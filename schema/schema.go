package schema

import (
	"time"

	"github.com/google/uuid"
)

type Root struct {
	Version string
	Self    string
	Vessels Vessels
	Sources Sources
}

// ====
// Delta(s)
// ====

type Delta struct {
	Context string
	Updates []DeltaUpdate
}

type DeltaUpdate struct {
	Source    Source
	Timestamp time.Time
	Values    []DeltaValue
}

type DeltaValue struct {
	Path  string
	Value any
}

// ====
// Vessel(s)
// ====

type VesselID string

type Vessels map[VesselID]Vessel

func (vessels *Vessels) Add(vessel Vessel)  {
    vessels[vesse
}

type Vessel struct {
	Name       string
	UUID       uuid.UUID
	MMSI       string
	Properties map[string]any
}

func NewVessel() Vessel {
	return Vessel{
		UUID:       uuid.New(),
		Properties: make(map[string]any),
	}
}

type VesselProperty struct {
	Source    string
	Timestamp time.Time
	Value     any
	Values    map[string]VesselShortProperty
}

type VesselShortProperty struct {
	Timestamp time.Time
	Value     any
}

// ====
// Source(s)
// ====

type Sources map[string]Source

type Source struct {
	Name string
}
