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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"errors"

	lhdriver "github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/factory"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

//  TASK: 	Determine number of tip boxes of each type
// INPUT: 	instructions
//OUTPUT: 	arrays of tip boxes
func (lh *Liquidhandler) Tip_box_setup(request *LHRequest) *LHRequest {
	tip_box_type := (*request).Tip_Type
	if tip_box_type == nil || tip_box_type.ID == "" {
		wutil.Error(errors.New("tip_box_setup: No tip_box type defined"))
	}
	tip_boxes := (*request).Tips
	if len(tip_boxes) == 0 {
		tip_boxes = make([]*wtype.LHTipbox, 0)
	}

	// the instructions are generated at this point so we just need to go through and count the tips used
	// of each type

	instrx := request.Instructions

	ntips := make(map[string]int)

	for _, ins := range instrx {
		if ins.InstructionType() == lhdriver.LOD {
			ttype := ins.GetParameter("TIPTYPE").([]string)[0]
			ntips[ttype] += ins.GetParameter("MULTI").(int)
		}
	}

	for tiptype, ntip := range ntips {
		// need to make sure the names match up here
		tbt := factory.GetTipByType(tiptype)
		ntbx := ntip/tbt.NTips + 1
		for i := 0; i < ntbx; i++ {
			tbt2 := factory.GetTipByType(tiptype)
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
