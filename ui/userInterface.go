package ui

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed **/*.html
var web embed.FS

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
	t := template.Must(template.ParseFS(web, "pages/_base.html", "pages/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	}
}
