package migration

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func MustMigrate(db *sqlx.DB) {
	processor := newProcessor(db)
	if err := processor.ensureVersionsTable(); err != nil {
		log.Fatal(err)
	}

	toMigrate, err := processor.getMigrationsToMigrate()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("To migrate: %d", len(toMigrate))

	for _, m := range toMigrate {
		log.Infof("Migrating: %s", m.Version())

		tx, err := processor.db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		err = m.Up(tx)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		err = processor.markMigrationAsMigrated(m, tx)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		tx.Commit()
	}

	log.Infoln("Done migrating")
}
