// Package target provides the construction of a target machine from a
// collection of devices
package target

import (
	"errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

// Capabilities for Devices
//   - MinVol, MaxVol
//   - MinIncubateTime, MaxIncubateTime
//   - MinIncubateTemp, MaxIncubateTemp
// Capabilities for movers
//   - Cost(Device, Device) int

// TODO(ddn): Add target instruction description

var (
	noLh     = errors.New("no liquid handler found")
	noTarget = errors.New("no target configuration found")
)

type targetKey int

const theTargetKey targetKey = 0

func GetTarget(ctx context.Context) (*Target, error) {
	v, ok := ctx.Value(theTargetKey).(*Target)
	if !ok {
		return nil, noTarget
	}
	return v, nil
}

func WithTarget(parent context.Context, t *Target) context.Context {
	return context.WithValue(parent, theTargetKey, t)
}

// TODO(ddn): Replace with a more efficient data structure (interval tree)

// An interval or union thereof
type Interval struct {
	values []struct{ A, B float64 }
}

// The nil interval does not contain any points
func (a Interval) Nil() bool {
	return len(a.values) == 0
}

func (a Interval) Contains(x, y float64) bool {
	for _, v := range a.values {
		if v.A <= x && y <= v.B {
			return true
		}
	}
	return false
}

func (a Interval) Add(x Interval) *Interval {
	var values []struct{ A, B float64 }
	for _, v := range a.values {
		values = append(values, v)
	}
	for _, v := range x.values {
		values = append(values, v)
	}
	return &Interval{values: values}
}

// Create the interval [a, b]
func NewInterval(a, b float64) *Interval {
	return &Interval{
		values: []struct{ A, B float64 }{struct{ A, B float64 }{A: a, B: b}},
	}
}

type Request struct {
	MixVol Interval
	Temp   Interval
	Time   Interval
}

type Device interface {
	Can(reqs ...Request) bool // Can device handle this request
	MoveCost(from Device) int // A non-negative cost to move to this device
}

// Target machine for execution.
//
// NB(ddn): API is in flux while the abstractions for targets are being worked
// out (29-01-2016).
type Target struct {
	devices []Device
}

func New() *Target {
	// TODO(ddn): Add Generic Manual Equipment
	return &Target{}
}

func (a *Target) Can(reqs ...Request) (r []Device) {
	for _, d := range a.devices {
		if d.Can(reqs...) {
			r = append(r, d)
		}
	}
	return
}

func (a *Target) AddDevice(d Device) {
	a.devices = append(a.devices, d)
}

// XXX(ddn): remove after compile refactor is done
func (a *Target) Mix(insts []*wtype.LHInstruction) error {
	for _, d := range a.devices {
		if mixer, ok := d.(Mixer); !ok {
			continue
		} else if _, err := mixer.PrepareMix(insts); err != nil {
			return err
		} else {
			return nil
		}
	}
	return noLh
}
