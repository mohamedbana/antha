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

package liquidhandling_test

import (
    "testing"
    "fmt"
    "strings"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
    "github.com/antha-lang/antha/microArch/simulator"
    lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
)

//
// Code for specifying a VLH
//

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
	tips := factory.GetTipList()
	for _, tt := range tips {
		tb := factory.GetTipByType(tt)
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
    ret := []string{}
    for _,va := range a {
        c := false
        for _,vb := range b {
            if va == vb {
                c = true
            }
        }
        if !c {
            ret = append(ret, va)
        }
    }
    return ret
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
 * ####################################### Default Types
 */

func default_lhplate() *wtype.LHPlate {
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

func default_lhtipbox() *wtype.LHTipbox {
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

func default_lhtipwaste() *wtype.LHTipwaste {
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

func default_lhproperties() *liquidhandling.LHProperties {
    valid_props := LHPropertiesParams{
        "Device Name",
        "Device Manufacturer",
        []LayoutParams{
            LayoutParams{"tipbox_1" ,     0.0,   0.0,   0.0},
            LayoutParams{"tipbox_2" ,   200.0,   0.0,   0.0},
            LayoutParams{"input_1" ,    400.0,   0.0,   0.0},
            LayoutParams{"input_2" ,      0.0, 200.0,   0.0},
            LayoutParams{"output_1" ,   200.0, 200.0,   0.0},
            LayoutParams{"output_2" ,   400.0, 200.0,   0.0},
            LayoutParams{"tipwaste" ,     0.0, 400.0,   0.0},
            LayoutParams{"wash" ,       200.0, 400.0,   0.0},
            LayoutParams{"waste" ,      400.0, 400.0,   0.0},
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
        []string{"tipbox_1","tipbox_2",},             //Tip_preferences
        []string{"input_1", "input_2",},             //Input_preferences
        []string{"output_1","output_2",},             //Output_preferences
        []string{"tipwaste",},                         //Tipwaste_preferences
        []string{"wash",},                         //Wash_preferences
        []string{"waste",},                         //Waste_preferences
    }

    return makeLHProperties(&valid_props)
}

func independent_lhproperties() *liquidhandling.LHProperties {
    valid_props := LHPropertiesParams{
        "Device Name",
        "Device Manufacturer",
        []LayoutParams{
            LayoutParams{"tipbox_1" ,     0.0,   0.0,   0.0},
            LayoutParams{"tipbox_2" ,   200.0,   0.0,   0.0},
            LayoutParams{"input_1" ,    400.0,   0.0,   0.0},
            LayoutParams{"input_2" ,      0.0, 200.0,   0.0},
            LayoutParams{"output_1" ,   200.0, 200.0,   0.0},
            LayoutParams{"output_2" ,   400.0, 200.0,   0.0},
            LayoutParams{"tipwaste" ,     0.0, 400.0,   0.0},
            LayoutParams{"wash" ,       200.0, 400.0,   0.0},
            LayoutParams{"waste" ,      400.0, 400.0,   0.0},
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
                    true,                       //independent
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
                        true,                           //independent
                        0,                              //orientation
                        0,                              //head
                    },
                },
            },
        },
        []string{"tipbox_1","tipbox_2",},             //Tip_preferences
        []string{"input_1", "input_2",},             //Input_preferences
        []string{"output_1","output_2",},             //Output_preferences
        []string{"tipwaste",},                         //Tipwaste_preferences
        []string{"wash",},                         //Wash_preferences
        []string{"waste",},                         //Waste_preferences
    }

    return makeLHProperties(&valid_props)
}

func default_vlh() *lh.VirtualLiquidHandler {
    vlh := lh.NewVirtualLiquidHandler(default_lhproperties())
    return vlh
}

/*
 * ######################################## InstructionParams
 */

type TestRobotInstruction interface {
    Apply(*lh.VirtualLiquidHandler)
}

//Initialize
type Initialize struct {}

func (self *Initialize) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.Initialize()
}

//Finalize
type Finalize struct {}

func (self *Finalize) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.Finalize()
}

//AddPlateTo
type AddPlateTo struct {
    position        string
    plate           interface{}
    name            string
}

func (self *AddPlateTo) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.AddPlateTo(self.position, self.plate, self.name)
}

//LoadTips
type LoadTips struct {
    channels        []int
    head            int
    multi           int
    platetype       []string
    position        []string
    well            []string
}

func (self *LoadTips) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.LoadTips(self.channels, self.head, self.multi, self.platetype, self.position, self.well)
}

//UnloadTips
type UnloadTips struct {
    channels        []int
    head            int
    multi           int
    platetype       []string
    position        []string
    well            []string
}

func (self *UnloadTips) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.UnloadTips(self.channels, self.head, self.multi, self.platetype, self.position, self.well)
}

//Move
type Move struct {
    deckposition        []string
    wellcoords          []string
    reference           []int
    offsetX             []float64
    offsetY             []float64
    offsetZ             []float64
    plate_type          []string
    head                int
}

func (self *Move) Apply(vlh *lh.VirtualLiquidHandler) {
    vlh.Move(self.deckposition, self.wellcoords, self.reference, self.offsetX, self.offsetY, self.offsetZ, self.plate_type, self.head)
}

/*
 * ######################################## Setup
 */

type SetupFn func(*lh.VirtualLiquidHandler)

func removeTipboxTips(tipbox_loc string, wells []string) *SetupFn {
    var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
        tipbox := vlh.GetPlateAt(tipbox_loc).(*wtype.LHTipbox)
        for _,well := range wells {
            wc := wtype.MakeWellCoords(well)
            tipbox.Tips[wc.X][wc.Y] = nil
        }
    }
    return &ret
}

func preloadAdaptorTips(head int, tipbox_loc string, channels []int) *SetupFn {
    var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
        adaptor := vlh.GetAdaptorState(head)
        tipbox := vlh.GetPlateAt(tipbox_loc).(*wtype.LHTipbox)

        for _,ch := range channels {
            adaptor.GetChannel(ch).SetTip(tipbox.Tiptype.Dup());
        }
    }
    return &ret
}

func fillTipwaste(tipwaste_loc string, count int) *SetupFn {
    var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
        tipwaste := vlh.GetPlateAt(tipwaste_loc).(*wtype.LHTipwaste)
        tipwaste.Contents += count
    }
    return &ret
}

/*
 * ######################################## Assertions (about the final state)
 */

type AssertionFn func(string, *testing.T, *lh.VirtualLiquidHandler)

//tipboxAssertion assert that the tipbox has tips missing in the given locations only
func tipboxAssertion(tipbox_loc string, missing_tips []string) *AssertionFn {
    var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
        mmissing_tips := make(map[string]bool)
        for _,tl := range missing_tips {
            mmissing_tips[tl] = true
        }

        tipbox := vlh.GetPlateAt(tipbox_loc).(*wtype.LHTipbox)
        errors := []string{}
        for y := 0; y < tipbox.Nrows; y++ {
            for x:= 0; x < tipbox.Ncols; x++ {
                wc := wtype.WellCoords{x,y}
                wcs:= wc.FormatA1()
                if hta, mt := tipbox.HasTipAt(wc), mmissing_tips[wcs]; hta && mt {
                    errors = append(errors, fmt.Sprintf("Unexpected tip missing at %s", wcs))
                } else if !hta && !mt {
                    errors = append(errors, fmt.Sprintf("Unexpected tip present at %s", wcs))
                }
            }
        }
        if len(errors) > 0 {
            t.Errorf("TipboxAssertion failed in test \"%s\", tipbox at \"%s\":\n%s", name, tipbox_loc, strings.Join(errors, "\n"))
        }
    }
    return &ret
}

//adaptorAssertion assert that the adaptor has tips in the given positions
func adaptorAssertion(head int, tip_channels []int) *AssertionFn {
    var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
        mtips := make(map[int]bool)
        for _,tl := range tip_channels {
            mtips[tl] = true
        }

        adaptor := vlh.GetAdaptorState(head)
        errors := []string{}
        for ch := 0; ch < adaptor.GetChannelCount(); ch++ {
            if itl, et := adaptor.GetChannel(ch).HasTip(), mtips[ch]; itl && !et {
                errors = append(errors, fmt.Sprintf("Unexpected tip on channel %v", ch))
            } else if !itl && et {
                errors = append(errors, fmt.Sprintf("Expected tip on channel %v", ch))
            }
        }
        if len(errors) > 0 {
            t.Errorf("AdaptorAssertion failed in test \"%s\", Head%v:\n%s", name, head, strings.Join(errors, "\n"))
        }
    }
    return &ret
}

//tipwasteAssertion assert the number of tips which should be in the tipwaste
func tipwasteAssertion(tipwaste_loc string, expected_contents int) *AssertionFn {
    var ret AssertionFn = func(name string, t *testing.T, vlh *lh.VirtualLiquidHandler) {
        tipwaste := vlh.GetPlateAt(tipwaste_loc).(*wtype.LHTipwaste)
        
        if tipwaste.Contents != expected_contents {
            t.Errorf("TipwasteAssertion failed in test \"%s\" at location %s: expected %v tips, got %v", 
                name, tipwaste_loc, expected_contents, tipwaste.Contents)
        }
    }
    return &ret
}

/*
 * ######################################## SimulatorTest
 */

type SimulatorTest struct {
    Name            string
    Props           *liquidhandling.LHProperties
    Setup           []*SetupFn
    Instructions    []TestRobotInstruction
    ExpectedErrors  []string
    Assertions      []*AssertionFn
}

func (self *SimulatorTest) run(t *testing.T) {
    
    if self.Props == nil {
        self.Props = default_lhproperties()
    }
    vlh := lh.NewVirtualLiquidHandler(self.Props)

    //do setup
    if self.Setup != nil {
        for _,setup_fn := range self.Setup {
            (*setup_fn)(vlh)
        }
    }

    //run the instructions
    if self.Instructions != nil {
        for _,inst := range self.Instructions {
            inst.Apply(vlh)
        }
    }

    //check errors
    if self.ExpectedErrors != nil {
        compare_errors(t, self.Name, self.ExpectedErrors, vlh.GetErrors())
    } else {
        compare_errors(t, self.Name, []string{}, vlh.GetErrors())
    }

    //check assertions
    if self.Assertions != nil {
        for _,a := range self.Assertions {
            (*a)(self.Name, t, vlh)
        }
    }
}

