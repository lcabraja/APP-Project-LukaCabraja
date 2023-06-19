package switchboard

type CustomValueCalculator struct {
	handler func(map[string]interface{}) (string, error)
}

func (cvc *CustomValueCalculator) Result(dependencies map[string]interface{}) (string, error) {
	return cvc.handler(dependencies)
}

func NewCustomValueCalculator(handler func(map[string]interface{}) (string, error)) ValueCalculator {
	return &CustomValueCalculator{
		handler: handler,
	}
}
