package models

//go:generate easytags $GOFILE json,db,db_ops,db_type,db_default,db_index

type Userer interface {
	TableName() string
}

type User struct {
	Email    string `json:"email" db:"email" db_type:"text" db_ops:"create,update" db_default:"not null" db_index:"unique"`
	Phone    string `json:"phone" db:"phone" db_type:"text" db_ops:"create,update" db_default:"not null" db_index:"unique"`
	Password string `json:"password" db:"password" db_type:"text" db_ops:"create,update" db_default:"not null" db_index:"unique"`
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

func (u *User) SetPhone(phone string) {
	u.Phone = phone
}

func (u *User) GetPhone() string {
	return u.Phone
}
