package main

import (
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client"
	"github.com/bubaew95/yandex-diplom-2/internal/application/client/pages"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/joho/godotenv"
	"log"
)

var (
	version   = "0.1"
	buildDate string
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
	fmt.Printf("Версия: %s\nДата сборки: %s\n", version, buildDate)

	cfg := config.NewConfig()
	tui, err := pages.NewTUI(cfg)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	defer func(Client *client.Client) {
		fmt.Printf("Client is shutting down")
		err := Client.Stop()
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}(tui.Client)

	if err := tui.Run(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
