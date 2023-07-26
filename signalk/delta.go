package signalk

import (
	"errors"
	"time"

	"signalk/tree"
)

var (
	ErrDeltaUnsupportedUpdate = errors.New("delta update is unsupported")
	ErrDeltaPathInvalid       = errors.New("delta update path is invalid")
)

type Delta struct {
	Context tree.Path
	Updates []DeltaUpdate
}

type DeltaUpdate struct {
	Source    string
	SourceRef string
	Timestamp time.Time
	Values    []DeltaUpdateValues
	Meta      []DeltaUpdateMeta
}
type DeltaUpdateValues struct {
	Path  tree.Path
	Value any
}
type DeltaUpdateMeta struct {
	Path  tree.Path
	Value any
}
