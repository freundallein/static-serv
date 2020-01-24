package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	ErrNoOptions = errors.New("no options provided")
)

// Options - fileserver parameters
type Options struct {
	Port    string
	RootDir string
	Prefix  string
}

// Server - main control struct
type Server struct {
	options *Options
}

// New - service constructor
func New(options *Options) (*Server, error) {
	if options == nil {
		return nil, ErrNoOptions
	}
	return &Server{options: options}, nil
}

// Run - start fileserver
func (srv *Server) Run() error {
	log.Printf("[server] Start listening on :%s\n", srv.options.Port)
	staticDIr := http.Dir(srv.options.RootDir)
	fileserver := http.FileServer(staticDIr)
	chain := MiddlewareChain(
		fileserver,
		RestrictListing,
		Cache(60 * time.Second),  // Cached items expiration timeout
		GetMethodOnly,
		AccessLog,
	)
	handler := http.StripPrefix(srv.options.Prefix, chain)

	mux := http.NewServeMux()
	mux.Handle("/static/", handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	addr := fmt.Sprintf("0.0.0.0:%s", srv.options.Port)
	serv := &http.Server{
		Handler:        mux,
		Addr:           addr,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := serv.ListenAndServe()
	return err
}
