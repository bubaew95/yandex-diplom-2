package token

import (
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	User model.User
}

const TokenExp = time.Hour * 3
const SecretKey = "sdgsg!35$#%TSGsdhdfhsd436093598!@$#%"

func EncodeJWTToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		User: user,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWTToken(tokenString string) (model.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return model.User{}, model.Error(err.Error(), http.StatusInternalServerError)
	}

	if !token.Valid {
		return model.User{}, model.Error("Token is not valid", http.StatusInternalServerError)
	}

	return claims.User, nil
}
