package http

import (
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	grpcServer "github.com/bubaew95/yandex-diplom-2/internal/application/server/grpc"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type HTTPServer interface {
	Start()
	Stop()
}

type httpServer struct {
	server  *grpc.Server
	config  config.Config
	srv     grpcServer.Service
	chError chan string
}

func NewServer(srv grpcServer.Service, cfg config.Config) HTTPServer {
	return &httpServer{
		config:  cfg,
		srv:     srv,
		chError: make(chan string),
	}
}

func (s *httpServer) Start() {
	s.server = s.listenGRPC()

	select {
	case result := <-s.chError:
		logger.Log.Fatal(result)
	default:
		close(s.chError)
	}
}

func (s *httpServer) listenGRPC() *grpc.Server {
	listener, err := net.Listen("tcp", ":3232")
	if err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(grpcServer.LoginInterceptor()))
	pb.RegisterGoKeeperServer(server, grpcServer.NewServer(s.srv))

	logger.Log.Info("Run rpc server")
	if err := server.Serve(listener); err != nil {
		s.chError <- fmt.Sprintf("Rpc server error: %s", err)
	}

	return server
}

func (s *httpServer) Stop() {
	logger.Log.Info("Stopping grpc server")
	s.server.Stop()
}
