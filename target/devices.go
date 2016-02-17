package target

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type MixResult struct {
	Request *lh.LHRequest
	Data    []byte
}

// Devices that can mix liquids in batched execution
type Mixer interface {
	Device
	PrepareMix(mixes []*wtype.LHInstruction) (*MixResult, error)
}

// Devices that can be configured
type Shaper interface {
	Device
	Shape() interface{}
}

// Devices that can move things between devices
type Mover interface {
	Device
	Move(from, to Device) error
}

type Incubator interface {
	Device
	Incubate() error
}
