package signalk

import (
	"errors"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/google/uuid"

	"signalk/tree"
)

type VesselProperty struct {
	Source    string                         `json:"source"`
	Value     any                            `json:"value"`
	Timestamp time.Time                      `json:"time"`
	Values    map[string]VesselPropertyValue `json:"values,omitempty"`
}

type VesselPropertyValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     any       `json:"value"`
}

type VesselProperties map[string]VesselProperty

type Vessel struct {
	ID         VesselID
	Name       string
	Properties VesselProperties
}

func NewVessel() *Vessel {
	return &Vessel{
		ID:         CreateVesselUUID(uuid.New()),
		Properties: make(VesselProperties),
	}
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(v.ID.String(), v.ID.typ.String())
	obj.Set(v.Name, "name")
	for path, property := range v.Properties {
		obj.Set(property, path)
	}
	return obj.MarshalJSON()
}

func (v *Vessel) GetByPath(path tree.Path) (any, tree.Path, error) {
	if path.IsEmpty() {
		return v, nil, nil
	}

	next, rest := path.FirstOut()
	if next.String() == "name" {
		return v.Name, rest, nil
	}

	potentialProperty := path.String()
	property, ok := v.Properties[potentialProperty]
	if ok {
		return property, nil, nil
	}

	return nil, nil, ErrDeltaPathInvalid
}

func (v *Vessel) SetProperty(path tree.Path, timestamp time.Time, source string, value any) error {
	if v.doesPropertyConflict(path) {
		return errors.New("property conflict")
	}
	property, propertyAlreadyExists := v.Properties[path.String()]
	if !propertyAlreadyExists {
		property = VesselProperty{
			Timestamp: timestamp,
			Source:    source,
			Value:     value,
			Values:    make(map[string]VesselPropertyValue),
		}
	}
	property.Values[source] = VesselPropertyValue{
		Timestamp: timestamp,
		Value:     value,
	}
	v.Properties[path.String()] = property
	return nil
}

func (v *Vessel) doesPropertyConflict(propertyPath tree.Path) bool {
	var fullPath tree.Path
	for _, part := range propertyPath {
		fullPath.Add(part)
		if _, conflict := v.Properties[fullPath.String()]; conflict {
			return true
		}
	}
	return false
}
