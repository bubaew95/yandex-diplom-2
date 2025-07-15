package middleware

import (
	"context"
	"github.com/bubaew95/yandex-diplom-2/internal/application/server/model"
	"github.com/bubaew95/yandex-diplom-2/pkg/response"
	"github.com/bubaew95/yandex-diplom-2/pkg/token"
	"net/http"
)

func AuthMiddleware() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			cookieToken, err := request.Cookie("token")
			if err != nil {
				response.WriteResponse(writer, http.StatusUnauthorized, model.Error("Authorization denied", http.StatusUnauthorized))
				return
			}

			user, err := token.DecodeJWTToken(cookieToken.Value)
			if err != nil {
				response.WriteResponse(writer, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(request.Context(), "user", user)
			nRequest := request.WithContext(ctx)
			
			next.ServeHTTP(writer, nRequest)
		})
	}
}
