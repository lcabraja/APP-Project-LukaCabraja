package map_manager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ValueType int

const (
	UnknownType ValueType = iota

	StringType
	JsonType
	IntType
	FloatType
	BoolType
	DurationType
)

var names = []string{
	"UnknownType",
	"StringType",
	"JsonType",
	"IntType",
	"FloatType",
	"BoolType",
	"DurationType",
}

func (vt ValueType) String() string {
	if vt < StringType || vt > DurationType {
		return names[0]
	}

	return names[vt]
}

func (vt ValueType) Validate(data interface{}) error {
	err := fmt.Errorf("provided value [%-16v]:16] not of type %s", data, vt)
	switch vt {
	case StringType:
		switch data.(type) {
		case string:
			return nil
		default:
			return err
		}
	case JsonType:
		switch data.(type) {
		case string:
			if json.Valid([]byte(data.(string))) {
				return nil
			}
		case []byte:
			if json.Valid(data.([]byte)) {
				return nil
			}
		default:
			return err
		}
	case IntType:
		switch data.(type) {
		case int:
			return nil
		default:
			return err
		}
	case FloatType:
		switch data.(type) {
		case float64:
			return nil
		default:
			return err
		}
	case BoolType:
		switch data.(type) {
		case bool:
			return nil
		default:
			return err
		}
	case DurationType:
		switch data.(type) {
		case time.Duration:
			return nil
		default:
			return err
		}
	}
	return err
}

func GetType(name string) ValueType {
	for i, n := range names[1:] {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(n)) {
			return ValueType(i)
		}
	}

	return UnknownType
}
