package liquidhandling

import (
	//"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/factory"
)

func GetLiquidHandlerForTest() *Liquidhandler {
	gilson := makeGilson()
	return Init(gilson)
}

func GetLHRequestForTest() *LHRequest {
	req := NewLHRequest()
	return req
}

func GetComponentForTest(name string, vol wunit.Volume) *wtype.LHComponent {
	c := factory.GetComponentByType(name)
	c.SetVolume(vol)
	return c
}

func TestNoInstructions(t *testing.T) {
	req := GetLHRequestForTest()
	lh := GetLiquidHandlerForTest()
	err := lh.MakeSolutions(req)

	if err.Error() != "9 (LH_ERR_OTHER) :  : Nil plan requested: no Mix Instructions present" {
		t.Fatal(fmt.Sprint("Unexpected error: ", err.Error()))
	}
}

func TestDeckSpace1(t *testing.T) {

	lh := GetLiquidHandlerForTest()

	for i := 0; i < len(lh.Properties.Tip_preferences); i++ {
		tb := factory.GetTipBoxByTip(lh.Properties.Tips[0])
		err := lh.Properties.AddTipBox(tb)
		if err != nil {
			t.Fatal(fmt.Sprintf("Should be able to fill deck with tipboxes, only managed %d", i+1))
		}
	}

	tb := factory.GetTipBoxByTip(lh.Properties.Tips[0])
	err := lh.Properties.AddTipBox(tb)
	if err.Error() != "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add tip box" {
		t.Fatal(fmt.Sprint("Expected error : 1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add tip box\n", " got: ", err.Error()))
	}
}

func TestDeckSpace2(t *testing.T) {
	lh := GetLiquidHandlerForTest()

	for i := 0; i < len(lh.Properties.Input_preferences); i++ {
		plate := factory.GetPlateByType("pcrplate_skirted")
		err := lh.Properties.AddPlate(lh.Properties.Input_preferences[i], plate)
		if err != nil {
			t.Fatal(fmt.Sprintf("position %s is full, should be empty", lh.Properties.Input_preferences[i]))
		}
	}

	plate := factory.GetPlateByType("pcrplate_skirted")
	err := lh.Properties.AddPlate(lh.Properties.Input_preferences[0], plate)

	if err.Error() != "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add plate to full position position_4" {
		t.Fatal(fmt.Sprint("1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add plate to full position position_4\n", " got: ", err.Error()))
	}

}

func TestNoTips(t *testing.T) {

}

func TestNotImplemented(t *testing.T) {

}

func TestDriv(t *testing.T) {

}
