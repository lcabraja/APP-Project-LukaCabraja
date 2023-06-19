package action_handler

type ActionHandler interface {
	Handle(map[string]interface{})
}
