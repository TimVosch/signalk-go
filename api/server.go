package api

import (
	"encoding/json"
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

func Send(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	r := s.r
	r.Get("/", s.httpGetGeneralInformation())
	r.Route("/v1/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	})
}

func (s *Server) httpGetGeneralInformation() http.HandlerFunc {
	type Endpoint struct {
		Version string `json:"version"`
		HTTP    string `json:"signalk-http"`
		WS      string `json:"signalk-ws"`
		TCP     string `json:"signalk-tcp"`
	}
	type ServerInfo struct {
		ID      string `json:"id"`
		Version string `json:"version"`
	}
	type Hello struct {
		Endpoints map[string]Endpoint `json:"endpoints"`
		Server    ServerInfo          `json:"server"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		Send(w, Hello{
			Endpoints: map[string]Endpoint{
				"v1": {
					Version: "1.0.0",
					HTTP:    "http://127.0.0.1:3000/signalk/v1/api",
				},
			},
			Server: ServerInfo{
				ID:      "signalk-go-server",
				Version: "0.0.0",
			},
		})
	}
}
