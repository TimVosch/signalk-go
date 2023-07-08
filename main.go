package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/cors"

	"signalk/signalk"
)

var (
	self = &signalk.Vessel{
		ID: signalk.CreateVesselUUID(uuid.New()),
	}
	root = &signalk.Root{
		Version: "0.0.1",
		Self:    signalk.CreatePath("vessels", self.ID.String()),
		Vessels: signalk.VesselList{
			self.ID: self,
		},
	}
)

func main() {
	// Seed testdata
	if err := SeedData(); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding error: %v\n", err)
	}
	// Run server
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func SeedData() error {
	// Random stuff for testing
	err := root.ApplyDelta(signalk.Delta{
		Context: signalk.CreatePath(),
		Updates: []signalk.DeltaUpdate{
			{
				Values: []signalk.DeltaUpdateValues{
					{
						Path:  signalk.CreatePath("vessels.self.navigation.position"),
						Value: []byte(`{"latitude": 15.222, "longitude":3.251}`),
					},
					{
						Path:  signalk.CreatePath("vessels.urn:mrn:imo:mmsi:12345678.navigation.position"),
						Value: []byte(`{"latitude": 15.222, "longitude":9.221}`),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func Run() error {
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(root)
	})

	return http.ListenAndServe(":3000", cors.AllowAll().Handler(r))
}
