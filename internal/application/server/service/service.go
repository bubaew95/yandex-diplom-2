package service

import (
	"context"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/bubaew95/yandex-diplom-2/pkg/token"
	"go.uber.org/zap"
	"net/http"
)

type Repository interface {
	CreateUser(ctx context.Context, r *model.RegistrationRequest) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, r *model.LoginRequest) (model.User, error)
}

type Service struct {
	repo Repository
	cfg  config.Config
}

func NewService(repo Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s Service) AddUser(ctx context.Context, r *model.RegistrationRequest) (model.RegistrationResponse, error) {
	if r.Password != r.RePassword {
		return model.RegistrationResponse{}, model.Error("password does not match", http.StatusInternalServerError)
	}

	hash, err := crypto.EncodeHash(r.Password)
	if err != nil {
		return model.RegistrationResponse{}, err
	}

	r.Password = hash
	userID, err := s.repo.CreateUser(ctx, r)
	if err != nil {
		return model.RegistrationResponse{}, err
	}

	user := model.User{
		ID:        userID,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}

	return model.RegistrationResponse{
		User: user,
	}, nil
}

func (s Service) Login(ctx context.Context, r *model.LoginRequest) (model.AuthResponse, error) {
	passwordHash, err := crypto.EncodeHash(r.Password)
	if err != nil {
		return model.AuthResponse{}, model.ErrorResponse{
			Message: err.Error(),
		}
	}

	r.Password = passwordHash

	user, err := s.repo.FindUserByEmail(ctx, r)
	if err != nil {
		logger.Log.Debug("login failed", zap.Error(err))
		return model.AuthResponse{}, model.Error(err.Error(), http.StatusUnauthorized)
	}

	jwt, err := token.EncodeJWTToken(user)
	if err != nil {
		return model.AuthResponse{}, model.Error(err.Error(), http.StatusUnauthorized)
	}

	if user.Password != r.Password {
		return model.AuthResponse{}, model.Error("not correct login or password", http.StatusUnauthorized)
	}

	return model.AuthResponse{
		Token: jwt,
	}, nil
}
