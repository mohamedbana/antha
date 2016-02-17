package human

import (
	"errors"
	"io"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/target"
)

const (
	HumanByHumanCost = 50  // Cost of manually moving from another human device
	HumanByXCost     = 100 // Cost of manually moving from any non-human device
)

var (
	tbd = errors.New("to be implemented")
)

var (
	_ target.Device = &Human{}
	_ target.Mixer  = &Human{}
	_ target.Mover  = &Human{}
)

type Human struct {
}

func (a *Human) Can(...target.Request) bool {
	return true
}

func (a *Human) MoveCost(from target.Device) int {
	if _, ok := from.(*Human); ok {
		return HumanByHumanCost
	}
	return HumanByXCost
}

func (a *Human) Move(from, to target.Device) error {
	return tbd
}

func (a *Human) PrepareMix(mixes []*wtype.LHInstruction) (*target.MixResult, error) {
	return nil, nil
}

type Opt struct {
	In  io.Reader
	Out io.Writer
}

func New(opt Opt) *Human {
	return nil
}
