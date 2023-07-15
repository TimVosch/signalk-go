package signalk

import (
	"errors"

	"github.com/Jeffail/gabs/v2"

	"signalk/tree"
)

var _ tree.Node = (*RootNode)(nil)

type RootNode struct {
	self    tree.Path
	version string
	vessels tree.Node
	sources tree.Node
	atons   tree.Node
	shore   tree.Node
}

func createRootNode() *RootNode {
	return &RootNode{
		version: "1.0.0",
	}
}

func (root RootNode) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	obj.Set(root.version, "version")
	return obj.MarshalJSON()
}

func (r *RootNode) Parent() tree.Node {
	return nil
}

func (r *RootNode) SetParent(_ tree.Node) {
	// Do nothing
}

func (r *RootNode) Key() string {
	return ""
}

func (r *RootNode) Path() tree.Path {
	return tree.CreatePath()
}

func (r *RootNode) GetChild(key string) (tree.Node, error) {
	switch key {
	case "vessels":
		return r.vessels, nil
	case "sources":
		return r.sources, nil
	case "version":
		// TODO: Should ?
		return tree.CreateLeaf(r.version), nil
	case "atons":
		return r.atons, nil
	}
	return nil, errors.New("error!")
}

func (r *RootNode) SetValue(v any) error {
	return errors.New("root is not leaf")
}

func (r *RootNode) GetValue() any {
	return errors.New("root is not leaf")
}

func (r *RootNode) AddChild(node tree.Node) error {
	return errors.New("no children allowed")
}
