package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtkey = []byte("minha_chave_secreta")

func GenerateToken(userID int, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtkey)
}
