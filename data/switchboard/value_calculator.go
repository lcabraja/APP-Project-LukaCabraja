package switchboard

type ValueCalculator interface {
	Result(map[string]interface{}) (string, error)
}
