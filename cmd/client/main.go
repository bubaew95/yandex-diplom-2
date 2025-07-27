package main

import (
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client/pages"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := logger.Load(); err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("No .env file found")
	}
}

func main() {
	cfg := config.NewConfig()
	tui, err := pages.NewTUI(cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	if err := tui.Run(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
