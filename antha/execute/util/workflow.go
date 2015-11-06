// /antharun/workflow.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package util

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/antha-lang/antha/antha/component/lib"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/internal/github.com/twinj/uuid"
)

type Workflow struct {
	flow.Graph
	InPorts  []flow.FullPortName
	OutPorts []flow.FullPortName
}

type Errors []error

func (e Errors) Error() string {
	var buf bytes.Buffer
	for _, v := range e {
		buf.WriteString(v.Error())
	}
	return buf.String()
}

type WorkflowRun struct {
	ID       string
	Outs     map[flow.FullPortName]chan execute.ThreadParam
	Ins      map[flow.FullPortName]chan execute.ThreadParam
	Errors   chan error
	Done     chan bool
	Messages chan execute.ThreadParam
}

func doSliceAppend(to interface{}, what reflect.Value) {
	av := reflect.ValueOf(to).Elem()
	av.Set(reflect.Append(av, what))
}

//jman change id to string to avoid vendorized uuid struct to be a problem
// id is actually the block id, as that is waht groups an execution.
func NewWorkflowRun(id execute.BlockID, wf *Workflow, cf *Config) (*WorkflowRun, error) {
	params := make(map[flow.FullPortName][]interface{})
	var first = true
	for _, parameters := range cf.Parameters {
		tid := execute.ThreadID(uuid.NewV4().String())
		for _, port := range wf.InPorts {
			if first {
				//first run
				params[port] = make([]interface{}, 0)
			}
			//			param, ok := cf.Parameters[port.Proc]
			param, ok := parameters[port.Proc]
			if !ok {
				return nil, fmt.Errorf("required parameter not found %v", port)
			}

			// Need to unpack fields to map entries; use reflection to avoid
			// copying pointer based structures
			pv := reflect.ValueOf(param)
			fv := pv.Elem().FieldByName(port.Port)
			if !fv.IsValid() {
				return nil, fmt.Errorf("required parameter not found %v", port)
			}
			var xx []interface{}
			doSliceAppend(&xx, fv)
			p := execute.ThreadParam{
				Value:   xx[0],
				ID:      tid, //TODO this must be generated for every input
				BlockID: id,
				Error:   false,
			}

			params[port] = append(params[port], p) //TODO ugly hack :(
		}
		first = false
	}

	ins := make(map[flow.FullPortName]chan execute.ThreadParam)
	outs := make(map[flow.FullPortName]chan execute.ThreadParam)
	for _, port := range wf.InPorts {
		ins[port] = make(chan execute.ThreadParam)
		wf.Graph.SetInPort(fmt.Sprintf("%s.%s", port.Proc, port.Port), ins[port])
	}
	for _, port := range wf.OutPorts {
		outs[port] = make(chan execute.ThreadParam)
		wf.Graph.SetOutPort(fmt.Sprintf("%s.%s", port.Proc, port.Port), outs[port])
	}

	messages := make(chan execute.ThreadParam)
	errors := make(chan error)
	done := make(chan bool)

	// Point of no return... start running workflow
	flow.RunNet(&wf.Graph)

	sg := sync.WaitGroup{}
	for _, ch := range outs {
		sg.Add(1)
		go func(ch chan execute.ThreadParam) {
			for v := range ch {
				messages <- v
			}
			sg.Done()
		}(ch)
	}

	go func() {
		defer func() {
			for _, ch := range ins {
				close(ch)
			}
		}()

		for port, vch := range params {
			for _, v := range vch {
				ins[port] <- v.(execute.ThreadParam) //param//instantiate the parameter up at the beggining
			}
		}
	}()

	go func() {
		sg.Wait()
		defer close(messages)
		defer close(errors)

		<-wf.Graph.Wait()
		if err := execute.TakeErrors(); err != nil {
			errors <- err
		}
		done <- true
	}()

	return &WorkflowRun{
		ID:       string(id.ThreadID),
		Outs:     outs,
		Ins:      ins,
		Messages: messages,
		Errors:   errors,
		Done:     done,
	}, nil
}

// Asynchronously runs a workflow for one sample
func (w *Workflow) Run(cf *Config) (*WorkflowRun, error) {
	//	id := uuid.NewV4()// let's replace this with the value coming as jobid
	if id := cf.Config["JOBID"]; id == nil {
		return nil, errors.New("missing job id")
	} else if count := cf.Config["OUTPUT_COUNT"]; count == nil {
		return nil, errors.New("missing output count")
	} else if sid, ok := id.(string); !ok {
		return nil, errors.New("unexpected id type")
	} else if fcount, ok := count.(float64); !ok {
		return nil, errors.New("unexpected count type")
	} else {
		bid := execute.BlockID{ThreadID: execute.ThreadID(sid), OutputCount: int(fcount)}

		//	ctx := execution.GetContext()
		//	ctx.ConfigService.SetConfig(bid, cf.Config)
		if wfr, err := NewWorkflowRun(bid, w, cf); err != nil {
			return nil, err
		} else {
			return wfr, nil
		}
	}
}

func NewWorkflow(js []byte) (*Workflow, error) {
	//g := flow.ParseJSON(js)
	g, err := flow.ParseJSON(js)
	if err != nil {
		return nil, err
	}

	wf := &Workflow{
		Graph: *g,
	}
	for _, port := range wf.GetUnboundInPorts() {
		wf.InPorts = append(wf.InPorts, port)
		wf.Graph.MapInPort(fmt.Sprintf("%s.%s", port.Proc, port.Port), port.Proc, port.Port)
	}
	for _, port := range wf.GetUnboundOutPorts() {
		wf.OutPorts = append(wf.OutPorts, port)
		wf.Graph.MapOutPort(fmt.Sprintf("%s.%s", port.Proc, port.Port), port.Proc, port.Port)
	}
	return wf, nil
}

func init() {
	cs := lib.GetComponents()
	for _, c := range cs {
		flow.Register(c.Name, c.Constructor)
	}
}
