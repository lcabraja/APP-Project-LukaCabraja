package log

import "fmt"

type LogLevel int

const (
	DEV LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var names = []string{
	"DEV",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

func (ll LogLevel) String() string {
	if ll < DEV || ll > FATAL {
		return fmt.Sprintf("LogLevel [%d] not found", ll)
	}
	return names[ll]
}

var colors = []Color{
	Magenta,
	Cyan,
	Blue,
	Yellow,
	Red,
	Black,
}

func (ll LogLevel) Color() Color {
	if ll < DEV || ll > FATAL {
		return Reset
	}
	return colors[ll]
}

var background = []Color{
	BackgroundReset,
	BackgroundReset,
	BackgroundReset,
	BackgroundReset,
	BackgroundReset,
	BackgroundWhite,
}

func (ll LogLevel) Background() Color {
	if ll < DEV || ll > FATAL {
		return Reset
	}
	return background[ll]
}
