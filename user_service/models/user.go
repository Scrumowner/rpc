package models

//go:generate easytags $GOFILE json,sql

type Userer interface {
	TableName() string
}

type User struct {
	Email    string `json:"email" db:"email" db_type:"text"`
	Password string `json:"password" db:"password" db_type:"text"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) GetPassword() string {
	return u.Password
}
