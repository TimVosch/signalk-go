package signalk

import (
	"github.com/Jeffail/gabs/v2"
)

type VesselList map[VesselID]*Vessel

func (l VesselList) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	for _, v := range l {
		obj.Set(v, v.ID.String())
	}
	return obj.MarshalJSON()
}

type Vessel struct {
	// Either MMSI, UUID or URL
	ID   VesselID
	Name string

	// Data groups
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(v.ID.String(), v.ID.typ.String())
	return obj.MarshalJSON()
}
