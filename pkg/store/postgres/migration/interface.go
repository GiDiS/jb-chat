package migration

import "database/sql"

//noinspection GoNameStartsWithPackageName
type MigrationInterface interface {
	Version() string
	Up(tx *sql.Tx) error
}
