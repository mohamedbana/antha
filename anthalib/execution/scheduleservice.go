// scheduler/scheduleservice.go: Part of the Antha language
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

// data structure for holding channels to communicate
// with the scheduler
type ScheduleService struct {
	RequestsIn  chan ScheduleRequest
	RequestsOut chan ScheduleRequest
}

// start up the scheduler and make the communication channels
// the channels have room to buffer 5 requests before blocking
func (ss *ScheduleService) Init() {
	ss.RequestsIn = make(chan ScheduleRequest, 5)
	ss.RequestsOut = make(chan ScheduleRequest, 5)

	go func() {
		scheduleDaemon(ss)
	}()
}

// communicate a request to the scheduler and return the reply
func (ss *ScheduleService) RequestSchedule(rin ScheduleRequest) ScheduleRequest {
	ss.RequestsIn <- rin
	rout := <-ss.RequestsOut
	return rout
}

// a process for monitoring the request queue
func scheduleDaemon(ss *ScheduleService) {
	for {
		rin := <-ss.RequestsIn
		ss.RequestsOut <- rin
	}
}
