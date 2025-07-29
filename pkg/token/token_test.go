package token

import (
	"encoding/base64"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTTokenEncodeDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		user           model.User
		modifyToken    func(string) string
		wantUser       model.User
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:      "успешное кодирование и декодирование",
			user:      model.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"},
			wantUser:  model.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"},
			expectErr: false,
		},
		{
			name: "повреждённый токен",
			user: model.User{ID: 1},
			modifyToken: func(tok string) string {
				return tok[:10] + "corrupted"
			},
			expectErr:      true,
			expectedErrMsg: "token contains an invalid number of segments",
		},
		{
			name: "подпись с другим секретом",
			user: model.User{ID: 2},
			modifyToken: func(_ string) string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					},
					User: model.User{ID: 2},
				})
				tokStr, _ := token.SignedString([]byte("WRONG_SECRET"))
				return tokStr
			},
			expectErr:      true,
			expectedErrMsg: "signature is invalid",
		},
		{
			name: "просроченный токен",
			user: model.User{ID: 3},
			modifyToken: func(_ string) string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
					},
					User: model.User{ID: 3},
				})
				tokStr, _ := token.SignedString([]byte(SecretKey))
				return tokStr
			},
			expectErr:      true,
			expectedErrMsg: "token is expired",
		},
		{
			name: "неподдерживаемый алгоритм подписи",
			user: model.User{ID: 4},
			modifyToken: func(_ string) string {
				// Подделываем токен с неподдерживаемым алгоритмом
				header := `{"alg":"none","typ":"JWT"}`
				payload := `{"User":{"id":4},"exp":4070880000}` // exp в далеком будущем
				// Закодируем
				h := base64.RawURLEncoding.EncodeToString([]byte(header))
				p := base64.RawURLEncoding.EncodeToString([]byte(payload))
				// Не добавляем подпись вообще
				return fmt.Sprintf("%s.%s.", h, p)
			},
			expectErr:      true,
			expectedErrMsg: "unexpected signing method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok, err := EncodeJWTToken(tt.user)
			require.NoError(t, err)
			require.NotEmpty(t, tok)

			if tt.modifyToken != nil {
				tok = tt.modifyToken(tok)
			}

			got, err := DecodeJWTToken(tok)

			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUser, got)
			}
		})
	}
}

func TestDecodeJWTToken(t *testing.T) {
	t.Parallel()

	user := model.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Иван",
		LastName:  "Иванов",
	}

	type testCase struct {
		name           string
		tokenGenerator func() string
		wantUser       model.User
		wantErr        bool
		expectedErrMsg string
	}

	tests := []testCase{
		{
			name: "валидный токен",
			tokenGenerator: func() string {
				tok, err := EncodeJWTToken(user)
				require.NoError(t, err)
				return tok
			},
			wantUser: user,
			wantErr:  false,
		},
		{
			name: "токен с невалидной подписью",
			tokenGenerator: func() string {
				tok, err := EncodeJWTToken(user)
				require.NoError(t, err)
				parts := strings.Split(tok, ".")
				parts[2] = "invalidsignature" // подменяем подпись
				return strings.Join(parts, ".")
			},
			wantErr:        true,
			expectedErrMsg: "signature is invalid", // ← исправлено
		},
		{
			name: "просроченный токен",
			tokenGenerator: func() string {
				tok := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					},
					User: user,
				})
				signed, err := tok.SignedString([]byte(SecretKey))
				require.NoError(t, err)
				return signed
			},
			wantErr:        true,
			expectedErrMsg: "token is expired",
		},
		{
			name: "неподдерживаемый алгоритм",
			tokenGenerator: func() string {
				// Ручной токен с alg=none
				header := `{"alg":"none","typ":"JWT"}`
				payload := fmt.Sprintf(`{"User":{"id":%d,"email":"%s"},"exp":%d}`,
					user.ID, user.Email, time.Now().Add(time.Hour).Unix())
				h := base64.RawURLEncoding.EncodeToString([]byte(header))
				p := base64.RawURLEncoding.EncodeToString([]byte(payload))
				return fmt.Sprintf("%s.%s.", h, p)
			},
			wantErr:        true,
			expectedErrMsg: "unexpected signing method",
		},
		{
			name: "некорректный токен (мало сегментов)",
			tokenGenerator: func() string {
				return "abc.def" // не хватает подписи
			},
			wantErr:        true,
			expectedErrMsg: "token is malformed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := tt.tokenGenerator()

			gotUser, err := DecodeJWTToken(tok)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUser, gotUser)
			}
		})
	}
}
