package movement

import (
	"github.com/antha-lang/antha/anthalib/wtype"
)

// This package deals with functions required to move things from
// one physical location to another

// the service side of this will be dealt with in the execution layer

type ConcreteMover struct {
	AnthaObject                       // has an ID, Name and Inst
	Manufacturer string               // should be wrapped in AnthaDevice
	Slots        []wtype.VariableSlot // places where it can hold things
}

func (cm *ConcreteMover) Grab(e Entity) bool {
	for _, s := range cm.Slots {
		if s.CanHold(e) {
			s.Add(e)
			return true
		}
	}

	return false
}

func (cm *ConcreteMover) Drop(s *Slot) Entity {
}
