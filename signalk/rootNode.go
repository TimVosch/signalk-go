package signalk

import (
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

func (root RootNode) MarshalJSON() ([]byte, error) {
	obj := gabs.New()
	return obj.MarshalJSON()
}

func (r *RootNode) Parent() tree.Node {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) SetParent(_ tree.Node) {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) Key() string {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) Path() tree.Path {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) GetChild(key string) (tree.Node, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) SetValue(v any) error {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) GetValue() any {
	panic("not implemented") // TODO: Implement
}

func (r *RootNode) AddChild(node tree.Node) error {
	panic("not implemented") // TODO: Implement
}
