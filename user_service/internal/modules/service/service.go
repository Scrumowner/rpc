package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"user/internal/models"
	"user/internal/modules/storage"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(dbx *sqlx.DB) *UserService {
	return &UserService{storage: storage.NewUserStorage(dbx)}
}

func (s *UserService) GetList(ctx context.Context) ([]models.User, error) {
	users, err := s.storage.GetListFromStorage(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (s *UserService) GetUser(user *models.User) (*models.User, error) {
	ctx := context.Background()
	resp, err := s.storage.GetUserFromStorage(ctx, user)
	if err != nil {
		return &models.User{}, err
	}
	return resp, nil
}

func (s *UserService) SetUser(user *models.User) error {
	ctx := context.Background()
	err := s.storage.SetUserToStorage(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
