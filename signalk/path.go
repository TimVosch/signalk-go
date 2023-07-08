package signalk

import (
	"encoding/json"
	"strings"
)

type Path []string

func CreatePath(pathStrings ...string) Path {
	path := make(Path, 0)
	for _, pathString := range pathStrings {
		pathString = strings.Trim(pathString, " \t")
		for _, pathPart := range strings.Split(pathString, ".") {
			if pathPart == "" {
				continue
			}
			path = append(path, strings.Trim(pathPart, " \t"))
		}
	}
	return path
}

func (path Path) String() string {
	return strings.Join(path, ".")
}

func (path Path) IsEmpty() bool {
	return len(path) == 0
}

func (path Path) Append(other Path) Path {
	return append(path, other...)
}

func (path Path) Prepend(other Path) Path {
	return append(other, path...)
}

func (path Path) FirstOut() (Path, Path) {
	if path.IsEmpty() {
		return Path{}, Path{}
	}
	if len(path) < 2 {
		return path[:1], Path{}
	}
	return path[:1], path[1:]
}

func (path Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(path.String())
}
