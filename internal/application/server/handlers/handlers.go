package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"github.com/bubaew95/yandex-diplom-2/pkg/response"
	"go.uber.org/zap"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Service --filename=servicemock_test.go --inpackage
type Service interface {
	AddUser(ctx context.Context, r *model.RegistrationRequest) (model.RegistrationResponse, error)
	Login(ctx context.Context, r *model.LoginRequest) (model.AuthResponse, error)
}

type Handler struct {
	service Service
	config  *config.Config
}

func NewHandler(service Service, config *config.Config) *Handler {
	return &Handler{service: service, config: config}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), &loginRequest)
	if err != nil {
		response.WriteResponse(w, http.StatusOK, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    user.Token,
		Path:     "/",
		HttpOnly: true,
	})

	response.WriteResponse(w, http.StatusOK, user)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest model.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.service.AddUser(r.Context(), &registerRequest)
	if err != nil {
		if errors.Is(err, model.UserAlreadyExistsError) {
			response.WriteResponse(w, http.StatusFound, model.Error(err.Error(), http.StatusFound))
			return
		}

		logger.Log.Debug("error adding user", zap.Error(err))
		response.WriteResponse(w, http.StatusInternalServerError,
			model.Error(err.Error(), http.StatusInternalServerError),
		)
		return
	}

	response.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Println(user)
}
