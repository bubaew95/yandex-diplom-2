package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	grpcServer "github.com/bubaew95/yandex-diplom-2/internal/application/server/handlers/grpc"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/service"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

type HTTPServer interface {
	Start()
	Stop()
}

type httpServer struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	router     *chi.Mux
	config     config.Config
	repo       service.Repository
	chError    chan string
}

func NewServer(r *chi.Mux, repo service.Repository, cfg config.Config) HTTPServer {
	return &httpServer{
		config:  cfg,
		router:  r,
		repo:    repo,
		chError: make(chan string),
	}
}

func (s *httpServer) Start() {
	if s.config.EnableGRPC {
		s.grpcServer = s.listenGRPC()
	} else {
		s.httpServer = s.listenHTTP()
	}

	select {
	case result := <-s.chError:
		logger.Log.Fatal(result)
	default:
		close(s.chError)
	}
}

func (s *httpServer) listenHTTP() *http.Server {
	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.router,
	}

	if s.config.EnableHTTPS {
		logger.Log.Info("Running https server", zap.String("port", s.config.ServerAddress))
		go func() {
			if err := server.Serve(autocert.NewListener(s.config.ServerAddress)); err != nil {
				s.chError <- fmt.Sprintf("Failed to start https(tsl) server: %s", err)
			}
		}()
	} else {
		logger.Log.Info("Running https server", zap.String("port", s.config.Port))
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.chError <- fmt.Sprintf("Failed to start http server: %s", err)
			}
		}()
	}

	return server
}

func (s *httpServer) listenGRPC() *grpc.Server {
	listener, err := net.Listen("tcp", ":3232")
	if err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	srv := service.NewService(s.repo, s.config)

	//grpc.UnaryInterceptor(grpc.AuthInterceptor()),
	server := grpc.NewServer()
	pb.RegisterGoKeeperServer(server, grpcServer.NewServer(srv))

	logger.Log.Info("Run rpc server")
	if err := server.Serve(listener); err != nil {
		s.chError <- fmt.Sprintf("Rpc server error: %s", err)
	}

	return server
}

func (s *httpServer) Stop() {
	if s.config.EnableGRPC && s.grpcServer != nil {
		logger.Log.Info("Stopping grpc server")
		s.grpcServer.Stop()
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			logger.Log.Fatal("Failed to stop http server", zap.Error(err))
		}
	}
}
