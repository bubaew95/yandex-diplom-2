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

// HTTPServer описывает интерфейс управления сервером (запуск и остановка).
type HTTPServer interface {
	Start() // Запускает gRPC-сервер
	Stop()  // Останавливает gRPC-сервер
}

// httpServer представляет реализацию gRPC-сервера с авторизационным интерсептором.
//
// Хранит:
//   - config: настройки (порт),
//   - srv: реализация бизнес-логики (service-layer),
//   - server: экземпляр gRPC-сервера,
//   - chError: канал для передачи ошибок запуска.
type httpServer struct {
	server  *grpc.Server
	config  config.Config
	srv     grpcServer.Service
	chError chan string
}

// NewServer создаёт и возвращает новый экземпляр httpServer,
// инициализируя зависимости и конфигурацию.
//
// Используется как точка входа для старта gRPC-сервера.
func NewServer(srv grpcServer.Service, cfg config.Config) HTTPServer {
	return &httpServer{
		config:  cfg,
		srv:     srv,
		chError: make(chan string),
	}
}

// Start запускает gRPC-сервер на указанном порту из конфигурации.
//
// Использует LoginInterceptor для авторизации.
func (s *httpServer) Start() {
	s.server = s.listenGRPC()

	select {
	case result := <-s.chError:
		logger.Log.Fatal(result)
	default:
		close(s.chError)
	}
}

// listenGRPC инициализирует gRPC-сервер, настраивает слушателя,
// регистрирует сервер GoKeeper и запускает `Serve()`.
//
// Ошибки передаются в канал s.chError.
func (s *httpServer) listenGRPC() *grpc.Server {
	listener, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		logger.Log.Fatal("Rpc server error", zap.Error(err))
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(grpcServer.LoginInterceptor(s.config)))
	pb.RegisterGoKeeperServer(server, grpcServer.NewServer(s.srv))

	logger.Log.Info("Run rpc server. Port: " + s.config.Port)
	if err := server.Serve(listener); err != nil {
		s.chError <- fmt.Sprintf("Rpc server error: %s", err)
	}

	return server
}

// Stop корректно останавливает gRPC-сервер.
func (s *httpServer) Stop() {
	logger.Log.Info("Stopping grpc server")
	s.server.Stop()
}
