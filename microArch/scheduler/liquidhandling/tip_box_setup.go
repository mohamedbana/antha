// anthalib//liquidhandling/tip_box__setup.go: Part of the Antha language
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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	lhdriver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

//  TASK: 	Determine number of tip boxes of each type
// INPUT: 	instructions
//OUTPUT: 	arrays of tip boxes
func (lh *Liquidhandler) Tip_box_setup(request *LHRequest) *LHRequest {
	tip_boxes := (*request).Tips
	if len(tip_boxes) == 0 {
		tip_boxes = make([]*wtype.LHTipbox, 0)
	}

	// the instructions have been generated at this point so we just need to go through and count the tips used
	// of each type
	instrx := request.Instructions
	ntips := make(map[string]int)

	// aide memoire: ultimately these tip types derive from the LHParams object which was passed into
	// the call to generating concrete instructions
	// there is a "tips" field which is just an array of LHTip objects. The Name field from one of these
	// is passed through to TIPTYPE
	// tips are ultimately chosen based on which one has a minimum volume which is closest to the required minimum volume
	// now this behaviour must be a lot more complicated in general since there are cases in which
	// it is permissible to use more than one tip type for the same request and in fact it would only be a case of
	// generating more or fewer transfers
	// this can be supported but it's a problem for the scheduler to optimize against preferences for number of operations,
	// time limits, space usage etc. as well as hard constraints on the above
	//

	for _, ins := range instrx {
		if ins.InstructionType() == lhdriver.LOD {
			ttype := ins.GetParameter("TIPTYPE").([]string)[0]
			ntips[ttype] += ins.GetParameter("MULTI").(int)
		}
	}

	h := make(map[string]int, 3)

	for tiptype, ntip := range ntips {
		// need to make sure the names match up here
		tx := strings.Split(tiptype, "_")
		actualtiptype := tx[0]
		h[actualtiptype] += ntip
	}

	for actualtiptype, ntip := range h {
		logger.Debug(fmt.Sprintln("TIPS OF TYPE ", actualtiptype, " USED: ", ntip))

		logger.Info(fmt.Sprintf("Block %s Tips of type %s used: %d", request.BlockID, actualtiptype, ntip))

		tbt := factory.GetTipByType(actualtiptype)
		ntbx := (ntip-1)/tbt.NTips + 1

		for i := 0; i < ntbx; i++ {
			tbt2 := factory.GetTipByType(actualtiptype)
			tip_boxes = append(tip_boxes, tbt2)
		}
	}

	(*request).Tips = tip_boxes

	// need to fix the tip situation in the properties structure
	lh.Properties.RemoveTipBoxes()
	for _, tb := range tip_boxes {
		lh.Properties.AddTipBox(tb)
	}

	return request
}
