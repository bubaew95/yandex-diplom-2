package main

import (
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/handlers/http"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/handlers/http/middleware"
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

	repo := repository.NewRepository(db)
	srv := service.NewService(repo, cfg)
	handler := handlers.NewHandler(srv, &cfg)

	route.Use(middleware.GZipMiddleware)
	route.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware())

			r.Route("/text", func(r chi.Router) {
				r.Post("/", handler.AddText)
				r.Put("/{id}", handler.EditText)
				r.Delete("/{id}", handler.DeleteText)
			})

			r.Route("/cards", func(r chi.Router) {
				r.Post("/", handler.AddCard)
				r.Put("/{id}", handler.EditCard)
				r.Delete("/{id}", handler.DeleteCard)
			})
		})
	})

	server := httpServer.NewServer(route, repo, cfg)
	server.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch

	logger.Log.Info("Shutting down server")
	server.Stop()
}
