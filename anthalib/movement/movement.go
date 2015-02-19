package movement

// This package deals with functions required to move things from
// one physical location to another

// the service side of this will be dealt with in the execution layer

type ConcreteMover struct {
	AnthaObject                  // has an ID, Name and Inst
	Manufacturer string          // should be wrapped in AnthaDevice
	Slots        []*VariableSlot // places where it can hold things
}
