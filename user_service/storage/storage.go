package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user/models"
)

type UserStorageInt interface {
	GetUserFromStorage(ctx context.Context, email string) (models.User, error)
	SetUserToStorage(ctx context.Context, user models.User) error
	GetListFromStorage(ctx context.Context) ([]models.User, error)
}

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(dbx *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: dbx,
	}
}

func (s *UserStorage) GetUserFromStorage(ctx context.Context, email string) (models.User, error) {
	user := models.User{}
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserStorage) SetUserToStorage(ctx context.Context, user models.User) error {
	_, err := s.db.NamedExecContext(
		ctx,
		"INSERT INTO users (email, password) VALUES (:email, :password)",
		map[string]interface{}{
			"email":    user.Email,
			"password": user.Password})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetListFromStorage(ctx context.Context) ([]models.User, error) {
	users := []models.User{}
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM users ORDER BY email ASC")
	if err != nil {
		return []models.User{}, err
	}
	return users, nil

}
