package storage

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	cryptography "hugoproxy-main/proxy/tools"
)

var TokenAuth = jwtauth.New("HS256", []byte("mysecretkey"), nil)

type Userer interface {
	Create(User) error
	GetToken(User) (string, error)
}

func NewStorage() Userer {
	return &Users{usr: make(map[int]User)}
}
func (u *Users) Create(user User) error {
	hashedPass, err := cryptography.GeneratePassword(user.Password)
	if err != nil {
		return err
	}
	if _, ok := u.usr[user.ID]; ok {
		return fmt.Errorf("User with this id are registred")
	}
	user.Password = hashedPass
	u.usr[user.ID] = user
	return nil

}

func (u *Users) GetToken(user User) (string, error) {
	if usr, ok := u.usr[user.ID]; !ok && cryptography.ComparePassowrd(usr.Password, user.Password) != true {
		return "", fmt.Errorf("Invalid password")
	}
	_, token, _ := TokenAuth.Encode(map[string]interface{}{"id": user.Password})
	return fmt.Sprintf(`{"token": "Bearer %s"}`, token), nil
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users struct {
	usr map[int]User
}
