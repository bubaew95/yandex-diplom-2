package infra

import (
	"database/sql"
	"github.com/bubaew95/yandex-diplom-2/config"
)

type DataBase struct {
	*sql.DB
}

func NewDB(c *config.Config) (*DataBase, error) {
	db, err := connectDB(c)
	if err != nil {
		return nil, err
	}

	return &DataBase{db}, nil
}

func connectDB(c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", c.DSN)
	if err != nil {
		return nil, err
	}

	//db.SetConnMaxLifetime(time.Minute * time.Duration(c.Database.ConnMaxLifeTimeInMinute))
	//db.SetMaxOpenConns(c.Database.MaxOpenConns)
	//db.SetMaxIdleConns(c.Database.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (d DataBase) Close() error {
	return d.DB.Close()
}
