package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"signalk/api"
	"signalk/signalk"
	"signalk/ui"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func Run() error {
	svc := signalk.NewService()
	uiServer := ui.NewServer()
	apiTransport := api.NewServer(svc)
	// Setup HTTP router
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Mount("/signalk", apiTransport)
	r.Mount("/", uiServer)

	http.ListenAndServe(":3000", cors.AllowAll().Handler(r))
	return nil
}
