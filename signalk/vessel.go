package signalk

import "github.com/Jeffail/gabs/v2"

type Vessel struct {
	ID     VesselID       `json:"-"`
	Name   string         `json:"name,omitempty"`
	UUID   VesselUUID     `json:"uuid,omitempty"`
	MMSI   VesselMMSI     `json:"mmsi,omitempty"`
	Values VesselDataTree `json:"-"`
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	if v.Name != "" {
		obj.Set(v.Name, "name")
	}
	if v.UUID != "" {
		obj.Set(v.UUID, "uuid")
	}
	if v.MMSI != "" {
		obj.Set(v.MMSI, "mmsi")
	}
	for path, dataEntry := range v.Values.Flatten() {
		obj.SetP(dataEntry, path)
	}
	return obj.MarshalJSON()
}

func CreateVessel(id VesselID) Vessel {
	return Vessel{
		Name:   "",
		ID:     id,
		UUID:   id.UUID,
		MMSI:   id.MMSI,
		Values: NewVesselDataTree(),
	}
}
