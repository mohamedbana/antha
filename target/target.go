// Package target provides the construction of a target machine from a
// collection of devices
package target

import (
	"errors"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

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

// Target machine for execution.
type Target struct {
	devices []Device
	runners map[string][]*Runner
}

func New() *Target {
	return &Target{
		runners: make(map[string][]*Runner),
	}
}

func (a *Target) canCompile(d Device, reqs ...ast.Request) bool {
	for _, req := range reqs {
		if !d.CanCompile(req) {
			return false
		}
	}
	return true
}

func (a *Target) CanCompile(reqs ...ast.Request) (r []Device) {
	for _, d := range a.devices {
		if a.canCompile(d, reqs...) {
			r = append(r, d)
		}
	}
	return
}

func (a *Target) CanRun(ftype string) []*Runner {
	return a.runners[ftype]
}

func (a *Target) Runners() (rs []*Runner) {
	for _, d := range a.devices {
		if r, ok := d.(*Runner); ok {
			rs = append(rs, r)
		}
	}
	return
}

func (a *Target) AddDevice(d Device) error {
	a.devices = append(a.devices, d)
	if r, ok := d.(*Runner); ok {
		ftypes, err := r.types()
		if err != nil {
			return err
		}
		for _, ftype := range ftypes {
			a.runners[ftype] = append(a.runners[ftype], r)
		}
	}

	return nil

}
