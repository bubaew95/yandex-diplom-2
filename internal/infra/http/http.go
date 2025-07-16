package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"time"
)

type HTTPServer interface {
	Start()
	Stop()
}

type httpServer struct {
	server *http.Server
	config config.Config
}

func NewServer(r *chi.Mux, cfg config.Config) HTTPServer {
	return &httpServer{
		server: &http.Server{
			Addr:    cfg.Port,
			Handler: r,
		},
		config: cfg,
	}
}

func (s *httpServer) Start() {
	logger.Log.Info("Starting http server", zap.String("address", s.config.Port))
	chError := make(chan string)

	if s.config.EnableHTTPS {
		chError = s.startTsl(chError)
	} else {
		chError = s.startHttp(chError)
	}

	select {
	case result := <-chError:
		logger.Log.Fatal(result)
	}
}

func (s *httpServer) startHttp(chError chan string) chan string {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chError <- fmt.Sprintf("Failed to start http server: %s", err)
		}
	}()

	return chError
}

func (s *httpServer) startTsl(chError chan string) chan string {
	logger.Log.Info("Running https server", zap.String("port", s.config.ServerAddress))
	go func() {
		if err := s.server.Serve(autocert.NewListener(s.config.ServerAddress)); err != nil {
			chError <- fmt.Sprintf("Failed to start https(tsl) server: %s", err)
		}
	}()

	return chError
}

func (s *httpServer) Stop() {
	logger.Log.Info("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Failed to stop http server", zap.Error(err))
	}
}
