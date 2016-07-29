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
/*
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
*/
func tipTestLayout() *SetupFn {
    var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
        vlh.Initialize()
        vlh.AddPlateTo("tipbox_1",  default_lhtipbox(), "tipbox1")
        vlh.AddPlateTo("tipbox_2",  default_lhtipbox(), "tipbox2")
        vlh.AddPlateTo("tipwaste", default_lhtipwaste(), "tipwaste")
    }
    return &ret
}
/*
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
                    []string{"H12"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"H12"}),
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
                removeTipboxTips("tipbox_1", []string{"H12"}),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},   //location
                    []string{"G12"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"H12","G12"}),
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
                    []string{"F12","G12","H12"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"F12","G12","H12"}),
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
            "OK - 3 tips (independent)",
            independent_lhproperties(),
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,4,7}, //channels
                    0,                      //head
                    3,                      //multi
                    []string{"tipbox","tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1","tipbox_1"},   //location
                    []string{"A1","E1","H1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A1","E1","H1"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0,4,7}),
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
                    []string{"A12","B12","C12","D12","E12","F12","G12","H12"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{"A12","B12","C12","D12","E12","F12","G12","H12"}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{0,1,2,3,4,5,6,7}),
                tipwasteAssertion("tipwaste", 0),
            },
        }, 
        SimulatorTest{
            "invalid channel 8",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{8},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load tip to channel 8 of 8-channel adaptor",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "invalid channel -1",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{-1},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load tip to channel -1 of 8-channel adaptor",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "duplicate channels",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1,2,3,4,5,6,3}, //channels
                    0,                      //head
                    8,                      //multi
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"},   //location
                    []string{"A12","B12","C12","D12","E12","F12","G12","H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Channel3 appears more than once",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "unknown head",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    1,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Request for unknown head 1",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "unknown head -1",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    -1,                     //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Request for unknown head -1",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "mismatching multi",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                     //head
                    2,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: channels, platetype, position, well should be of length multi=2",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "multiple locations",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1},             //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_2"},    //location
                    []string{"G12","H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load tips from multiple locations",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "multiple platetypes",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1},             //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipbox","tipwaste"},     //platetype
                    []string{"tipbox_1","tipbox_1"},    //location
                    []string{"G12","H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: platetype should be equal",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "wrong plate type",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipwaste"},     //tipbox
                    []string{"tipwaste"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: No tipbox found at location \"tipwaste\"",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "well out of range",
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
                    []string{"tipbox_1"},   //location
                    []string{"H13"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Request for well H13, but tipbox size is [12x8]",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "invalid well",
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
                    []string{"tipbox_1"},   //location
                    []string{"not_a_well"}, //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Couldn't parse well \"not_a_well\"",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "Loading collision",
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
                    []string{"tipbox_1"},   //location
                    []string{"G12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load G12->channel0 due to tip at H12 (Head0 is not independent)",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "multi loading collision",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,7},             //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1"},   //location
                    []string{"A12","H12"},  //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load A12->channel0, H12->channel7 due to tips at B12,C12,D12,E12,F12,G12 (Head0 is not independent)",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "non-contiguous wells",
            independent_lhproperties(),
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0,1},             //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipbox","tipbox"},     //tipbox
                    []string{"tipbox_1","tipbox_1"},   //location
                    []string{"F12","H12"},  //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load F12->channel0, H12->channel1, tip spacing doesn't match channel spacing",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "tip missing",
            nil,
            []*SetupFn{
                tipTestLayout(),
                removeTipboxTips("tipbox_1", []string{"H12"}),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load H12->channel0 as H12 is empty",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "tip already loaded",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &LoadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},    //location
                    []string{"H12"},        //well
                },
            },
            []string{       //errors
                "(err) LoadTips: Cannot load tips while adaptor already contains 1 tip",
            },
            nil,            //assertions
        },
    }

    for _,test := range tests {
        test.run(t)
    }
}




func Test_UnloadTips(t *testing.T) {

    tests := []SimulatorTest{
        SimulatorTest{
            "OK - single tip",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipwaste"},     //tipbox
                    []string{"tipwaste"},   //location
                    []string{"A1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{}),
                tipwasteAssertion("tipwaste", 1),
            },
        },
        SimulatorTest{
            "OK - 8 tips",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0,1,2,3,4,5,6,7}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0,1,2,3,4,5,6,7},               //channels
                    0,                      //head
                    8,                      //multi
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste"},     //tipbox
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste"},   //location
                    []string{"A1","A1","A1","A1","A1","A1","A1","A1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{}),
                tipwasteAssertion("tipwaste", 8),
            },
        },
        SimulatorTest{
            "OK - independent tips",
            independent_lhproperties(),
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0,1,2,3,4,5,6,7}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0,2,4,6},               //channels
                    0,                      //head
                    4,                      //multi
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste"},     //tipbox
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste"},   //location
                    []string{"A1","A1","A1","A1"},        //well
                },
            },
            nil,            //errors
            []*AssertionFn{ //assertions
                tipboxAssertion("tipbox_1", []string{}),
                tipboxAssertion("tipbox_2", []string{}),
                adaptorAssertion(0, []int{1,3,5,7}),
                tipwasteAssertion("tipwaste", 4),
            },
        },
        SimulatorTest{
            "can only unload all tips",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0,1,2,3,4,5,6,7}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0,2,4,6},               //channels
                    0,                      //head
                    4,                      //multi
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste"},     //tipbox
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste"},   //location
                    []string{"A1","A1","A1","A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: Cannot unload tips from Head0(channels 0,2,4,6) due to other tips on the adaptor (independent is false)",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "wrong multi",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0},               //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipwaste"},     //tipbox
                    []string{"tipwaste"},   //location
                    []string{"A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: channels, platetype, position, well should be of length multi=2",
            },
            nil,  
        },
        SimulatorTest{
            "wrong location",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipbox"},     //tipbox
                    []string{"tipbox_1"},   //location
                    []string{"A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: No tipwaste found at location \"tipbox_1\"",
            },
            nil,  
        },
        SimulatorTest{
            "wrong well",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipwaste"},     //tipbox
                    []string{"tipwaste"},   //location
                    []string{"B1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: Tipwaste at \"tipwaste\" has no well B1",
            },
            nil,  
        },
        SimulatorTest{
            "multiple locations",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0,1}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0,1},               //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipwaste","tipwaste"},     //tipbox
                    []string{"tipwaste","tipwaste2"},   //location
                    []string{"A1","A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: Cannot unload tips to multiple locations",
            },
            nil,  
        },
        SimulatorTest{
            "multiple locations",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0,7},               //channels
                    0,                      //head
                    2,                      //multi
                    []string{"tipwaste","tipwaste"},     //tipbox
                    []string{"tipwaste","tipwaste"},   //location
                    []string{"A1","A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: Cannot unload tip from channel 7 as no tip is loaded there",
            },
            nil,  
        },
        SimulatorTest{
            "multiple locations",
            nil,
            []*SetupFn{
                tipTestLayout(),
                preloadAdaptorTips(0, "tipbox_1", []int{0}),
                fillTipwaste("tipwaste", 700),
            },
            []TestRobotInstruction{
                &UnloadTips{
                    []int{0},               //channels
                    0,                      //head
                    1,                      //multi
                    []string{"tipwaste"},   //tipbox
                    []string{"tipwaste"},   //location
                    []string{"A1"},        //well
                },
            },
            []string{            //errors
                "(err) UnloadTips: Tipwaste at \"tipwaste\" is overfull",
            },
            nil,  
        },

    }

    for _,test := range tests {
        test.run(t)
    }
}


func Test_Move(t *testing.T) {

    tests := []SimulatorTest{
        SimulatorTest{
            "OK",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste"}, //deckposition
                    []string{"A1","A1","A1","A1","A1","A1","A1","A1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"","","","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"","","","A1","B1","C1","D1","E1"}, //wellcoords
                    []int{0,0,0,1,1,1,1,1}, //reference (first 3 should be ignored)
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"","","","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            nil,            //errors
            nil,            //assertions
        },
        SimulatorTest{
            "unknown location",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox7","tipbox7","tipbox7","tipbox7","tipbox7","tipbox7","tipbox7","tipbox7"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste","tipwaste"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Unknown location \"tipbox7\"",
                "(err) Move: Location \"tipbox_1\" has no plate of type \"tipwaste\"",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "unknown head",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    1, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    -1, //head
                },
            },
            []string{       //errors
                "(err) Move: Unknown head 1",
                "(err) Move: Unknown head -1",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "invalid wellcoords",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"B1","C1","D1","E1","F1","G1","H1","I1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"B1","C1","D1","E1","F1","G1","H1","not_a_well"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Request for well I1 in tipbox at \"tipbox_1\", size [8x12]",
                "(err) Move: Couldn't parse well \"not_a_well\"",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "Invalid reference",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{-1,-1,-1,-1,-1,-1,-1,-1,}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{2,2,2,2,2,2,2,2,}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{0,0,0,0,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Invalid reference -1",
                "(err) Move: Invalid reference 2",
                "(err) Move: References must be equal as adaptor is not independent",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "offsets differ",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,3.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,1.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,1.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Offsets cannot differ between channels when independent is false",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "layout mismatch",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B2","C1","D2","E1","F2","G1","H2"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Requested wells do not match adaptor layout",
            },
            nil,            //assertions
        },
        SimulatorTest{
            "crashes",
            nil,
            []*SetupFn{
                tipTestLayout(),
            },
            []TestRobotInstruction{
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{1,1,1,1,1,1,1,1}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{-50.,-50.,-50.,0.,-50.,-50.,-50.,-50.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{0,0,0,0,0,0,0,0}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{-0.1,-0.1,-0.1,-0.1,-0.1,-0.1,-0.1,-0.1,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
                &Move{
                    []string{"tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1","tipbox_1"}, //deckposition
                    []string{"A1","B1","C1","D1","E1","F1","G1","H1"}, //wellcoords
                    []int{0,0,0,0,0,0,0,0}, //reference
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetX
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetY
                    []float64{0.,0.,0.,0.,0.,0.,0.,0.,}, //offsetZ
                    []string{"tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox","tipbox"}, //plate_type
                    0, //head
                },
            },
            []string{       //errors
                "(err) Move: Crash predicted, tip intersects with well base",
                "(err) Move: Crash predicted, tip intersects with well base",
                "(err) Move: Crash predicted, tip intersects with well side",
            },
            nil,            //assertions
        },
    }

    for _,test := range tests {
        test.run(t)
    }
}

*/
