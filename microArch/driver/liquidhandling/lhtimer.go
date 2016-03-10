package liquidhandling

import "time"

// records timing info
// preliminary implementation assumes all instructions of a given
// type have the same timing, TimeFor is expressed in terms of the instruction
// however so it will be possible to modify this behaviour in future

type LHTimer struct {
	Times []time.Duration
}

func NewTimer() *LHTimer {
	var t LHTimer
	t.Times = make([]time.Duration, 50)
	return &t
}

func (t *LHTimer) TimeFor(r RobotInstruction) time.Duration {
	var d time.Duration
	if r.InstructionType() > 0 && r.InstructionType() < len(t.Times) {
		d = t.Times[r.InstructionType()]
	} else {
	}
	return d
}
