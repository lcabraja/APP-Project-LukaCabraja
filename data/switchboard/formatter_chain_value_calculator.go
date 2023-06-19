package switchboard

import (
	"errors"
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/formatter"
)

type FormatterChainValueCalculator struct {
	formatters []formatter.DataFormatter
}

func (fcvc *FormatterChainValueCalculator) Result(dependencies map[string]interface{}) (string, error) {
	switch len(dependencies) {
	case 0:
		return "", errors.New("no dependencies provided")
	case 1:
		var (
			value interface{}
			err   error
		)
		for _, v := range dependencies {
			value = v
			for _, formatter := range fcvc.formatters {
				value, err = formatter.Format(value)
				if err != nil {
					return "", err
				}
			}
			return fmt.Sprintf("%v", value), err
		}
	}
	return "", errors.New("too many dependencies provided, ")
}

func NewFormatterChainValueCalculator(formatters ...formatter.DataFormatter) ValueCalculator {
	return &FormatterChainValueCalculator{
		formatters: formatters,
	}
}
