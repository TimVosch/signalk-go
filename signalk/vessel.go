package signalk

import (
	"encoding/json"
	"time"
)

type Vessel struct {
	ID     VesselID   `json:"-"`
	Name   string     `json:"name,omitempty"`
	UUID   VesselUUID `json:"uuid,omitempty"`
	MMSI   VesselMMSI `json:"mmsi,omitempty"`
	Values VesselData `json:".,omitempty"`
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	if v.Name != "" {
		m["name"] = v.Name
	}
	if v.UUID != "" {
		m["uuid"] = v.UUID
	}
	if v.MMSI != "" {
		m["mmsi"] = v.MMSI
	}
	v.Values.encodeToMap(m)
	return json.Marshal(m)
}

type VesselData []VesselDataEntry

func (vd VesselData) encodeToMap(m map[string]any) {
	setPath := func(path Path, v any) {
		tmp := m
		parts := path.Parts()
		for _, part := range parts[:len(parts)-2] {
			partMap, ok := tmp[part]
			if !ok {
				partMap = map[string]any{}
			}
			tmp[part] = partMap
		}
		tmp[parts[len(parts)-1]] = v
	}

	for _, entry := range vd {
		setPath(entry.Path, entry)
	}
}

func (vd VesselData) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	vd.encodeToMap(m)
	return json.Marshal(m)
}

type VesselDataEntry struct {
	Path      Path      `json:"-"`
	SourceRef string    `json:"$source,omitempty"`
	Value     DataValue `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`

	// Optional used if there are multiple sources generating this value
	Values map[string]VesselData `json:"values,omitempty"`
}

func CreateVessel(id VesselID) Vessel {
	return Vessel{
		Name:   "",
		ID:     id,
		UUID:   id.UUID,
		MMSI:   id.MMSI,
		Values: VesselData{},
	}
}
