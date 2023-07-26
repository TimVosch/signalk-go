package tree

import (
	"encoding/json"
	"errors"
)

type NodeBrancher interface {
	GetChild(key string) (Node, error)
	AddChild(node Node) error
}

type NodeLeafer interface {
	SetValue(v any) error
	GetValue() any
}

var _ Node = (*GenericNode[any])(nil)

type GenericNode[T any] struct {
	key    string
	parent Node
	impl   T
}

func CreateGenericNode[T any](v T) *GenericNode[T] {
	return &GenericNode[T]{
		impl: v,
	}
}

func (g *GenericNode[T]) Parent() Node {
	return g.parent
}

func (g *GenericNode[T]) SetParent(parent Node) {
	g.parent = parent
}

func (g *GenericNode[T]) Key() string {
	return g.key
}

func (g *GenericNode[T]) Path() Path {
	var path Path
	if g.parent != nil {
		path = g.parent.Path()
	}
	return path.Add(g.Key())
}

func (g *GenericNode[T]) GetChild(key string) (Node, error) {
	brancher, ok := any(g.impl).(NodeBrancher)
	if ok {
		return brancher.GetChild(key)
	}
	return nil, errors.New("Generic node does not implement Brancher")
}

func (g *GenericNode[T]) SetValue(v any) error {
	leafer, ok := any(g.impl).(NodeLeafer)
	if ok {
		return leafer.SetValue(v)
	}
	return errors.New("Generic node does not implement Leafer")
}

func (g *GenericNode[T]) GetValue() any {
	leafer, ok := any(g.impl).(NodeLeafer)
	if ok {
		return leafer.GetValue()
	}
	return errors.New("Generic node does not implement Leafer")
}

func (g *GenericNode[T]) AddChild(node Node) error {
	brancher, ok := any(g.impl).(NodeBrancher)
	if ok {
		return brancher.AddChild(node)
	}
	return errors.New("Generic node does not implement Brancher")
}

func (g GenericNode[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(any(g.impl))
}
