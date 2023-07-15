package tree

import (
	"encoding/json"
	"errors"
	"fmt"
)

var errChildNotFound = errors.New("child not found")

/*

Tree requirements:
 - A field must (optionally) validated by a FieldValueValidator
 - Must be able to get/set a field by path (create if not exist)
 - Certain fields, preconfigured by path, must have predefined FieldValueValidators
*/

type Node interface {
	Parent() Node
	SetParent(Node)
	Key() string
	Path() Path
	GetChild(key string) (Node, error)
	SetValue(v any) error
	GetValue() any
	AddChild(node Node) error
}

var (
	_ Node = (*Branch)(nil)
	_ Node = (*Leaf[any])(nil)
)

type (
	NodeGroup map[string]Node
	// Branch represents a group of nodes at the root or as part of the tree
	Branch struct {
		parent Node
		key    string
		value  NodeGroup
	}
	// Leaf represents and end node of the tree, this node has a value
	Leaf[T any] struct {
		parent Node
		key    string
		value  T
	}
)

func CreateLeaf[T any](v T) *Leaf[T] {
	return &Leaf[T]{
		value: v,
	}
}

func (branch *Branch) Parent() Node {
	return branch.parent
}

func (leaf *Leaf[T]) Parent() Node {
	return leaf.parent
}

func (branch *Branch) SetParent(node Node) {
	branch.parent = node
}

func (leaf *Leaf[T]) SetParent(node Node) {
	leaf.parent = node
}

func (branch *Branch) Key() string {
	return branch.key
}

func (leaf *Leaf[T]) Key() string {
	return leaf.key
}

func (branch *Branch) Path() Path {
	if branch.Parent() == nil {
		return CreatePath(branch.Key())
	}
	return branch.Parent().Path().Add(branch.Key())
}

func (leaf *Leaf[T]) Path() Path {
	if leaf.Parent() == nil {
		return CreatePath(leaf.Key())
	}
	return leaf.Parent().Path().Add(leaf.Key())
}

func (branch *Branch) GetChild(key string) (Node, error) {
	child, ok := branch.value[key]
	if !ok {
		return nil, errChildNotFound
	}
	return child, nil
}

func (leaf *Leaf[T]) GetChild(key string) (Node, error) {
	return nil, errors.New("leaf node has no children")
}

func (branch *Branch) AddChild(node Node) error {
	branch.value[node.Key()] = node
	node.SetParent(branch)
	return nil
}

func (leaf *Leaf[T]) AddChild(node Node) error {
	return errors.New("leaf node has no children")
}

func (branch *Branch) SetValue(v any) error {
	return errors.New("setting values on branch node is not implemented")
}

func (leaf *Leaf[T]) SetValue(v any) error {
	vT, ok := v.(T)
	if !ok {
		return errors.New("value is invalid type")
	}
	leaf.value = vT
	return nil
}

func (branch *Branch) GetValue() any {
	return branch.value
}

func (leaf *Leaf[T]) GetValue() any {
	return leaf.value
}

func (branch Branch) MarshalJSON() ([]byte, error) {
	return json.Marshal(branch.value)
}

func (leaf Leaf[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(leaf.value)
}

type Tree struct {
	root Node
}

func (tree Tree) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.root.GetValue())
}

func Create() *Tree {
	return CreateWith(&Branch{key: "", value: make(NodeGroup)})
}

func CreateWith(node Node) *Tree {
	return &Tree{
		root: node,
	}
}

func (tree *Tree) Get(path Path) (Node, error) {
	return tree.get(path, false)
}

func (tree *Tree) get(path Path, create bool) (Node, error) {
	var err error
	var node Node

	next, rest := path.FirstOut()
	node, err = tree.root.GetChild(next.String())
	if errors.Is(err, errChildNotFound) {
		fmt.Printf("Root is missing key: %v\n", next)
		if create {
			return tree.buildAndGet(tree.root, next.Append(rest))
		}
		return nil, errChildNotFound
	}
	if err != nil {
		fmt.Printf("Root error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Root found\n")
	for !rest.IsEmpty() {
		next, rest = rest.FirstOut()
		nextNode, err := node.GetChild(next.String())
		if errors.Is(err, errChildNotFound) {
			fmt.Printf("Missing %s on %s\n", next, node.Key())
			if create {
				return tree.buildAndGet(node, next.Append(rest))
			}
			return nil, err
		}
		fmt.Printf("found %s on %s\n", next.String(), node.Key())
		node = nextNode
		// If there is more path, then keep traversing
		if !rest.IsEmpty() {
			fmt.Printf("More to come\n")
			continue
		}
		// Otherwise we should be at a leaf and set its value
		return node, nil
	}
	return node, nil
}

func (tree *Tree) buildAndGet(node Node, path Path) (Node, error) {
	fmt.Printf("buildAndGet for %s, for path %s\n", node.Key(), path)
	var newNode Node
	var err error

	next, rest := path.FirstOut()
	for {
		// Create a new node, leaf if last path of the path, otherwise branch
		if rest.IsEmpty() {
			fmt.Printf("Building leaf node: %s\n", next)
			newNode = &Leaf[any]{
				key:    next.String(),
				parent: node,
				value:  nil, // TODO: What?
			}
		} else {
			fmt.Printf("Building branch node: %s\n", next)
			newNode = &Branch{
				key:    next.String(),
				parent: node,
				value:  NodeGroup{},
			}
		}

		// Append new node to branch
		fmt.Printf("Adding %s to %s\n", newNode.Key(), node.Key())
		err = node.AddChild(newNode)
		if err != nil {
			return nil, err
		}
		// If this is the last key (thus a leaf) stop traversing
		if rest.IsEmpty() {
			fmt.Printf("was last node, quitting: %s\n", next)
			return newNode, nil
		}
		// otherwise make our newNode our main node and traverse a step further
		node = newNode
		next, rest = rest.FirstOut()
	}
}

func (tree *Tree) Set(path Path, value any) error {
	node, err := tree.get(path, true)
	if err != nil {
		return err
	}
	return node.SetValue(value)
}
