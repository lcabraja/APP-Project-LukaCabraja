package stopwatch

import (
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"time"
)

func NewTimeoutStopwatch(name string, parent *Stopwatch) *Stopwatch {
	to := NewStopwatch(name, 15, -time.Minute*2, 0)
	parent.Subscribe(&StopwatchEventSubscription{
		Event: ClockStarted,
		Handler: func(e StopwatchEvent, t time.Duration, sw *Stopwatch) {
			to.Start()
		},
	})
	parent.Subscribe(&StopwatchEventSubscription{
		Event: ClockResumed,
		Handler: func(e StopwatchEvent, t time.Duration, sw *Stopwatch) {
			to.Resume()
		},
	})
	parent.Subscribe(&StopwatchEventSubscription{
		Event: ClockPaused,
		Handler: func(e StopwatchEvent, t time.Duration, sw *Stopwatch) {
			to.Pause()
		},
	})
	if parent.IsStarted() {
		to.Start()
	}

	if !parent.IsPaused() {
		to.Resume()
	}

	log.Dev("parent, timeout =", parent.IsPaused(), to.IsPaused())
	return to
}
