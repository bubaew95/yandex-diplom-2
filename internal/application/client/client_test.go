package client

import (
	"context"
	"errors"
	"github.com/bubaew95/yandex-diplom-2/internal/model"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/bubaew95/yandex-diplom-2/internal/proto" // проверь, может быть pb
	"github.com/stretchr/testify/assert"
)

func TestClient_Login(t *testing.T) {
	t.Run("успешная авторизация", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()

		mockClient.On("Login", ctx, &proto.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}).Return(&proto.TokenResponse{Token: "jwt-token"}, nil)

		client := &Client{
			State:        &State{},
			KeeperClient: mockClient,
		}

		token, err := client.Login(ctx, "test@example.com", "password")
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token", token)
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка авторизации", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()

		mockClient.On("Login", ctx, &proto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpass",
		}).Return(nil, errors.New("unauthorized"))

		client := &Client{
			State:        &State{},
			KeeperClient: mockClient,
		}

		token, err := client.Login(ctx, "test@example.com", "wrongpass")
		assert.Error(t, err)
		assert.Empty(t, token)
		mockClient.AssertExpectations(t)
	})
}

func TestClient_GetAllData(t *testing.T) {
	t.Run("успешное получение данных", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()

		mockClient.On("FindAll", mock.Anything, &proto.EmptyRequest{}).
			Return(&proto.DataList{List: []*proto.DataResponse{}}, nil)

		c := &Client{
			State:        &State{Token: "test-token"},
			KeeperClient: mockClient,
		}

		resp, err := c.GetAllData(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка получения данных", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()

		mockClient.On("FindAll", mock.Anything, &proto.EmptyRequest{}).
			Return(nil, errors.New("db error"))

		c := &Client{
			State:        &State{Token: "test-token"},
			KeeperClient: mockClient,
		}

		resp, err := c.GetAllData(ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockClient.AssertExpectations(t)
	})
}

func TestClient_Add(t *testing.T) {
	t.Run("успешное добавление", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		input := &model.Data{Text: "test", Type: "text"}

		mockClient.On("Add", mock.Anything, &proto.DataRequest{Text: "test", Type: "text"}).
			Return(&proto.DataResponse{}, nil)

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Add(ctx, input)
		assert.NoError(t, err)
		assert.True(t, success)
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка при добавлении", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		input := &model.Data{Text: "bad", Type: "text"}

		mockClient.On("Add", mock.Anything, &proto.DataRequest{Text: "bad", Type: "text"}).
			Return(nil, errors.New("insert error"))

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Add(ctx, input)
		assert.Error(t, err)
		assert.False(t, success)
		mockClient.AssertExpectations(t)
	})
}

func TestClient_Edit(t *testing.T) {
	t.Run("успешное редактирование", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		input := &model.Data{Text: "updated"}

		mockClient.On("Edit", mock.Anything, &proto.DataEditRequest{Id: 1, Text: "updated"}).
			Return(&proto.DataResponse{}, nil)

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Edit(ctx, 1, input)
		assert.NoError(t, err)
		assert.True(t, success)
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка при редактировании", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		input := &model.Data{Text: "bad"}

		mockClient.On("Edit", mock.Anything, &proto.DataEditRequest{Id: 2, Text: "bad"}).
			Return(nil, errors.New("update error"))

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Edit(ctx, 2, input)
		assert.Error(t, err)
		assert.False(t, success)
		mockClient.AssertExpectations(t)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("успешное удаление", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		idReq := &proto.IdRequest{Id: 1}

		mockClient.On("Delete", mock.Anything, idReq).
			Return(&proto.SuccessResponse{Success: true}, nil)

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Delete(ctx, idReq)
		assert.NoError(t, err)
		assert.True(t, success)
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка при удалении", func(t *testing.T) {
		mockClient := &MockGoKeeperClient{}
		ctx := context.Background()
		idReq := &proto.IdRequest{Id: 2}

		mockClient.On("Delete", mock.Anything, idReq).
			Return(nil, errors.New("delete error"))

		c := &Client{
			State:        &State{Token: "tok"},
			KeeperClient: mockClient,
		}

		success, err := c.Delete(ctx, idReq)
		assert.Error(t, err)
		assert.False(t, success)
		mockClient.AssertExpectations(t)
	})
}
