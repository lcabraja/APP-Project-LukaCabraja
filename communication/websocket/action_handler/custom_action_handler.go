package action_handler

type CustomActionHandler struct {
	handler func(map[string]interface{})
}

func NewCustomActionHandler(handler func(map[string]interface{})) *CustomActionHandler {
	return &CustomActionHandler{
		handler: handler,
	}
}

func (cah *CustomActionHandler) Handle(jsonData map[string]interface{}) {
	cah.handler(jsonData)
}
