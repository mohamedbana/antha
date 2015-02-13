// execution/stockservice.go: Part of the Antha language
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

// a map structure for defining requests to the stock manager
type StockRequest map[string] interface{}

// structure to hold communication channels with the stock service
type StockService struct{
	RequestsIn  chan StockRequest
	RequestsOut chan StockRequest
}

// Initialize the communication service 
// The channels have the capacity to hold 5 requests before blocking
func (ss *StockService)Init(){
	ss.RequestsIn =make(chan StockRequest, 5)
	ss.RequestsOut=make(chan StockRequest, 5)

	go func() {
		stockDaemon(ss)
	}()
}

// Make a request to the stock agent and return the result
func (ss *StockService)RequestStock(rin StockRequest)StockRequest{
	ss.RequestsIn<-rin
	rout:=<-ss.RequestsOut
	return rout
}

// a process for monitoring the stock request queue
func stockDaemon(ss *StockService){
	for{
		rin:=<-ss.RequestsIn
		// deal with request

		rin["inst"]=GetUUID()
		// obviously this needs to be a little more complicated!

		ss.RequestsOut<-rin
	}
}
