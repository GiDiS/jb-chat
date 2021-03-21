package migration

import (
	"github.com/GiDiS/jb-chat/pkg/store/postgres/migration/migrations"
)

var migrationsList []MigrationInterface

func init() {
	migrationsList = []MigrationInterface{
		migrations.Migration000Init{},
		migrations.Migration001AddChanMsgsCount{},
		migrations.Migration001AddMemberLastSeen{},
	}
}
