// anthalib//liquidhandling/newexecutionplanner.go: Part of the Antha language
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

package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

// robot here should be a copy... this routine will be destructive of state
func ImprovedExecutionPlanner(request *LHRequest, robot *liquidhandling.LHProperties) *LHRequest {
	logger.Info("Improved execution planner YEAH")
	/*
		// this volume correction needs to be removed asap
		// essentially its purpose is to account for extra volume lost
		// while clinging to outside of tips
		volume_correction := 0.5
	*/

	// 1 -- set output order, this is based on dependencies
	//set_output_order(request)
	// this now happens waaaaaay at the beginning

	// 2 -- we might optimize at this point: for instance grouping components
	//      or generating stages of execution
	/*
		newoutputorder := make([]string, 0, 1)
		optimize_runs(request, request.InstructionChain, newoutputorder)
		request.Output_order = newoutputorder
	*/
	// 3 -- generate top-level instructions

	for _, insID := range request.Output_order {
		request.InstructionSet.Add(ConvertInstruction(request.LHInstructions[insID], robot))
	}

	// 4 -- make the low-level instructions

	inx := request.InstructionSet.Generate(request.Policies, robot)
	instrx := make([]liquidhandling.TerminalRobotInstruction, len(inx))
	for i := 0; i < len(inx); i++ {
		fmt.Println(liquidhandling.InsToString(inx[i]))
		instrx[i] = inx[i].(liquidhandling.TerminalRobotInstruction)
	}
	request.Instructions = instrx

	return request
}
