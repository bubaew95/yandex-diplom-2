package handlers

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	type response struct {
		data model.RegistrationResponse
		err  error
	}

	type request struct {
		data string
		code int
	}

	tests := []struct {
		name     string
		request  request
		response response
	}{
		{
			name: "created success",
			request: request{
				data: `{ "first_name": "test", "last_name": "test", "email": "tesd33d14@mail.ru", "password": "1234", "re_password": "1234" }`,
				code: http.StatusCreated,
			},
			response: response{
				data: model.RegistrationResponse{
					User: model.User{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "test@mailr",
						Password:  "1234",
					},
				},
				err: nil,
			},
		},
		{
			name: "email already exists",
			request: request{
				data: `{ "first_name": "test", "last_name": "test", "email": "tesd33d14@mail.ru", "password": "1234", "re_password": "1234" }`,
				code: http.StatusInternalServerError,
			},
			response: response{
				data: model.RegistrationResponse{},
				err: model.ErrorResponse{
					Errors: map[string]string{
						"message": "user with this email already exists",
					},
				},
			},
		},
		{
			name: "password doesn't match",
			request: request{
				data: `{ "first_name": "test", "last_name": "test", "email": "tesd33d14@mail.ru", "password": "12341", "re_password": "1234" }`,
				code: http.StatusInternalServerError,
			},
			response: response{
				data: model.RegistrationResponse{},
				err: model.ErrorResponse{
					Errors: map[string]string{
						"message": "password does not match",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockService(t)
			handler := NewHandler(mockService, &config.Config{})

			router := chi.NewRouter()
			router.Post("/api/v1/register", handler.Register)

			ts := httptest.NewServer(router)
			defer ts.Close()

			var registerRequest model.RegistrationRequest
			err := json.Unmarshal([]byte(tt.request.data), &registerRequest)
			require.NoError(t, err)

			mockService.On("AddUser", mock.Anything, &registerRequest).Return(tt.response.data, tt.response.err)

			resp, err := http.Post(ts.URL+"/api/v1/register", "application/json", strings.NewReader(tt.request.data))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.request.code, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		})
	}
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
	}{
		{
			name: "success",
			data: `{ "login": "test", "password": "test" }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockService(t)
			handler := NewHandler(mockService, &config.Config{})
			router := chi.NewRouter()
			router.Post("/api/v1/login", handler.Login)
			ts := httptest.NewServer(router)
			defer ts.Close()
		})
	}
}
