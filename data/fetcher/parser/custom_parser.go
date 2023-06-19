package parser

import (
	"github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"
)

type CustomParser struct {
	handler func([]byte, *map_manager.MapManager)
	mm      *map_manager.MapManager
	vt      map_manager.ValueType
	key     string
}

func NewCustomParser(handler func([]byte, *map_manager.MapManager), mm *map_manager.MapManager) *CustomParser {
	return &CustomParser{
		handler: handler,
		mm:      mm,
	}
}

func (cp *CustomParser) Parse(body []byte) error {
	cp.handler(body, cp.mm)
	return nil
}
