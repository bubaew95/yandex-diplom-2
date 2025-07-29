package token

import (
	"fmt"
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

// TokenExp задаёт срок действия JWT-токена: 3 часа.
const TokenExp = time.Hour * 3

// SecretKey используется для подписи и проверки JWT-токенов.
const SecretKey = "sdgsg!35$#%TSGsdhdfhsd436093598!@$#%"

// EncodeJWTToken генерирует JWT-токен для заданного пользователя.
//
// Токен содержит:
//   - поле `User` со всей структурой model.User,
//   - срок действия (ExpiresAt).
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

// DecodeJWTToken парсит и проверяет JWT-токен, извлекая пользователя.
//
// Возвращает:
//   - model.User — если токен валиден,
//   - error — если токен невалиден или подпись некорректна.
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
