package response

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diplom-2/internal/logger"
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
