package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

func CookieMiddleware(isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &responseRecorder{ResponseWriter: w, headers: http.Header{}}
			next.ServeHTTP(rec, r)

			for key, values := range rec.Header() {
				if strings.ToLower(key) == "grpc-metadata-set-cookie" {
					for _, value := range values {
						parts := strings.SplitN(value, "=", 2)
						if len(parts) != 2 {
							continue
						}
						cookieName := parts[0]
						cookieValue := url.QueryEscape(parts[1])

						http.SetCookie(w, &http.Cookie{
							Name:     cookieName,
							Value:    cookieValue,
							HttpOnly: true,
							Secure:   isProd,
							Path:     "/",
						})
					}

					rec.Header().Del(key)
				}
			}

			if rec.statusCode == 0 {
				rec.statusCode = http.StatusOK
			}

			for key, values := range rec.headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(rec.statusCode)
			w.Write(rec.body)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	headers    http.Header
	statusCode int
	body       []byte
}

func (r *responseRecorder) Header() http.Header {
	return r.headers
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	r.body = append(r.body, body...)
	return len(body), nil
}
