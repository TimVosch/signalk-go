package signalk

import (
	"errors"
	"time"
)

var (
	ErrDeltaUnsupportedUpdate = errors.New("delta update is unsupported")
	ErrDeltaPathInvalid       = errors.New("delta update path is invalid")
)

type Delta struct {
	Context Path
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
	Path  Path
	Value any
}
type DeltaUpdateMeta struct {
	Path  Path
	Value any
}
