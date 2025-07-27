package service

import (
	"context"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/bubaew95/yandex-diplom-2/pkg/token"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, r *model.RegistrationDTO) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, r *model.LoginDTO) (model.User, error)

	AddText(ctx context.Context, r *model.Data, userID int64) (int64, error)
	Edit(ctx context.Context, r *model.TextRequest, userID int64) (int64, error)
	Delete(ctx context.Context, userID int64, ID int64) error
	FindAll(ctx context.Context, userID int64) ([]*pb.DataResponse, error)
}

type Service struct {
	repo Repository
	cfg  config.Config
}

func NewService(repo Repository, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s Service) AddUser(ctx context.Context, r *model.RegistrationDTO) (*model.AuthResponse, error) {
	if r.Password != r.RePassword {
		return nil, model.PasswordNotMatchError
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	r.Password = string(hash)
	userID, err := s.repo.CreateUser(ctx, r)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:        userID,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}

	jwt, err := token.EncodeJWTToken(user)
	if err != nil {
		return nil, model.AuthorizationError
	}

	return &model.AuthResponse{
		Token: jwt,
	}, nil
}
func (s Service) Login(ctx context.Context, r *model.LoginDTO) (model.AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, r)
	if err != nil {
		logger.Log.Debug("login failed", zap.Error(err))
		return model.AuthResponse{}, model.AuthorizationError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
	if err != nil {
		return model.AuthResponse{}, model.LoginAndPasswordError
	}

	jwt, err := token.EncodeJWTToken(user)
	if err != nil {
		return model.AuthResponse{}, model.AuthorizationError
	}

	return model.AuthResponse{
		Token: jwt,
	}, nil
}

func (s Service) Add(ctx context.Context, r *model.Data) (model.TextResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	dataID, err := s.repo.AddText(ctx, r, user.ID)
	if err != nil {
		return model.TextResponse{}, err
	}

	return model.TextResponse{
		ID:     dataID,
		Text:   r.Text,
		UserID: user.ID,
	}, nil
}

func (s Service) Edit(ctx context.Context, r *model.TextRequest) (model.TextResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	_, err := s.repo.Edit(ctx, r, user.ID)
	if err != nil {
		return model.TextResponse{}, err
	}

	return model.TextResponse{
		ID:     r.ID,
		Text:   r.Text,
		UserID: user.ID,
	}, nil
}

func (s Service) Delete(ctx context.Context, ID int64) error {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return s.repo.Delete(ctx, user.ID, ID)
}

func (s Service) FindAll(ctx context.Context) ([]*pb.DataResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return s.repo.FindAll(ctx, user.ID)
}
