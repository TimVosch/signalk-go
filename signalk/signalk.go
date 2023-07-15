package signalk

import "signalk/tree"

type Service struct {
	root *tree.Tree
}

func NewService() *Service {
	root := tree.CreateWith(rootNode)
	return &Service{}
}
