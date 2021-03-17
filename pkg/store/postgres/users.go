package postgres

import (
	"context"
	"fmt"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store"
	"github.com/GiDiS/jb-chat/pkg/store/postgres/utils"
	"github.com/jmoiron/sqlx"
	"strings"
)

type usersPostgresStore struct {
	db *sqlx.DB
}

func NewUsersPostgresStore(db *sqlx.DB) *usersPostgresStore {
	return &usersPostgresStore{
		db: db,
	}
}

func (s *usersPostgresStore) Register(ctx context.Context, user models.User) (uid models.Uid, err error) {
	err = utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		uid, err = s.create(ctx, tx, user)
		return err
	})
	return
}

func (s *usersPostgresStore) Save(ctx context.Context, user models.User) (uid models.Uid, err error) {
	err = utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if user.UserId > 0 {
			uid, err = s.update(ctx, tx, user)
		} else {
			uid, err = s.create(ctx, tx, user)
		}
		return err
	})
	return
}

func (s *usersPostgresStore) SetStatus(ctx context.Context, uid models.Uid, status models.UserStatus) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "UPDATE users SET status = $1 WHERE uid = $2 ", status, uid)
		return err
	})
}

func (s *usersPostgresStore) GetByEmail(ctx context.Context, email string) (models.User, error) {
	return s.FindOne(ctx, store.UserSearchCriteria{Emails: []string{email}})
}

func (s *usersPostgresStore) GetByUid(ctx context.Context, uid models.Uid) (models.User, error) {
	return s.FindOne(ctx, store.UserSearchCriteria{Uids: []models.Uid{uid}})
}

func (s *usersPostgresStore) FindOne(ctx context.Context, filter store.UserSearchCriteria) (models.User, error) {
	matched, err := s.Find(ctx, filter)
	if err != nil {
		return models.User{}, err
	}
	if len(matched) > 0 {
		return matched[0], nil
	}
	return models.User{}, store.ErrUserNotFound
}

func (s *usersPostgresStore) Find(ctx context.Context, filter store.UserSearchCriteria) (list []models.User, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		list, err = s.find(ctx, tx, filter)
		return err
	})
	return
}

func (s *usersPostgresStore) FindActive(ctx context.Context, limits models.Limits) ([]models.User, error) {
	return s.Find(ctx, store.UserSearchCriteria{
		Statuses: []models.UserStatus{models.UserStatusOnline, models.UserStatusAway},
		Limits:   limits,
	})
}

func (s *usersPostgresStore) Estimate(ctx context.Context, filter store.UserSearchCriteria) (uint64, error) {
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

func (s *usersPostgresStore) create(ctx context.Context, tx *sqlx.Tx, user models.User) (models.Uid, error) {
	existed, err := s.FindOne(ctx, store.UserSearchCriteria{Emails: []string{user.Email}})
	if err == nil {
		return existed.UserId, fmt.Errorf("user: %s, err: %w", user.Email, store.ErrUserAlreadyRegistered)
	} else if err != store.ErrUserNotFound {
		return models.NoUser, err
	}

	var uid models.Uid
	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (nickname, title, email, avatar_url) VALUES ($1, $2, $3, $4) RETURNING uid`,
		user.Nickname, user.Title, user.Email, user.AvatarUrl,
	).Scan(&uid)

	return uid, err
}

func (s *usersPostgresStore) update(ctx context.Context, tx *sqlx.Tx, user models.User) (models.Uid, error) {
	existed, err := s.FindOne(ctx, store.UserSearchCriteria{Emails: []string{user.Email}})
	if err == store.ErrUserNotFound {
		return models.NoUser, err
	} else if err != nil {
		return models.NoUser, err
	} else if existed.UserId == models.NoUser {
		return models.NoUser, store.ErrUserNotFound
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE users u SET nickname = $1, title = $2, avatar_url = $3 WHERE uid = $4",
		user.Nickname, user.Title, user.AvatarUrl, user.UserId,
	)

	return user.UserId, err
}

func (s *usersPostgresStore) find(ctx context.Context, tx *sqlx.Tx, filter store.UserSearchCriteria) (list []models.User, err error) {
	query, args, err := s.buildFind([]string{"uid", "nickname", "title", "email", "avatar_url"}, filter, true)
	if err != nil {
		return
	}

	query = tx.Rebind(query)

	list = make([]models.User, 0)
	err = tx.SelectContext(ctx, &list, query, args...)
	return list, err
}

func (s *usersPostgresStore) buildFind(fields []string, filter store.UserSearchCriteria, useLimit bool) (string, []interface{}, error) {
	var (
		where = make([]string, 0)
		args  = make([]interface{}, 0)
		err   error
	)

	if len(filter.Uids) > 0 {
		where, args, err = utils.In(where, args, "uid", filter.Uids)
		if err != nil {
			return "", nil, err
		}
	}
	if len(filter.Nicknames) > 0 {
		where, args, err = utils.In(where, args, "nickname", filter.Nicknames)
		if err != nil {
			return "", nil, err
		}
	}

	if len(filter.Emails) > 0 {
		where, args, err = utils.In(where, args, "email", filter.Emails)
		if err != nil {
			return "", nil, err
		}
	}

	if len(filter.Statuses) > 0 {
		where, args, err = utils.In(where, args, `status`, filter.Statuses)
		if err != nil {
			return "", nil, err
		}
	}

	if filter.WithAvatars {
		where = append(where, "coalesce(avatar_url, \"\") != \"\"")
	}

	query := `
		SELECT ` + strings.Join(fields, ", ") + `
		FROM users u
	`

	if len(where) > 0 {
		query += fmt.Sprintf("\tWHERE (%s)", strings.Join(where, ") \n\t\tAND ("))
	}

	if useLimit {
		query += "\n\tORDER BY status !='online', title, nickname"
	}

	if useLimit && (filter.Limits.Offset > 0 || filter.Limits.Limit > 0) {
		query += fmt.Sprintf("\n\tLIMIT %d OFFSET %d", filter.Limits.Limit, filter.Limits.Offset)
	}

	return query, args, nil
}
