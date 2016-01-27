// Package target provides the construction of a target machine from a
// collection of equipment
package target

import (
	"errors"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
)

// TODO(ddn): Move equipment.Equipment here

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

// Target machine for execution.
//
// NB(ddn): API is in flux while the abstractions for targets are being worked
// out (29-01-2016).
type Target struct {
	equips []equipment.Equipment
}

func New() *Target {
	// TODO(ddn): Add Generic Manual Equipment
	return &Target{}
}

func (a *Target) AddLiquidHandler(e equipment.Equipment) error {
	a.equips = append(a.equips, e)
	return nil
}

func (a *Target) GetLiquidHandler() (equipment.Equipment, error) {
	if len(a.equips) == 0 {
		return nil, noLh
	}
	return a.equips[0], nil
}
