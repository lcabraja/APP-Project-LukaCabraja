package formatter

import (
	"fmt"
	"time"
)

type TimeFormatter struct {
	format string
}

func NewTimeFormatter(format string) *TimeFormatter {
	return &TimeFormatter{
		format: format,
	}
}

func (tf *TimeFormatter) Format(value interface{}) (interface{}, error) {
	duration, ok := value.(time.Duration)
	if !ok {
		return nil, fmt.Errorf("value %s is not a time.Duration\n", value)
	}

	return FormatTime(tf.format, duration), nil
}


func FormatTime(format string, duration time.Duration) string {
	t := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(duration)
	return t.Format(format)
}
