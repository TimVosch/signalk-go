package signalk

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_                       json.Marshaler = (*DataValue)(nil)
	ErrDataValueTypeInvalid                = errors.New("type of DataValue is invalid")
)

type DataValueType uint8

const (
	DataValueNumerical DataValueType = iota
	DataValueString
	DataValueObject
)

type DataValue struct {
	Type      DataValueType
	Numerical float64
	String    string
	Object    map[string]any
}

func (value DataValue) MarshalJSON() ([]byte, error) {
	switch value.Type {
	case DataValueNumerical:
		return json.Marshal(value.Numerical)
	case DataValueString:
		return json.Marshal(value.String)
	case DataValueObject:
		return json.Marshal(value.Object)
	}
	return nil, fmt.Errorf("DataValue marshal error: %w", ErrDataValueTypeInvalid)
}

func DataValueFromNumerical(v float64) DataValue {
	return DataValue{
		Type:      DataValueNumerical,
		Numerical: v,
	}
}

func DataValueFromString(v string) DataValue {
	return DataValue{
		Type:   DataValueString,
		String: v,
	}
}

func DataValueFromObject(v map[string]any) DataValue {
	return DataValue{
		Type:   DataValueObject,
		Object: v,
	}
}
