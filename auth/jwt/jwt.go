package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	time "time"
)

const Secret = "mysycret"

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}
type JwtAuther struct{}

func NewJwtAuther() *JwtAuther {
	return &JwtAuther{}
}

func (a *JwtAuther) GenerateToken(email, password string) (string, error) {
	time := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		Username: email,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return strToken, nil
}
func (a *JwtAuther) CheckToken(token string) (string, string, bool) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return "", "", false
	} else if !tkn.Valid {
		return "", "", false
	}
	return claims.Username, claims.Password, true
}
