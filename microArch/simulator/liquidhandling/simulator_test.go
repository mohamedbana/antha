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
    lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
)


func TestUnknownLocations(t *testing.T) {
    tests := make([]SimulatorTest, 0)

    lhp := default_lhproperties()
    lhp.Tip_preferences = append(lhp.Tip_preferences, "undefined_tip_pref")
    tests = append(tests, SimulatorTest{"Undefined Tip_preference", lhp, nil, nil,
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tip_pref\" referenced in tip preferences"},
        nil})

    lhp = default_lhproperties()
    lhp.Input_preferences = append(lhp.Tip_preferences, "undefined_input_pref")
    tests = append(tests, SimulatorTest{"passing undefined Input_preference", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_input_pref\" referenced in input preferences"},
        nil})

    lhp = default_lhproperties()
    lhp.Output_preferences = append(lhp.Tip_preferences, "undefined_output_pref")
    tests = append(tests, SimulatorTest{"passing undefined Output_preference", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_output_pref\" referenced in output preferences"},
        nil})

    lhp = default_lhproperties()
    lhp.Tipwaste_preferences = append(lhp.Tip_preferences, "undefined_tipwaste_pref")
    tests = append(tests, SimulatorTest{"passing undefined Tipwaste_preference", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tipwaste_pref\" referenced in tipwaste preferences"},
        nil})

    lhp = default_lhproperties()
    lhp.Wash_preferences = append(lhp.Tip_preferences, "undefined_wash_pref")
    tests = append(tests, SimulatorTest{"passing undefined Wash_preference", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_wash_pref\" referenced in wash preferences"},
        nil})

    lhp = default_lhproperties()
    lhp.Waste_preferences = append(lhp.Tip_preferences, "undefined_waste_pref")
    tests = append(tests, SimulatorTest{"passing undefined Waste_preference", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_waste_pref\" referenced in waste preferences"},
        nil})

    for _,test := range tests {
        test.run(t)
    }
}

func TestMissingPrefs(t *testing.T) {
    tests := make([]SimulatorTest, 0)

    lhp := default_lhproperties()
    lhp.Tip_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing Tip_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No tip preferences specified"},
        nil})

    lhp = default_lhproperties()
    lhp.Input_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing Input_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No input preferences specified"},
        nil})

    lhp = default_lhproperties()
    lhp.Output_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing Output_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No output preferences specified"},
        nil})

    lhp = default_lhproperties()
    lhp.Tipwaste_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing TipWaste_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No tipwaste preferences specified"},
        nil})

    lhp = default_lhproperties()
    lhp.Wash_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing Wash_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No wash preferences specified"},
        nil})

    lhp = default_lhproperties()
    lhp.Waste_preferences = make([]string, 0)
    tests = append(tests, SimulatorTest{"passing missing Waste_preferences", lhp, nil, nil, 
        []string{"(warn) NewVirtualLiquidHandler: No waste preferences specified"},
        nil})

    for _,test := range tests {
        test.run(t)
    }
}

func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
    test := SimulatorTest{"Create Valid VLH", nil, nil, nil, nil, nil}
    test.run(t)
}

func TestVLH_AddPlateTo(t *testing.T) {
    tests := []SimulatorTest{
        SimulatorTest{
            "OK",       //name
            nil,        //default params
            nil,        //no setup
            []TestRobotInstruction{
                &Initialize{},
                &AddPlateTo{"tipbox_1", default_lhtipbox(), "tipbox1"},
                &AddPlateTo{"tipbox_2", default_lhtipbox(), "tipbox2"},
                &AddPlateTo{"input_1", default_lhplate(), "input1"},
                &AddPlateTo{"input_2", default_lhplate(), "input2"},
                &AddPlateTo{"output_1", default_lhplate(), "output1"},
                &AddPlateTo{"output_2", default_lhplate(), "output2"},
                &AddPlateTo{"tipwaste", default_lhtipwaste(), "tipwaste"},
            },
            nil,        //no errors
            nil,        //no assertions
        },
        SimulatorTest{
            "non plate type",       //name
            nil,                    //default params
            nil,                    //no setup
            []TestRobotInstruction{
                &Initialize{},
                &AddPlateTo{"tipbox_1", "my plate's gone stringy", "not_a_plate"},
            },
            []string{"(err) AddPlateTo: Cannot add plate \"not_a_plate\" of type string to location \"tipbox_1\""},
            nil,        //no assertions
        },
        SimulatorTest{
            "location full",        //name
            nil,                    //default params
            nil,                    //no setup
            []TestRobotInstruction{
                &Initialize{},
                &AddPlateTo{"tipbox_1", default_lhtipbox(), "p0"},
                &AddPlateTo{"tipbox_1", default_lhtipbox(), "p1"},
            },
            []string{"(err) AddPlateTo: Cannot add plate \"p1\" to location \"tipbox_1\" which is already occupied by plate \"p0\""},
            nil,        //no assertions
        },
//        SimulatorTest{   -- We'll probably want a test along these lines at some point, but Preferences aren't very strict at the moment
//            "wrong plate type",     //name
//            nil,                    //default params
//            nil,                    //no setup
//            []TestRobotInstruction{
//                &Initialize{},
//                &AddPlateTo{"tipbox_1", default_lhplate(), "tipbox"},
//            },
//            []string{"(warn) AddPlateTo: Added type Plate to location \"tipbox_1\", when preferences requested Tipbox"},
//            nil,        //no assertions
//        },
    }

    for _,test := range tests {
        test.run(t)
    }
}




// ########################################################################################################################
// ########################################################## Tip Loading/Unloading
// ########################################################################################################################

func tipTestLayout() *SetupFn {
    var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
        vlh.Initialize()
        vlh.AddPlateTo("tipbox_1",  default_lhtipbox(), "tipbox1")
        vlh.AddPlateTo("tipbox_2",  default_lhtipbox(), "tipbox2")
        vlh.AddPlateTo("tipwaste", default_lhtipwaste(), "tipwaste")
    }
    return &ret
}

func TestLoadTips(t *testing.T) {
    tests := []SimulatorTest{
        SimulatorTest{
            "OK - single tip",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H11"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"H11"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0}),
                tipwasteAssertion("tipwaste", 0),
            },
        },
        SimulatorTest{
            "OK - single tip (alt)",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{7},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"A1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A1"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{7}),
                tipwasteAssertion("tipwaste", 0),
            },
        },
        SimulatorTest{
            "OK - single tip above space",
            nil,
            []*SetupFn{
                tipTestLayout(),
                removeTipboxTips("tipbox_1", []string{"H11"}),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},   //location
                    []string{"G11"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"H11","G11"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0}),
                tipwasteAssertion("tipwaste", 0),
            },
        }, 
        SimulatorTest{
            "OK - single tip above space (alt)",
            nil,
            []*SetupFn{
                tipTestLayout(),
                removeTipboxTips("tipbox_1", []string{"A1"}),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{7},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"B1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A1","B1"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{7}),
                tipwasteAssertion("tipwaste", 0),
            },
        },
        SimulatorTest{
            "OK - 3 tips",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1,2}, //channels
                    0,                      //head
                    3,                      //multi
                    []string{"tipbox","tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1","tipbox_1"},   //location
                    []string{"F11","G11","H11"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"F11","G11","H11"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0,1,2}),
                tipwasteAssertion("tipwaste", 0),
            },
        },
        SimulatorTest{
            "OK - 3 tips (alt)",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{5,6,7}, //channels
                    0,                      //head
                    3,                      //multi
                    []string{"tipbox","tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1","tipbox_1"},   //location
                    []string{"A1","B1","C1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A1","B1","C1"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{5,6,7}),
                tipwasteAssertion("tipwaste", 0),
            },
        },
        SimulatorTest{
            "OK - 8 tips",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1,2,3,4,5,6,7}, //channels
                    0,                      //head
                    8,                      //multi
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"},   //location
                    []string{"A11","B11","C11","D11","E11","F11","G11","H11"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A11","B11","C11","D11","E11","F11","G11","H11"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0,1,2,3,4,5,6,7}),
                tipwasteAssertion("tipwaste", 0),
            },
        }, 
    }

    for _,test := range tests {
        test.run(t)
    }
}


/*
func get_tip_test_vlh(independent bool) *VirtualLiquidHandler {
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
                    independent,                //independent
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
                        independent,                    //independent
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
    independent         bool
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
    independent     bool
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            "independently load three non-contiguous tips at once",
            true,                       //independent
            LoadTipsParams{
                []int{0,4,6,},     //channels
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
                    "A1",
                    "E1",
                    "G1",
                },
            },
            nil,                        //setup
            nil,                        //expected_error
            []TipLoc{                   //missing_tips
                TipLoc{0,0},
                TipLoc{4,0},
                TipLoc{6,0},
            },
            []int{0,4,6,},    //loaded_tips
        },
        LoadTipsTest{
            "loading a single tip above a missing tip",
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
            false,                      //independent
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
        vlh := get_tip_test_vlh(test.independent)
        test.run(t, vlh)
    }

}

func Test_UnloadTips(t *testing.T) {

    tests := []UnloadTipsTest {
        UnloadTipsTest{
            "unload a tip",     //testName        string
            false,              //independent
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
            false,              //independent
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
            "independently unload 2 of 8 tips",     //testName        string
            true,              //independent
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
            nil,
            []int{1,2,3,4,5,6},         //remaining_tips  []int
            2,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "can only unload all tips",     //testName        string
            false,              //independent
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
            false,              //independent
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
            false,              //independent
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
        },
        UnloadTipsTest{
            "wrong well",       //testName        string
            false,              //independent
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
        },
        UnloadTipsTest{
            "wrong platetype",     //testName        string
            false,              //independent
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
            false,              //independent
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
            false,              //independent
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
            []string{"(err) UnloadTips: Cannot unload tip from channel 7 as no tip is loaded there"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "unload too many",     //testName        string
            false,              //independent
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
            []string{"(err) UnloadTips: Cannot unload tips from channels 7,4 as no tips are loaded there"},
            nil,                        //remaining_tips  []int
            0,                          //tips_in_waste   int
        },
        UnloadTipsTest{
            "waste is full",     //testName        string
            false,              //independent
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
                fillWaste(700),
            },
            []string{"(err) UnloadTips: Tipwaste at \"tipwaste_loc\" is overfull"},//expected_errors []string
            nil,                            //remaining_tips  []int
            1,                              //tips_in_waste   int
        },
    }

    for _,test := range tests {
        vlh := get_tip_test_vlh(test.independent)
        test.run(t, vlh)
    }
}

*/
