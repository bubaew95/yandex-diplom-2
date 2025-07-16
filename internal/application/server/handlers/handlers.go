package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bubaew95/yandex-diplom-2/config"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/bubaew95/yandex-diplom-2/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Service --filename=servicemock_test.go --inpackage
type Service interface {
	AddUser(ctx context.Context, r *model.RegistrationRequest) (model.RegistrationResponse, error)
	Login(ctx context.Context, r *model.LoginRequest) (model.AuthResponse, error)
	AddText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	EditText(ctx context.Context, r *model.TextRequest) (model.TextResponse, error)
	DeleteText(ctx context.Context, ID int64) error
	AddBinary(ctx context.Context, r *model.BinaryRequest) (model.BinaryResponse, error)
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
		response.WriteErrors(w, err)
		return
	}

	user, err := h.service.Login(r.Context(), &loginRequest)
	if err != nil {
		response.WriteErrors(w, err)
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
		response.WriteErrors(w, err)
		return
	}

	if valid := registerRequest.Validate(); len(valid) != 0 {
		response.WriteResponse(w, http.StatusUnprocessableEntity, model.ValidationError(valid))
		return
	}

	user, err := h.service.AddUser(r.Context(), &registerRequest)
	if err != nil {
		response.WriteErrors(w, err)
		return
	}

	response.WriteResponse(w, http.StatusCreated, user)
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Println(user)
}

func (h *Handler) AddText(w http.ResponseWriter, r *http.Request) {
	var textRequest model.TextRequest
	if err := json.NewDecoder(r.Body).Decode(&textRequest); err != nil {
		response.WriteErrors(w, err)
		return
	}

	text, err := h.service.AddText(r.Context(), &textRequest)
	if err != nil {
		response.WriteErrors(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, text)
}

func (h *Handler) EditText(w http.ResponseWriter, r *http.Request) {
	var textRequest model.TextRequest
	if err := json.NewDecoder(r.Body).Decode(&textRequest); err != nil {
		response.WriteErrors(w, err)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteErrors(w, err)
		return
	}

	textRequest.ID = id
	text, err := h.service.EditText(r.Context(), &textRequest)
	if err != nil {
		response.WriteErrors(w, err)
		return
	}

	response.WriteResponse(w, http.StatusOK, text)
}

func (h *Handler) DeleteText(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteErrors(w, err)
		return
	}

	if err := h.service.DeleteText(r.Context(), id); err != nil {
		response.WriteErrors(w, err)
		return
	}

	response.WriteResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) AddBinary(w http.ResponseWriter, r *http.Request) {
	//var binaryRequest model.BinaryRequest

}
