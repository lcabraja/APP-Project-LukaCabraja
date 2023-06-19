package stopwatch

type StopwatchEvent int

const (
	ClockPulse        StopwatchEvent = iota
	ClockPulseRunning
	ClockPulsePaused

	SecondReached
	SectionReached
	SectionPassed

	ClockStarted
	ClockPaused
	ClockPausedManually
	ClockPausedAutomatically
	ClockResumed
	ClockToggled
	ClockFinished
	ClockOffset
	ClockIncrement
	ClockDecrement

	ClockReset

	EventFired
	ExceptPulse
)

func (swe StopwatchEvent) String() string {
	names := []string{
		"ClockPulse",
		"SecondReached",
		"SectionReached",
		"SectionPassed",

		"ClockStarted",
		"ClockPaused",
		"ClockPausedManually",
		"ClockPausedAutomatically",
		"ClockResumed",
		"ClockToggled",
		"ClockFinished",
		"ClockOffset",
		"ClockIncrement",
		"ClockDecrement",

		"ClockReset",

		"EventFired",
		"ExceptPulse",
	}

	if swe < ClockPulse || swe > ExceptPulse {
		return "UnknownStopwatchEvent"
	}

	return names[swe]
}
