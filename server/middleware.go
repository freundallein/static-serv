package server

import (
	"log"
	"net/http"
	"strings"
)

// Middleware - http middleware
type Middleware func(http.Handler) http.Handler

// MiddlewareChain - chain multiple middlewares
func MiddlewareChain(handler http.Handler, midllewares ...Middleware) http.Handler {
	if len(midllewares) < 1 {
		return handler
	}
	wrapped := handler
	for i := 0; i < len(midllewares); i++ {
		wrapped = midllewares[i](wrapped)
	}
	return wrapped
}

// RestrictListing - prevent directory listing
func RestrictListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AccessLog - log cliet requests
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientIP string
		clientIP = r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}
		log.Printf("[server] %s requested %s %s\n", clientIP, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// GetMethodOnly - return 405 for all except GET
func GetMethodOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
