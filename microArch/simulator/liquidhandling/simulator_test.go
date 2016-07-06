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
    "testing"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	. "github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/simulator"
)

type LayoutParams struct {
    Name    string
    Xpos    float64
    Ypos    float64
    Zpos    float64
}

type UnitParams struct {
    Value   float64
    Unit    string
}

type ChannelParams struct {
    Name            string
    Minvol          UnitParams
    Maxvol          UnitParams
    Minrate         UnitParams
    Maxrate         UnitParams
    multi           int
    Independent     bool
    Orientation     int
    Head            int
}

func makeLHChannelParameter(cp ChannelParams) *wtype.LHChannelParameter {
    return wtype.NewLHChannelParameter(cp.Name,
                                       wunit.NewVolume(cp.Minvol.Value, cp.Minvol.Unit),
                                       wunit.NewVolume(cp.Maxvol.Value, cp.Maxvol.Unit),
                                       wunit.NewFlowRate(cp.Minrate.Value, cp.Minrate.Unit),
                                       wunit.NewFlowRate(cp.Maxrate.Value, cp.Maxrate.Unit),
                                       cp.multi,
                                       cp.Independent,
                                       cp.Orientation,
                                       cp.Head)
}

type AdaptorParams struct {
    Name      string
    Mfg       string
    Channel   ChannelParams
}

func makeLHAdaptor(ap AdaptorParams) *wtype.LHAdaptor {
    return wtype.NewLHAdaptor(ap.Name,
                              ap.Mfg,
                              makeLHChannelParameter(ap.Channel))
}

type HeadParams struct {
    Name        string
    Mfg         string
    Channel     ChannelParams
    Adaptor     AdaptorParams
}

func makeLHHead(hp HeadParams) *wtype.LHHead {
    ret := wtype.NewLHHead(hp.Name, hp.Mfg, makeLHChannelParameter(hp.Channel))
    ret.Adaptor = makeLHAdaptor(hp.Adaptor)
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


/*
 *######################################### Test Data
 */

var valid_props = LHPropertiesParams{
    "Device Name",
    "Device Manufacturer",
    []LayoutParams{
        LayoutParams{"position1" ,   0.0,   0.0,   0.0},
        LayoutParams{"position2" , 100.0,   0.0,   0.0,},
        LayoutParams{"position3" , 200.0,   0.0,   0.0},
        LayoutParams{"position4" ,   0.0, 100.0,   0.0},
        LayoutParams{"position5" , 100.0, 100.0,   0.0},
        LayoutParams{"position6" , 200.0, 100.0,   0.0},
        LayoutParams{"position7" ,   0.0, 200.0,   0.0},
        LayoutParams{"position8" , 100.0, 200.0,   0.0},
        LayoutParams{"position9" , 200.0, 200.0,   0.0},
    },
    []HeadParams{
        HeadParams{
            "Head0 Name",
            "Head0 Manufacturer",
            ChannelParams{
                "Head0 ChannelParams",      //Name
                UnitParams{0.1, "ul"},      //min volume
                UnitParams{1.,  "ml"},      //max volume
                UnitParams{0.1, "ml/min"},  //min flowrate
                UnitParams{10., "ml/min",}, //max flowrate
                8,                          //multi
                false,                      //independent
                0,                          //orientation
                0,                          //head
            },
            AdaptorParams{
                "Head0 Adaptor",
                "Head0 Adaptor Manufacturer",
                ChannelParams{
                    "Head0 Adaptor ChannelParams",  //Name
                    UnitParams{0.1, "ul"},          //min volume
                    UnitParams{1.,  "ml"},          //max volume
                    UnitParams{0.1, "ml/min"},      //min flowrate
                    UnitParams{10., "ml/min",},     //max flowrate
                    8,                              //multi
                    false,                          //independent
                    0,                              //orientation
                    0,                              //head
                },
            },
        },
    },
    []string{"position1","position2",},             //Tip_preferences
    []string{"position3","position4",},             //Input_preferences
    []string{"position5","position6",}, //Output_preferences
    []string{"position7",},             //Tipwaste_preferences
    []string{"position8",},             //Wash_preferences
    []string{"position9",},             //Waste_preferences
}


/*
 *######################################### Testing Begins
 */


func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
    lhp := makeLHProperties(valid_props)
    vlh := NewVirtualLiquidHandler(lhp)

    errors, max_severity := vlh.GetErrors()
    if len(errors) > 0 {
        t.Error("Unexpected Error: %v", errors)
    } else if max_severity != simulator.SeverityNone {
        t.Error("Severty should be SeverityNone, instead got: %v", max_severity)
    }
}


