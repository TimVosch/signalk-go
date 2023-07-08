package signalk

import (
	"encoding/json"
)

var (
	_ json.Marshaler = (*VesselList)(nil)
	_ json.Marshaler = (*VesselList)(nil)
)

type VesselList []Vessel

func (vm VesselList) MarshalJSON() ([]byte, error) {
	jsonMap := make(map[string]Vessel)
	for _, v := range vm {
		jsonMap[v.ID.String()] = v
	}
	return json.Marshal(jsonMap)
}

type SourceList []Source

func (sl SourceList) MarshalJSON() ([]byte, error) {
	jsonMap := make(map[string]Source)
	for _, v := range sl {
		jsonMap[v.Name] = v
	}
	return json.Marshal(jsonMap)
}

// FullFormat model
type FullFormat struct {
	Version string     `json:"version"`
	Self    Path       `json:"self"`
	Vessels VesselList `json:"vessels"`
	Sources SourceList `json:"sources"`
}
