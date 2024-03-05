package migrator

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"user/internal/models"
)

type Migrator struct {
	db *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{
		db: db,
	}
}

func (m *Migrator) Migrate(users ...models.Userer) {
	for _, user := range users {
		tableName := user.TableName()
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (email TEXT UNIQUE NOT NULL , phone TEXT UNIQUE NOT NULL ,password TEXT NOT NULL );", tableName)
		m.db.MustExec(query)

	}
}
