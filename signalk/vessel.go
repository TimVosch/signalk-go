package signalk

import (
	"errors"

	"github.com/Jeffail/gabs/v2"

	"signalk/tree"
)

var _ tree.NodeBrancher = (*vesselList)(nil)

type vesselListNode = *tree.GenericNode[vesselList]

type vesselList struct {
	vessels map[VesselID]vesselNode
}

func createVesselListNode() vesselListNode {
	return tree.CreateGenericNode(vesselList{
		vessels: make(map[VesselID]vesselNode),
	})
}

func (v *vesselList) GetChild(key string) (tree.Node, error) {
	id, err := ParseVesselID(key)
	if err != nil {
		return nil, err
	}
	node, ok := v.vessels[id]
	if !ok {
		return nil, ErrPathInvalid
	}
	return node, nil
}

func (v *vesselList) AddChild(node tree.Node) error {
	panic("not implemented") // TODO: Implement
}

//
//
//

type vesselNode = *tree.GenericNode[vessel]

type vessel struct {
	// Either MMSI, UUID or URL
	ID   VesselID
	Name string

	// Data groups
}

func (v vessel) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(v.ID.String(), v.ID.typ.String())
	return obj.MarshalJSON()
}

func (n *vessel) GetChild(key string) (tree.Node, error) {
	switch key {
	}
	return nil, ErrPathInvalid
}

func (n *vessel) AddChild(node tree.Node) error {
	return errors.New("adding children directly to a vessel is currently not implemented")
}
