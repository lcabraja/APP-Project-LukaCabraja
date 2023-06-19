package switchboard

import "github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"

type Dependency struct {
	mm  *map_manager.MapManager
	vt  map_manager.ValueType
	key string
}

func (d Dependency) isValid() bool {
	return d.mm != nil && d.vt != map_manager.UnknownType && d.key != ""
}

func NewDependency(mm *map_manager.MapManager, vt map_manager.ValueType, key string) *Dependency {
	return &Dependency{
		mm:  mm,
		vt:  vt,
		key: key,
	}
}
