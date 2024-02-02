package migrator

import (
	"fmt"
	_ "github.com/lib/pq"
	"hugoproxy-main/proxy/adapter"
)

type SQLGenerator interface {
	CreateTableSQL(table adapter.Tabler) string
}

type SQLiteGenerator struct{}

func NewSQLiteGenerator() SQLGenerator {
	return &SQLiteGenerator{}
}

func (sg *SQLiteGenerator) CreateTableSQL(table adapter.Tabler) string {
	tableName := table.TableName()
	tableFields := table.TableFieldAndType()

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, tableFields)
}
