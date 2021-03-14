package server

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server struct points to server instance.
type Server struct {
	srv *http.Server
}

// ServerOptions is functional options
type ServerOptions func(s *Server) error

// GetServer function returns pointer to server instance.
func GetServer(opts ...ServerOptions) (*Server, error) {
	s := &Server{
		srv: &http.Server{},
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// WithAddr returns method to specify server address.
func WithAddr(addr string) ServerOptions {
	return func(s *Server) error {
		s.srv.Addr = addr
		return nil
	}
}

// WithErrLogger returns method to invoke custom error logger.
func WithErrLogger(l *log.Logger) ServerOptions {
	return func(s *Server) error {
		s.srv.ErrorLog = l
		return nil
	}
}

// WithRouter returns method to invoke custom handler.
func WithRouter(router *mux.Router) ServerOptions {
	return func(s *Server) error {
		s.srv.Handler = router
		return nil
	}
}

// WithTLS returns method to invoke custom certificate
func WithTLS(cert *tls.Certificate) ServerOptions {
	return func(s *Server) error {
		s.srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}
		return nil
	}
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

// StartServerTLS returns ListenAndServeTLS method to start https server.
func (s *Server) StartServerTLS() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("Server missing handler")
	}

	return s.srv.ListenAndServeTLS("", "")
}

// CloseServer returns method to close server connection.
func (s *Server) CloseServer() error {
	return s.srv.Close()
}
