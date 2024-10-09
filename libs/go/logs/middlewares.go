package logs

import (
	"net/http"
)

func ContextWithLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(toContext(r.Context(), &loggerCollector{}))
		next.ServeHTTP(w, r)
	})
}
