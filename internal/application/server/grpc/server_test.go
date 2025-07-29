package grpc

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	pb "github.com/bubaew95/yandex-diplom-2/internal/proto"
	"github.com/bubaew95/yandex-diplom-2/pkg/crypto"
	"github.com/bubaew95/yandex-diplom-2/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLoginInterceptor(t *testing.T) {
	interceptor := LoginInterceptor()

	validUser := model.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	validToken, _ := token.EncodeJWTToken(validUser)

	tests := []struct {
		name        string
		method      string
		md          metadata.MD
		expectCode  codes.Code
		expectCalls bool
	}{
		{
			name:       "Публичный метод — Registration",
			method:     "/gokeeper.GoKeeper/Registration",
			md:         nil,
			expectCode: codes.OK,
		},
		{
			name:       "Публичный метод — Login",
			method:     "/gokeeper.GoKeeper/Login",
			md:         nil,
			expectCode: codes.OK,
		},
		{
			name:       "Метаданные отсутствуют",
			method:     "/gokeeper.GoKeeper/Add",
			md:         nil,
			expectCode: codes.Unauthenticated,
		},
		{
			name:       "Токен отсутствует",
			method:     "/gokeeper.GoKeeper/Add",
			md:         metadata.New(nil),
			expectCode: codes.Unauthenticated,
		},
		{
			name:   "Неверный формат токена",
			method: "/gokeeper.GoKeeper/Add",
			md: metadata.New(map[string]string{
				"token": "NotBearerFormat",
			}),
			expectCode: codes.Unauthenticated,
		},
		{
			name:   "Недействительный токен",
			method: "/gokeeper.GoKeeper/Add",
			md: metadata.New(map[string]string{
				"token": "Bearer invalidtoken",
			}),
			expectCode: codes.Unauthenticated,
		},
		{
			name:   "Корректный токен",
			method: "/gokeeper.GoKeeper/Add",
			md: metadata.New(map[string]string{
				"token": "Bearer " + validToken,
			}),
			expectCode: codes.OK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.md != nil {
				ctx = metadata.NewIncomingContext(ctx, tc.md)
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: tc.method,
			}

			var handlerCalled bool
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				handlerCalled = true
				if tc.expectCode == codes.OK {
					user, ok := ctx.Value(crypto.KeyUser).(model.User)
					if tc.method != "/gokeeper.GoKeeper/Login" && tc.method != "/gokeeper.GoKeeper/Registration" {
						assert.True(t, ok)
						assert.Equal(t, validUser.Email, user.Email)
					}
				}
				return "ok", nil
			}

			resp, err := interceptor(ctx, nil, info, handler)

			if tc.expectCode == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, "ok", resp)
				assert.True(t, handlerCalled)
			} else {
				st, _ := status.FromError(err)
				assert.Equal(t, tc.expectCode, st.Code())
			}
		})
	}
}

func TestRegistration(t *testing.T) {
	t.Parallel()

	t.Run("Успешная регистрация", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)

		req := &pb.RegistrationRequest{
			Email:      "test@example.com",
			Password:   "password",
			RePassword: "password",
			FirstName:  "Test",
			LastName:   "User",
		}
		mockService.On("AddUser", mock.Anything, mock.Anything).
			Return(&model.AuthResponse{Token: "jwt-token"}, nil)

		resp, err := srv.Registration(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token", resp.Token)
	})

	t.Run("Пароли не совпадают", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.RegistrationRequest{
			Email:      "test@example.com",
			Password:   "pass1",
			RePassword: "pass2",
			FirstName:  "Test",
			LastName:   "User",
		}

		_, err := srv.Registration(context.Background(), req)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Пользователь уже существует", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.RegistrationRequest{
			Email:      "exist@example.com",
			Password:   "password",
			RePassword: "password",
			FirstName:  "Test",
			LastName:   "User",
		}
		mockService.On("AddUser", mock.Anything, mock.Anything).
			Return(nil, model.UserAlreadyExistsError)

		_, err := srv.Registration(context.Background(), req)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.AlreadyExists, st.Code())
	})

	t.Run("Внутренняя ошибка", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.RegistrationRequest{
			Email:      "err@example.com",
			Password:   "password",
			RePassword: "password",
			FirstName:  "Test",
			LastName:   "User",
		}
		mockService.On("AddUser", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		_, err := srv.Registration(context.Background(), req)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestLogin(t *testing.T) {
	t.Run("Успешный вход", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		mockService.On("Login", mock.Anything, mock.Anything).
			Return(model.AuthResponse{Token: "jwt-token"}, nil)

		resp, err := srv.Login(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token", resp.Token)
	})

	t.Run("Неверный логин/пароль", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrong",
		}
		mockService.On("Login", mock.Anything, mock.Anything).
			Return(model.AuthResponse{}, model.LoginAndPasswordError)

		_, err := srv.Login(context.Background(), req)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("Внутренняя ошибка", func(t *testing.T) {
		mockService := NewMockService(t)
		srv := NewServer(mockService)
		req := &pb.LoginRequest{
			Email:    "err@example.com",
			Password: "password",
		}
		mockService.On("Login", mock.Anything, mock.Anything).
			Return(model.AuthResponse{}, errors.New("db fail"))

		_, err := srv.Login(context.Background(), req)
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestAdd(t *testing.T) {
	mockService := NewMockService(t)
	srv := NewServer(mockService)

	ctx := context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1})
	t.Run("Успешное добавление", func(t *testing.T) {
		req := &pb.DataRequest{
			Text: "secret",
			Type: "text",
		}

		mockService.On("Add", mock.Anything, mock.Anything).
			Return(model.TextResponse{ID: 1, Text: "secret", UserID: 1}, nil).Once()

		resp, err := srv.Add(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.Id)
	})

	t.Run("Внутренняя ошибка", func(t *testing.T) {
		mockService.On("Add", mock.Anything, mock.Anything).
			Return(model.TextResponse{}, errors.New("internal error")).Once()

		_, err := srv.Add(ctx, &pb.DataRequest{Text: "fail", Type: "text"})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestEdit(t *testing.T) {
	mockService := NewMockService(t)
	srv := NewServer(mockService)

	ctx := context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1})
	t.Run("Успешное редактирование", func(t *testing.T) {
		req := &pb.DataEditRequest{Id: 1, Text: "new text"}
		mockService.On("Edit", mock.Anything, mock.Anything).
			Return(model.TextResponse{ID: 1, Text: "new text", UserID: 1}, nil).Once()

		resp, err := srv.Edit(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "new text", resp.Text)
	})

	t.Run("Внутренняя ошибка", func(t *testing.T) {
		mockService.On("Edit", mock.Anything, mock.Anything).
			Return(model.TextResponse{}, errors.New("fail")).Once()

		_, err := srv.Edit(ctx, &pb.DataEditRequest{Id: 1, Text: "fail"})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestDelete(t *testing.T) {
	mockService := NewMockService(t)
	srv := NewServer(mockService)

	ctx := context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1})
	t.Run("Успешное удаление", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, int64(1)).
			Return(nil).Once()

		resp, err := srv.Delete(ctx, &pb.IdRequest{Id: 1})
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("Запись не найдена", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, int64(2)).
			Return(model.NotFoundError).Once()

		_, err := srv.Delete(ctx, &pb.IdRequest{Id: 2})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("Нет доступа", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, int64(3)).
			Return(model.AccessDeniedError).Once()

		_, err := srv.Delete(ctx, &pb.IdRequest{Id: 3})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Внутренняя ошибка", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, int64(4)).
			Return(errors.New("fail")).Once()

		_, err := srv.Delete(ctx, &pb.IdRequest{Id: 4})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
	})
}

func TestFindAll(t *testing.T) {
	mockService := NewMockService(t)
	srv := NewServer(mockService)

	ctx := context.WithValue(context.Background(), crypto.KeyUser, model.User{ID: 1})
	t.Run("Успешный возврат", func(t *testing.T) {
		mockService.On("FindAll", mock.Anything).
			Return([]*pb.DataResponse{{Id: 1, Text: "text", UserId: 1}}, nil).Once()

		resp, err := srv.FindAll(ctx, &pb.EmptyRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.List, 1)
	})

	t.Run("Ошибка получения", func(t *testing.T) {
		mockService.On("FindAll", mock.Anything).
			Return(nil, errors.New("fail")).Once()

		_, err := srv.FindAll(ctx, &pb.EmptyRequest{})
		st, _ := status.FromError(err)
		assert.Equal(t, codes.Internal, st.Code())
	})
}
