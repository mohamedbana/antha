package liquidhandling

import (
	//"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"fmt"
	"testing"
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

	for i := 0; i < lh.Properties.Nposns; i++ {

	}
}

func TestNoTips(t *testing.T) {

}

func TestNotImplemented(t *testing.T) {

}

func TestDriv(t *testing.T) {

}
