package infra

import (
	"database/sql"
	"github.com/bubaew95/yandex-diplom-2/config"
	"time"
)

// DataBase представляет обёртку над *sql.DB,
// используемую для управления соединением с PostgreSQL через драйвер pgx.
type DataBase struct {
	*sql.DB
}

// NewDB создаёт новый экземпляр DataBase, устанавливая подключение к PostgreSQL.
//
// Использует параметры подключения из конфигурации `config.Config`.
// Устанавливает параметры пула соединений (макс. время жизни, макс. число соединений и т.д.).
func NewDB(c *config.Config) (*DataBase, error) {
	db, err := connectDB(c)
	if err != nil {
		return nil, err
	}

	return &DataBase{db}, nil
}

// connectDB устанавливает подключение к базе данных PostgreSQL через драйвер "pgx",
// настраивает параметры пула соединений и проверяет доступность БД через Ping().
func connectDB(c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", c.DSN)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(3))
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close закрывает соединение с базой данных.
func (d DataBase) Close() error {
	return d.DB.Close()
}
