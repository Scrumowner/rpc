package tool

import (
	"auth/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	time "time"
)

const (
	AccessToken = iota
	RefreshToken
)

type Claims struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}
type JwtAuther struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

func NewJwtAuther(config *config.Config) *JwtAuther {
	return &JwtAuther{
		AccessSecret: []byte(config.Secret),
	}
}

func (a *JwtAuther) GenerateToken(email, password string, kind int) (string, error) {

	claims := &Claims{
		Email: email,
		Phone: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	var secret []byte
	switch kind {
	case AccessToken:
		secret = a.AccessSecret
	case RefreshToken:
		secret = a.RefreshSecret
	default:
		return "", jwt.ErrInvalidKey
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return strToken, nil
}
func (a *JwtAuther) CheckToken(token string, kind int) (Claims, bool) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		var secret []byte
		switch kind {
		case AccessToken:
			secret = a.AccessSecret
		case RefreshToken:
			secret = a.RefreshSecret
		default:
			return "", jwt.ErrSignatureInvalid
		}
		_ = secret

		return secret, nil

	})
	if err != nil {
		return Claims{}, false
	} else if !tkn.Valid {
		return Claims{}, false
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, false
	}
	return Claims{
		Email: claims["email"].(string),
		Phone: claims["phone"].(string),
	}, true
}
