// /anthalib/liquidhandling/rationalise_outputs.go: Part of the Antha language
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
	"reflect"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/logger"
)

func Rationalise_Outputs(lhr *LHRequest, lhp *liquidhandling.LHProperties, outputs ...*wtype.LHComponent) *LHRequest {

	op := make(map[string]*wtype.LHPlate)

	for i := 0; i < len(outputs); i++ {
		// need to check that the outputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic
		n := reflect.TypeOf(outputs[i].Container())
		if !((n.Kind() == reflect.Ptr && n.Elem().Name() == "wtype.LHWell") || n.Name() == "liquidhandling.wtype.LHWell") {
			logger.Fatal("Rationalise_outputs: cannot use non-microplate formats in liquid handler... yet!")
			panic("Rationalise_outputs: cannot use non-microplate formats in liquid handler... yet!")
		}

		op[outputs[i].LContainer.Plate.ID] = outputs[i].LContainer.Plate

	}

	lhr.Output_plates = op

	return lhr
}
