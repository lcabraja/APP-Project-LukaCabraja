package switchboard

type SwitchboardEvent int

const (
	UnknownEvent SwitchboardEvent = iota
	KeyChanged
	KeyAdded
	KeyRemoved
)

func (se SwitchboardEvent) String() string {
	switch se {
	case KeyChanged:
		return "KeyChanged"
	case KeyAdded:
		return "KeyAdded"
	case KeyRemoved:
		return "KeyRemoved"
	default:
		return "UnknownEvent"
	}
}
