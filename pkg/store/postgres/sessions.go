package postgres

import (
	"context"
	"database/sql"
	"github.com/GiDiS/jb-chat/pkg/models"
	"github.com/GiDiS/jb-chat/pkg/store/postgres/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

type sessionsPostgresStore struct {
	db *sqlx.DB
}

func NewSessionsPostgresStore(db *sqlx.DB) *sessionsPostgresStore {
	return &sessionsPostgresStore{
		db: db,
	}
}

func (s *sessionsPostgresStore) GetUserSessions(ctx context.Context, uid models.Uid) (list []models.Session, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		list = make([]models.Session, 0)
		if err := tx.SelectContext(ctx, &list, "SELECT * FROM sessions WHERE uid = $1", uid); err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *sessionsPostgresStore) GetSessions(ctx context.Context) (list []models.Session, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		list = make([]models.Session, 0)
		if err := tx.SelectContext(ctx, &list, "SELECT * FROM sessions"); err != nil {
			return err
		}
		return nil
	})
	return
}

func (s *sessionsPostgresStore) GetSession(ctx context.Context, sid string) (sess models.Session, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		dbErr := tx.GetContext(ctx, &sess, "SELECT * FROM sessions WHERE sid = $1", sid)
		if dbErr == sql.ErrNoRows {
			now := time.Now()
			sess = models.Session{
				SessionId: sid,
				Started:   now,
				Updated:   now,
				IsOnline:  false,
			}
			return nil
		}
		return dbErr
	})
	return
}

func (s *sessionsPostgresStore) SetSession(ctx context.Context, sid string, sess models.Session) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var _sid string
		err := tx.QueryRow("SELECT sid FROM sessions WHERE sid = $1", sid).Scan(&_sid)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			if sess.Started.IsZero() {
				sess.Started = time.Now()
			}
			if sess.Updated.IsZero() {
				sess.Updated = sess.Started
			}

			_, err = tx.NamedExecContext(ctx, `
				INSERT INTO sessions
				(sid, uid, service, app_id, app_token, token, expired, started, expires, is_online) 
				VALUES 
			    (:sid, :uid, :service, :app_id, :app_token, :token, :expired, :started, :expires, :is_online)
			`, sess)
		} else {
			sess.Updated = time.Now()
			_, err = tx.NamedExecContext(ctx, `
				UPDATE sessions SET
                    uid = :uid, 
                    service = :service, 
                    app_id = :app_id, 
                    app_token = :app_token, 
                    token = :token, 
                    expired = :expired, 
                    started = :started, 
                    expires = :expires,
                    is_online = :is_online
				WHERE 				
					sid = :sid 
			`, sess)
		}
		return err
	})
}

func (s *sessionsPostgresStore) ClearSession(ctx context.Context, sid string) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sessions WHERE sid = $1", sid)
		return err
	})
}

func (s *sessionsPostgresStore) SetOnline(ctx context.Context, sid string, uid models.Uid, isOnline bool) error {
	return utils.InWriteTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "UPDATE sessions SET is_online = $1 WHERE sid = $2", isOnline, sid)
		return err
	})
}

func (s *sessionsPostgresStore) IsOnline(ctx context.Context, uid models.Uid) (isOnline bool, err error) {
	sess, err := s.GetOnline(ctx, uid)
	if err != nil {
		return false, err
	}
	return len(sess) > 0, nil
}

func (s *sessionsPostgresStore) GetOnline(ctx context.Context, uid models.Uid) (list []string, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		sess := make([]struct {
			Sid string `db:"sid"`
		}, 0)

		queryErr := tx.SelectContext(ctx, &sess, "SELECT sid FROM sessions WHERE uid = $1 and is_online = true", uid)
		if queryErr != nil {
			return queryErr
		}
		list = make([]string, 0, len(sess))
		for _, s := range sess {
			list = append(list, s.Sid)
		}
		return nil
	})
	return
}

func (s *sessionsPostgresStore) GetAllOnline(ctx context.Context) (users map[models.Uid][]string, err error) {
	err = utils.InReadOnlyTransactionX(s.db, ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		sess := make([]struct {
			Sid string     `db:"sid"`
			Uid models.Uid `db:"uid"`
		}, 0)

		queryErr := tx.SelectContext(ctx, &sess, "SELECT sid, uid FROM sessions WHERE is_online = true")
		if queryErr != nil {
			return queryErr
		}
		users = make(map[models.Uid][]string, len(sess))
		for _, s := range sess {
			if _, ok := users[s.Uid]; !ok {
				users[s.Uid] = make([]string, 0)
			}
			users[s.Uid] = append(users[s.Uid], s.Sid)
		}
		return nil
	})
	return
}
