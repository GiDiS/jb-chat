package migrations

import "database/sql"

type Migration001AddChanMsgsCount struct{}

func (Migration001AddChanMsgsCount) Version() string {
	return "001AddChanMsgsCount"
}

func (m Migration001AddChanMsgsCount) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE channels ADD COLUMN messages_count INTEGER DEFAULT 0 NOT NULL 
	`)
	return err
}
