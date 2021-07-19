package logger

import (
	"log"
	"net/http"
)

// responseWriter is a wrapper around http.ResponseWriter that provides
// extra information on the request.
type responseWriter struct {
	http.ResponseWriter
	status int // Records the status code of the request.
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Middleware adds logging to the given http.Handler.
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		h.ServeHTTP(rw, r)

		c := statusColor(rw.status)
		log.Printf(
			"%s %s -> %s %s",
			r.Method, r.URL.String(), c(rw.status), c(http.StatusText(rw.status)),
		)
	})
}
