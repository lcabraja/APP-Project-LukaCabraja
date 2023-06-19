package formatter

import (
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"time"
)

type CustomDurationFormatter struct {
	formatter func(time.Duration) string
}

func NewCustomDurationFormatter(formatter func(d time.Duration) string) *CustomDurationFormatter {
	return &CustomDurationFormatter{
		formatter: formatter,
	}
}

func (tf *CustomDurationFormatter) Format(value interface{}) (interface{}, error) {
	duration, ok := value.(time.Duration)
	if !ok {
		return nil, fmt.Errorf("value %s is not a time.Duration\n", value)
	}

	defer func() {
		if r := recover(); r != nil {
			log.E("Recovered from panic:", r)
		}
	}()
	return tf.formatter(duration), nil
}
