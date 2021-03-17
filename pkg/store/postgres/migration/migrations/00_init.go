package migrations

import (
	"database/sql"
	"fmt"
)

type Migration000Init struct{}

func (Migration000Init) Version() string {
	return "000Init"
}

func (m Migration000Init) Up(tx *sql.Tx) error {
	queries := make([]string, 0)
	queries = append(queries, m.upChannels(tx)...)
	queries = append(queries, m.upChannelMembers(tx)...)
	queries = append(queries, m.upMessages(tx)...)
	queries = append(queries, m.upSessions(tx)...)
	queries = append(queries, m.upUsers(tx)...)

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("query: %s, err: %w", q, err)
		}
	}

	return nil
}

func (Migration000Init) upChannels(tx *sql.Tx) []string {
	queries := make([]string, 0)
	queries = append(queries, `
		CREATE TABLE channels (
			cid INTEGER GENERATED ALWAYS AS IDENTITY,
			type VARCHAR NOT NULL,
			title VARCHAR NOT NULL,
			created TIMESTAMP NOT NULL,
			owner_uid INTEGER DEFAULT 0,
			last_msg_id INTEGER DEFAULT 0,
			last_msg_at TIMESTAMP WITH TIME ZONE NULL,
			members_count INTEGER default 0,
			PRIMARY KEY (cid)
		);
	`)
	queries = append(queries, `CREATE INDEX channels_owner_uid ON channels (owner_uid);`)
	queries = append(queries, `CREATE INDEX channels_title ON channels (title);`)

	return queries
}

func (Migration000Init) upChannelMembers(tx *sql.Tx) []string {
	queries := make([]string, 0)
	queries = append(queries, `
		CREATE TABLE channel_members (
			cid INTEGER NOT NULL,
			uid INTEGER NOT NULL,
			created TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (cid, uid)
		);
`)
	queries = append(queries, `CREATE INDEX channel_members_uid ON channel_members (uid);`)

	return queries
}

func (Migration000Init) upMessages(tx *sql.Tx) []string {
	queries := make([]string, 0)
	queries = append(queries, `
		CREATE TABLE messages (
			mid INTEGER GENERATED ALWAYS AS IDENTITY,
			cid INTEGER NOT NULL,
			uid INTEGER NOT NULL,
			pid INTEGER NOT NULL DEFAULT 0,
			created TIMESTAMP WITH TIME ZONE  NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted TIMESTAMP WITH TIME ZONE  NULL DEFAULT NULL,
			body TEXT,
			is_thread BOOLEAN,
			likes TEXT,
			PRIMARY KEY (mid)
		);
`)
	queries = append(queries, `CREATE INDEX messages_cid ON messages (cid);`)
	queries = append(queries, `CREATE INDEX messages_uid ON messages (uid);`)
	queries = append(queries, `CREATE INDEX messages_pid ON messages (pid);`)
	return queries
}

func (Migration000Init) upSessions(tx *sql.Tx) []string {
	queries := make([]string, 0)
	queries = append(queries, `
		CREATE TABLE SESSIONS (
			sid varchar NOT NULL ,
			uid integer NOT NULL ,
			service varchar NOT NULL DEFAULT '',
			app_id varchar NOT NULL DEFAULT '',
			app_token varchar NOT NULL DEFAULT '',
			token varchar NOT NULL DEFAULT '',
			expired BOOLEAN NOT NULL DEFAULT FALSE,
			is_online BOOLEAN NOT NULL DEFAULT FALSE,
			started TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
			PRIMARY KEY (sid)
		);
`)
	queries = append(queries, `CREATE INDEX sessions_uid ON sessions (uid);`)
	queries = append(queries, `CREATE INDEX sessions_is_online ON sessions (is_online);`)
	return queries
}

func (Migration000Init) upUsers(tx *sql.Tx) []string {
	queries := make([]string, 0)
	queries = append(queries, `
		CREATE TABLE users (
			uid INTEGER GENERATED ALWAYS AS IDENTITY,
			nickname varchar NOT NULL,
			title varchar NOT NULL,
			email varchar NOT NULL,
			avatar_url varchar NOT NULL,
			status varchar NOT NULL DEFAULT 'unknown',
			PRIMARY KEY (uid)
		);
`)
	queries = append(queries, `CREATE INDEX users_uid ON users (uid);`)
	queries = append(queries, `CREATE INDEX users_email ON users (email);`)
	queries = append(queries, `CREATE INDEX users_nickname ON users (nickname);`)
	queries = append(queries, `CREATE INDEX users_status ON users (status);`)
	return queries
}
