package signalk

import (
	"errors"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/google/uuid"

	"signalk/tree"
)

type Root struct {
	Version string    `json:"version"`
	Self    tree.Path `json:"self"`
	Vessels *Vessels  `json:"vessels"`
}

func (r *Root) CreateSelf() {
	vessel := &Vessel{
		ID: CreateVesselUUID(uuid.New()),
	}
	r.Vessels.Add(vessel)
	r.Self = tree.CreatePath(vessel.ID.String())
}

func (r *Root) GetByPath(path tree.Path) (any, tree.Path, error) {
	if path.IsEmpty() {
		return r, nil, nil
	}

	next, rest := path.FirstOut()
	switch next.String() {
	case "version":
		return r.Version, rest, nil
	case "self":
		return r.Self, rest, nil
	case "vessels":
		return r.Vessels, rest, nil
	default:
		return nil, rest, errors.New("path not found")
	}
}

type Vessels struct {
	root    *Root
	vessels map[VesselID]*Vessel
}

// Required because map[VesselID] is not marshable by default
func (v Vessels) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	for _, vessel := range v.vessels {
		obj.Set(vessel, vessel.ID.String())
	}
	return obj.MarshalJSON()
}

func (v Vessels) GetByPath(path tree.Path) (any, tree.Path, error) {
	if path.IsEmpty() {
		return v, nil, nil
	}

	next, rest := path.FirstOut()

	// Special case
	if next.String() == "self" {
		return v.GetByPath(v.root.Self.Append(rest))
	}

	vID, err := ParseVesselID(next.String())
	if err != nil {
		return nil, nil, err
	}

	vessel, exists := v.vessels[vID]
	if !exists {
		return nil, nil, errors.New("vessel not found")
	}

	return vessel, rest, nil
}

func (v Vessels) Add(vessel *Vessel) error {
	v.vessels[vessel.ID] = vessel
	return nil
}

type Vessel struct {
	ID         VesselID
	Name       string
	Navigation *VesselNavigation
}

func (v Vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(v.ID.String(), v.ID.typ.String())
	obj.Set(v.Name, "name")
	if v.Navigation != nil {
		obj.Set(v.Navigation, "navigation")
	}
	return obj.MarshalJSON()
}

func (v *Vessel) GetByPath(path tree.Path) (any, tree.Path, error) {
	if path.IsEmpty() {
		return v, nil, nil
	}

	next, rest := path.FirstOut()
	switch next.String() {
	case "name":
		return v.Name, rest, nil
	case "navigation":
		if v.Navigation == nil {
			return nil, nil, errors.New("not found")
		}
		return v.Navigation, rest, nil
	default:
		return nil, nil, errors.New("not found")
	}
}

type VesselNavigation struct {
	Position [3]float64
}

type Service struct {
	root *Root
}

func NewService() *Service {
	root := &Root{
		Version: "1.0.0",
	}

	// Instantiate vessels node
	root.Vessels = &Vessels{
		root:    root,
		vessels: make(map[VesselID]*Vessel),
	}

	// Create self vessel
	root.CreateSelf()

	// Add debugging vessel
	root.Vessels.Add(&Vessel{
		ID:   CreateVesselUUID(uuid.New()),
		Name: "Test",
	})

	return &Service{root}
}

func (s *Service) GetPath(path tree.Path) (any, error) {
	value, restPath, err := s.root.GetByPath(path)
	if err != nil {
		return nil, err
	}

	// Check if we can traverse further
	if !restPath.IsEmpty() {
		for {
			node, isTraversable := value.(tree.Traversable)
			if !isTraversable {
				break
			}
			value, restPath, err = node.GetByPath(restPath)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Got %v of type %T with pointer %p == nil ?: %v\n", value, value, value, value == nil)
			if restPath.IsEmpty() {
				break
			}
		}
	}
	return value, nil
}
