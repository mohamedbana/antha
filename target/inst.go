package target

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type Inst interface {
	Device() Device
	DependsOn() []Inst
	SetDependsOn([]Inst)
}

type Mix struct {
	Dev        Device
	Depends    []Inst
	Request    *lh.LHRequest
	Properties liquidhandling.LHProperties
	Files      []byte
}

func (a *Mix) Device() Device {
	return a.Dev
}

func (a *Mix) DependsOn() []Inst {
	return a.Depends
}

func (a *Mix) SetDependsOn(x []Inst) {
	a.Depends = x
}

type Manual struct {
	Dev     Device
	Details string
	Depends []Inst
}

func (a *Manual) DependsOn() []Inst {
	return a.Depends
}

func (a *Manual) Device() Device {
	return a.Dev
}

func (a *Manual) SetDependsOn(x []Inst) {
	a.Depends = x
}

// Virtual instruction to hang dependencies on
type Wait struct {
	Depends []Inst
}

func (a *Wait) Device() Device {
	return nil
}

func (a *Wait) DependsOn() []Inst {
	return a.Depends
}

func (a *Wait) SetDependsOn(x []Inst) {
	a.Depends = x
}

// TODO: refine with more accurate representation
type Move struct {
	Dev     Device
	Depends []Inst
	Comps   []*wtype.LHComponent
}

func (a *Move) Device() Device {
	return a.Dev
}

func (a *Move) DependsOn() []Inst {
	return a.Depends
}

func (a *Move) SetDependsOn(x []Inst) {
	a.Depends = x
}
