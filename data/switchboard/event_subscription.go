package switchboard

type SwitchboardEventSubscription struct {
	Event   []SwitchboardEvent
	Keys    []string
	Handler func(SwitchboardEvent, string, *Switchboard)
}
