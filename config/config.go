// Package config предоставляет структуру и функции для конфигурации
// параметров приложения, включая порт, строку подключения (DSN)
// и токен авторизации.
package config

import (
	"flag"
	"os"
	"time"
)

// token представляет данные токена авторизации,
// включая секретный ключ и время истечения.
type token struct {
	Secret string    // Secret — секретный ключ токена
	Exp    time.Time // Exp — время истечения токена
}

// Config представляет конфигурацию приложения,
// включая порт, строку подключения к базе данных (DSN)
// и данные токена авторизации.
type Config struct {
	Port  string // Port — порт, на котором будет слушать приложение
	DSN   string // DSN — строка подключения к базе данных
	Token token  // Token — данные токена авторизации
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
	flag.Parse()

	// Если задана переменная окружения PORT, она переопределяет значение из флага.
	if envAddress := os.Getenv("PORT"); envAddress != "" {
		*port = envAddress
	}

	// Если задана переменная окружения DSN, она переопределяет значение из флага.
	if envDSN := os.Getenv("DSN"); envDSN != "" {
		*dsn = envDSN
	}

	return &Config{
		Port: *port,
		DSN:  *dsn,
	}
}
