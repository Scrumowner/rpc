package service

import (
	"hugoproxy-main/proxy/storage"
)

type AuthServiceer interface {
	Register(r UserService) (string, error)
	Login(r UserService) (string, error)
}

type AuthService struct {
	storage storage.Userer
}

func NewAuthService(storage storage.Userer) AuthServiceer {
	return &AuthService{storage: storage}
}

func (auth *AuthService) Register(r UserService) (string, error) {
	var user storage.User = storage.User{ID: r.ID, Name: r.Name, Phone: r.Phone, Email: r.Email, Password: r.Password}
	err := auth.storage.Create(user)
	if err != nil {
		return "", err
	}
	return "Registration successful", nil

}

func (auth *AuthService) Login(r UserService) (string, error) {
	var user storage.User = storage.User{ID: r.ID, Name: r.Name, Phone: r.Phone, Email: r.Email, Password: r.Password}
	token, err := auth.storage.GetToken(user)
	if err != nil {
		return "", err
	}
	return token, nil

}

type UserService struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
