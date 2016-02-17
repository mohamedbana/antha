// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"encoding/json"
	"errors"

	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
	"github.com/antha-lang/antha/workflow"
)

var (
	cannotConfigure = errors.New("cannot configure liquid handler")
)

type Options struct {
	WorkflowData []byte         // JSON data describing workflow
	Workflow     *workflow.Desc // Or workflow directly
	ParamData    []byte         // JSON data describing parameters
	Params       *RawParams     // Or parameters directly
	Target       *target.Target // Target machine configuration
	Id           string         // Job Id
}

// Simple entrypoint for one-shot execution of workflows.
func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData, FromDesc: opt.Workflow})
	if err != nil {
		return nil, err
	}

	var params *RawParams
	if opt.Params != nil {
		params = opt.Params
	} else if opt.ParamData != nil {
		if err := json.Unmarshal(opt.ParamData, &params); err != nil {
			return nil, err
		}
	}

	if _, err := setParams(parent, params, w); err != nil {
		return nil, err
	}

	ctx := target.WithTarget(WithId(parent, opt.Id), opt.Target)

	if err := w.Run(trace.WithResolver(ctx, resolveIntrinsics)); err != nil {
		return nil, err
	}

	return w, nil
}
