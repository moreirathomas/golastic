package httputil

import (
	"fmt"
	"log"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides
// extra information on the request.
type ResponseWriter struct {
	http.ResponseWriter
	Status int // Records the status code of the request.
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.Status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	return rw.ResponseWriter.Write(data)
}

// RequestLogger adds logging to the given http.Handler.
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		h.ServeHTTP(rw, r)

		c := statusColor(rw.Status)
		log.Printf("%s %s -> %s %s", r.Method, r.URL.String(), c(rw.Status), c(http.StatusText(rw.Status)))
	})
}

var (
	success   = color("\033[1;32m%s\033[0m")
	clientErr = color("\033[1;33m%s\033[0m")
	serverErr = color("\033[1;31m%s\033[0m")
)

func color(a string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(a, fmt.Sprint(args...))
	}
	return sprint
}

func statusColor(statusCode int) func(...interface{}) string {
	switch {
	case statusCode >= 500:
		return serverErr
	case statusCode >= 400:
		return clientErr
	default:
		return success
	}
}
