package middleware

import (
	"github.com/bubaew95/yandex-diplom-2/internal/compress"
	"net/http"
	"strings"
)

// GZipMiddleware - middleware, производит сжатие запроса.
func GZipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if isContentEncoding(r) {
			cr, err := compress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		if isAcceptEncoding(r) && isAccessContentType(r) {
			cw := compress.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

func isAcceptEncoding(r *http.Request) bool {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}

func isContentEncoding(r *http.Request) bool {
	contentEncoding := r.Header.Get("Content-Encoding")
	return strings.Contains(contentEncoding, "gzip")
}

func isAccessContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-type")
	return contentType == "application/json" || contentType == "text/html"
}
