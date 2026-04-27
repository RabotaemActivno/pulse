package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

type Config struct {
	Port string `default:"8080" envconfig:"HTTP_PORT"`
}

func New(handler http.Handler, c Config) *Server {

	httpserver := &http.Server{
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Addr:         net.JoinHostPort("", c.Port),
	}

	s := &Server{server: httpserver}

	go s.start()

	// TODO make log about server started

	return s
}

func (s *Server) start() {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {

	}
}

func (s *Server) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		// TODO make logs
	}

	// TODO make logs
}
