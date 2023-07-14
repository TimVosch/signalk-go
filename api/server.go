package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	r chi.Router
}

func NewServer() *Server {
	srv := &Server{
		r: chi.NewRouter(),
	}
	srv.setupRoutes()

	return srv
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	r := s.r
	r.Get("/", s.httpIndex())
}

func (s *Server) httpIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	}
}
