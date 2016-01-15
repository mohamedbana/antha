// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"fmt"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	eq "github.com/antha-lang/antha/microArch/equipment"
	em "github.com/antha-lang/antha/microArch/equipmentManager"
	"github.com/antha-lang/antha/workflow"
)

type idKey int

const theIdKey idKey = 0

func getId(ctx context.Context) string {
	v, ok := ctx.Value(theIdKey).(string)
	if !ok {
		return ""
	}
	return v
}

type Options struct {
	WorkflowData  []byte
	ParamData     []byte
	FromEM        em.EquipmentManager // Use equipment handler to find liquid handler
	FromEquipment eq.Equipment        // Use equipment as liquid handler
	Id            string
}

func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData})
	if err != nil {
		return nil, err
	}

	cd, err := setParams(parent, opt.ParamData, w)
	if err != nil {
		return nil, fmt.Errorf("cannot set initial parameters: %s", err)
	}

	var lh eq.Equipment
	switch {
	case opt.FromEquipment != nil:
		lh = opt.FromEquipment
	case opt.FromEM != nil:
		if l, err := getLhFromEm(opt.FromEM); err != nil {
			return nil, err
		} else {
			lh = l
		}
	case lh == nil:
		return nil, noLh
	}

	ctx, done, err := newLHContext(context.WithValue(parent, theIdKey, opt.Id), lh, cd)
	if done != nil {
		defer done()
	}
	if err != nil {
		return nil, fmt.Errorf("cannot initialize liquid handler: %s", err)
	}

	if err := w.Run(ctx); err != nil {
		return nil, err
	}

	return w, nil
}
