package map_manager

type MapManagerEventSubscription struct {
	Event   MapManagerEvent
	Handler func(MapManagerEvent, ValueType, string, *MapManager)
}
