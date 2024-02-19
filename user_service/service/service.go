package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user/models"
	"user/storage"
)

type UserServiceInt interface {
	GetUser(email string) (models.User, error)
	SetUser(user models.User) error
	GetList() ([]models.User, error)
}
type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(dbx *sqlx.DB) *UserService {
	return &UserService{storage: storage.NewUserStorage(dbx)}
}

func (s *UserService) GetUser(email string) (models.User, error) {
	ctx := context.Background()
	user, err := s.storage.GetUserFromStorage(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) SetUser(user models.User) error {
	ctx := context.Background()
	err := s.storage.SetUserToStorage(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetList() ([]models.User, error) {
	ctx := context.Background()
	users, err := s.storage.GetListFromStorage(ctx)
	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}
