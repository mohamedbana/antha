package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/workflow"
)

var (
	badFormat = errors.New("bad format")
)

// Flatten []map[string]RawMessage to RawParams
func flattenSList(pdata []byte, param *execute.RawParams) error {
	var params []map[string]json.RawMessage
	if err := json.Unmarshal(pdata, &params); err != nil {
		return err
	} else if len(params) == 0 {
		return badFormat
	}

	var firstKeys map[string]bool
	for idx, p := range params {
		keys := make(map[string]bool)
		for k := range p {
			keys[k] = true
		}

		// Check that keys are the same between values
		if len(keys) == 0 {
			return badFormat
		} else if len(firstKeys) == 0 {
			firstKeys = keys
		} else if len(keys) != len(firstKeys) {
			return badFormat
		} else {
			for k := range firstKeys {
				if !keys[k] {
					return badFormat
				}
			}
		}

		key := fmt.Sprintf("Input%d", idx)
		param.Parameters[key] = p
	}
	return nil
}

// Flatten []RawParams to RawParams
func flattenList(pdata []byte, param *execute.RawParams) error {
	var params []execute.RawParams
	if err := json.Unmarshal(pdata, &params); err != nil {
		return err
	} else if len(params) == 0 {
		return badFormat
	}

	for idx, p := range params {
		if p.Config != nil {
			return fmt.Errorf("cannot extend workflow with config parameters")
		}
		if len(p.Parameters) == 0 {
			return badFormat
		}
		for k, v := range p.Parameters {
			key := fmt.Sprintf("%d%s", idx, k)
			param.Parameters[key] = v
		}
	}
	return nil
}

// Parse parameters and workflow. If there are multiple input parameters and
// workflow is just one element, modify workflow to take multiple parameters.
func tryExpand(wdata, pdata []byte) (desc *workflow.Desc, param *execute.RawParams, err error) {
	getFirstProcess := func(desc *workflow.Desc) *workflow.Process {
		for _, p := range desc.Processes {
			return &p
		}
		return nil
	}

	if err = json.Unmarshal(wdata, &desc); err != nil {
		return
	} else if err = json.Unmarshal(pdata, &param); err == nil {
		// Not list of parameters, use data as-is
		return
	}

	// Update parameters
	param = &execute.RawParams{
		Parameters: make(map[string]map[string]json.RawMessage),
	}
	if len(desc.Connections) > 0 {
		err = fmt.Errorf("cannot extend workflow with connections")
	} else if len(desc.Processes) != 1 {
		err = fmt.Errorf("can only expand workflows with one process")
	} else if err = flattenList(pdata, param); err == nil {
	} else if err = flattenSList(pdata, param); err == nil {
		// ^ This should be last as most parameters will match
	} else {
		err = fmt.Errorf("exhausted methods for expanding workflow")
	}
	if err != nil {
		return
	}

	// Update workflow
	process := getFirstProcess(desc)
	desc = &workflow.Desc{
		Processes: make(map[string]workflow.Process),
	}
	for k := range param.Parameters {
		desc.Processes[k] = *process
	}

	return
}
