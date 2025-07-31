// Package config предоставляет структуру и функции для конфигурации
// параметров приложения, включая порт, строку подключения (DSN)
// и токен авторизации.
package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// Token представляет данные токена авторизации,
// включая секретный ключ и время истечения.
type Token struct {
	Secret string        // Secret — секретный ключ токена
	Exp    time.Duration // Exp — время истечения токена
}

// Config представляет конфигурацию приложения,
// включая порт, строку подключения к базе данных (DSN)
// и данные токена авторизации.
type Config struct {
	Port      string // Port — порт, на котором будет слушать приложение
	DSN       string // DSN — строка подключения к базе данных
	SecretKey string // SecretKey - секрет ключ шифрования
	Token     Token  // Token — данные токена авторизации
}

// NewConfig создает новый экземпляр Config.
// Приоритет загрузки конфигурации:
//  1. Значения, переданные через флаги командной строки (-p для порта, -d для DSN).
//  2. Переменные окружения PORT и DSN (перезаписывают значения из флагов).
//
// Возвращает указатель на структуру Config.
func NewConfig() *Config {
	port := flag.String("p", "", "port to listen on")
	dsn := flag.String("d", "", "dsn to connect to")
	tokenSecret := flag.String("t", "", "Token secret")
	tokenExpiration := flag.Int("e", 1, "Token expiration")
	secretKey := flag.String("s", "", "Secret key")
	flag.Parse()

	// Если задана переменная окружения PORT, она переопределяет значение из флага.
	if envAddress := os.Getenv("PORT"); envAddress != "" {
		*port = envAddress
	}

	// Если задана переменная окружения DSN, она переопределяет значение из флага.
	if envDSN := os.Getenv("DSN"); envDSN != "" {
		*dsn = envDSN
	}

	if envTokenSecretEnv := os.Getenv("TOKEN_SECRET_KEY"); envTokenSecretEnv != "" {
		*tokenSecret = envTokenSecretEnv
	}

	if tokenExpirationEnv := os.Getenv("TOKEN_EXPIRE_AT"); tokenExpirationEnv != "" {
		exp, _ := strconv.Atoi(tokenExpirationEnv)
		*tokenExpiration = exp
	}

	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		*secretKey = envSecretKey
	}

	return &Config{
		Port:      *port,
		DSN:       *dsn,
		SecretKey: *secretKey,
		Token: Token{
			Secret: *tokenSecret,
			Exp:    time.Duration(*tokenExpiration) * time.Hour,
		},
	}
}
