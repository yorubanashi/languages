package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	config *Config
	logger *log.Logger
	mux    *http.Server
}

func New(config *Config, logger *log.Logger) *Server {
	mux := http.NewServeMux()
	return &Server{
		config: config,
		logger: logger,
		mux:    &http.Server{Addr: config.Server.Address, Handler: mux},
	}
}

func (s *Server) Start() {
	s.logger.Println(fmt.Sprintf("Starting server on %s...", s.mux.Addr))

	go func() {
		if err := s.mux.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.logger.Println(fmt.Sprintf("Error during ListenAndServer: %v", err))
			}
		}
	}()
}

func (s *Server) Stop() {
	// TODO: Set a timeout here
	s.logger.Println("Server shutting down...")
	s.mux.Shutdown(context.Background())
	s.logger.Println("Server successfully shut down!")
}
