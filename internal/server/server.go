package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type HandlerFunc func(ctx context.Context, decode func(interface{}) error) (interface{}, error)

type Server struct {
	// User-provided structs
	config *Config
	logger *log.Logger

	// Underlying components
	mux   *http.Server
	cache map[string]interface{}
}

func New(config *Config, logger *log.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
		mux:    &http.Server{Addr: config.Server.Address, Handler: nil},
		cache:  make(map[string]interface{}),
	}
}

func (s *Server) Register() {
	mux := http.NewServeMux()
	for route, handler := range s.songRoutes() {
		mux.HandleFunc(route, translateHandler(handler))
	}

	mux.HandleFunc("/svelte", translateHandler(s.svelteWalkHandler))
	s.mux.Handler = mux
}

func (s *Server) Start() {
	s.logger.Printf("Starting server on %s...\n", s.mux.Addr)

	go func() {
		if err := s.mux.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.logger.Printf("Error during ListenAndServer: %v\n", err)
			}
		}
	}()

	// TODO: Should this just be its own binary?
	if s.config.StartOptions.Index {
		go func() { s.indexAll() }()
	}
}

func (s *Server) Stop() {
	timeout := time.Duration(s.config.Server.Timeouts.Shutdown) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.logger.Println("Server shutting down...")
	if err := s.mux.Shutdown(ctx); err != nil {
		s.logger.Println(err)
	} else {
		s.logger.Println("Server successfully shut down!")
	}
}
