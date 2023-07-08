package signalk

import (
	"encoding/json"
	"strings"
)

type Path struct {
	values []string
}

func CreatePath(values ...string) Path {
	return Path{values}
}

func (path Path) String() string {
	return strings.Join(path.values, ".")
}

func (path Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(path.String())
}

func (path Path) Child(value string) Path {
	if value == "" {
		return path
	}
	path.values = append(path.values, value)
	return path
}

func (path Path) Parts() []string {
	return path.values
}
