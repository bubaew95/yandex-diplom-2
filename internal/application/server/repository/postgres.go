package repository

import (
	"context"
	"database/sql"
	"errors"
	infra "github.com/bubaew95/yandex-diplom-2/internal/infra/database"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Repository struct {
	db *infra.DataBase
}

func NewRepository(db *infra.DataBase) *Repository {
	return &Repository{db: db}
}

func (s *Repository) CreateUser(ctx context.Context, r *model.RegistrationDTO) (int64, error) {
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
func (s *Repository) FindUserByEmail(ctx context.Context, r *model.LoginDTO) (model.User, error) {
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

//Text table

func (s *Repository) AddText(ctx context.Context, r *model.Data, userID int64) (int64, error) {
	sqlQuery := `INSERT INTO data (text, user_id, type) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	row := s.db.QueryRowContext(ctx, sqlQuery, r.Text, userID, r.Type)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}
func (s *Repository) Edit(ctx context.Context, r *model.TextRequest, userID int64) (int64, error) {
	data, err := s.GetText(ctx, r.ID)
	if err != nil {
		return -1, err
	}

	if userID != data.UserID {
		return -1, model.AccessDeniedError
	}

	if data.Text == r.Text {
		return -1, model.DataNotChangedError
	}

	_, err = s.db.ExecContext(ctx, `UPDATE data SET text = $1 WHERE id = $2`, r.Text, r.ID)
	if err != nil {
		return -1, err
	}

	return r.ID, nil
}
func (s *Repository) Delete(ctx context.Context, userID int64, ID int64) error {
	textData, err := s.GetText(ctx, ID)
	if err != nil {
		return err
	}

	if textData.IsDeleted == true {
		return model.NotFoundError
	}

	if textData.UserID != userID {
		return model.AccessDeniedError
	}

	_, err = s.db.ExecContext(ctx, `UPDATE data SET is_deleted = $1 WHERE id = $2`, true, ID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Repository) GetText(ctx context.Context, ID int64) (model.TextResponse, error) {
	var text model.TextResponse

	sqlQuery := `SELECT id, text, user_id, is_deleted FROM data WHERE id = $1`
	row := s.db.QueryRowContext(ctx, sqlQuery, ID)
	if err := row.Scan(&text.ID, &text.Text, &text.UserID, &text.IsDeleted); err != nil {
		return model.TextResponse{}, err
	}

	return text, nil
}
func (s *Repository) FindAll(ctx context.Context, userID int64) ([]*pb.DataResponse, error) {
	sqlQuery := `SELECT id, text, type, user_id, is_deleted FROM data WHERE user_id = $1 AND is_deleted = false ORDER BY id DESC`
	rows, err := s.db.QueryContext(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*pb.DataResponse, 0)
	for rows.Next() {
		var text pb.DataResponse
		if err := rows.Scan(&text.Id, &text.Text, &text.Type, &text.UserId, &text.IsDeleted); err != nil {
			return nil, err
		}

		list = append(list, &text)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
