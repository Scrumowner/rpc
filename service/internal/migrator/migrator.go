package migrator

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"rpc/service/internal/db"
	"strings"
)

type Migrator struct {
	sql     *sqlx.DB
	scanner *db.Scanner
}

func NewMigrator(sql *sqlx.DB) *Migrator {
	return &Migrator{
		sql:     sql,
		scanner: db.NewScanner(),
	}
}

func (m *Migrator) Migrate(tables ...interface{}) {
	for _, table := range tables {
		info := m.scanner.GetFieldsTypes(table)
		var query strings.Builder
		query.WriteString("CREATE TABLE IF NOT EXISTS ")
		query.WriteString(reflect.TypeOf(table).Name())
		query.WriteString(" (")
		for j, column := range info.Names {
			tp := fmt.Sprintf("%s ", column)
			for _, t := range info.Types[column] {
				tp += fmt.Sprintf("%s ", t)
			}
			query.WriteString(tp)
			if j < len(info.Names)-1 {
				query.WriteString(", ")
			}
		}
		query.WriteString(")")
		fmt.Println(query.String())
	}
}
