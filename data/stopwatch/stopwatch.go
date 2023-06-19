package stopwatch

import (
	"context"
	"encoding/json"
	"github.com/lcabraja/APP-Project-LukaCabraja/data"
	"time"

	"github.com/google/uuid"
)

type Stopwatch struct {


	name string


	ctx                 context.Context
	cancel              context.CancelFunc
	started             bool
	paused              bool
	pausedAutomatically bool


	firstStart      time.Time
	prevTime        time.Time
	cancelled       time.Time
	elapsed         time.Duration
	durationPaused  time.Duration
	lastWholeSecond time.Duration
	manualOffset    time.Duration


	frameRate    int
	minimumValue time.Duration
	startingTime time.Duration
	maximumValue time.Duration
	pauses       []time.Duration
	pauseIdx     int


	subscriptions map[string]*StopwatchEventSubscription `json:"-"`
}


func NewStopwatch(name string, frameRate int, minimumValue, maximumValue time.Duration, pauses ...time.Duration) *Stopwatch {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewStopwatchWithContext(name, ctx, cancel, frameRate, minimumValue, maximumValue, pauses...)
	return s
}

func NewStopwatchWithContext(name string, ctx context.Context, cancel context.CancelFunc, frameRate int, minimumValue, maximumValue time.Duration, pauses ...time.Duration) *Stopwatch {
	return &Stopwatch{
		name: name,

		ctx:    ctx,
		cancel: cancel,
		paused: true,

		elapsed: minimumValue,

		frameRate:    frameRate,
		minimumValue: minimumValue,
		maximumValue: maximumValue,
		pauses:       pauses,

		subscriptions: make(map[string]*StopwatchEventSubscription),
	}
}

func (sw *Stopwatch) fireEvent(e StopwatchEvent) {
	for _, sub := range sw.subscriptions {
		if sub.Event == e || sub.Event == EventFired || (sub.Event == ExceptPulse && e > 2) {
			go sub.Handler(e, sw.Elapsed(), sw)
		}
	}
}

func (sw *Stopwatch) Subscribe(sub *StopwatchEventSubscription) string {
	uuid := uuid.NewString()
	sw.subscriptions[uuid] = sub
	return uuid
}

func (sw *Stopwatch) Unsubscribe(uuid string) {
	delete(sw.subscriptions, uuid)
}

func (sw *Stopwatch) SubscriptionsJson() ([]byte, error) {
	if jsonData, err := json.Marshal(sw.subscriptions); err != nil {
		return nil, err
	} else {
		return jsonData, nil
	}
}

func (sw *Stopwatch) SubscriptionsCount() int {
	return len(sw.subscriptions)
}

func (sw *Stopwatch) GetName() string {
	return sw.name
}

func (sw *Stopwatch) GetType() data.ResourceType {
	return data.StopwatchType
}

func (sw *Stopwatch) Start() {
	if !sw.started {
		sw.started = true
		sw.firstStart = time.Now()
		sw.prevTime = time.Now()
		go sw.worker()
		sw.fireEvent(ClockStarted)
	}
}

func (sw *Stopwatch) worker() {
	for {
		select {
		case <-sw.ctx.Done():
			return
		default:
			sw.work()
			time.Sleep(time.Second / time.Duration(sw.frameRate))
		}
	}
}

func (sw *Stopwatch) work() {
	sw.fireEvent(ClockPulse)
	if sw.paused {
		sw.durationPaused += time.Since(sw.prevTime)
		sw.prevTime = time.Now()
		sw.fireEvent(ClockPulsePaused)
	} else {
		sw.fireEvent(ClockPulseRunning)

		if sw.Elapsed() > sw.maximumValue {
			sw.pause(sw.maximumValue)
			sw.fireEvent(ClockFinished)
			return
		}

		if sw.pauseIdx < len(sw.pauses) && sw.Elapsed() > sw.pauses[sw.pauseIdx] {
			sw.pause(sw.pauses[sw.pauseIdx])
			sw.pauseIdx += 1
			return
		}

		if sw.Elapsed()-sw.lastWholeSecond > time.Second {
			sw.lastWholeSecond += time.Second
			sw.fireEvent(SecondReached)
		}

		sw.elapsed += time.Since(sw.prevTime)
		sw.prevTime = time.Now()
	}
}


func (sw *Stopwatch) Stop() {
	sw.cancelled = time.Now()
	sw.cancel()
	sw.paused = true
	sw.pausedAutomatically = false
	sw.fireEvent(ClockPaused)
	sw.fireEvent(ClockFinished)
}

func (sw *Stopwatch) IsStarted() bool {
	return sw.started
}

func (sw *Stopwatch) IsPaused() bool {
	return sw.paused
}

func (sw *Stopwatch) Pause() {
	sw.paused = true
	sw.fireEvent(ClockPausedManually)
	sw.fireEvent(ClockPaused)
	sw.fireEvent(ClockToggled)
}

func (sw *Stopwatch) pause(pause time.Duration) {
	sw.paused = true
	sw.pausedAutomatically = true

	sw.manualOffset = 0
	sw.elapsed = pause
	sw.lastWholeSecond = time.Duration(pause.Seconds()) * time.Second

	sw.fireEvent(ClockPausedAutomatically)
	sw.fireEvent(ClockPaused)
	sw.fireEvent(SecondReached)
	sw.fireEvent(SectionReached)
}

func (sw *Stopwatch) resume(resume time.Duration) {
	sw.paused = false
	sw.pausedAutomatically = false

	sw.manualOffset = 0
	sw.elapsed = resume

}

func (sw *Stopwatch) Resume() {
	if !sw.started {
		sw.Start()
	}

	if sw.paused {
		sw.prevTime = time.Now()
		sw.paused = false
		if sw.pausedAutomatically {
			sw.pausedAutomatically = false
			sw.fireEvent(SectionPassed)
		}
		sw.fireEvent(ClockResumed)
		sw.fireEvent(ClockToggled)
	}
}

func (sw *Stopwatch) Toggle() {
	if !sw.IsPaused() {
		sw.Pause()
	} else {
		sw.Resume()
	}
}

func (sw *Stopwatch) Reset() {
	ctx, cancel := context.WithCancel(context.Background())
	sw.ctx = ctx
	sw.cancel = cancel
	sw.firstStart = time.Now()
	sw.prevTime = time.Now()
	sw.fireEvent(ClockReset)
}

func (sw *Stopwatch) Offset(d time.Duration, unit time.Duration) {
	sw.fireEvent(ClockOffset)

	newTime := sw.elapsed + sw.manualOffset + d*unit
	if newTime > sw.maximumValue {
		sw.pause(sw.maximumValue)
		return
	}
	if sw.pauseIdx < len(sw.pauses) && newTime > sw.pauses[sw.pauseIdx] {
		sw.pause(sw.pauses[sw.pauseIdx])
		sw.pauseIdx += 1
		return
	}
	if newTime < sw.minimumValue {
		sw.elapsed = sw.minimumValue
		sw.manualOffset = 0
		sw.pauseIdx = 0
		return
	}

	if d < 0 && sw.pauseIdx > 0 && sw.pauses[sw.pauseIdx-1] > newTime {
		sw.pauseIdx -= 1
	}

	sw.manualOffset += d * unit
	sw.lastWholeSecond = time.Duration(sw.Elapsed().Seconds()) * time.Second

	if d > 0 {
		sw.fireEvent(ClockIncrement)
	} else {
		sw.fireEvent(ClockDecrement)
	}
}

func (sw *Stopwatch) Elapsed() time.Duration {
	return sw.elapsed + sw.manualOffset
}

func (sw *Stopwatch) Json(additional map[string]interface{}) ([]byte, error) {
	var nextPause time.Duration

	if sw.pauseIdx < len(sw.pauses) {
		nextPause = sw.pauses[sw.pauseIdx]
	} else {
		nextPause = sw.maximumValue
	}

	var swMap = map[string]interface{}{
		"stopwatch": sw.name,
		"elapsed":   sw.Elapsed(),
		"min":       sw.minimumValue,
		"max":       sw.maximumValue,
		"paused":    sw.IsPaused(),
		"pauses":    sw.pauses,
		"nextPause": nextPause,
	}

	for k, v := range additional {
		if _, present := swMap[k]; !present {
			swMap[k] = v
		}
	}

	data, err := json.Marshal(swMap)

	if err != nil {
		return nil, err
	}
	return data, nil
}
