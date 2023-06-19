package data

type ResourceType int

const (
	StopwatchType ResourceType = iota
	MapManagerType
	FetcherType
	FormatterType

	UnknownType
)

func (rt ResourceType) String() string {
	switch rt {
	case StopwatchType:
		return "StopwatchType"
	case MapManagerType:
		return "MapManagerType"
	case FetcherType:
		return "FetcherType"
	case FormatterType:
		return "FormatterType"
	}
	return "Unknown"
}
