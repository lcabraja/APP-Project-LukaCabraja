package fetcher

type FetchEvent int

const (
	UnknownEvent FetchEvent = iota
	DataReceived
	DataChanged
	ErrorReceived
	ManualParse
)

func (fe FetchEvent) String() string {
	switch fe {
	case DataReceived:
		return "DataReceived"
	case DataChanged:
		return "DataChanged"
	case ErrorReceived:
		return "ErrorReceived"
	case ManualParse:
		return "ManualParse"
	default:
		return "UnknownEvent"
	}
}
