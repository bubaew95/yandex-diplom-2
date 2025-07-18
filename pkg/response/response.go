package response

import (
	"encoding/json"
	"errors"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			logger.Log.Debug("Error json encoded")
		}
	}
}

func WriteErrors(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := err.Error()

	switch {
	case errors.Is(err, model.UserNotFoundError):
		statusCode = http.StatusNotFound
	case errors.Is(err, model.UserAlreadyExistsError):
		statusCode = http.StatusConflict
	case errors.Is(err, model.DataNotChangedError):
		statusCode = http.StatusNotModified
	case errors.Is(err, model.AccessDeniedError):
		statusCode = http.StatusForbidden
	case errors.Is(err, model.NotFoundError):
		statusCode = http.StatusNotFound
	case errors.Is(err, model.LoginAndPasswordError):
	case errors.Is(err, model.AuthorizationError):
		statusCode = http.StatusUnauthorized
	}

	logger.Log.Debug(message, zap.Error(err))
	WriteResponse(w, statusCode, model.Error(message, statusCode))
}
