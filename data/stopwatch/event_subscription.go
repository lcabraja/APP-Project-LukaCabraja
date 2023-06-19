package stopwatch

import "time"

type StopwatchEventSubscription struct {
	Event   StopwatchEvent
	Handler func(StopwatchEvent, time.Duration, *Stopwatch)
}
