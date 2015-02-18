// execution/execution.go: Part of the Antha language
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

//import "fmt"

// A structure which defines the execution context
// - holds pointers to the services
type ExecutionService struct{
	StockMgr *StockService
	SampleTracker *SampleTrackerService
	Logger *LogService
	Scheduler *ScheduleService
	GarbageCollector *GarbageCollectionService
	// and also to the config
	Config *AnthaConfig
}

// globally accessible execution context variable
var executionContext ExecutionService
var servicesStarted bool = false
var mutex sync.Mutex

// not a true function: has the side-effect of starting runtime services
func StartRuntime(){
	sm:=NewStockService()
	st:=NewSampleTrackerService()
	l:=NewLogService()
	s:=NewScheduleService()
	gc:=NewGarbageCollectionService()
	ac:=NewAnthaConfig()
	executionContext=ExecutionService{sm, st, l, s, gc, ac}
}

// accessor for above to enforce singleton status

func GetContext()*ExecutionService{
	mutex.Lock()
	defer mutex.Unlock()

	if(!servicesStarted){
		StartRuntime()
		servicesStarted=true
	}

	return &executionContext
}



