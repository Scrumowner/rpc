package decoder

import (
	"github.com/golang-jwt/jwt/v5"
)

type Decoder struct {
	Secret []byte
}

func NewDecoder(secret string) *Decoder {
	return &Decoder{
		Secret: []byte(secret),
	}
}

func (d *Decoder) Decdoe(token string) (email string, phone string) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return d.Secret, nil
	})
	if err != nil {
		return "", ""
	}
	mp, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", ""
	}
	return mp["email"].(string), mp["phone"].(string)
}
