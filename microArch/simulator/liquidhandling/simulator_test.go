// /anthalib/simulator/liquidhandling/simulator_test.go: Part of the Antha language
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
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	. "github.com/antha-lang/antha/microArch/factory"
)

type LayoutParams struct {
    Name    string
    Xpos    float32
    Ypos    float32
    Zpos    float32
}

type ChannelParams struct {
    Name            string
    Minvol          string
    Maxvol          string
    Minrate         string
    Maxrate         string
    multi           int
    Independent     bool
    Orientation     int
    Head            int
}

func makeLHChannelParameter(cp ChannelParams) *wtype.LHChannelParameter {
    return wtype.NewLHChannelParameter(cp.Name,
                                       wtype.NewVolume(cp.Minvol),
                                       wtype.NewVolume(cp.Maxvol),
                                       wtype.NewFlowRate(cp.Minrate),
                                       wtype.NewFlowRate(cp.Maxrate),
                                       cp.multi,
                                       cp.Independent,
                                       cp.Orientation,
                                       cp.Head)
}

type AdaptorParams struct {
    ChannelParams
    Name string
    Mfg  string
}

func makeLHAdaptor(ap AdaptorParams) *wtype.LHAdaptor {
    return wtype.NewLHAdaptor(ap.Name,
                              ap.Mfg,
                              makeLHChannelParameter(AdaptorParams.ChannelParams))
}

type HeadParams struct {
    ChannelParams
    Name        string
    Mfg         string
    Adaptor     AdaptorParams
}

func makeLHHead(hp HeadParams) *wtype.LHHead {
    ret := wtype.NewLHHead(hp.Name, hp.Mfg, makeLHChannelParameter(HeadParams.ChannelParams))
    ret.Adaptor = makeLHAdaptor(cp.Adaptor)
    return ret
}

type LHPropertiesParams struct {
    Name                    string
    Mfg                     string
    Layouts                 []LayoutParams
    Heads                   []HeadParams
    Tip_preferences         []string
	Input_preferences       []string
	Output_preferences      []string
	Tipwaste_preferences    []string
	Wash_preferences        []string
	Waste_preferences       []string
}
    

func AddAllTips(lhp *liquidhandling.LHProperties) *liquidhandling.LHProperties {
	tips := GetTipList()
	for _, tt := range tips {
		tb := GetTipByType(tt)
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}
	return lhp
}


func makeLHProperties(p LHPropertiesParams) *liquidhandling.LHProperties {


	layout := make(map[string]wtype.Coordinates)
    for _, lp := range p.Layouts {
        layout[lp.Name] = wtype.Coordinates{lp.Xpos, lp.Ypos, lp.Zpos}
    }
        
	lhp := liquidhandling.NewLHProperties(len(layout), p.Name, p.Mfg, "discrete", "disposable", layout)

	AddAllTips(lhp)

    lhp.Heads = make([]*wtype.LHHead, 0)
    for _, hp := range p.Heads {
        lhp.Heads = append(lhp.Heads, makeLHHead(hp))
    }

    return lhp
}
