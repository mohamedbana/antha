// /anthalib/execution/services_ddn.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package execution

import (
	"errors"
)

type WorkflowId interface{}
type DeviceId interface{}

var dummy interface{}

///////////////////////

// Manages service deployment
type ServiceDeployer struct{}

///////////////////////

// Logs service actions
type Logger struct{}

///////////////////////

// Authenticates and validates service requests
type Auth struct{}

///////////////////////

// Presents a friendlier interface for the web frontend
type WebFrontEndCondenser struct{}

///////////////////////

// Manages the consumable resources currently used in experiments
type SampleManager struct{}

///////////////////////

// Manages physical devices
type DeviceManager struct{}

///////////////////////

// Manages disposal of physical waste
type SampleDisposer struct{}

///////////////////////

// Manages the consumable resources not currently used in experiments
type StockManager struct{}

///////////////////////

// Schedules workflows onto devices
type WorkflowScheduler struct{}
type Amount interface{}
type ConsumableId interface{}
type Deadline interface{}

type WorkflowStatus interface {
	Feasible() bool
	Done() bool
}
type Consumable interface {
	ConsumableId() ConsumableId
	Amount() Amount
}
type WorkflowPlan interface {
	WorkflowId() WorkflowId
	Deadline() Deadline
	Devices() []DeviceId
	Consumables() []Consumable
}

// Estimates resources required (i.e., time, stock, devices) to execute a workflow
func (s *WorkflowScheduler) Plan(script string) (WorkflowPlan, error) {
	return nil, nil
}

func (m *WorkflowScheduler) Run(WorkflowId) error {
	return nil
}

func (m *WorkflowScheduler) Release(WorkflowId) error {
	return nil
}

// Checks the status of a workflow
func (m *WorkflowScheduler) Check(WorkflowId) (WorkflowStatus, error) {
	return nil, nil
}

// Subscribes to updates to workflow
func (m *WorkflowScheduler) Subscribe(WorkflowId) (chan WorkflowStatus, error) {
	return nil, nil
}

///////////////////////

func PlanExperiment(script string) (WorkflowId, error) {
	_, err := dummy.(*WorkflowScheduler).Plan(script)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func ReleaseExperiment(wid WorkflowId) error {
	err := dummy.(*WorkflowScheduler).Release(wid)
	return err
}

func RunExperiment(wid WorkflowId) error {
	// Check if workflow is feasible
	ws, err := dummy.(*WorkflowScheduler).Check(wid)
	if err != nil {
		return err
	}
	if !ws.Feasible() {
		return errors.New("Unfeasible workflow")
	}
	err = dummy.(*WorkflowScheduler).Run(wid)
	if err != nil {
		return err
	}

	// Wait for experiment to finish
	ch, err := dummy.(*WorkflowScheduler).Subscribe(wid)
	if err != nil {
		return err
	}
Wait:
	for {
		switch msg := <-ch; {
		case msg.Done():
			break Wait
		}
	}
	return nil
}
