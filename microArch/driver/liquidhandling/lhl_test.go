package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/factory"
	"testing"
)

func makeGilson() *LHProperties {
	// gilson pipetmax

	layout := make(map[string]wtype.Coordinates)
	i := 0
	x0 := 3.886
	y0 := 3.513
	z0 := -82.035
	xi := 149.86
	yi := 95.25
	xp := x0
	yp := y0
	zp := z0
	for y := 0; y < 3; y++ {
		xp = x0
		for x := 0; x < 3; x++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{xp, yp, zp}
			layout[posname] = crds
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := NewLHProperties(9, "Pipetmax", "Gilson", "discrete", "disposable", layout)
	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"}
	//lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_7"}

	//lhp.Tip_preferences = []string{"position_9", "position_6", "position_3", "position_5", "position_2"} //jmanart i cut it down to 5, as it was hardcoded in the liquidhandler getInputs call before

	// original preferences
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"}
	lhp.Output_preferences = []string{"position_8", "position_9", "position_6", "position_5", "position_3", "position_1"}

	// use these new preferences for gel loading: this is needed because outplate overlaps inplate otherwise so move inplate to position 5 rather than 4 (pos 4 deleted)
	//lhp.Input_preferences = []string{"position_5", "position_6", "position_9", "position_8", "position_3"}
	//lhp.Output_preferences = []string{"position_9", "position_8", "position_7", "position_6", "position_5", "position_3"}

	lhp.Wash_preferences = []string{"position_8"}
	lhp.Tipwaste_preferences = []string{"position_1", "position_7"}
	lhp.Waste_preferences = []string{"position_9"}
	//	lhp.Tip_preferences = []int{2, 3, 6, 9, 5, 8, 4, 7}
	//	lhp.Input_preferences = []int{24, 25, 26, 29, 28, 23}
	//	lhp.Output_preferences = []int{10, 11, 12, 13, 14, 15}
	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(250, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", "GilsonPipetmax", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor
	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(20, "ul")
	newminspd := wunit.NewFlowRate(0.1, "ml/min")
	newmaxspd := wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", "GilsonPipetmax", newminvol, newmaxvol, newminspd, newmaxspd, 8, false, wtype.LHVChannel, 1)
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	lhp.Heads = append(lhp.Heads, hvhead)
	lhp.Heads = append(lhp.Heads, lvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, hvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, lvhead)

	return lhp
}

func makeTestGilson() *LHProperties {
	params := makeGilson()

	params.AddTipWaste(factory.GetTipwasteByType("Gilsontipwaste"))
	params.AddTipBox(factory.GetTipboxByType("DL10 Tip Rack (PIPETMAX 8x20)"))
	params.AddTipBox(factory.GetTipboxByType("DF200 Tip Rack (PIPETMAX 8x200)"))

	return params
}

func getTestBlowout(robot *LHProperties) RobotInstruction {
	v := wunit.NewVolume(10.0, "ul")
	ch, _ := ChooseChannel(v, robot)
	bi := NewBlowInstruction()
	bi.Multi = 1
	bi.What = append(bi.What, "soup")
	bi.PltTo = append(bi.PltTo, "position_4")
	bi.WellTo = append(bi.WellTo, "A1")
	bi.Volume = append(bi.Volume, v)
	bi.TPlateType = append(bi.TPlateType, "pcrplate_skirted_riser40")
	bi.TVolume = append(bi.TVolume, wunit.ZeroVolume())
	bi.Prms = ch
	bi.Head = ch.Head
	return bi
}

func TestBlowWithTipChange(t *testing.T) {
	t.Skip()
	robot := makeTestGilson()
	bi := getTestBlowout(robot)
	pol, _ := GetLHPolicyForTest()

	rule := NewLHPolicyRule("TESTRULE1")
	rule.AddCategoryConditionOn("LIQUIDCLASS", "soup")
	pols := make(LHPolicy, 2)
	pols["POST_MIX"] = 5
	pols["POST_MIX_VOLUME"] = 100.0
	pol.AddRule(rule, pols)

	set := NewRobotInstructionSet(bi)

	ris, err := set.Generate(pol, robot)

	if err != nil {
		t.Fatal(err)
	}

	expectedIns := []int{MOV, DSP, MOV, ULD, MOV, LOD, MOV, MMX, MOV, BLO}

	if len(ris) != len(expectedIns) {
		t.Fatal(fmt.Sprintf("Error: Expected %d instructions, got %d", len(expectedIns), len(ris)))
	}

	for i, ins := range ris {
		if ins.InstructionType() != expectedIns[i] {
			t.Fatal(fmt.Sprintf("Error generating high mix volume blow: expected %s got %s", Robotinstructionnames[expectedIns[i]], Robotinstructionnames[ins.InstructionType()]))
		}
	}
}

func TestBlowNoTipChange(t *testing.T) {
	robot := makeTestGilson()
	bi := getTestBlowout(robot)
	pol, _ := GetLHPolicyForTest()

	rule := NewLHPolicyRule("TESTRULE1")
	rule.AddCategoryConditionOn("LIQUIDCLASS", "soup")
	pols := make(LHPolicy, 2)
	pols["POST_MIX"] = 5
	pols["POST_MIX_VOLUME"] = 10.0
	pol.AddRule(rule, pols)

	set := NewRobotInstructionSet(bi)

	ris, err := set.Generate(pol, robot)

	if err != nil {
		t.Fatal(err)
	}
	expectedIns := []int{MOV, DSP, MOV, MMX, MOV, BLO}

	if len(ris) != len(expectedIns) {
		t.Fatal(fmt.Sprintf("Error: Expected %d instructions, got %d", len(expectedIns), len(ris)))
	}

	for i, ins := range ris {
		if ins.InstructionType() != expectedIns[i] {
			t.Fatal(fmt.Sprintf("Error generating low mix volume blow: expected %s got %s", Robotinstructionnames[expectedIns[i]], Robotinstructionnames[ins.InstructionType()]))
		}
	}
}
