package signalk

import (
	"encoding/json"
	"strings"
)

type Path struct {
	values []string
}

func (path Path) Empty() bool {
	return len(path.values) == 0
}

func CreatePath(values ...string) Path {
	path := Path{
		values: make([]string, 0, len(values)),
	}
	for _, part := range values {
		path = path.Child(part)
	}
	return path
}

func (path Path) String() string {
	return strings.Join(path.values, ".")
}

func (path Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(path.String())
}

func (path Path) Pop() (string, Path) {
	if len(path.values) == 0 {
		return "", Path{}
	}
	next := path.values[0]
	if len(path.values) == 1 {
		return next, Path{}
	}
	path.values = path.values[1:]
	return next, path
}

func (path Path) Child(value string) Path {
	for _, part := range strings.Split(value, ".") {
		path = path.addSingleChild(part)
	}
	return path
}

func (path Path) addSingleChild(value string) Path {
	trimmed := strings.Trim(value, " \t")
	if trimmed == "" {
		return path
	}
	path.values = append(path.values, trimmed)
	return path
}

func (path Path) Parts() []string {
	return path.values
}
