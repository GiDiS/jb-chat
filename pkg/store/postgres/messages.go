package postgres

import (
	"context"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/store/postgres/utils"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type messagesPostgresStore struct {
	db *sqlx.DB
}

func NewMessagesPostgresStore(db *sqlx.DB) *messagesPostgresStore {
	return &messagesPostgresStore{
		db: db,
	}
}

func (s *messagesPostgresStore) Create(ctx context.Context, message models.Message) (models.MessageId, error) {
	var msgId models.MessageId
	err := utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO messages (cid, uid, pid, created, deleted, body, is_thread, likes) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8) RETURNING mid
		`, message.ChannelId, message.UserId, message.ParentId, message.Created, message.Deleted,
			message.Body, message.IsThread, packLikes(message.Likes),
		).Scan(&msgId)
		return err
	})
	return msgId, err
}

func (s *messagesPostgresStore) Delete(ctx context.Context, mid models.MessageId) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM messages WHERE mid = $1", mid)
		return err
	})
}

func (s *messagesPostgresStore) MarkAsThread(ctx context.Context, mid models.MessageId, isThread bool) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "UPDATE messages SET is_thread = $1 WHERE mid = $2", isThread, mid)
		return err
	})
}

func (s *messagesPostgresStore) Get(ctx context.Context, mid models.MessageId) (models.Message, error) {
	var messages, err = s.Find(ctx, store.MessagesSearchCriteria{
		Ids: []models.MessageId{mid}, Limits: models.Limits{Limit: 1},
	})
	if err != nil {
		return models.Message{}, err
	}
	if len(messages) > 0 {
		return messages[0], nil
	}
	return models.Message{}, store.ErrMessageNotFound
}

func (s *messagesPostgresStore) Find(ctx context.Context, filter store.MessagesSearchCriteria) (list []models.Message, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		query, args, buildErr := s.buildFind([]string{"*"}, filter, true)
		if buildErr != nil {
			return buildErr
		}

		query = tx.Rebind(query)
		list = make([]models.Message, 0)
		if queryErr := tx.SelectContext(ctx, &list, query, args...); queryErr != nil {
			return queryErr
		}

		return nil
	})
	return
}

func (s *messagesPostgresStore) Estimate(ctx context.Context, filter store.MessagesSearchCriteria) (uint64, error) {
	var total uint64
	err := utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		query, args, err := s.buildFind([]string{"count(*)"}, filter, false)
		if err != nil {
			return err
		}
		query = tx.Rebind(query)
		if err := s.db.GetContext(ctx, &total, query, args...); err != nil {
			return err
		}
		return nil
	})
	return total, err
}

func (s *messagesPostgresStore) buildFind(fields []string, filter store.MessagesSearchCriteria, useLimit bool) (string, []interface{}, error) {
	var (
		where = make([]string, 0)
		args  = make([]interface{}, 0)
		err   error
	)

	if len(filter.Ids) > 0 {
		where, args, err = utils.In(where, args, "mid", filter.Ids)
		if err != nil {
			return "", nil, err
		}
	}

	if filter.ChannelId > 0 {
		where = append(where, "cid = ?")
		args = append(args, filter.ChannelId)
	}

	if filter.ParentId > 0 {
		where = append(where, "pid = ?")
		args = append(args, filter.ParentId)
	}

	if filter.Search != "" {
		where = append(where, "body ILIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", filter.Search))
	}

	query := `
		SELECT ` + strings.Join(fields, ", ") + `
		FROM messages m
	`

	if len(where) > 0 {
		query += fmt.Sprintf("\tWHERE (%s)", strings.Join(where, ") \n\t\tAND ("))
	}

	if useLimit {
		query += "\n\tORDER BY created"
	}

	if useLimit && (filter.Limits.Offset > 0 || filter.Limits.Limit > 0) {
		query += fmt.Sprintf("\n\tLIMIT %d OFFSET %d", filter.Limits.Limit, filter.Limits.Offset)
	}

	return query, args, nil
}

func packLikes(likes []models.Uid) string {
	list := make([]string, 0, len(likes))
	for _, uid := range likes {
		list = append(list, uid.String())
	}
	return "[" + strings.Join(list, ",") + "]"
}

func unpackLikes(likes string) []models.Uid {
	list := make([]models.Uid, 0)
	likes = strings.Trim(likes, "[]")
	for _, s := range strings.Split(likes, ",") {
		uid, _ := strconv.Atoi(s)
		list = append(list, models.Uid(uid))
	}
	return list
}
