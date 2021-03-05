package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server struct points to server instance.
type Server struct {
	srv *http.Server
}

// GetServer function returns pointer to server instance.
func GetServer() *Server {
	return &Server{
		srv: &http.Server{},
	}
}

// WithAddr returns method to specify server address.
func (s *Server) WithAddr(addr string) *Server {
	s.srv.Addr = addr
	return s
}

// WithErrLogger returns method to invoke custom error logger.
func (s *Server) WithErrLogger(l *log.Logger) *Server {
	s.srv.ErrorLog = l
	return s
}

// WithRouter returns method to invoke custom handler.
func (s *Server) WithRouter(router *mux.Router) *Server {
	s.srv.Handler = router
	return s
}

// StartServer returns ListenAndServe method to start server.
func (s *Server) StartServer() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("Server missing handler")
	}

	return s.srv.ListenAndServe()
}

// CloseServer returns method to close server connection.
func (s *Server) CloseServer() error {
	return s.srv.Close()
}