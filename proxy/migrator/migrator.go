package migrator

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Migratorer interface {
	Migrate() error
}

type Migrator struct {
	db           *sqlx.DB
	sqlGenerator SQLGenerator
}

func NewMigrator(db *sqlx.DB) Migratorer {
	return &Migrator{
		db:           db,
		sqlGenerator: NewSQLiteGenerator(),
	}
}

func (m *Migrator) Migrate() error {
	var querys []string = []string{
		"CREATE TABLE IF NOT EXISTS address (query TEXT, result TEXT, lat TEXT , lon TEXT)",
		"CREATE TABLE IF NOT EXISTS geo (lat TEXT , lng TEXT, result TEXT, r_lat TEXT , r_lon TEXT )",
	}
	for _, query := range querys {
		_, _ = m.db.Exec(query)
	}
	return nil
}
