package signalk

import "time"

// DeltaFormat model
type DeltaFormat struct {
	Context Path          `json:"context"`
	Updates []DeltaUpdate `json:"updates"`
}

type DeltaUpdate struct {
	Source    Source              `json:"source"`
	Timestamp time.Time           `json:"timestamp"`
	Values    []DeltaUpdateObject `json:"values"`
}

type DeltaUpdateObject struct {
	Path  Path      `json:"path"`
	Value DataValue `json:"value"`
}
