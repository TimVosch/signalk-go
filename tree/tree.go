package tree

type Traversable interface {
	GetByPath(path Path) (any, Path, error)
}

type Settable interface {
	SetValue(v any) error
}
