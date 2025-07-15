package main

import (
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/handlers"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/handlers/middleware"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/repository"
	service "github.com/bubaew95/yandex-diplom-2/internal/application/server/service"
	infra "github.com/bubaew95/yandex-diplom-2/internal/infra/database"
	httpServer "github.com/bubaew95/yandex-diplom-2/internal/infra/http"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	route := chi.NewRouter()

	db, err := infra.NewDB(&cfg)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		logger.Log.Fatal("Repository init failed", zap.Error(err))
	}

	srv := service.NewService(repo, cfg)
	handler := handlers.NewHandler(srv, &cfg)

	route.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware())
			r.Get("/sync", handler.Sync)
		})
	})

	server := httpServer.NewServer(route, cfg)
	server.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch

	server.Stop()
}
