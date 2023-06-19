package websocket

type WebsocketEventSubscription struct {
	Event   WebsocketEvent
	Handler func(WebsocketEvent, []byte, *WebsocketHandler)
}
