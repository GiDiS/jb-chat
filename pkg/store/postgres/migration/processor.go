package migration

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Processor struct {
	migrations map[string]MigrationInterface
	db         *sqlx.DB
}

func newProcessor(db *sqlx.DB) Processor {
	p := Processor{
		db:         db,
		migrations: make(map[string]MigrationInterface),
	}

	for _, migration := range migrationsList {
		p.addMigration(migration)
	}

	return p
}

func (p *Processor) addMigration(m MigrationInterface) {
	if _, ok := p.migrations[m.Version()]; ok {
		log.Fatalf("can not add migration with version %s: is already in map", m.Version())
	}

	p.migrations[m.Version()] = m
}

func (p Processor) ensureVersionsTable() error {
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version varchar(100) not null,
			primary key (version)
		) 
	`)

	return err
}

func (p Processor) getMigratedVersions() ([]string, error) {
	var versions []string
	err := p.db.Select(&versions, `SELECT version FROM migrations`)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (p Processor) getMigrationsToMigrate() ([]MigrationInterface, error) {
	migrated, err := p.getMigratedVersions()
	if err != nil {
		return nil, err
	}

	var allIDs []string
	for id := range p.migrations {
		allIDs = append(allIDs, id)
	}

	idsToMigrate := diff(allIDs, migrated)
	if len(idsToMigrate) == 0 {
		return nil, nil
	}

	sort.Strings(idsToMigrate)
	var out []MigrationInterface
	for _, id := range idsToMigrate {
		out = append(out, p.migrations[id])
	}

	return out, nil
}

func (p Processor) isMigrated(m MigrationInterface) (bool, error) {
	var count int

	err := p.db.QueryRow(`SELECT COUNT(*) FROM migrations WHERE version = $1`, m.Version()).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p Processor) markMigrationAsMigrated(m MigrationInterface, tx *sql.Tx) error {
	migrated, err := p.isMigrated(m)
	if err != nil {
		return err
	}

	if migrated {
		return nil
	}

	_, err = tx.Exec(`INSERT INTO migrations VALUES ($1)`, m.Version())
	return err
}

func diff(list1, list2 []string) []string {
	var out []string

	for _, item1 := range list1 {
		found := false
		for _, item2 := range list2 {
			if item1 == item2 {
				found = true
				break
			}
		}

		if !found {
			out = append(out, item1)
		}
	}

	return out
}
