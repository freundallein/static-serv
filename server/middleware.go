package server

import (
	"github.com/freundallein/static-serv/cache"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
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

// Cache - caches files
func Cache(expiration time.Duration) func(next http.Handler) http.Handler {
	store := cache.New(expiration)
	go store.GarbageCollect(60 * time.Second)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uri := r.URL.Path
			response, ok := store.Get(uri)
			if ok {
				for key, value := range response.Header() {
					w.Header().Set(key, strings.Join(value, ","))
				}
				w.Write(response.Data())
				log.Printf("[server] %s (cached)\n", r.URL.Path)
				return
			}
			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)
			result := recorder.Result()
			if result.StatusCode < 400 {
				response := cache.NewItem(recorder)
				store.Set(uri, response)
			}
			for key, value := range result.Header {
				w.Header().Set(key, strings.Join(value, ","))
			}
			w.Write(recorder.Body.Bytes())
		})
	}
}
