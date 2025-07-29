package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	infra "github.com/bubaew95/yandex-diplom-2/internal/infra/database"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func setup(t *testing.T) (Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := Repository{db: &infra.DataBase{DB: db}}

	return repo, mock, func() { db.Close() }
}

func TestGetUserByEmail(t *testing.T) {
	repo, mock, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name      string
		email     string
		mockSetup func()
		wantFound bool
		wantErr   error
	}{
		{
			name:  "Пользователь уже существует",
			email: "exists@example.com",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs("exists@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantFound: true,
			wantErr:   nil,
		},
		{
			name:  "Пользователь не найден",
			email: "missing@example.com",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs("missing@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantFound: false,
			wantErr:   nil,
		},
		{
			name:  "Ошибка",
			email: "error@example.com",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs("error@example.com").
					WillReturnError(errors.New("db error"))
			},
			wantFound: false,
			wantErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			found, err := repo.GetUserByEmail(ctx, tt.email)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantFound, found)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateUser(t *testing.T) {
	repo, mock, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	sampleReg := &model.RegistrationDTO{
		User: model.User{
			Email:     "new@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "hashed",
		},
	}

	tests := []struct {
		name      string
		setupMock func()
		wantID    int64
		wantErr   error
	}{
		{
			name: "Пользователь уже существует",
			setupMock: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs(sampleReg.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantID:  -1,
			wantErr: model.UserAlreadyExistsError,
		},
		{
			name: "Пользователь успешно создан",
			setupMock: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs(sampleReg.Email).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sampleReg.Email, sampleReg.FirstName, sampleReg.LastName, sampleReg.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
			},
			wantID:  42,
			wantErr: nil,
		},
		{
			name: "Ошибка создания пользователя",
			setupMock: func() {
				mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).
					WithArgs(sampleReg.Email).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(sampleReg.Email, sampleReg.FirstName, sampleReg.LastName, sampleReg.Password).
					WillReturnError(errors.New("insert error"))
			},
			wantID:  -1,
			wantErr: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			id, err := repo.CreateUser(ctx, sampleReg)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantID, id)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindUserById(t *testing.T) {
	repo, mock, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	loginDTO := &model.LoginDTO{
		Email:    "user@example.com",
		Password: "1234",
	}

	tests := []struct {
		name      string
		email     string
		mockSetup func()
		wantUser  model.User
		wantErr   error
	}{
		{
			name:  "Пользователь найден",
			email: loginDTO.Email,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, email, first_name, last_name, password FROM users WHERE email = \$1`).
					WithArgs(loginDTO.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "password"}).
						AddRow(int64(1), "user@example.com", "Ivan", "Ivanov", "hashed-password"))
			},
			wantUser: model.User{
				ID:        1,
				Email:     "user@example.com",
				FirstName: "Ivan",
				LastName:  "Ivanov",
				Password:  "hashed-password",
			},
			wantErr: nil,
		},
		{
			name:  "Пользователь не найден",
			email: "missing@example.com",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, email, first_name, last_name, password FROM users WHERE email = \$1`).
					WithArgs("missing@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser: model.User{},
			wantErr:  model.UserNotFoundError,
		},
		{
			name:  "Ошибка базы данных",
			email: "fail@example.com",
			mockSetup: func() {
				mock.ExpectQuery(`SELECT id, email, first_name, last_name, password FROM users WHERE email = \$1`).
					WithArgs("fail@example.com").
					WillReturnError(errors.New("ошибка подключения"))
			},
			wantUser: model.User{},
			wantErr:  errors.New("ошибка подключения"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := repo.FindUserByEmail(ctx, &model.LoginDTO{Email: tt.email})

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddText(t *testing.T) {
	repo, mock, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name      string
		input     *model.Data
		userID    int64
		mockSetup func()
		wantID    int64
		wantErr   bool
	}{
		{
			name: "Успешное добавление текста",
			input: &model.Data{
				Text: "тестовое сообщение",
				Type: "text",
			},
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO data `).
					WithArgs("тестовое сообщение", int64(1), "text").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))
			},
			wantID:  42,
			wantErr: false,
		},
		{
			name: "Ошибка выполнения запроса",
			input: &model.Data{
				Text: "ошибочный текст",
				Type: "text",
			},
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO data`).
					WithArgs("ошибочный текст", int64(2), "text").
					WillReturnError(errors.New("ошибка БД"))
			},
			wantID:  -1,
			wantErr: true,
		},
		{
			name: "Ошибка чтения id",
			input: &model.Data{
				Text: "некорректный id",
				Type: "text",
			},
			userID: 3,
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO data`).
					WithArgs("некорректный id", int64(3), "text").
					WillReturnRows(sqlmock.NewRows([]string{"invalid_column"}).AddRow("???"))
			},
			wantID:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			id, err := repo.AddText(ctx, tt.input, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestEdit(t *testing.T) {
	tests := []struct {
		name         string
		input        *model.TextRequest
		userID       int64
		getTextMock  *sqlmock.Rows
		getTextError error
		updateError  error
		expectedID   int64
		expectedErr  error
	}{
		{
			name: "успешное обновление текста",
			input: &model.TextRequest{
				ID:   1,
				Text: "новый текст",
			},
			userID: 1,
			getTextMock: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "старый текст", 1, false),
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "пользователь не владелец",
			input: &model.TextRequest{
				ID:   1,
				Text: "новый текст",
			},
			userID: 99,
			getTextMock: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "старый текст", 1, false),
			expectedID:  -1,
			expectedErr: model.AccessDeniedError,
		},
		{
			name: "текст не изменился",
			input: &model.TextRequest{
				ID:   1,
				Text: "старый текст",
			},
			userID: 1,
			getTextMock: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "старый текст", 1, false),
			expectedID:  -1,
			expectedErr: model.DataNotChangedError,
		},
		{
			name: "ошибка получения записи",
			input: &model.TextRequest{
				ID:   1,
				Text: "что угодно",
			},
			userID:       1,
			getTextError: errors.New("ошибка запроса"),
			expectedID:   -1,
			expectedErr:  errors.New("ошибка запроса"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := setup(t)
			defer cleanup()

			if tt.getTextMock != nil {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.input.ID).
					WillReturnRows(tt.getTextMock)
			} else if tt.getTextError != nil {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.input.ID).
					WillReturnError(tt.getTextError)
			}

			if tt.expectedErr == nil || errors.Is(tt.expectedErr, model.AccessDeniedError) || errors.Is(tt.expectedErr, model.DataNotChangedError) {
				// если дошло до UPDATE
				if tt.expectedErr == nil {
					mock.ExpectExec(`UPDATE data SET text = \$1 WHERE id = \$2`).
						WithArgs(tt.input.Text, tt.input.ID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				} else if tt.updateError != nil {
					mock.ExpectExec(`UPDATE data SET text = \$1 WHERE id = \$2`).
						WithArgs(tt.input.Text, tt.input.ID).
						WillReturnError(tt.updateError)
				}
			}

			id, err := repo.Edit(context.Background(), tt.input, tt.userID)

			assert.Equal(t, tt.expectedID, id)
			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputID     int64
		userID      int64
		getTextRow  *sqlmock.Rows
		getTextErr  error
		expectExec  bool
		execErr     error
		expectedErr error
	}{
		{
			name:    "успешное удаление",
			inputID: 1,
			userID:  42,
			getTextRow: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "text", 42, false),
			expectExec:  true,
			expectedErr: nil,
		},
		{
			name:    "данные уже удалены",
			inputID: 1,
			userID:  42,
			getTextRow: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "text", 42, true),
			expectedErr: model.NotFoundError,
		},
		{
			name:    "нет доступа к удалению (другой пользователь)",
			inputID: 1,
			userID:  99,
			getTextRow: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "text", 42, false),
			expectedErr: model.AccessDeniedError,
		},
		{
			name:        "ошибка при получении текста",
			inputID:     1,
			userID:      42,
			getTextErr:  errors.New("db error"),
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := setup(t)
			defer cleanup()

			if tt.getTextRow != nil {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.inputID).
					WillReturnRows(tt.getTextRow)
			} else if tt.getTextErr != nil {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.inputID).
					WillReturnError(tt.getTextErr)
			}

			if tt.expectExec {
				mock.ExpectExec(`UPDATE data SET is_deleted = \$1 WHERE id = \$2`).
					WithArgs(true, tt.inputID).
					WillReturnError(tt.execErr).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := repo.Delete(context.Background(), tt.userID, tt.inputID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputID     int64
		mockRows    *sqlmock.Rows
		mockError   error
		expected    model.TextResponse
		expectedErr error
	}{
		{
			name:    "успешное получение текста",
			inputID: 1,
			mockRows: sqlmock.NewRows([]string{"id", "text", "user_id", "is_deleted"}).
				AddRow(1, "example text", 42, false),
			expected: model.TextResponse{
				ID:        1,
				Text:      "example text",
				UserID:    42,
				IsDeleted: false,
			},
		},
		{
			name:        "текст не найден",
			inputID:     2,
			mockError:   sql.ErrNoRows,
			expectedErr: sql.ErrNoRows,
		},
		{
			name:        "ошибка базы данных",
			inputID:     3,
			mockError:   errors.New("db error"),
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := setup(t)
			defer cleanup()

			if tt.mockRows != nil {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.inputID).
					WillReturnRows(tt.mockRows)
			} else {
				mock.ExpectQuery(`SELECT id, text, user_id, is_deleted FROM data WHERE id = \$1`).
					WithArgs(tt.inputID).
					WillReturnError(tt.mockError)
			}

			result, err := repo.GetText(context.Background(), tt.inputID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      int64
		mockRows    *sqlmock.Rows
		mockError   error
		expected    []*pb.DataResponse
		expectedErr error
	}{
		{
			name:   "успешное получение списка данных",
			userID: 1,
			mockRows: sqlmock.NewRows([]string{"id", "text", "type", "user_id", "is_deleted"}).
				AddRow(int64(1), "text-1", "text", int64(1), false).
				AddRow(int64(2), "text-2", "binary", int64(1), false),
			expected: []*pb.DataResponse{
				{Id: 1, Text: "text-1", Type: "text", UserId: 1, IsDeleted: false},
				{Id: 2, Text: "text-2", Type: "binary", UserId: 1, IsDeleted: false},
			},
		},
		{
			name:        "ошибка выполнения запроса",
			userID:      2,
			mockError:   errors.New("db error"),
			expectedErr: errors.New("db error"),
		},
		{
			name:   "ошибка сканирования строки",
			userID: 3,
			mockRows: sqlmock.NewRows([]string{"id", "text", "type", "user_id", "is_deleted"}).
				AddRow("invalid", "text", "type", 3, false), // id должен быть int64
			expectedErr: errors.New("sql: Scan error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, cleanup := setup(t)
			defer cleanup()

			query := `SELECT id, text, type, user_id, is_deleted FROM data WHERE user_id = \$1 AND is_deleted = false ORDER BY id DESC`
			if tt.mockRows != nil {
				mock.ExpectQuery(query).
					WithArgs(tt.userID).
					WillReturnRows(tt.mockRows)
			} else {
				mock.ExpectQuery(query).
					WithArgs(tt.userID).
					WillReturnError(tt.mockError)
			}

			result, err := repo.FindAll(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
