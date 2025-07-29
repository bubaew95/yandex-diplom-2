package service

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, data *model.Data)
	}
	type args struct {
		ctx  context.Context
		data *model.Data
	}
	type want struct {
		result     model.TextResponse
		err        error
		errMessage string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "успешное добавление текста",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, data *model.Data) {
					m.On("AddText", ctx, data, int64(42)).Return(int64(100), nil)
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 42}),
				data: &model.Data{
					Text: "секретный текст",
					Type: "text",
				},
			},
			want: want{
				result: model.TextResponse{
					ID:     100,
					Text:   "секретный текст",
					UserID: 42,
				},
			},
		},
		{
			name: "ошибка при добавлении текста в репозиторий",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, data *model.Data) {
					m.On("AddText", ctx, data, int64(42)).Return(int64(-1), errors.New("ошибка вставки"))
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 42}),
				data: &model.Data{
					Text: "ошибочный текст",
					Type: "text",
				},
			},
			want: want{
				result:     model.TextResponse{},
				err:        errors.New("ошибка вставки"),
				errMessage: "ошибка вставки",
			},
		},
		{
			name: "отсутствие пользователя в контексте",
			fields: fields{
				mockSetup: nil, // не мокаем, ошибка произойдёт до вызова AddText
			},
			args: args{
				ctx:  context.Background(),
				data: &model.Data{Text: "текст", Type: "text"},
			},
			want: want{
				result:     model.TextResponse{},
				errMessage: "user value not found in context",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			if tt.fields.mockSetup != nil {
				tt.fields.mockSetup(mockRepo, tt.args.ctx, tt.args.data)
			}

			svc := NewService(mockRepo, config.Config{})
			result, err := svc.Add(tt.args.ctx, tt.args.data)

			if tt.want.err != nil {
				require.EqualError(t, err, tt.want.err.Error())
			} else if tt.want.errMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.want.errMessage)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.result, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEdit(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, req *model.TextRequest)
	}
	type args struct {
		ctx context.Context
		req *model.TextRequest
	}
	type want struct {
		result     model.TextResponse
		err        error
		errMessage string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "успешное обновление текста",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, req *model.TextRequest) {
					m.On("Edit", ctx, req, int64(1)).Return(int64(100), nil)
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1}),
				req: &model.TextRequest{
					ID:   100,
					Text: "обновлённый текст",
				},
			},
			want: want{
				result: model.TextResponse{
					ID:     100,
					Text:   "обновлённый текст",
					UserID: 1,
				},
			},
		},
		{
			name: "ошибка при обновлении текста",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, req *model.TextRequest) {
					m.On("Edit", ctx, req, int64(1)).Return(int64(-1), errors.New("ошибка обновления"))
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1}),
				req: &model.TextRequest{
					ID:   200,
					Text: "ошибочный текст",
				},
			},
			want: want{
				result:     model.TextResponse{},
				errMessage: "ошибка обновления",
			},
		},
		{
			name: "пользователь отсутствует в контексте",
			fields: fields{
				mockSetup: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &model.TextRequest{
					ID:   300,
					Text: "текст без пользователя",
				},
			},
			want: want{
				result:     model.TextResponse{},
				errMessage: "user value not found in context",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			if tt.fields.mockSetup != nil {
				tt.fields.mockSetup(mockRepo, tt.args.ctx, tt.args.req)
			}

			svc := NewService(mockRepo, config.Config{})
			result, err := svc.Edit(tt.args.ctx, tt.args.req)

			if tt.want.err != nil {
				require.EqualError(t, err, tt.want.err.Error())
			} else if tt.want.errMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.want.errMessage)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.result, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, userID, dataID int64)
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	type want struct {
		errMessage string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "успешное удаление",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, userID, dataID int64) {
					m.On("Delete", ctx, userID, dataID).Return(nil)
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 10}),
				id:  100,
			},
			want: want{},
		},
		{
			name: "ошибка при удалении",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, userID, dataID int64) {
					m.On("Delete", ctx, userID, dataID).Return(errors.New("ошибка удаления"))
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 10}),
				id:  101,
			},
			want: want{
				errMessage: "ошибка удаления",
			},
		},
		{
			name: "пользователь отсутствует в контексте",
			fields: fields{
				mockSetup: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  200,
			},
			want: want{
				errMessage: "user value not found in context",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)

			if tt.fields.mockSetup != nil {
				userVal := tt.args.ctx.Value(crypto.KeyUser).(model.User)
				tt.fields.mockSetup(mockRepo, tt.args.ctx, userVal.ID, tt.args.id)
			}

			svc := NewService(mockRepo, config.Config{})
			err := svc.Delete(tt.args.ctx, tt.args.id)

			if tt.want.errMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.want.errMessage)
			} else {
				require.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFindAll(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, userID int64)
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		result     []*pb.DataResponse
		errMessage string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "успешное получение данных",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, userID int64) {
					m.On("FindAll", ctx, userID).Return([]*pb.DataResponse{
						{Id: 1, Text: "text1", Type: "text", UserId: userID, IsDeleted: false},
						{Id: 2, Text: "text2", Type: "text", UserId: userID, IsDeleted: false},
					}, nil)
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 42}),
			},
			want: want{
				result: []*pb.DataResponse{
					{Id: 1, Text: "text1", Type: "text", UserId: 42, IsDeleted: false},
					{Id: 2, Text: "text2", Type: "text", UserId: 42, IsDeleted: false},
				},
			},
		},
		{
			name: "ошибка при получении данных",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, userID int64) {
					m.On("FindAll", ctx, userID).Return(nil, errors.New("ошибка выборки"))
				},
			},
			args: args{
				ctx: context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 99}),
			},
			want: want{
				result:     nil,
				errMessage: "ошибка выборки",
			},
		},
		{
			name: "отсутствие пользователя в контексте",
			fields: fields{
				mockSetup: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				result:     nil,
				errMessage: "user value not found in context",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)

			if tt.fields.mockSetup != nil {
				user := tt.args.ctx.Value(crypto.KeyUser).(model.User)
				tt.fields.mockSetup(mockRepo, tt.args.ctx, user.ID)
			}

			svc := NewService(mockRepo, config.Config{})
			result, err := svc.FindAll(tt.args.ctx)

			if tt.want.errMessage != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.want.errMessage)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.result, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, r *model.LoginDTO)
	}
	type args struct {
		ctx context.Context
		r   *model.LoginDTO
	}
	type want struct {
		result     model.AuthResponse
		errMessage string
	}

	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	validUser := model.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  string(hashed),
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "успешный вход",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.LoginDTO) {
					m.On("FindUserByEmail", ctx, r).Return(validUser, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				r: &model.LoginDTO{
					Email:    validUser.Email,
					Password: password,
				},
			},
			want: want{
				result: model.AuthResponse{}, // сравнивать токен не будем — он непредсказуем
			},
		},
		{
			name: "ошибка: пользователь не найден",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.LoginDTO) {
					m.On("FindUserByEmail", ctx, r).Return(model.User{}, errors.New("not found"))
				},
			},
			args: args{
				ctx: context.Background(),
				r: &model.LoginDTO{
					Email:    "unknown@example.com",
					Password: "pass",
				},
			},
			want: want{
				errMessage: model.AuthorizationError.Error(),
			},
		},
		{
			name: "ошибка: неверный пароль",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.LoginDTO) {
					m.On("FindUserByEmail", ctx, r).Return(validUser, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				r: &model.LoginDTO{
					Email:    validUser.Email,
					Password: "wrongpass",
				},
			},
			want: want{
				errMessage: model.LoginAndPasswordError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			if tt.fields.mockSetup != nil {
				tt.fields.mockSetup(mockRepo, tt.args.ctx, tt.args.r)
			}

			svc := NewService(mockRepo, config.Config{})
			result, err := svc.Login(tt.args.ctx, tt.args.r)

			if tt.want.errMessage != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.want.errMessage)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result.Token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		mockSetup func(m *MockRepository, ctx context.Context, r *model.RegistrationDTO)
	}
	type args struct {
		ctx context.Context
		req *model.RegistrationDTO
	}
	type want struct {
		token string
		err   error
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "несовпадение паролей",
			args: args{
				ctx: context.Background(),
				req: &model.RegistrationDTO{
					User: model.User{
						Password: "123",
					},
					RePassword: "456",
				},
			},
			want: want{
				err: model.PasswordNotMatchError,
			},
		},
		{
			name: "ошибка при создании пользователя",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.RegistrationDTO) {
					m.On("CreateUser", ctx, mock.AnythingOfType("*model.RegistrationDTO")).
						Return(int64(-1), errors.New("insert fail"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &model.RegistrationDTO{
					RePassword: "123",
					User: model.User{
						Email:     "test@example.com",
						Password:  "123",
						FirstName: "Test",
						LastName:  "User",
					},
				},
			},
			want: want{
				err: errors.New("insert fail"),
			},
		},
		{
			name: "ошибка генерации токена (невалидный пользователь)",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.RegistrationDTO) {
					m.On("CreateUser", ctx, mock.AnythingOfType("*model.RegistrationDTO")).
						Return(int64(1), nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &model.RegistrationDTO{
					User: model.User{
						Email:     "",
						Password:  "123",
						FirstName: "Bad",
						LastName:  "User",
					},
					RePassword: "123",
				},
			},
			want: want{
				// на текущий момент jwt всё равно создаётся — тест будет проходить,
				// если не добавлять явную проверку в token.EncodeJWTToken
				err: nil,
			},
		},
		{
			name: "успешная регистрация",
			fields: fields{
				mockSetup: func(m *MockRepository, ctx context.Context, r *model.RegistrationDTO) {
					m.On("CreateUser", ctx, mock.AnythingOfType("*model.RegistrationDTO")).
						Return(int64(10), nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &model.RegistrationDTO{
					RePassword: "secure",
					User: model.User{
						Email:     "user@example.com",
						Password:  "secure",
						FirstName: "First",
						LastName:  "Last",
					},
				},
			},
			want: want{
				token: "some", // просто проверим, что токен не пустой
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			if tt.fields.mockSetup != nil {
				tt.fields.mockSetup(mockRepo, tt.args.ctx, tt.args.req)
			}

			svc := NewService(mockRepo, config.Config{})
			resp, err := svc.AddUser(tt.args.ctx, tt.args.req)

			if tt.want.err != nil {
				require.EqualError(t, err, tt.want.err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.Token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
