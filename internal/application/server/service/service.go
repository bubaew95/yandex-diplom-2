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
)

type Repository interface {
	CreateUser(ctx context.Context, r *model.RegistrationDTO) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, r *model.LoginDTO) (model.User, error)

	AddText(ctx context.Context, r *model.TextRequest, userID int64) (int64, error)
	EditText(ctx context.Context, r *model.TextRequest, userID int64) (int64, error)
	DeleteText(ctx context.Context, userID int64, ID int64) error
	FindAllText(ctx context.Context, userID int64) ([]*pb.TextResponse, error)

	AddCard(ctx context.Context, r *model.CardRequest, userID int64) (int64, error)
	EditCard(ctx context.Context, r *model.CardRequest, userID int64) (int64, error)
	DeleteCard(ctx context.Context, userID int64, ID int64) error

	AddBinary(ctx context.Context, r *model.BinaryRequest, userID int64) (int64, error)
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

	hash, err := crypto.EncodeHash(r.Password)
	if err != nil {
		return nil, err
	}

	r.Password = hash
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
		return model.AuthResponse{}, model.AuthorizationError
	}

	if user.Password != r.Password {
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

func (s Service) Add(ctx context.Context, r *model.TextRequest) (model.TextResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	hashText, err := crypto.EncodeHash(r.Text)
	if err != nil {
		return model.TextResponse{}, err
	}

	r.Text = hashText
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
func (s Service) EditText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	hashText, err := crypto.EncodeHash(r.Text)
	if err != nil {
		return model.TextResponse{}, err
	}
	r.Text = hashText

	_, err = s.repo.EditText(ctx, r, user.ID)
	if err != nil {
		return model.TextResponse{}, err
	}

	return model.TextResponse{
		ID:     r.ID,
		Text:   r.Text,
		UserID: user.ID,
	}, nil
}
func (s Service) DeleteText(ctx context.Context, ID int64) error {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return s.repo.DeleteText(ctx, user.ID, ID)
}
func (s Service) FindAllText(ctx context.Context) ([]*pb.TextResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return s.repo.FindAllText(ctx, user.ID)
}

func (s Service) AddCard(ctx context.Context, r *model.CardRequest) (model.CardResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)
	dataID, err := s.repo.AddCard(ctx, r, user.ID)
	if err != nil {
		return model.CardResponse{}, err
	}

	return model.CardResponse{
		ID:     dataID,
		Number: r.Number,
		UserID: user.ID,
	}, nil
}
func (s Service) EditCard(ctx context.Context, r *model.CardRequest) (model.CardResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	_, err := s.repo.EditCard(ctx, r, user.ID)
	if err != nil {
		return model.CardResponse{}, err
	}

	return model.CardResponse{
		ID:     r.ID,
		Number: r.Number,
		UserID: user.ID,
	}, nil
}
func (s Service) DeleteCard(ctx context.Context, ID int64) error {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return s.repo.DeleteCard(ctx, user.ID, ID)
}

func (s Service) AddBinary(ctx context.Context, r *model.BinaryRequest) (model.BinaryResponse, error) {
	user := ctx.Value(crypto.KeyUser).(model.User)

	return model.BinaryResponse{
		UserID: user.ID,
	}, nil
}
