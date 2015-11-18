package liquidhandling

import (
	"time"
)

func GetTimerFor(model, mnfr string) *LHTimer {
	timers := makeTimers()
	_, ok := timers[model+mnfr]
	if ok {
		return timers[model+mnfr]
	}

	return nil
}

func makeTimers() map[string]*LHTimer {
	timers := make(map[string]*LHTimer, 2)

	timers["GilsonPipetmax"] = makeGilsonPipetmaxTimer()

	return timers
}

func makeGilsonPipetmaxTimer() *LHTimer {
	// first instance an intermediate level of detail is probably the best we can do
	t := NewTimer()
	t.Times[7], _ = time.ParseDuration("8s")
	t.Times[8], _ = time.ParseDuration("6s")
	t.Times[19], _ = time.ParseDuration("4s")
	t.Times[20], _ = time.ParseDuration("4s")
	return t
}
