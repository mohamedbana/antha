package inject

import (
	"fmt"
	"golang.org/x/net/context"
)

type RunFunc func(context.Context, Value) (Value, error)

type Runner interface {
	Run(context.Context, Value) (Value, error)
}

type FuncRunner struct {
	RunFunc
}

func (a *FuncRunner) Run(ctx context.Context, value Value) (Value, error) {
	return a.RunFunc(ctx, value)
}

type CheckedRunner struct {
	RunFunc
	In  interface{}
	Out interface{}
}

func (a *CheckedRunner) Run(ctx context.Context, value Value) (Value, error) {
	inT := a.In
	if err := AssignableTo(value, inT); err != nil {
		return nil, fmt.Errorf("input value not assignable to %T: %s", inT, err)
	}

	out, err := a.RunFunc(ctx, value)

	if err != nil {
		return out, err
	}

	outT := a.Out
	if err != AssignableTo(out, outT) {
		return nil, fmt.Errorf("output value not assignable to %T: %s", outT, err)
	}

	return out, err
}
