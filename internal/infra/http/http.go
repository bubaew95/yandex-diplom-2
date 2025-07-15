package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
			Addr:    cfg.Address,
			Handler: r,
		},
		config: cfg,
	}
}

func (s *httpServer) Start() {
	logger.Log.Info("Starting http server", zap.String("address", s.config.Address))
	chError := make(chan string)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chError <- fmt.Sprintf("Failed to start http server: %s", err)
		}
	}()

	select {
	case result := <-chError:
		logger.Log.Fatal(result)
	}
}

func (s *httpServer) Stop() {
	logger.Log.Info("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Failed to stop http server", zap.Error(err))
	}
}
