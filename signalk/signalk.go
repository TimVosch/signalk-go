package signalk

import "signalk/tree"

type Service struct {
	root *tree.Tree
}

func NewService() *Service {
	root := tree.CreateWith(createRootNode())
	return &Service{
		root: root,
	}
}

func (s *Service) GetPath(path tree.Path) (tree.Node, error) {
	return s.root.Get(path)
}
