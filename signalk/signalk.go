package signalk

import (
	"errors"
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

type Root struct {
	Version string     `json:"version"`
	Self    Path       `json:"self"`
	Vessels VesselList `json:"vessels"`
}

type VesselList map[VesselID]*Vessel

func (l VesselList) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	for _, v := range l {
		obj.Set(v, v.ID.String())
	}
	return obj.MarshalJSON()
}

func (l VesselList) applyDeltaUpdateValue(root *Root, path Path, value any) error {
	if path.IsEmpty() {
		return fmt.Errorf("%w: values cannot be merged into vessel list", ErrDeltaUnsupportedUpdate)
	}

	// Find the vessel to which this delta is pointed
	vesselIDPart, rest := path.FirstOut()
	vesselIDStr := vesselIDPart.String()

	// Case: if vesselIDPart is "self" then find out local boat id
	if vesselIDStr == "self" {
		vesselIDStr = root.Self[len(root.Self)-1]
	}

	// Parse path part to vessel ID
	vesselID, err := ParseVesselID(vesselIDStr)
	if err != nil {
		return fmt.Errorf("error parsing vessel id: %w", err)
	}
	vessel, found := l[vesselID]
	if !found {
		vessel = CreateVessel(vesselID)
		l[vesselID] = vessel
	}

	return vessel.applyDeltaUpdateValue(root, rest, value)
}

type Vessel struct {
	// Either MMSI, UUID or URL
	ID   VesselID
	Name string

	// Data groups
	Navigation *NavigationData
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(v.ID.String(), v.ID.typ.String())
	if v.Navigation != nil {
		obj.Set(v.Navigation, "navigation")
	}
	return obj.MarshalJSON()
}

func (v *Vessel) applyDeltaUpdateValue(root *Root, path Path, value any) error {
	if path.IsEmpty() {
		// TODO: Merge
		return errors.New("merging in vessel is not implemented")
	}

	next, rest := path.FirstOut()
	switch next.String() {
	case "navigation":
		if v.Navigation == nil {
			v.Navigation = &NavigationData{}
		}
		return v.Navigation.applyDeltaUpdateValue(root, rest, value)
	}

	return fmt.Errorf("%w: for %s in vessel", ErrDeltaPathInvalid, path)
}

func CreateVessel(id VesselID) *Vessel {
	return &Vessel{
		ID: id,
	}
}
