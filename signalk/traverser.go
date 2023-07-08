package signalk

import (
	"errors"
	"fmt"
)

var ErrNoSuchKey = errors.New("request key does not exist")

type Traverser interface {
	Next(key string) (Traverser, error)
	Value() any
}

type FullFormatTraverser struct {
	root FullFormat
}

func (t *FullFormatTraverser) Next(key string) (Traverser, error) {
	switch key {
	case "version":
		return &literalTraverser{t.root.Version}, nil
	case "self":
		return &literalTraverser{t.root.Self}, nil
	case "vessels":
		return &vesselListTraverser{t.root.Vessels}, nil
	case "sources":
		return &sourceListTraverser{t.root.Sources}, nil
	}
	return nil, ErrNoSuchKey
}

func (t *FullFormatTraverser) Value() any {
	return t.root
}

type literalTraverser struct {
	literal any
}

func (t *literalTraverser) Next(key string) (Traverser, error) {
	return nil, ErrNoSuchKey
}

func (t *literalTraverser) Value() any {
	return t.literal
}

type vesselListTraverser struct {
	vessels VesselList
}

func (t *vesselListTraverser) Next(key string) (Traverser, error) {
	id, err := VesselIDFromString(key)
	if err != nil {
		return nil, fmt.Errorf("cannot get vessel with id '%s' from vessels: %w", key, err)
	}
	for _, vessel := range t.vessels {
		if vessel.ID == id {
			return &vesselTraverser{vessel}, nil
		}
	}
	return nil, ErrNoSuchKey
}

func (t *vesselListTraverser) Value() any {
	return t.vessels
}

type vesselDataEntryTraverser struct {
	data VesselDataEntry
}

func (t *vesselDataEntryTraverser) Next(key string) (Traverser, error) {
	panic("not implemented")
}

func (t *vesselDataEntryTraverser) Value() any {
	return t.data
}

type vesselTraverser struct {
	vessel Vessel
}

func (t *vesselTraverser) Next(key string) (Traverser, error) {
	switch key {
	case "self":
		// root.self
	case "id":
		return &literalTraverser{t.vessel.ID}, nil
	case "name":
		return &literalTraverser{t.vessel.Name}, nil
	case "uuid":
		return &literalTraverser{t.vessel.UUID.String()}, nil
	case "mmsi":
		return &literalTraverser{t.vessel.MMSI.String()}, nil
	}
	dataEntry, dataTree, err := t.vessel.Values.Get(CreatePath(key))
	if err != nil {
		return nil, err
	}
	if dataEntry != nil {
		return &literalTraverser{dataEntry}, nil
	}
	return &literalTraverser{dataTree}, nil
}

func (t *vesselTraverser) Value() any {
	return t.vessel
}

type sourceListTraverser struct {
	sources SourceList
}

func (t *sourceListTraverser) Next(key string) (Traverser, error) {
	panic("not implemented") // TODO: Implement
}

func (t *sourceListTraverser) Value() any {
	return t.sources
}
