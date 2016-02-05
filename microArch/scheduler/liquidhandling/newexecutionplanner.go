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

import "github.com/antha-lang/antha/microArch/driver/liquidhandling"

func ImprovedExecutionPlanner(request *LHRequest, parameters *liquidhandling.LHProperties) *LHRequest {
	/*
		// this volume correction needs to be removed asap
		// essentially its purpose is to account for extra volume lost
		// while clinging to outside of tips
		volume_correction := 0.5
	*/

	// 1 -- set output order, this is based on dependencies
	set_output_order(request)

	// 2 -- we might optimize at this point: for instance grouping components
	//      or generating stages of execution

	// optimize_runs(request)

	// 3 -- generate top-level instructions

	for _, insID := range request.Output_order {
		request.InstructionSet.Add(ConvertInstruction(request.LHInstructions[insID]))
	}

	return request
}
