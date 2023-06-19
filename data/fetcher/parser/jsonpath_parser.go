package parser

import (
	"encoding/json"
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"
	"github.com/oliveagle/jsonpath"
)

type JsonPathParser struct {
	jsonPath string
	mm       *map_manager.MapManager
	vt       map_manager.ValueType
	key      string
}

func NewJsonPathParser(jsonPath string, mm *map_manager.MapManager, vt map_manager.ValueType, key string) *JsonPathParser {
	return &JsonPathParser{jsonPath: jsonPath,
		mm:  mm,
		vt:  vt,
		key: key}
}

func (jpp *JsonPathParser) Parse(body []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	value, err := jsonpath.JsonPathLookup(data, jpp.jsonPath)
	if err != nil {
		return err
	}

	if err := jpp.vt.Validate(value); err != nil {
		if jpp.vt == map_manager.StringType {
			jpp.mm.SetString(jpp.key, fmt.Sprintf("%v", value))
		}
		return err
	}

	if err := jpp.mm.Set(jpp.vt, jpp.key, value); err != nil {
		return err
	}

	return nil
}
