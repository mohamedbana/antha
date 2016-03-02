// /anthalib/driver/liquidhandling/funcs.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// random helper functions

func MinMinVol(channels []*wtype.LHChannelParameter) wunit.Volume {
	mmv := wunit.NewVolume(9999999.0, "ul")

	for _, c := range channels {
		if c.Minvol.LessThan(mmv) {
			mmv = c.Minvol
		}
	}

	return mmv
}

/*
func ChooseChannel(vol *wunit.Volume, channels []*wtype.LHChannelParameter) *wtype.LHChannelParameter {
	// greedy greedy
	var r *wtype.LHChannelParameter = nil

	d := 999999.0
	v := vol.RawValue()

	for _, c := range channels {
		min := c.Minvol.ConvertTo(vol.Unit())

		// can't do less than the stated minimum
		// choose the channel which has the minimum
		// nearest the required volume; change this in future
		if v > min {
			df := v - min
			if df < d {
				d = v - min
				r = c
			}
		}
	}

	return r
}
*/
