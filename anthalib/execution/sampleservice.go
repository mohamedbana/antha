// execution/sampleservice.go: Part of the Antha language
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

// data structure defining sample requests
type SampleRequest map[string] interface{}

// struct to hold channels for communicating with 
// the sample tracker service
type SampleTrackerService struct{
	RequestsIn  chan SampleRequest
	RequestsOut chan SampleRequest
}

// initialize the sample tracker agent
// the channels have capacity to buffer 5 requests before blocking
func (ss *SampleTrackerService)Init(){
	ss.RequestsIn =make(chan SampleRequest, 5)
	ss.RequestsOut=make(chan SampleRequest, 5)

	go func() {
		sampleDaemon(ss)
	}()
}

// send a request to the sample tracking service - e.g. update status or 
// register a new sample
func (ss *SampleTrackerService)TrackSample(rin SampleRequest)SampleRequest{
	ss.RequestsIn<-rin
	rout:=<-ss.RequestsOut
	return rout
}

// daemon process to keep the queue moving
func sampleDaemon(ss *SampleTrackerService){
	for{
		rin:=<-ss.RequestsIn
		// deal with request
		ss.RequestsOut<-rin
	}
}
