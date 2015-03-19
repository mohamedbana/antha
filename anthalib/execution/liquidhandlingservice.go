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
)

// the liquid handler holds channels for communicating
// with the liquid handling service provider
type LiquidHandlingService struct {
	RequestsIn  chan liquidhandling.LHRequest
	RequestsOut chan liquidhandling.LHRequest
}

// Initialize the liquid handling service
// the channels are given a default capacity for 5 items before
// blocking
func (lhs *LiquidHandlingService) Init() {
	lhs.RequestsIn = make(chan liquidhandling.LHRequest, 5)
	lhs.RequestsOut = make(chan liquidhandling.LHRequest, 5)

	go func() {
		liquidhandlingDaemon(lhs)
	}()
}

// send an item to the waste stream
func (lhs *LiquidHandlingService) RequestLiquidHandling(rin liquidhandling.LHRequest) liquidhandling.LHRequest {
	lhs.RequestsIn <- rin
	rout := <-lhs.RequestsOut
	return rout
}

// Daemon for passing requests through to the service
func liquidhandlingDaemon(lhs *LiquidHandlingService) {
	for {
		rin := <-lhs.RequestsIn
		// do something
		lhs.RequestsOut <- rin
	}
}
