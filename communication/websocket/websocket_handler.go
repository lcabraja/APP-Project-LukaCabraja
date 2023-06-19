package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/lcabraja/APP-Project-LukaCabraja/communication/websocket/action_handler"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {

	route       string
	connections map[string]*websocket.Conn
	mu          sync.Mutex


	upgrader websocket.Upgrader

	ah            map[string][]action_handler.ActionHandler
	subscriptions map[string]*WebsocketEventSubscription
}

func NewWebsocketHandler(route string, console bool) *WebsocketHandler {
	wh := &WebsocketHandler{
		route:       route,
		connections: make(map[string]*websocket.Conn),
		mu:          sync.Mutex{},

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},

		ah:            make(map[string][]action_handler.ActionHandler),
		subscriptions: make(map[string]*WebsocketEventSubscription),
	}

	if console {
		wh.Subscribe(&WebsocketEventSubscription{
			Event: ConnectionEstablished,
			Handler: func(event WebsocketEvent, bytes []byte, handler *WebsocketHandler) {
				log.Devf("[%s] Connection established [%s]\n", handler.Route(), bytes)
			},
		})
		wh.Subscribe(&WebsocketEventSubscription{
			Event: ConnectionClosed,
			Handler: func(event WebsocketEvent, bytes []byte, handler *WebsocketHandler) {
				log.Devf("Connection closed [%s]\n", bytes)
			},
		})
	}

	return wh
}

func (wh *WebsocketHandler) fireEvent(e WebsocketEvent, msg []byte) {
	for _, sub := range wh.subscriptions {
		if sub.Event == e || sub.Event == EventFired {
			go sub.Handler(e, msg, wh)
		}
	}
}

func (wh *WebsocketHandler) RegisterActionHandler(action string, ah ...action_handler.ActionHandler) {
	wh.ah[action] = append(wh.ah[action], ah...)
}

func (wh *WebsocketHandler) Subscribe(sub *WebsocketEventSubscription) string {
	uuid := uuid.NewString()
	wh.subscriptions[uuid] = sub
	return uuid
}

func (wh *WebsocketHandler) Unsubscribe(uuid string) {
	delete(wh.subscriptions, uuid)
}

func (wh *WebsocketHandler) Route() string {
	return wh.route
}

func (wh *WebsocketHandler) BroadcastMessage(msg []byte) {
	for id := range wh.connections {
		go wh.SendMessage(id, msg)
	}
}

func (wh *WebsocketHandler) SendMessage(uuid string, msg []byte) {
	wh.mu.Lock()
	conn, ok := wh.connections[uuid]
	wh.mu.Unlock()

	if !ok {
		log.E("Connection not found")
		return
	}

	wh.mu.Lock()
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.E(err.Error())
	}
	wh.mu.Unlock()
	wh.fireEvent(MessageSent, msg)
}

func (wh *WebsocketHandler) HttpHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := wh.upgrader.Upgrade(w, r, nil)
	uuid := uuid.NewString()
	wh.fireEvent(ConnectionInitiated, []byte(uuid))
	if err != nil {
		log.E(err.Error())
		wh.fireEvent(ConnectionBroken, []byte(err.Error()))
		return
	}

	defer wh.fireEvent(ConnectionClosed, []byte(uuid))
	defer conn.Close()

	wh.handleWebsocket(conn, uuid)
}

func (wh *WebsocketHandler) addConnection(uuid string, conn *websocket.Conn) {
	wh.mu.Lock()
	wh.connections[uuid] = conn
	wh.mu.Unlock()
}

func (wh *WebsocketHandler) removeConnection(uuid string) {
	wh.mu.Lock()
	delete(wh.connections, uuid)
	wh.mu.Unlock()
}

func (wh *WebsocketHandler) handleWebsocket(conn *websocket.Conn, uuid string) {
	wh.addConnection(uuid, conn)
	wh.fireEvent(ConnectionEstablished, []byte(uuid))
	defer wh.removeConnection(uuid)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.E(err.Error())
			wh.fireEvent(ErrorReceived, []byte(err.Error()))
			wh.fireEvent(ConnectionBroken, []byte(err.Error()))
			return
		}

		wh.fireEvent(MessageReceived, msg)
		if msgType == websocket.TextMessage {
			wh.fireEvent(TextMessageReceived, msg)
			wh.tryHandleAction(msg, uuid)
		}
	}
}

func (wh *WebsocketHandler) tryHandleAction(msg []byte, uuid string) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(msg, &jsonData); err != nil {
		log.E(err.Error())
		return
	}

	var (
		action interface{}
		ok     bool
	)
	if action, ok = jsonData["action"]; !ok {
		log.E("no action found in message\n")
		return
	}

	if a, ok := action.(string); ok {
		wh.fireEvent(ActionReceived, msg)
		wh.handleAction(jsonData, a, uuid)
	}
}

func (wh *WebsocketHandler) handleAction(msg map[string]interface{}, action, uuid string) {
	if _, ok := wh.ah[action]; !ok {
		log.E("no handler for action [%s]\n", uuid)
		return
	}
	for _, handler := range wh.ah[action] {
		log.Devf("handling [%s] for connection [%s]\n", action, uuid)
		go handler.Handle(msg)
	}

}
