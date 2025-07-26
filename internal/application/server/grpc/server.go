package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/bubaew95/yandex-diplom-2/pkg/token"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Service --filename=servicemock_test.go --inpackage
type Service interface {
	AddUser(ctx context.Context, r *model.RegistrationDTO) (*model.AuthResponse, error)
	Login(ctx context.Context, r *model.LoginDTO) (model.AuthResponse, error)

	Add(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	EditText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	DeleteText(ctx context.Context, ID int64) error
	FindAllText(ctx context.Context) ([]*pb.TextResponse, error)

	AddCard(ctx context.Context, r *model.CardRequest) (model.CardResponse, error)
	EditCard(ctx context.Context, r *model.CardRequest) (model.CardResponse, error)
	DeleteCard(ctx context.Context, ID int64) error

	AddBinary(ctx context.Context, r *model.BinaryRequest) (model.BinaryResponse, error)
}

type Server struct {
	pb.UnimplementedGoKeeperServer
	service Service
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}

func LoginInterceptor() grpc.UnaryServerInterceptor {
	publicMethods := map[string]struct{}{
		"/gokeeper.GoKeeper/Registration": {},
		"/gokeeper.GoKeeper/Login":        {},
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not found")
		}

		tkn := ""
		if vals := md.Get("token"); len(vals) > 0 {
			tkn = vals[0]
		}

		if tkn == "" || !strings.HasPrefix(tkn, "Bearer") {
			logger.Log.Debug("invalid token", zap.String("token", tkn))
			return nil, status.Error(codes.Unauthenticated, "authorization error")
		}

		user, err := token.DecodeJWTToken(tkn[7:])
		if err != nil {
			logger.Log.Debug("token decode error", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		nCtx := context.WithValue(ctx, crypto.KeyUser, user)
		fmt.Println(user)

		return handler(nCtx, req)
	}
}

func (s *Server) Registration(ctx context.Context, r *pb.RegistrationRequest) (*pb.TokenResponse, error) {
	regData := model.RegistrationDTO{
		User: model.User{
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Email:     r.Email,
			Password:  r.Password,
		},
		RePassword: r.RePassword,
	}

	if valid := regData.Validate(); len(valid) != 0 {
		logger.Log.Debug("invalid regData", zap.Any("regData", valid))
		return nil, status.Error(codes.InvalidArgument, "invalid registration request")
	}

	jwt, err := s.service.AddUser(ctx, &regData)
	if err != nil {
		if errors.Is(err, model.UserAlreadyExistsError) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TokenResponse{
		Token: jwt.Token,
	}, nil
}

func (s *Server) Login(ctx context.Context, r *pb.LoginRequest) (*pb.TokenResponse, error) {
	user, err := s.service.Login(ctx, &model.LoginDTO{
		Email:    r.Email,
		Password: r.Password,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TokenResponse{
		Token: user.Token,
	}, nil
}

func (s *Server) Add(ctx context.Context, r *pb.TextRequest) (*pb.TextResponse, error) {
	data, err := s.service.Add(ctx, &model.TextRequest{
		Text: r.Text,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TextResponse{
		Id:     data.ID,
		Text:   data.Text,
		UserId: data.UserID,
	}, nil
}
func (s *Server) EditText(ctx context.Context, r *pb.TextEditRequest) (*pb.TextResponse, error) {
	data, err := s.service.EditText(ctx, &model.TextRequest{
		ID:   r.Id,
		Text: r.Text,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TextResponse{
		Id:     data.ID,
		Text:   data.Text,
		UserId: data.UserID,
	}, nil
}
func (s *Server) DeleteText(ctx context.Context, r *pb.IdRequest) (*pb.SuccessResponse, error) {
	if err := s.service.DeleteText(ctx, r.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SuccessResponse{
		Success: true,
	}, nil
}
func (s *Server) FindAllText(ctx context.Context, r *pb.DataRequest) (*pb.TextList, error) {
	list, err := s.service.FindAllText(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TextList{
		List: list,
	}, nil
}

func (s *Server) AddCard(ctx context.Context, r *pb.CardRequest) (*pb.CardResponse, error) {
	return &pb.CardResponse{}, nil
}
func (s *Server) EditCard(ctx context.Context, r *pb.CardEditRequest) (*pb.CardResponse, error) {
	return &pb.CardResponse{}, nil
}
func (s *Server) DeleteCard(ctx context.Context, r *pb.IdRequest) (*pb.SuccessResponse, error) {
	return &pb.SuccessResponse{}, nil
}
