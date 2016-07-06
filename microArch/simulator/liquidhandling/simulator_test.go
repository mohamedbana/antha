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
    "fmt"
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

    lhp.Tip_preferences = p.Tip_preferences         
	lhp.Input_preferences = p.Input_preferences
	lhp.Output_preferences = p.Output_preferences
	lhp.Tipwaste_preferences = p.Tipwaste_preferences
	lhp.Wash_preferences = p.Wash_preferences
	lhp.Waste_preferences = p.Waste_preferences

    return lhp
}


/*
 *######################################### Test Data
 */

func get_valid_props() *LHPropertiesParams {
    valid_props := LHPropertiesParams{
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

    return &valid_props
}

type LHPropertyTest struct {
    Properties      *LHPropertiesParams
    ErrorStrings    []string
}

func get_unknown_locations() []LHPropertyTest {
    ret := make([]LHPropertyTest, 0)

    lhp := get_valid_props()
    lhp.Tip_preferences = append(lhp.Tip_preferences, "undefined_tip_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_tip_pref\" referenced in tip preferences"}})

    lhp = get_valid_props()
    lhp.Input_preferences = append(lhp.Tip_preferences, "undefined_input_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_input_pref\" referenced in input preferences"}})

    lhp = get_valid_props()
    lhp.Output_preferences = append(lhp.Tip_preferences, "undefined_output_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_output_pref\" referenced in output preferences"}})

    lhp = get_valid_props()
    lhp.Tipwaste_preferences = append(lhp.Tip_preferences, "undefined_tipwaste_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_tipwaste_pref\" referenced in tipwaste preferences"}})

    lhp = get_valid_props()
    lhp.Wash_preferences = append(lhp.Tip_preferences, "undefined_wash_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_wash_pref\" referenced in wash preferences"}})

    lhp = get_valid_props()
    lhp.Waste_preferences = append(lhp.Tip_preferences, "undefined_waste_pref")
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) Undefined location \"undefined_waste_pref\" referenced in waste preferences"}})

    return ret
}

func get_missing_prefs() []LHPropertyTest {
    ret := make([]LHPropertyTest, 0)

    lhp := get_valid_props()
    lhp.Tip_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No tip preferences specified"}})

    lhp = get_valid_props()
    lhp.Input_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No input preferences specified"}})

    lhp = get_valid_props()
    lhp.Output_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No output preferences specified"}})

    lhp = get_valid_props()
    lhp.Tipwaste_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No tipwaste preferences specified"}})

    lhp = get_valid_props()
    lhp.Wash_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No wash preferences specified"}})

    lhp = get_valid_props()
    lhp.Waste_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{lhp, []string{"(warn) No waste preferences specified"}})

    return ret
}

/*
 * ######################################## utils
 */

//test that the worst reported error severity is the worst
func test_worst(t *testing.T, errors []*simulator.SimulationError, worst simulator.ErrorSeverity) {
    s := simulator.SeverityNone
    for _, err := range errors {
        if err.Severity() > s {
            s = err.Severity()
        }
    }

    if s != worst {
        t.Error("Expected maximum severity %v, actual maximum severity %v", worst, s)
    }
}


func run_prop_test(t *testing.T, pt LHPropertyTest) {
    vlh := NewVirtualLiquidHandler(makeLHProperties(*pt.Properties))
    errors, worst := vlh.GetErrors()
    test_worst(t, errors, worst)

    actual_errors := make([]string, 0)
    for _,err := range errors {
        actual_errors = append(actual_errors, err.Error())
    }

    compare_errors(t, pt.ErrorStrings, actual_errors)
}

//return subset of a not in b
func get_not_in(a, b []string) []string {
    for _,vb := range b {
        for i,va := range a {
            if va == vb {
                a = append(a[:i], a[i+1:]...)
                break
            }
        }
    }
    return a
}



func compare_errors(t *testing.T, expected []string, actual []string) {
    // maybe sort alphabetically?
    
    missing := get_not_in(expected, actual)
    extra := get_not_in(actual, expected)

    if len(missing) != 0 {
        t.Error(fmt.Sprintf("missing expected error(s): %v", missing))
    }
    if len(extra) != 0 {
        t.Error(fmt.Sprintf("got extra error(s): %v", extra))
    }
}

/*
 *######################################### Testing Begins
 */


func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
    lhp := makeLHProperties(*get_valid_props())
    vlh := NewVirtualLiquidHandler(lhp)

    errors, max_severity := vlh.GetErrors()
    if len(errors) > 0 {
        t.Error("Unexpected Error: %v", errors)
    } else if max_severity != simulator.SeverityNone {
        t.Error("Severty should be SeverityNone, instead got: %v", max_severity)
    }
}

func TestNewVirtualLiquidHandler_UnknownLocation(t *testing.T) {
    tests := get_unknown_locations()
    for _,test := range tests {
        run_prop_test(t, test)
    }
}

func TestNewVirtualLiquidHandler_MissingPrefs(t *testing.T) {
    tests := get_missing_prefs()
    for _,test := range tests {
        run_prop_test(t, test)
    }
}


