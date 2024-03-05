package adapter

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"user/internal/models"
)

type SqlAdapter struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func NewSqlAdapter(db *sqlx.DB) *SqlAdapter {
	var builder squirrel.StatementBuilderType
	return &SqlAdapter{
		db: db,
		sq: builder,
	}
}

func (s *SqlAdapter) CreateUser(user *models.User) error {
	query := fmt.Sprintf("INSERT INTO users (email,phone,password) VALUES ($1,$2,$3)")
	_, err := s.db.Exec(query, user.Email, user.Phone, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlAdapter) GetUser(user *models.User) (*models.User, error) {
	query := fmt.Sprintf("SELECT email,phone,password FROM users WHERE email=$1 AND phone=$2;")
	resp := []models.User{}
	err := s.db.Select(&resp, query, user.Email, user.Phone)
	if err != nil {
		return nil, fmt.Errorf("Can't exec query SELECT", err)
	}
	us := resp[0]
	return &us, nil

}

func (s *SqlAdapter) GetList() ([]models.User, error) {
	query := fmt.Sprintf("SELECT * FROM users ;")
	resUsers := []models.User{}
	err := s.db.Select(&resUsers, query)
	if err != nil {
		return nil, fmt.Errorf("Can't exec query SELECT ALL", err)
	}
	return resUsers, err
}
