// logger/logservice.go: Part of the Antha language
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

// structure defining a log service - holds
// channels for communicating with the service
type LogService struct {
	RequestsIn  chan LogRequest
	RequestsOut chan LogRequest
}

// initialize the log service
// the channels have capacity to buffer 5 requests by default
func (ls *LogService) Init() {
	ls.RequestsIn = make(chan LogRequest, 5)
	ls.RequestsOut = make(chan LogRequest, 5)

	go func() {
		logDaemon(ls)
	}()
}

// function for communicating with the log daemon
func (ls *LogService) RequestLog(rin LogRequest) LogRequest {
	ls.RequestsIn <- rin
	rout := <-ls.RequestsOut
	return rout
}

// daemon process for monitoring the queue
func logDaemon(ls *LogService) {
	for {
		rin := <-ls.RequestsIn
		ls.RequestsOut <- rin
	}
}
