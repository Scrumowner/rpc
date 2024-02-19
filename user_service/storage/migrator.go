package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"user/models"
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
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (email TEXT ,password TEXT);", tableName)
		m.db.MustExec(query)
	}
}
