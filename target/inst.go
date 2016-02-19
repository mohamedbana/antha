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

type MixInst struct {
	Dev        Device
	Depends    []Inst
	Request    *lh.LHRequest
	Properties liquidhandling.LHProperties
	Files      []byte
}

func (a *MixInst) Device() Device {
	return a.Dev
}

func (a *MixInst) DependsOn() []Inst {
	return a.Depends
}

func (a *MixInst) SetDependsOn(x []Inst) {
	a.Depends = x
}

type ManualInst struct {
	Dev     Device
	Details string
	Depends []Inst
}

func (a *ManualInst) DependsOn() []Inst {
	return a.Depends
}

func (a *ManualInst) Device() Device {
	return a.Dev
}

func (a *ManualInst) SetDependsOn(x []Inst) {
	a.Depends = x
}

// Virtual instruction to hang dependencies on
type WaitInst struct {
	Depends []Inst
}

func (a *WaitInst) Device() Device {
	return nil
}

func (a *WaitInst) DependsOn() []Inst {
	return a.Depends
}

func (a *WaitInst) SetDependsOn(x []Inst) {
	a.Depends = x
}

// TODO: refine with more accurate representation
type MoveToInst struct {
	Dev     Device
	Depends []Inst
	Comps   []*wtype.LHComponent
}

func (a *MoveToInst) Device() Device {
	return a.Dev
}

func (a *MoveToInst) DependsOn() []Inst {
	return a.Depends
}

func (a *MoveToInst) SetDependsOn(x []Inst) {
	a.Depends = x
}
