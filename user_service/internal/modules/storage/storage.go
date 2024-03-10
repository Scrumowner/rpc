package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"user/internal/infra/db/adapter"
	"user/internal/models"
)

type UserStorageInt interface {
	GetUserFromStorage(ctx context.Context, email string) (models.User, error)
	SetUserToStorage(ctx context.Context, user models.User) error
	GetListFromStorage(ctx context.Context) ([]models.User, error)
}

type UserStorage struct {
	adapter *adapter.SqlAdapter
}

func NewUserStorage(dbx *sqlx.DB) *UserStorage {
	return &UserStorage{
		adapter: adapter.NewSqlAdapter(dbx),
	}
}

func (s *UserStorage) SetUserToStorage(ctx context.Context, user *models.User) error {
	err := s.adapter.CreateUser(user)
	if err != nil {
		return fmt.Errorf("Error from storage", err)
	}
	return nil
}
func (s *UserStorage) GetUserFromStorage(ctx context.Context, user *models.User) (*models.User, error) {
	userFromStorage, err := s.adapter.GetUser(user)
	if err != nil {
		return nil, fmt.Errorf("Error from storage", err)
	}

	return userFromStorage, nil
}

func (s *UserStorage) GetListFromStorage(ctx context.Context) ([]models.User, error) {
	users, err := s.adapter.GetList()
	if err != nil {
		return nil, fmt.Errorf("Error from storage", err)
	}
	return users, nil
}
