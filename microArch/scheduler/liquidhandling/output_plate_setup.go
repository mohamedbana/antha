// anthalib//liquidhandling/output_plate_setup.go: Part of the Antha language
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
	"errors"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

//  TASK: 	define output plates
// INPUT: 	"output_platetype", "outputs"
//OUTPUT: 	"output_plates"      -- these each have components in wells
//		"output_assignments" -- map with arrays of assignment strings, i.e. {tea: [plate1:A:1, plate1:A:2...] }etc.
func output_plate_setup(request *LHRequest) *LHRequest {
	logger.Debug("in output plate setup")
	//(map[string]*wtype.LHPlate, map[string][]string) {
	output_platetype := (*request).Output_platetypes
	if output_platetype == nil || len(output_platetype) == 0 {
		logger.Debug(fmt.Sprint(output_platetype, " ", len(output_platetype)))
		wutil.Error(errors.New("plate_setup: No output plate type defined"))
	}

	// this might be OK if we have pre-defined output
	if (*request).Output_major_group_layouts == nil {
		wutil.Error(errors.New("plate setup: Output major groups undefined"))
	}

	output_plates := (*request).Output_plates

	if len(output_plates) == 0 {
		output_plates = make(map[string]*wtype.LHPlate, len(request.Output_major_group_layouts))
	}

	//

	// just assign based on number of groups

	opl := request.Output_plate_layout

	for i := 0; i < len(request.Output_major_group_layouts); i++ {
		//p := wtype.New_Plate(request.Output_platetype)
		//p := factory.GetPlateByType(request.Output_platetype.Type)
		mgl := request.Output_major_group_layouts[i]
		//p := factory.GetPlateByType(output_platetype[mgl])
		p := factory.GetPlateByType(output_platetype[0].Type)
		output_plates[p.ID] = p
		opl[i] = p.ID
		name := fmt.Sprintf("Output_plate_%d", i+1)
		p.PlateName = name
	}

	(*request).Output_plate_layout = opl
	(*request).Output_plates = output_plates
	return request
}
