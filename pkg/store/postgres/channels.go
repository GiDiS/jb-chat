package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/store/postgres/utils"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type channelsPostgresStore struct {
	db *sqlx.DB
}

func NewChannelsPostgresStore(db *sqlx.DB) *channelsPostgresStore {
	return &channelsPostgresStore{
		db: db,
	}
}

func (s *channelsPostgresStore) CreateDirect(ctx context.Context, uidA models.Uid, uidB models.Uid) (models.ChannelId, error) {
	var cid = models.NoChannel
	err := utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		channel := models.Channel{
			Title:        models.DirectTitle(uidA, uidB),
			LastMsgId:    models.NoMessage,
			MembersCount: 0,
			Type:         models.ChannelTypeDirect,
		}

		var dbErr error
		cid, dbErr = s.create(ctx, tx, channel)
		if dbErr != nil {
			return fmt.Errorf("channel create failed: %w", dbErr)
		}

		if err := s.join(ctx, tx, cid, uidA); err != nil {
			return err
		}
		if err := s.join(ctx, tx, cid, uidB); err != nil {
			return err
		}
		return nil
	})

	return cid, err
}

func (s *channelsPostgresStore) GetDirect(ctx context.Context, uidA models.Uid, uidB models.Uid) (models.ChannelId, error) {
	channels, err := s.Find(ctx, store.ChannelsSearchCriteria{
		Type:  models.ChannelTypeDirect,
		Title: models.DirectTitle(uidA, uidB),
	})
	if err != nil {
		return models.NoChannel, err
	}
	if len(channels) > 0 {
		return channels[0].Cid, nil
	}
	return models.NoChannel, nil
}

func (s *channelsPostgresStore) CreatePublic(ctx context.Context, authorUid models.Uid, title string) (models.ChannelId, error) {
	var cid = models.NoChannel
	err := utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		channel := models.Channel{
			Title:        title,
			LastMsgId:    models.NoMessage,
			MembersCount: 0,
			Type:         models.ChannelTypePublic,
		}

		var dbErr error
		cid, dbErr = s.create(ctx, tx, channel)
		if dbErr != nil {
			return fmt.Errorf("channel create failed: %w", dbErr)
		}

		if authorUid != models.NoUser {
			if err := s.join(ctx, tx, cid, authorUid); err != nil {
				return fmt.Errorf("join failed: %w", err)
			}
		}

		return nil
	})

	return cid, err
}

func (s *channelsPostgresStore) Delete(ctx context.Context, cid models.ChannelId) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.delete(ctx, tx, cid); err != nil {
			return fmt.Errorf("channel delete failed: %w", err)
		}
		return nil
	})
}

func (s *channelsPostgresStore) Get(ctx context.Context, cid models.ChannelId) (models.Channel, error) {
	var channels, err = s.Find(ctx, store.ChannelsSearchCriteria{
		ChannelIds: []models.ChannelId{cid}, Limits: models.Limits{Limit: 1},
	})
	if err != nil {
		return models.Channel{}, err
	}
	if len(channels) > 0 {
		return channels[0], nil
	}
	return models.Channel{}, store.ErrChanNotFound
}

func (s *channelsPostgresStore) Find(ctx context.Context, filter store.ChannelsSearchCriteria) (list []models.Channel, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		fields := []string{"cid", "title", "type", "owner_uid", "last_msg_id", "last_msg_at", "members_count"}
		query, args, buildErr := s.buildFind(fields, filter, true)
		if err != nil {
			return buildErr
		}
		query = tx.Rebind(query)
		list = make([]models.Channel, 0)
		if queryErr := tx.SelectContext(ctx, &list, query, args...); queryErr != nil {
			return fmt.Errorf("channels find failed: %v", queryErr)
		}

		return nil
	})
	return
}

func (s *channelsPostgresStore) Estimate(ctx context.Context, filter store.ChannelsSearchCriteria) (uint64, error) {
	var total uint64
	err := utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		query, args, err := s.buildFind([]string{"count(*)"}, filter, false)
		if err != nil {
			return err
		}
		query = tx.Rebind(query)
		if err := s.db.GetContext(ctx, &total, query, args...); err != nil {
			return fmt.Errorf("channels estimate failed: %v", err)
		}
		return nil
	})
	return total, err
}

func (s *channelsPostgresStore) buildFind(fields []string, filter store.ChannelsSearchCriteria, useLimit bool) (string, []interface{}, error) {
	var (
		where = make([]string, 0)
		args  = make([]interface{}, 0)
		err   error
	)

	if len(filter.ChannelIds) > 0 {
		where, args, err = utils.In(where, args, "cid", filter.ChannelIds)
		if err != nil {
			return "", nil, err
		}
	}

	if filter.Title != "" {
		where = append(where, "title = ?")
		args = append(args, filter.Title)
	}

	if filter.Type != models.ChannelTypeUnknown {
		where = append(where, "type = ?")
		args = append(args, filter.Type)
	}

	if filter.HasMember > 0 {
		where = append(where, "EXISTS (SELECT m.cid FROM channel_members WHERE c.cid=m.cid AND m.uid = ?)")
		args = append(args, filter.HasMember)
	}

	query := `
		SELECT ` + strings.Join(fields, ", ") + `
		FROM channels c
	`

	if len(where) > 0 {
		query += fmt.Sprintf("\tWHERE (%s)", strings.Join(where, ") \n\t\tAND ("))
	}

	if useLimit {
		query += "\n\tORDER BY title"
	}

	if useLimit && (filter.Limits.Offset > 0 || filter.Limits.Limit > 0) {
		query += fmt.Sprintf("\n\tLIMIT %d OFFSET %d", filter.Limits.Limit, filter.Limits.Offset)
	}

	return query, args, nil
}

func (s *channelsPostgresStore) Members(ctx context.Context, cid models.ChannelId) ([]models.Uid, error) {
	var uids = make([]models.Uid, 0)
	err := utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := tx.SelectContext(ctx, &uids, `
			SELECT uid FROM channel_members WHERE cid = $1
		`, cid)
		if err != nil {
			return fmt.Errorf("channel members failed: %v", err)
		}
		return nil
	})
	return uids, err
}

func (s *channelsPostgresStore) MemberOf(ctx context.Context, uid models.Uid) ([]models.ChannelId, error) {
	var channels = make([]models.ChannelId, 0)
	err := utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := tx.SelectContext(ctx, &channels, `
			SELECT cid FROM channel_members WHERE uid = $1
		`, uid)
		if err != nil {
			return fmt.Errorf("channel member of failed: %v", err)
		}
		return nil
	})
	return channels, err
}

func (s *channelsPostgresStore) IsMember(ctx context.Context, cid models.ChannelId, uid models.Uid) (isMember bool, err error) {
	if err := utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		isMember, err = s.isMember(ctx, tx, cid, uid)
		return err
	}); err != nil {
		return false, err
	}
	return
}

func (s *channelsPostgresStore) Join(ctx context.Context, cid models.ChannelId, uid models.Uid) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.join(ctx, tx, cid, uid)
	})
}

func (s *channelsPostgresStore) Leave(ctx context.Context, cid models.ChannelId, uid models.Uid) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.leave(ctx, tx, cid, uid)
	})
}

func (s *channelsPostgresStore) create(ctx context.Context, tx *sqlx.Tx, c models.Channel) (models.ChannelId, error) {
	var cid = models.NoChannel
	if c.Created.IsZero() {
		c.Created = time.Now()
	}
	err := tx.QueryRowContext(ctx, `
		INSERT INTO channels (title, type, created, owner_uid) VALUES ($1, $2, $3, $4) RETURNING cid
		`, c.Title, c.Type, c.Created, c.OwnerUid,
	).Scan(&cid)
	return cid, err
}

func (s *channelsPostgresStore) delete(ctx context.Context, tx *sqlx.Tx, cid models.ChannelId) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM channels WHERE cid = $1", cid)
	return err
}

func (s *channelsPostgresStore) isMember(ctx context.Context, tx *sqlx.Tx, cid models.ChannelId, uid models.Uid) (bool, error) {
	var _uid = models.NoUser
	err := tx.GetContext(ctx, &_uid, `
		SELECT uid FROM channel_members WHERE (cid,uid) = ($1, $2)
		`, cid, uid,
	)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return _uid == uid, err
}

func (s *channelsPostgresStore) join(ctx context.Context, tx *sqlx.Tx, cid models.ChannelId, uid models.Uid) error {
	if isMember, err := s.isMember(ctx, tx, cid, uid); err != nil {
		return err
	} else if isMember {
		return nil
	}
	_, err := tx.ExecContext(ctx,
		"INSERT INTO channel_members (cid, uid, created) values ($1, $2, $3)",
		cid, uid, time.Now(),
	)
	return err
}

func (s *channelsPostgresStore) leave(ctx context.Context, tx *sqlx.Tx, cid models.ChannelId, uid models.Uid) error {
	if isMember, err := s.isMember(ctx, tx, cid, uid); err != nil {
		return err
	} else if !isMember {
		return nil
	}
	_, err := tx.ExecContext(ctx,
		"DELETE FROM channel_members WHERE (cid, uid) = ($1, $2)",
		cid, uid,
	)
	return err
}
