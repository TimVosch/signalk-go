package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/rs/cors"

	"signalk/signalk"
)

type SignalKHello struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Roles     []string  `json:"roles"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	self signalk.Vessel
	root signalk.FullFormat
)

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.CheckOrigin = func(r *http.Request) bool { return true }
	done := make(chan struct{})
	u.OnMessage(func(c *websocket.Conn, mt websocket.MessageType, b []byte) {
		log.Printf("WS Message: %v %v\n", mt, string(b))
	})
	u.OnOpen(func(c *websocket.Conn) {
		log.Printf("New WS: %v\n", c.RemoteAddr())
		hello := SignalKHello{
			Name:      "Test server",
			Version:   "1.0.0",
			Roles:     []string{"master", "main"},
			Timestamp: time.Now(),
		}
		jsonBytes, _ := json.Marshal(&hello)
		c.WriteMessage(websocket.TextMessage, jsonBytes)

		go func() {
			for {
				delta := signalk.DeltaFormat{
					Context: signalk.CreatePath("vessels", self.ID.String()),
					Updates: []signalk.DeltaUpdate{
						{
							Timestamp: time.Now(),
							Values: []signalk.DeltaUpdateObject{
								{
									Path:  signalk.CreatePath("navigation", "speedOverGround"),
									Value: signalk.DataValueFromNumerical(rand.Float64() * 10),
								},
							},
						},
					},
				}
				jsonBytes, err := json.Marshal(&delta)
				if err != nil {
					log.Printf("WS Encode error: %v\n", err)
				}
				c.WriteMessage(websocket.TextMessage, jsonBytes)

				select {
				case <-done:
					return
				case <-time.After(1 * time.Second):
				}
			}
		}()
	})
	u.OnClose(func(c *websocket.Conn, err error) {
		close(done)
		log.Printf("Closed WS: %v, %s\n", c.RemoteAddr(), err)
	})
	return u
}

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v", err)
	}
}

func Run() error {
	self = signalk.CreateVessel(signalk.VesselIDFromUUID(uuid.New()))
	// vessel := signalk.CreateVessel(signalk.VesselIDFromMMSI("123123123"))

	// Add some data to vessel
	self.Values.Add(signalk.CreatePath("navigation.speedThroughWater"), &signalk.VesselDataEntry{
		Value:     signalk.DataValueFromNumerical(2.94),
		Timestamp: time.Now(),
		SourceRef: "testsource",
	})
	self.Values.Add(signalk.CreatePath("navigation.courseOverGround"), &signalk.VesselDataEntry{
		Value:     signalk.DataValueFromNumerical(55),
		Timestamp: time.Now(),
		SourceRef: "testsource",
	})

	// Create root format
	root = signalk.FullFormat{
		Version: "1.0",
		Self:    signalk.CreatePath("vessels", self.ID.String()),
		Vessels: signalk.VesselList{
			self,
		},
		Sources: []signalk.Source{},
	}

	// Setup HTTP router
	r := chi.NewRouter()
	r.Get("/signalk/v1/stream", func(w http.ResponseWriter, r *http.Request) {
		u := newUpgrader()
		_, err := u.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WS error: %v\n", err)
			return
		}
	})

	baseAPI := "/signalk/v1/api"
	r.Route(baseAPI, func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			pathStr, _ := strings.CutPrefix(r.URL.String(), baseAPI)
			path := signalk.CreatePath(strings.Split(pathStr, "/")...)
			w.Write([]byte(path.String()))
		})
	})

	r.Get("/delta", func(w http.ResponseWriter, r *http.Request) {
		delta := signalk.DeltaFormat{
			Context: signalk.CreatePath("vessels", self.ID.String()),
			Updates: []signalk.DeltaUpdate{
				{
					Timestamp: time.Now(),
					Values: []signalk.DeltaUpdateObject{
						{
							Path:  signalk.CreatePath("navigation", "speedOverGround"),
							Value: signalk.DataValueFromNumerical(rand.Float64() * 10),
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(&delta)
	})

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&root)
	}))

	http.ListenAndServe(":3000", cors.AllowAll().Handler(r))
	return nil
}
