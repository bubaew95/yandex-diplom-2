package grpc

import (
	"context"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Service --filename=servicemock_test.go --inpackage
type Service interface {
	AddUser(ctx context.Context, r *model.RegistrationRequest) (model.RegistrationResponse, error)
	Login(ctx context.Context, r *model.LoginRequest) (model.AuthResponse, error)

	AddText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	EditText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	DeleteText(ctx context.Context, ID int64) error

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

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not found")
		}

		token := ""
		if vals := md.Get("token"); len(vals) > 0 {
			token = vals[0]
		}

		if token == "" || !crypto.IsInvalidUserID(&http.Cookie{Name: "user_id", Value: userID}) {
			rawID := crypto.GenerateUserID()
			encodedID, err := crypto.EncodeUserID(rawID)
			if err != nil {
				return nil, status.Error(codes.Internal, "user ID encoding failed")
			}
			userID = encodedID
		}

		ctx = context.WithValue(ctx, crypto.KeyUser, userID)
		return handler(ctx, req)
	}
}

func (s *Server) AddText(ctx context.Context, r *pb.TextRequest) (*pb.TextResponse, error) {
	ctxx := context.WithValue(ctx, "user", model.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
	})

	data, err := s.service.AddText(ctxx, &model.TextRequest{
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

func (s *Server) AddCard(ctx context.Context, r *pb.CardRequest) (*pb.CardResponse, error) {
	return &pb.CardResponse{}, nil
}
func (s *Server) EditCard(ctx context.Context, r *pb.CardEditRequest) (*pb.CardResponse, error) {
	return &pb.CardResponse{}, nil
}
func (s *Server) DeleteCard(ctx context.Context, r *pb.IdRequest) (*pb.SuccessResponse, error) {
	return &pb.SuccessResponse{}, nil
}
