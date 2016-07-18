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
    "strings"
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


func makeLHProperties(p *LHPropertiesParams) *liquidhandling.LHProperties {


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

type ShapeParams struct {
    name            string 
    lengthunit      string 
    h               float64
    w               float64
    d               float64
}

func makeShape(p *ShapeParams) *wtype.Shape {
    return wtype.NewShape(p.name, p.lengthunit, p.h, p.w, p.d)
}

type LHWellParams struct {
    platetype       string 
    plateid         string 
    crds            string
    vunit           string 
    vol             float64
    rvol            float64 
    shape           ShapeParams 
    bott            int 
    xdim            float64
    ydim            float64 
    zdim            float64 
    bottomh         float64 
    dunit           string
}

func makeLHWell(p *LHWellParams) *wtype.LHWell {
    return wtype.NewLHWell(p.platetype, 
                           p.plateid, 
                           p.crds, 
                           p.vunit, 
                           p.vol, 
                           p.rvol,
                           makeShape(&p.shape),
                           p.bott,
                           p.xdim,
                           p.ydim,
                           p.zdim,
                           p.bottomh,
                           p.dunit)
}

type LHPlateParams struct {
    platetype       string 
    mfr             string 
    nrows           int 
    ncols           int 
    height          float64
    hunit           string
    welltype        LHWellParams
    wellXOffset     float64 
    wellYOffset     float64 
    wellXStart      float64
    wellYStart      float64
    wellZStart      float64
}

func makeLHPlate(p *LHPlateParams) *wtype.LHPlate {
    return wtype.NewLHPlate(p.platetype,
                            p.mfr,
                            p.nrows,
                            p.ncols,
                            p.height,
                            p.hunit,
                            makeLHWell(&p.welltype),
                            p.wellXOffset,
                            p.wellYOffset,
                            p.wellXStart,
                            p.wellYStart,
                            p.wellZStart)
}

type LHTipParams struct {
    mfr         string
    ttype       string 
    minvol      float64
    maxvol      float64 
    volunit     string
}

func makeLHTip(p *LHTipParams) *wtype.LHTip {
    return wtype.NewLHTip(p.mfr,
                         p.ttype,
                         p.minvol,
                         p.maxvol,
                         p.volunit)
}

type LHTipboxParams struct {
    nrows           int 
    ncols           int 
    height          float64 
    manufacturer    string
    boxtype         string 
    tiptype         LHTipParams
    well            LHWellParams 
    tipxoffset      float64
    tipyoffset      float64
    tipxstart       float64
    tipystart       float64
    tipzstart       float64
}

func makeLHTipbox(p *LHTipboxParams) *wtype.LHTipbox {
    return wtype.NewLHTipbox(p.nrows,
                             p.ncols,
                             p.height,
                             p.manufacturer,
                             p.boxtype,
                             makeLHTip(&p.tiptype),
                             makeLHWell(&p.well),
                             p.tipxoffset,
                             p.tipyoffset,
                             p.tipystart,
                             p.tipxstart,
                             p.tipzstart)
}

type LHTipwasteParams struct {
    capacity        int 
    typ             string
    mfr             string 
    height          float64 
    w               LHWellParams 
    wellxstart      float64
    wellystart      float64
    wellzstart      float64
}

func makeLHTipWaste(p *LHTipwasteParams) *wtype.LHTipwaste {
    return wtype.NewLHTipwaste(p.capacity,
                               p.typ,
                               p.mfr,
                               p.height,
                               makeLHWell(&p.w),
                               p.wellxstart,
                               p.wellystart,
                               p.wellzstart)
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
        []string{"position5","position6",},             //Output_preferences
        []string{"position7",},                         //Tipwaste_preferences
        []string{"position8",},                         //Wash_preferences
        []string{"position9",},                         //Waste_preferences
    }

    return &valid_props
}

func get_valid_vlh() *VirtualLiquidHandler {
    vlh := NewVirtualLiquidHandler(makeLHProperties(get_valid_props()))
    vlh.Initialize()
    return vlh
}

type LHPropertyTest struct {
    desc            string
    Properties      *LHPropertiesParams
    ErrorStrings    []string
}

func (pt *LHPropertyTest) apply(t *testing.T) {
    vlh := NewVirtualLiquidHandler(makeLHProperties(pt.Properties))
    errors := vlh.GetErrors()
    test_worst(t, errors, vlh.GetErrorSeverity())
    compare_errors(t, pt.desc, pt.ErrorStrings, errors)
}

func get_unknown_locations() []LHPropertyTest {
    ret := make([]LHPropertyTest, 0)

    lhp := get_valid_props()
    lhp.Tip_preferences = append(lhp.Tip_preferences, "undefined_tip_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Tip_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tip_pref\" referenced in tip preferences"}})

    lhp = get_valid_props()
    lhp.Input_preferences = append(lhp.Tip_preferences, "undefined_input_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Input_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_input_pref\" referenced in input preferences"}})

    lhp = get_valid_props()
    lhp.Output_preferences = append(lhp.Tip_preferences, "undefined_output_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Output_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_output_pref\" referenced in output preferences"}})

    lhp = get_valid_props()
    lhp.Tipwaste_preferences = append(lhp.Tip_preferences, "undefined_tipwaste_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Tipwaste_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tipwaste_pref\" referenced in tipwaste preferences"}})

    lhp = get_valid_props()
    lhp.Wash_preferences = append(lhp.Tip_preferences, "undefined_wash_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Wash_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_wash_pref\" referenced in wash preferences"}})

    lhp = get_valid_props()
    lhp.Waste_preferences = append(lhp.Tip_preferences, "undefined_waste_pref")
    ret = append(ret, LHPropertyTest{"passing undefined Waste_preference", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_waste_pref\" referenced in waste preferences"}})

    return ret
}

func get_missing_prefs() []LHPropertyTest {
    ret := make([]LHPropertyTest, 0)

    lhp := get_valid_props()
    lhp.Tip_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing Tip_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No tip preferences specified"}})

    lhp = get_valid_props()
    lhp.Input_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing Input_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No input preferences specified"}})

    lhp = get_valid_props()
    lhp.Output_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing Output_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No output preferences specified"}})

    lhp = get_valid_props()
    lhp.Tipwaste_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing TipWaste_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No tipwaste preferences specified"}})

    lhp = get_valid_props()
    lhp.Wash_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing Wash_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No wash preferences specified"}})

    lhp = get_valid_props()
    lhp.Waste_preferences = make([]string, 0)
    ret = append(ret, LHPropertyTest{"passing missing Waste_preferences", lhp, 
        []string{"(warn) NewVirtualLiquidHandler: No waste preferences specified"}})

    return ret
}

func get_lhplate() *wtype.LHPlate {
    params := LHPlateParams {
        "test_plate_type",  // platetype       string 
        "test_plate_mfr",   // mfr             string 
        8,                  // nrows           int 
        12,                 // ncols           int 
        25.7,               // height          float64
        "mm",               // hunit           string
        LHWellParams{           // welltype
            "test_welltype",    // platetype       string 
            "test_wellid",      // plateid         string 
            "",                 // crds            string
            "ul",               // vunit           string 
            200,                // vol             float64
            5,                  // rvol            float64 
            ShapeParams{            // shape           ShapeParams struct {
               "test_shape",        // name            string 
               "mm",                // lengthunit      string 
               5.5,                 // h               float64
               5.5,                 // w               float64
               20.4,                // d               float64
            },
           wtype.LHWBV,         // bott            int 
           5.5,                 // xdim            float64
           5.5,                 // ydim            float64 
           20.4,                // zdim            float64 
           1.4,                 // bottomh         float64 
           "mm",                // dunit           string
        },
        9.,        // wellXOffset     float64 
        9.,        // wellYOffset     float64 
        0.,        // wellXStart      float64
        0.,        // wellYStart      float64
        18.5,      // wellZStart      float64
    }

    return makeLHPlate(&params)
}

func get_lhtipbox() *wtype.LHTipbox {
    params := LHTipboxParams{
        8,                      //nrows           int 
        12,                     //ncols           int 
        60.13,                  //height          float64 
        "test Tipbox mfg",      //manufacturer    string
        "tipbox",     //boxtype         string 
        LHTipParams {           //tiptype
            "test_tip mfg",         //mfr         string
            "test_tip type",        //ttype       string 
            50,                     //minvol      float64
            1000,                   //maxvol      float64 
            "ul",                   //volunit     string
        },
        LHWellParams{           // well
            "test_welltype",        // platetype       string 
            "test_wellid",          // plateid         string 
            "",                     // crds            string
            "ul",                   // vunit           string 
            1000,                   // vol             float64
            50,                     // rvol            float64 
            ShapeParams{            // shape           ShapeParams struct {
               "test_shape",            // name            string 
               "mm",                    // lengthunit      string 
               7.3,                     // h               float64
               7.3,                     // w               float64
               51.2,                    // d               float64
            },
           wtype.LHWBV,             // bott            int 
           7.3,                     // xdim            float64
           7.3,                     // ydim            float64 
           51.2,                    // zdim            float64 
           0.0,                     // bottomh         float64 
           "mm",                    // dunit           string
        },
        9.,                     //tipxoffset      float64
        9.,                     //tipyoffset      float64
        0.,                     //tipxstart       float64
        0.,                     //tipystart       float64
        0.,                     //tipzstart       float64
    }

    return makeLHTipbox(&params)
}

func get_lhtipwaste() *wtype.LHTipwaste {
    params := LHTipwasteParams {
        700,                    //capacity        int 
        "tipwaste",             //typ             string
        "testTipwaste mfr",     //mfr             string 
        92.0,                   //height          float64 
        LHWellParams{           // w               LHWellParams
            "test_tipwaste_well",   // platetype       string 
            "test_wellid",          // plateid         string 
            "",                     // crds            string
            "ul",                   // vunit           string 
            800000.0,                   // vol             float64
            800000.0,                     // rvol            float64 
            ShapeParams{            // shape           ShapeParams struct {
               "test_tipbox",           // name            string 
               "mm",                    // lengthunit      string 
               123.0,                   // h               float64
               80.0,                    // w               float64
               92.0,                    // d               float64
            },
           wtype.LHWBV,             // bott            int 
           123.0,                   // xdim            float64
           80.0,                    // ydim            float64 
           92.0,                    // zdim            float64 
           0.0,                     // bottomh         float64 
           "mm",                    // dunit           string
        },
        85.5,               //wellxstart      float64
        45.5,               //wellystart      float64
        0.0,                //wellzstart      float64
    }
    return makeLHTipWaste(&params)
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



func compare_errors(t *testing.T, desc string, expected []string, actual []*simulator.SimulationError) {
    string_errors := make([]string, 0)
    for _,err := range actual {
        string_errors = append(string_errors, err.Error())
    }
    // maybe sort alphabetically?
    
    missing := get_not_in(expected, string_errors)
    extra := get_not_in(string_errors, expected)

    errs := []string{}
    for _,s := range missing {
        errs = append(errs, fmt.Sprintf("--\"%v\"", s))
    }
    for _,s := range extra {
        errs = append(errs, fmt.Sprintf("++\"%v\"", s))
    }
    if len(missing) > 0 || len(extra) > 0 {
        t.Errorf("Errors didn't match in test \"%v\":\n%s", 
            desc, strings.Join(errs, "\n"))
    }
}

/*
 *######################################### Testing Begins
 */


func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
    lhp := makeLHProperties(get_valid_props())
    vlh := NewVirtualLiquidHandler(lhp)

    compare_errors(t, "NewVLH with valid props", []string{}, vlh.GetErrors())
}

func TestNewVirtualLiquidHandler_UnknownLocation(t *testing.T) {
    tests := get_unknown_locations()
    for _,test := range tests {
        test.apply(t)
    }
}

func TestNewVirtualLiquidHandler_MissingPrefs(t *testing.T) {
    tests := get_missing_prefs()
    for _,test := range tests {
        test.apply(t)
    }
}

func TestVLH_AddPlateTo_Valid(t *testing.T) {
    vlh := get_valid_vlh()
    //try adding a test LHplate to preferred inputs and outputs
    for i,loc := range []string{"position3","position4","position5","position6",} {
        vlh.AddPlateTo(loc, get_lhplate(), fmt.Sprintf("LHPlate_%v", i))
    }
    //try adding a test LHTipBox to preferred Tips
    for i,loc := range []string{"position1","position2",} {
        vlh.AddPlateTo(loc, get_lhtipbox(), fmt.Sprintf("LHTipbox_%v", i))
    }
    //try adding a test LHTipWaste to preferred TipWaste
    for i,loc := range []string{"position7"} {
        vlh.AddPlateTo(loc, get_lhtipwaste(), fmt.Sprintf("LHTipwaste_%v", i))
    }

    compare_errors(t, "adding valid plates", []string{}, vlh.GetErrors())
}

func TestVLH_AddPlateTo_NotPlateType(t *testing.T) {
    vlh := get_valid_vlh()
    //try adding something that's the wrong type
    vlh.AddPlateTo("position2", "my plate's gone stringy", "not_a_plate")

    compare_errors(t, "adding string plate", 
        []string{"(err) AddPlateTo: Cannot add plate \"not_a_plate\" of type string to location \"position2\""},
        vlh.GetErrors())
}

func TestVLH_AddPlateTo_locationFull(t *testing.T) {
    vlh := get_valid_vlh()
    
    //add a plate
    vlh.AddPlateTo("position1", get_lhplate(), "p0")
    //try to add another plate in the same location
    vlh.AddPlateTo("position1", get_lhplate(), "p1")

    compare_errors(t, "adding plate to full location", 
        []string{"(err) AddPlateTo: Cannot add plate \"p1\" to location \"position1\" which is already occupied by plate \"p0\""},
        vlh.GetErrors())
}

// ########################################################################################################################
// ########################################################## Tip Loading/Unloading
// ########################################################################################################################


func get_tip_test_vlh() *VirtualLiquidHandler {
    props := LHPropertiesParams{
        "Device Name",
        "Device Manufacturer",
        []LayoutParams{
            LayoutParams{"tip_loc" ,   0.0,   0.0,   0.0},
            LayoutParams{"tipwaste_loc" , 100.0,   0.0,   0.0,},
            LayoutParams{"input_loc" , 200.0,   0.0,   0.0},
            LayoutParams{"output_loc" ,   0.0, 100.0,   0.0},
            LayoutParams{"wash_loc" , 100.0, 100.0,   0.0},
            LayoutParams{"waste_loc" , 200.0, 100.0,   0.0},
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
        []string{"tip_loc",},               //Tip_preferences
        []string{"input_loc",},             //Input_preferences
        []string{"output_loc",},            //Output_preferences
        []string{"tipwaste_loc",},          //Tipwaste_preferences
        []string{"wash_loc",},              //Wash_preferences
        []string{"waste_loc",},             //Waste_preferences
    }

    vlh := NewVirtualLiquidHandler(makeLHProperties(&props))
    vlh.Initialize()
    vlh.AddPlateTo("tip_loc", get_lhtipbox(), "tipbox1")
    vlh.AddPlateTo("tipwaste_loc", get_lhtipwaste(), "tipwaste")
    return vlh
}

func removeTips(wells []string) *func(*VirtualLiquidHandler) {
    ret := func(vlh *VirtualLiquidHandler) {
        tipbox := vlh.properties.PlateLookup[vlh.properties.PosLookup["tip_loc"]].(*wtype.LHTipbox)
        for _,well := range wells {
            wc := wtype.MakeWellCoords(well)
            tipbox.Tips[wc.X][wc.Y] = nil
        }
    }
    return &ret
}

func preloadTips(channels []int, head int) *func(*VirtualLiquidHandler) {
    ret := func(vlh *VirtualLiquidHandler) {
        p := vlh.properties
        adaptor := p.Heads[head].Adaptor
        tipbox := p.PlateLookup[p.PosLookup["tip_loc"]].(*wtype.LHTipbox)

        for _,ch := range channels {
            adaptor.AddTip(ch, tipbox.Tiptype.Dup());
        }
    }
    return &ret
}

type LoadTipsParams struct {
    channels        []int
    head            int
    multi           int
    platetype       []string
    position        []string
    well            []string
}

func (s *LoadTipsParams) apply(vlh *VirtualLiquidHandler) {
    vlh.LoadTips(s.channels,
                 s.head,
                 s.multi,
                 s.platetype,
                 s.position,
                 s.well)
}

type LoadTipsTest struct {
    desc                string
    params              LoadTipsParams
    setup               *func(*VirtualLiquidHandler)
    expected_errors     []string
    missing_tips        []TipLoc
    loaded_tips         []int
}

func (test *LoadTipsTest) run(t *testing.T, vlh *VirtualLiquidHandler) {
    if test.setup != nil {
        (*test.setup)(vlh)
    }
    test.params.apply(vlh)
    
    //check that there are no errors/warnings in the vlh
    if test.expected_errors == nil {
        test.expected_errors = []string{}
    }
    compare_errors(t, test.desc, test.expected_errors, vlh.GetErrors())

    //don't check missing tips or loaded tips if we were expecting errors
    if len(test.expected_errors) > 0 {
        return
    }

    //get the tipbox
    props := vlh.properties
    tipbox := props.PlateLookup[props.PosLookup["tip_loc"]].(*wtype.LHTipbox)
    
    //check that tips are missing iff they're in missing_tips
    tip_errors := []string{}
    for i := 0; i < tipbox.Ncols; i++ {
        for j := 0; j < tipbox.Nrows; j++ {
            if tipbox.Tips[i][j] == nil {
                //must be in missing tips
                if !contains(test.missing_tips,j,i) {
                    tip_errors = append(tip_errors, fmt.Sprintf("--(%v,%v)", i,j))
                }
            } else {
                //mustn't be in missing tips
                if contains(test.missing_tips,j,i) {
                    tip_errors = append(tip_errors, fmt.Sprintf("++(%v,%v)", i,j))
                }
            }
        }
    }
    if len(tip_errors) > 0 {
        t.Errorf("In test \"%v\": Final tipbox layout incorrect: (missing/extra --/++)\n%s", 
            test.desc, strings.Join(tip_errors, "\n"))
    }

    //get the adaptor
    adaptor := props.Heads[0].Adaptor

    //test the tips were loaded in the right place
    verifyAdaptorTips(t, test.desc, test.loaded_tips, adaptor)
}

func verifyAdaptorTips(t *testing.T, test_desc string, expected_tips []int, adaptor *wtype.LHAdaptor) {
    tipExpected := make(map[int]bool)
    for _,ch := range expected_tips {
        tipExpected[ch] = true
    }
    tip_errors := []string{}
    for ch := 0; ch < adaptor.Params.Multi; ch++ {
        if tipExpected[ch] && !adaptor.IsTipLoaded(ch) {
            tip_errors = append(tip_errors, fmt.Sprintf("Missing tip at channel %v", ch))
        } else if !tipExpected[ch] && adaptor.IsTipLoaded(ch) {
            tip_errors = append(tip_errors, fmt.Sprintf("Extra tip at channel %v", ch))
        }
    }
    if len(tip_errors) > 0 { 
        t.Errorf("In test \"%v\": Unexpected adaptor tips:\n%s",
            test_desc, strings.Join(tip_errors, "\n"))
    }
    
}

type UnloadTipsParams struct {
    channels        []int
    head            int
    multi           int
    platetype       []string
    position        []string
    well            []string
}

func (s *UnloadTipsParams) apply(vlh *VirtualLiquidHandler) {
    vlh.UnloadTips(s.channels,
                   s.head,
                   s.multi,
                   s.platetype,
                   s.position,
                   s.well)
}

type UnloadTipsTest struct {
    testName        string
    params          UnloadTipsParams
    setup           []*func(*VirtualLiquidHandler)
    expected_errors []string
    remaining_tips  []int
    tips_in_waste   int
}

func (self *UnloadTipsTest) run(t *testing.T, vlh *VirtualLiquidHandler) {
    if self.setup != nil {
        for _,f := range self.setup {
            (*f)(vlh)
        }
    }
    self.params.apply(vlh)

    if self.expected_errors == nil {
        self.expected_errors = []string{}
    }
    compare_errors(t, self.testName, self.expected_errors, vlh.GetErrors())

    //Don't worry about tip locations if we were expecting errors anyway
    if len(self.expected_errors) > 0 {
        return
    }

    adaptor := vlh.properties.Heads[0].Adaptor
    verifyAdaptorTips(t, self.testName, self.remaining_tips, adaptor)

    if self.remaining_tips == nil {
        self.remaining_tips = []int{}
    }
    tipwaste := vlh.properties.PlateLookup[vlh.properties.PosLookup["tipwaste_loc"]].(*wtype.LHTipwaste)
    if self.tips_in_waste != tipwaste.Contents {
        t.Errorf("Incorrect tipwaste contents after \"%s\", expected %v tips, got %v tips",
                 self.testName,
                 self.tips_in_waste,
                 tipwaste.Contents)
    }

}


type TipLoc struct {
    row     int
    col     int
}

func contains(s []TipLoc, row, col int) bool {
    for _, l := range s {
        if l.row == row && l.col == col {
            return true
        }
    }
    return false
}

func Test_LoadTips(t *testing.T) {

    tests := []LoadTipsTest{
        LoadTipsTest{
            "correctly loading a single tip",
            LoadTipsParams{
                []int{0},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                        //setup
            nil,                        //expected_error
            []TipLoc{TipLoc{7,11}},     //missing_tips
            []int{0},                   //loaded_tips
        },
        LoadTipsTest{
            "correctly loading eight tips at once",
            LoadTipsParams{
                []int{0,1,2,3,4,5,6,7},     //channels
                0,                          //head
                8,                          //multi
                []string{                   //platetype
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                },
                []string{                   //position
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                },
                []string{                   //well
                    "A12",
                    "B12",
                    "C12",
                    "D12",
                    "E12",
                    "F12",
                    "G12",
                    "H12",
                },
            },
            nil,                        //setup
            nil,                        //expected_error
            []TipLoc{                   //missing_tips
                TipLoc{0,11},
                TipLoc{1,11},
                TipLoc{2,11},
                TipLoc{3,11},
                TipLoc{4,11},
                TipLoc{5,11},
                TipLoc{6,11},
                TipLoc{7,11},
            },
            []int{0,1,2,3,4,5,6,7,},    //loaded_tips
        },
        LoadTipsTest{
            "correctly loading three tips at once",
            LoadTipsParams{
                []int{0,1,2,},     //channels
                0,                          //head
                3,                          //multi
                []string{                   //platetype
                    "tipbox",
                    "tipbox",
                    "tipbox",
                },
                []string{                   //position
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                },
                []string{                   //well,
                    "F1",
                    "G1",
                    "H1",
                },
            },
            nil,                        //setup
            nil,                        //expected_error
            []TipLoc{                   //missing_tips
                TipLoc{5,0},
                TipLoc{6,0},
                TipLoc{7,0},
            },
            []int{0,1,2,},    //loaded_tips
        },
        LoadTipsTest{
            "loading a single tip above a missing tip",
            LoadTipsParams{
                []int{0},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"G12"},            //well
            },
            removeTips([]string{"H12"}),//setup
            nil,                        //expected_error
            []TipLoc{TipLoc{7,11},TipLoc{6,11}},     //missing_tips
            []int{0},                   //loaded_tips
        },
        // and now the error tests
        LoadTipsTest{
            "invalid channel value",
            LoadTipsParams{
                []int{8},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load tip to channel 8 of 8-channel adaptor"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "invalid channel value",
            LoadTipsParams{
                []int{-1},                  //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load tip to channel -1 of 8-channel adaptor"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "too many channels",
            LoadTipsParams{
                []int{0,1,2,3,4,5,6,7,7},   //channels
                0,                          //head
                9,                          //multi
                []string{                   //platetype
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                    "tipbox",
                },
                []string{                   //position
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                    "tip_loc",
                },
                []string{                   //well
                    "A12",
                    "B12",
                    "C12",
                    "D12",
                    "E12",
                    "F12",
                    "G12",
                    "H12",
                    "H12",
                },
            },
            nil,                            //setup
            []string{"(err) LoadTips: Channel7 appears more than once"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "invalid head",
            LoadTipsParams{
                []int{0},                  //channels
                1,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Request for invalid Head 1"},
            nil,
            nil,
        },
        LoadTipsTest{
            "invalid head",
            LoadTipsParams{
                []int{0},                  //channels
                -1,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Request for invalid Head -1"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "mismatching multi",
            LoadTipsParams{
                []int{0},                  //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: channels, platetype, position, well should be of length multi=2"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "mismatching location",
            LoadTipsParams{
                []int{0,1},                  //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox",
                         "tipbox"},         //platetype
                []string{"tip_loc", 
                         "tipwaste_loc"},        //position
                []string{"G12", "H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load tips from multiple locations"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "mismatching platetype",
            LoadTipsParams{
                []int{0,1},                  //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox",
                         "tipwaste"},         //platetype
                []string{"tip_loc", 
                         "tip_loc"},        //position
                []string{"G12", "H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: platetype should be equal"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "wrong platetype",
            LoadTipsParams{
                []int{0,1},                  //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox",
                         "tipbox"},         //platetype
                []string{"tipwaste_loc", 
                         "tipwaste_loc"},        //position
                []string{"G12", "H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load tips from location \"tipwaste_loc\", no tipbox found"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "mismatching multi",
            LoadTipsParams{
                []int{0,1},                  //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox",
                         "tipbox"},         //platetype
                []string{"tip_loc", 
                         "tip_loc"},        //position
                []string{"H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: well should be of length multi=2"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "invalid well",
            LoadTipsParams{
                []int{0},                  //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H13"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Request for well H13, but tipbox size is [12x8]"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "invalid well",
            LoadTipsParams{
                []int{0},                  //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"not_a_well"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Couldn't parse well \"not_a_well\""},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "loading collision",
            LoadTipsParams{
                []int{0},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"G12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load G12->channel0 due to tip at H12 (Head0 is not independent)"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
        LoadTipsTest{
            "loading collision",
            LoadTipsParams{
                []int{0,7},                   //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox","tipbox"},         //platetype
                []string{"tip_loc","tip_loc"},        //position
                []string{"A12", "H12"},            //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load A12->channel0, H12->channel7 due to tips at B12,C12,D12,E12,F12,G12 (Head0 is not independent)"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        }, 
        LoadTipsTest{
            "non-contiguous wells",
            LoadTipsParams{
                []int{0,1},                   //channels
                0,                          //head
                2,                          //multi
                []string{"tipbox","tipbox"},         //platetype
                []string{"tip_loc","tip_loc"},        //position
                []string{"F12", "H12"},     //well
            },
            nil,                            //setup
            []string{"(err) LoadTips: Cannot load F12->channel0, H12->channel1, tip spacing doesn't match channel spacing"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        }, 
        LoadTipsTest{
            "missing tip",
            LoadTipsParams{
                []int{0},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            removeTips([]string{"H12"}),     //setup
            []string{"(err) LoadTips: Cannot load H12->channel0 as H12 is empty"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        }, 
        LoadTipsTest{
            "tip already loaded",
            LoadTipsParams{
                []int{0},                   //channels
                0,                          //head
                1,                          //multi
                []string{"tipbox"},         //platetype
                []string{"tip_loc"},        //position
                []string{"H12"},            //well
            },
            preloadTips([]int{0}, 0),     //setup
            []string{"(err) LoadTips: Cannot load tips while adaptor already contains 1 tip"},
            nil,                            //missing_tips
            nil,                            //loaded_tips
        },
    }

    for _, test := range tests {
        vlh := get_tip_test_vlh()
        test.run(t, vlh)
    }

}

func Test_UnloadTips(t *testing.T) {

    tests := []UnloadTipsTest {
        UnloadTipsTest{
            "unload a tip",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0},                   //channels        []int
                0,                          //head            int
                1,                          //multi           int
                []string{"tipwaste"},       //platetype       []string
                []string{"tipwaste_loc"},   //position        []string
                []string{"A1"},             //well            []string
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0}, 0),
            },
            nil,                            //expected_errors []string
            nil,                            //remaining_tips  []int
            1,                              //tips_in_waste   int
        },
        UnloadTipsTest{
            "unload 8 tips",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,1,2,3,4,5,6,7},     //channels        []int
                0,                          //head            int
                8,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste",
                         "tipwaste",
                         "tipwaste",
                         "tipwaste",
                         "tipwaste",
                         "tipwaste",
                         "tipwaste"},
                []string{"tipwaste_loc",    //position        []string
                         "tipwaste_loc",
                         "tipwaste_loc",
                         "tipwaste_loc",
                         "tipwaste_loc",
                         "tipwaste_loc",
                         "tipwaste_loc",
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1",
                         "A1",
                         "A1",
                         "A1",
                         "A1",
                         "A1",
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,1,2,3,4,5,6,7}, 0),
            },
            nil,                            //expected_errors []string
            nil,                            //remaining_tips  []int
            8,                              //tips_in_waste   int
        },
        UnloadTipsTest{
            "can only unload all tips",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste"},
                []string{"tipwaste_loc",    //position        []string
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,1,2,3,4,5,6,7}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload tips from Head0(channels 0,7) due to other tips on the adaptor (independent is false)"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "multi is wrong",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                3,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste"},
                []string{"tipwaste_loc",    //position        []string
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,7}, 0),
            },
            []string{"(err) UnloadTips: channels, platetype, position, well should be of length multi=3"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "wrong location",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipbox",        //platetype       []string
                         "tipbox"},
                []string{"tip_loc",         //position        []string
                         "tip_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,7}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload tips at location \"tip_loc\", no tipwaste found"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },/*
        UnloadTipsTest{
            "wrong well",       //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste"},
                []string{"tipwaste_loc",    //position        []string
                         "tipwaste_loc"},
                []string{"B1",              //well            []string
                         "B1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,7}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload tips as plate at \"tipwaste_loc\" has no well \"B1\""},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },*/
        UnloadTipsTest{
            "wrong platetype",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipbox",        //platetype       []string
                         "tipbox"},
                []string{"tipwaste_loc",         //position        []string
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,7}, 0),
            },
            []string{"(err) UnloadTips: Requested plate type \"tipbox\" but plate at tipwaste_loc is of type \"tipwaste\""},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "mixed location",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste"},
                []string{"tipwaste_loc",         //position        []string
                         "tip_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,7}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload tips to multiple locations"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "unload too many",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7},                 //channels        []int
                0,                          //head            int
                2,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste"},
                []string{"tipwaste_loc",         //position        []string
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload 2 tips as adaptor only has one tip loaded"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "unload too many",     //testName        string
            UnloadTipsParams{   //params          UnloadTipsParams
                []int{0,7,4},                 //channels        []int
                0,                          //head            int
                3,                          //multi           int
                []string{"tipwaste",        //platetype       []string
                         "tipwaste",
                         "tipwaste"},
                []string{"tipwaste_loc",         //position        []string
                         "tipwaste_loc",
                         "tipwaste_loc"},
                []string{"A1",              //well            []string
                         "A1",
                         "A1"},
            },
            []*func(*VirtualLiquidHandler){ //setup           []*func(*VirtualLiquidHandler)
                preloadTips([]int{0,1}, 0),
            },
            []string{"(err) UnloadTips: Cannot unload 3 tips as adaptor only has 2 tips loaded"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
    }

    for _,test := range tests {
        vlh := get_tip_test_vlh()
        test.run(t, vlh)
    }
}
