package signalk

import (
	"errors"
	"fmt"
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

func (root *Root) ApplyDelta(delta Delta) error {
	context := delta.Context
	if context == nil {
		context = CreatePath("vessels.self")
	}

	// Update values from delta
	for _, update := range delta.Updates {
		for _, value := range update.Values {
			path := value.Path.Prepend(context)
			err := root.applyDeltaUpdateValue(path, value.Value)
			if err != nil {
				return err
			}
		}
	}

	// TODO: Update meta from delta

	return nil
}

// TODO: Rename to SetPath?
func (root *Root) applyDeltaUpdateValue(path Path, value any) error {
	if path.IsEmpty() {
		return fmt.Errorf("%w: values cannot be merged into root", ErrDeltaUnsupportedUpdate)
	}

	next, rest := path.FirstOut()
	switch next.String() {
	case "vessels":
		return root.Vessels.applyDeltaUpdateValue(root, rest, value)
	}
	return fmt.Errorf("%w: for %s in root", ErrDeltaPathInvalid, path)
}
