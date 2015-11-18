package liquidhandling

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

func makeGilsonPipetMaxTimer() *LHTimer {

}
