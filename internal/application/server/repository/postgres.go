package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	infra "github.com/bubaew95/yandex-diplom-2/internal/infra/database"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Repository struct {
	db *infra.DataBase
}

func NewRepository(db *infra.DataBase) (*Repository, error) {
	return &Repository{db: db}, nil
}

func (s *Repository) CreateUser(ctx context.Context, r *model.RegistrationRequest) (int64, error) {
	isUser, err := s.GetUserByEmail(ctx, r.Email)
	if err != nil {
		return -1, err
	}

	if isUser {
		return -1, model.UserAlreadyExistsError
	}

	sqlQuery := `INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	row := s.db.QueryRowContext(ctx, sqlQuery, r.Email, r.FirstName, r.LastName, r.Password)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Repository) GetUserByEmail(ctx context.Context, email string) (bool, error) {
	var id int64

	row := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", email)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("user not found", zap.String("email", email))
			return false, nil
		}

		logger.Log.Debug("error getting user by email", zap.String("email", email), zap.Error(err))
		return false, err
	}

	return true, nil
}

func (s *Repository) FindUserByEmail(ctx context.Context, r *model.LoginRequest) (model.User, error) {
	var user model.User
	sqlQuery := `SELECT id, email, first_name, last_name, password FROM users WHERE email = $1`
	row := s.db.QueryRowContext(ctx, sqlQuery, r.Email)

	if err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.UserNotFoundError
		}

		return model.User{}, err
	}

	return user, nil
}
