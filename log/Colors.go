package log

import "fmt"

type Color string

const (
	Black   Color = "\033[30m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	White   Color = "\033[37m"
)

const (
	BackgroundBlack   Color = "\033[40m"
	BackgroundRed     Color = "\033[41m"
	BackgroundGreen   Color = "\033[42m"
	BackgroundYellow  Color = "\033[43m"
	BackgroundBlue    Color = "\033[44m"
	BackgroundMagenta Color = "\033[45m"
	BackgroundCyan    Color = "\033[46m"
	BackgroundWhite   Color = "\033[47m"
	BackgroundReset   Color = "\033[49m"
)

const Reset Color = "\033[0m"

const (
	Bright    Color = "\033[1m"
	Dim       Color = "\033[2m"
	Underline Color = "\033[4m"
	Blink     Color = "\033[5m"
	Reverse   Color = "\033[7m"
	Hidden    Color = "\033[8m"
)

func dye(text string, color, background Color, reset bool) string {
	if reset {
		return fmt.Sprintf("%s%s%s%s", background, color, text, Reset)
	} else {
		return fmt.Sprintf("%s%s%s", background, color, text)
	}
}
