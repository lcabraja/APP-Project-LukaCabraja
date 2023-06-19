package map_manager

type MapManagerEvent int

const (
	RecordUpdated   MapManagerEvent = iota
	StringUpdated
	JsonUpdated
	IntUpdated
	FloatUpdated
	BoolUpdated
	DurationUpdated

	RecordChanged
	StringChanged
	JsonChanged
	IntChanged
	FloatChanged
	BoolChanged
	DurationChanged

	RecordRead
	StringRead
	JsonRead
	IntRead
	FloatRead
	BoolRead
	DurationRead

	EventFired
)

func (mme MapManagerEvent) String() string {
	names := []string{
		"RecordUpdated",
		"StringUpdated",
		"JsonUpdated",
		"IntUpdated",
		"FloatUpdated",
		"BoolUpdated",
		"DurationUpdated",

		"RecordChanged",
		"StringChanged",
		"JsonChanged",
		"IntChanged",
		"FloatChanged",
		"BoolChanged",
		"DurationChanged",

		"RecordRead",
		"JsonRead",
		"StringRead",
		"IntRead",
		"FloatRead",
		"BoolRead",
		"DurationRead",

		"EventFired",
	}

	if mme < RecordUpdated || mme > EventFired {
		return "UnknownMapManagerEvent"
	}

	return names[mme]
}
