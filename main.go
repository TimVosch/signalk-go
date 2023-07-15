package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"signalk/api"
	"signalk/ui"
)

func Run() error {
	uiServer := ui.NewServer()
	apiTransport := api.NewServer()
	// Setup HTTP router
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Mount("/signalk", apiTransport)
	r.Mount("/", uiServer)

	http.ListenAndServe(":3000", cors.AllowAll().Handler(r))
	return nil
}
