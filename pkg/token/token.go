package token

import (
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// Claims представляет JWT-полезную нагрузку (payload), содержащую:
// - стандартные claims (exp, iat, и т.п.),
// - информацию о пользователе.
type Claims struct {
	jwt.RegisteredClaims            // Встроенные поля: ExpiresAt, IssuedAt, и т.п.
	User                 model.User // Пользовательские данные (структура User)
}

type Token struct {
	cfg *config.Config
}

// NewToken - инциализация токена
func NewToken(cfg *config.Config) *Token {
	return &Token{
		cfg: cfg,
	}
}

// EncodeJWTToken генерирует JWT-токен для заданного пользователя.
//
// Токен содержит:
//   - поле `User` со всей структурой model.User,
//   - срок действия (ExpiresAt).
func (tkn *Token) EncodeJWTToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tkn.cfg.Token.Exp)),
		},
		User: user,
	})

	tokenString, err := token.SignedString([]byte(tkn.cfg.Token.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeJWTToken парсит и проверяет JWT-токен, извлекая пользователя.
//
// Возвращает:
//   - model.User — если токен валиден,
//   - error — если токен невалиден или подпись некорректна.
func (tkn *Token) DecodeJWTToken(tokenString string) (model.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tkn.cfg.Token.Secret), nil
	})

	if err != nil {
		return model.User{}, model.Error(err.Error(), http.StatusInternalServerError)
	}

	if !token.Valid {
		return model.User{}, model.Error("Token is not valid", http.StatusInternalServerError)
	}

	return claims.User, nil
}
