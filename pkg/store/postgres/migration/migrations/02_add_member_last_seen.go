package migrations

import "database/sql"

type Migration001AddMemberLastSeen struct{}

func (Migration001AddMemberLastSeen) Version() string {
	return "001AddMemberLastSeen"
}

func (m Migration001AddMemberLastSeen) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE channel_members 
		    ADD COLUMN last_seen_id INTEGER DEFAULT 0 NOT NULL, 
		    ADD COLUMN last_seen_at TIMESTAMP WITH TIME ZONE NULL 
	`)
	return err
}
