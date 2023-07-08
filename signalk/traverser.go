package signalk

import (
	"errors"
	"fmt"
)

var ErrNoSuchKey = errors.New("requested key does not exist")

type Traverser interface {
	Get(path Path) (any, error)
}

type FullFormatTraverser struct {
	root FullFormat
}

func NewTraverser(root FullFormat) FullFormatTraverser {
	return FullFormatTraverser{root}
}

func (t FullFormatTraverser) Next(path Path) (any, error) {
	fmt.Printf("FullTraverser getting: %s\n", path)
	next, rest := path.Pop()
	var traverser Traverser
	switch next {
	case "":
		return t.root, nil
	case "version":
		traverser = &literalTraverser{t.root.Version}
	case "self":
		traverser = &literalTraverser{t.root.Self}
	case "vessels":
		traverser = &vesselListTraverser{t.root, t.root.Vessels}
	case "sources":
		traverser = &sourceListTraverser{t.root.Sources}
	}
	if traverser == nil {
		return nil, fmt.Errorf("%w: missing '%s' in root", ErrNoSuchKey, next)
	}
	return traverser.Get(rest)
}

type literalTraverser struct {
	literal any
}

func (t literalTraverser) Get(path Path) (any, error) {
	return t.literal, nil
}

type vesselListTraverser struct {
	root    FullFormat
	vessels VesselList
}

func (t vesselListTraverser) Get(path Path) (any, error) {
	fmt.Printf("VesselListTraverser getting: %s\n", path)
	next, rest := path.Pop()
	switch next {
	case "":
		return t.vessels, nil
	case "self":
		if next == "self" {
			self, err := NewTraverser(t.root).Next(t.root.Self)
			if err != nil {
				return nil, err
			}
			return (&vesselTraverser{self.(Vessel)}).Next(rest)
		}
	}

	id, err := VesselIDFromString(next)
	if err != nil {
		return nil, fmt.Errorf("cannot get vessel with id '%s' from vessels: %w", next, err)
	}
	fmt.Printf("VesselListTraverser finding: %s\n", id)
	for _, vessel := range t.vessels {
		if vessel.ID == id {
			return (&vesselTraverser{vessel}).Next(rest)
		}
	}
	return nil, fmt.Errorf("%w: no vessel with id '%s'", ErrNoSuchKey, next)
}

type vesselDataEntryTraverser struct {
	data VesselDataEntry
}

func (t vesselDataEntryTraverser) Next(path Path) (any, error) {
	panic("not implemented")
}

type vesselTraverser struct {
	vessel Vessel
}

func (t vesselTraverser) Next(path Path) (any, error) {
	next, rest := path.Pop()
	switch next {
	case "":
		return t.vessel, nil
	case "id":
		return (&literalTraverser{t.vessel.ID}).Get(rest)
	case "name":
		return (&literalTraverser{t.vessel.Name}).Get(rest)
	case "uuid":
		return (&literalTraverser{t.vessel.UUID.String()}).Get(rest)
	case "mmsi":
		return (&literalTraverser{t.vessel.MMSI.String()}).Get(rest)
	}
	dataEntry, dataTree, err := t.vessel.Values.Get(path)
	if err != nil {
		return nil, err
	}
	if dataEntry != nil {
		return dataEntry, nil
	}
	return dataTree, nil
}

type sourceListTraverser struct {
	sources SourceList
}

func (t sourceListTraverser) Get(path Path) (any, error) {
	panic("not implemented") // TODO: Implement
}
