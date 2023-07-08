package signalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type VesselDataEntry struct {
	SourceRef string    `json:"$source,omitempty"`
	Value     DataValue `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`

	// Optional used if there are multiple sources generating this value
	Values map[string]VesselDataEntry `json:"values,omitempty"`
}

type VesselDataTree struct {
	children map[string]VesselDataTree
	parent   *VesselDataTree
	leaf     *VesselDataEntry
}

func NewVesselDataTree() VesselDataTree {
	return VesselDataTree{}
}

func (t VesselDataTree) MarshalJSON() ([]byte, error) {
	if t.leaf != nil {
		return json.Marshal(t.leaf)
	}
	return json.Marshal(t.children)
}

func (t *VesselDataTree) Add(path Path, entry *VesselDataEntry) error {
	// We are always adding a leaf, if the path end has children then error
	next, rest := path.Pop()

	// if next is "" then we're at the paths end
	if next == "" {
		if t.children != nil {
			return errors.New("cannot add leaf mid tree")
		}
		t.leaf = entry
		return nil
	}

	if t.children == nil {
		t.children = make(map[string]VesselDataTree)
	}
	child, ok := t.children[next]
	if !ok {
		child = NewVesselDataTree()
		child.parent = t
	}
	err := child.Add(rest, entry)
	t.children[next] = child
	return err
}

func (t *VesselDataTree) Get(path Path) (*VesselDataEntry, *VesselDataTree, error) {
	fmt.Printf("VesselDataTree getting: %s\n", path)
	next, rest := path.Pop()

	// If there is no next key then we're at the end of the path
	// if there is are children here, return that part of the tree
	// otherwise return the leaf
	if next == "" {
		if t.children != nil {
			return nil, t, nil
		}
		return t.leaf, nil, nil
	}

	child, ok := t.children[next]
	if !ok {
		return nil, nil, fmt.Errorf("%w: '%s' not found in vessel values", ErrNoSuchKey, next)
	}
	return child.Get(rest)
}

func (t *VesselDataTree) Flatten() map[string]*VesselDataEntry {
	// Make sure we're not a leaf
	if t.children == nil {
		return nil
	}
	m := make(map[string]*VesselDataEntry)
	for key := range t.children {
		child := t.children[key]
		child.flatten(CreatePath(key), m)
	}
	return m
}

func (t *VesselDataTree) flatten(path Path, m map[string]*VesselDataEntry) {
	if t.leaf != nil {
		m[path.String()] = t.leaf
		return
	}

	for key := range t.children {
		child := t.children[key]
		child.flatten(path.Child(key), m)
	}
}
