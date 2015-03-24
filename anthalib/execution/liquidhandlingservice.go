// execution/liquidhandlerservice.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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
	"github.com/antha-lang/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/execute"
	"sync"
)

// the liquid handler holds channels for communicating
// with the liquid handling service provider
type LiquidHandlingService struct {
	RequestsIn   chan *liquidhandling.LHRequest
	RequestsOut  map[string]chan *liquidhandling.LHRequest
	RequestQueue map[execute.ThreadID]*liquidhandling.LHRequest
	lock         *sync.Mutex
}

// Initialize the liquid handling service
// the channels are given a default capacity for 5 items before
// blocking
func (lhs *LiquidHandlingService) Init() {
	lhs.RequestsIn = make(chan *liquidhandling.LHRequest, 5)
	lhs.RequestsOut = make(map[string]chan *liquidhandling.LHRequest, 10)
	lhs.RequestQueue = make(map[execute.ThreadID]*liquidhandling.LHRequest)
	lhs.lock = new(sync.Mutex)
	go func() {
		liquidhandlingDaemon(lhs)
	}()
}

func (lhs *LiquidHandlingService) MakeMixRequest(solution *liquidhandling.LHSolution) *liquidhandling.LHRequest {
	lhs.lock.Lock()
	defer lhs.lock.Unlock()
	rq, ok := lhs.RequestQueue[execute.ThreadID(solution.ID)]

	if !ok {
		// if we don't have a request with this ID, make a new one
		rq = liquidhandling.NewLHRequest()
	}

	rq.Output_solutions[solution.ID] = solution
	lhs.RequestQueue[execute.ThreadID(solution.ID)] = rq

	return rq
}

// Daemon for passing requests through to the service
// when do these output channels get destroyed? Now I guess
func liquidhandlingDaemon(lhs *LiquidHandlingService) {
	for {
		rin := <-lhs.RequestsIn

		// handle request
		// what do we do?
		//

		rout, ok := lhs.RequestsOut[rin.ID]

		if !ok {
			panic("Liquidhandlingdaemon: No channel for request output")
		}

		delete(lhs.RequestsOut, rin.ID)

		rout <- rin
	}
}
