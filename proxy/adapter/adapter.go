package adapter

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Adapterer interface {
	BuildSelect(tableName string, condition Condition, fields ...string) (string, []interface{}, error)
	Create(ctx context.Context, entity Tabler, model interface{}) error
	Update(ctx context.Context, entity Tabler, model interface{}, condition Condition, opts ...interface{}) error
	Delete(ctx context.Context, entity Tabler, condition Condition, opts ...interface{}) error
}

func NewSqlAdapter(db *sqlx.DB) *SqlAdapter {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &SqlAdapter{Db: db, sqlBuilder: builder}
}

type SqlAdapter struct {
	sqlBuilder squirrel.StatementBuilderType
	Db         *sqlx.DB
}

type Condition struct {
	Equal       map[string]interface{}
	NotEqual    map[string]interface{}
	Order       []*Order
	LimitOffset *LimitOffset
	ForUpdate   bool
	Upsert      bool
}

type Order struct {
	Field string
	Asc   bool
}

type LimitOffset struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

func GetOrderDirection(asc bool) string {
	if asc {
		return "ASC"
	}
	return "DESC"

}

func (s *SqlAdapter) BuildSelect(tableName string, condition Condition, fields ...string) (string, []interface{}, error) {
	query := s.sqlBuilder.Select("*").From(tableName)
	if condition.Equal != nil {
		for column, val := range condition.Equal {
			query = query.Where(squirrel.Eq{column: val})
		}
	}
	if condition.NotEqual != nil {
		for column, val := range condition.NotEqual {
			query = query.Where(squirrel.NotEq{column: val})
		}
	}
	if condition.Order != nil {
		for _, order := range condition.Order {
			cond := GetOrderDirection(order.Asc)
			query = query.OrderBy(fmt.Sprintf("%s %s", order.Field, cond))
		}
	}
	if condition.LimitOffset != nil {
		query = query.Limit(uint64(condition.LimitOffset.Limit)).Offset(uint64(condition.LimitOffset.Offset))
	}
	if condition.ForUpdate {
		sql, args, err := query.ToSql()
		if err != nil {
			return "", nil, fmt.Errorf("Cant build select qurey to database ")
		}
		return sql + fmt.Sprintf(" FOR UPDATE"), args, err
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("Cant build select qurey to database ")
	}
	return sql, args, err
}

func (s *SqlAdapter) Create(ctx context.Context, entity Tabler, model interface{}) error {
	query := s.sqlBuilder.Insert(entity.TableName())
	info := GetStructInfo(&model)
	columns := info.Fields
	values := info.Pointers
	query = query.Columns(columns...).Values(values...)
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("Can't build insert query to database")
	}

	_, err = s.Db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Can't exec creat query to db")
	}
	return nil

}

func (s *SqlAdapter) Update(ctx context.Context, entity Tabler, model interface{}, condition Condition, opts ...interface{}) error {
	query := s.sqlBuilder.Update(entity.TableName())
	info := GetStructInfo(&model)
	columns := info.Fields
	values := info.Pointers
	for i := range columns {
		query = query.Set(columns[i], values[i])
	}
	if condition.Equal != nil {
		for column, val := range condition.Equal {
			query = query.Where(squirrel.Eq{column: val})
		}
	}
	if condition.NotEqual != nil {
		for column, val := range condition.NotEqual {
			query = query.Where(squirrel.NotEq{column: val})
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("Can't build update query to database: %v", err)
	}

	_, err = s.Db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Error executing update query: %v", err)
	}

	return nil
}

func (s *SqlAdapter) Delete(ctx context.Context, entity Tabler, condition Condition, opts ...interface{}) error {
	query := s.sqlBuilder.Update(entity.TableName())

	if condition.Equal != nil {
		for column, val := range condition.Equal {
			query = query.Where(squirrel.Eq{column: val})
		}
	}
	query = query.Set("deleted_at", time.Now().String())

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("Can't build delete query to database: %v", err)
	}

	_, err = s.Db.ExecContext(ctx, sql, args...)
	return err
}
