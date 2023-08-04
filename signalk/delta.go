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
	Context tree.Path     `json:"context"`
	Updates []DeltaUpdate `json:"updates"`
}

type DeltaUpdate struct {
	Source    string       `json:"source"`
	SourceRef string       `json:"source_ref"`
	Timestamp time.Time    `json:"timestamp"`
	Values    []DeltaValue `json:"values"`
	Meta      []DeltaMeta  `json:"meta"`
}
type DeltaValue struct {
	Path  tree.Path `json:"path"`
	Value any       `json:"value"`
}
type DeltaMeta struct {
	Path  tree.Path `json:"path"`
	Value any       `json:"value"`
}
