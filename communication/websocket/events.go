package websocket

type WebsocketEvent int

const (
	ConnectionInitiated   WebsocketEvent = iota
	ConnectionEstablished
	ConnectionBroken
	ConnectionClosed

	ErrorReceived
	MessageReceived
	TextMessageReceived
	ActionReceived
	ErrorSent
	MessageSent

	EventFired
)

func (wse WebsocketEvent) String() string {
	names := []string{
		"ConnectionInitiated",
		"ConnectionEstablished",
		"ConnectionBroken",
		"ConnectionClosed",

		"ErrorReceived",
		"MessageReceived",
		"TextMessageReceived",
		"ActionReceived",
		"ErrorSent",
		"MessageSent",

		"EventFired",
	}

	if wse < ConnectionInitiated || wse > EventFired {
		return "UnknownWebsocketEvent"
	}

	return names[wse]
}
