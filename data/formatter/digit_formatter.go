package formatter

import (
	"fmt"
	"math"
)

type DigitFormatter struct {
	digits int
	cutoff bool
}

func NewDigitFormatter(digits int, cutoff bool) (*DigitFormatter, error) {
	if digits < 1 {
		return nil, fmt.Errorf("digits must be greater than 0")
	}

	return &DigitFormatter{
		digits: digits,
		cutoff: cutoff,
	}, nil
}

func (df *DigitFormatter) Format(value interface{}) (interface{}, error) {
	intval, ok := value.(int)
	if !ok {
		return nil, fmt.Errorf("value %s is not an int", value)
	}

	if df.cutoff {
		intval = intval % int(math.Pow10(df.digits))
	}

	return fmt.Sprintf("%0*d", df.digits, intval), nil
}
