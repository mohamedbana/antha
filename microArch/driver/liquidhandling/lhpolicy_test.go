package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func getChannelForTest() *wtype.LHChannelParameter {
	return wtype.NewLHChannelParameter("ch", "gilson", wunit.NewVolume(20.0, "ul"), wunit.NewVolume(200.0, "ul"), wunit.NewFlowRate(0.0, "ml/min"), wunit.NewFlowRate(100.0, "ml/min"), 8, false, wtype.LHVChannel, 1)
}

func TestPolicyMerger(t *testing.T) {
	pft, _ := GetLHPolicyForTest()

	tp := TransferParams{
		What:    "PEG",
		Volume:  wunit.NewVolume(190.0, "ul"),
		Channel: getChannelForTest(),
	}

	ins1 := NewSuckInstruction()
	ins1.AddTransferParams(tp)

	p := pft.GetPolicyFor(ins1)

	if p["ASPZOFFSET"].(float64) != 2.5 {
		t.Fatal("ASPZOFFSET for PEG must be 2.5")
	}
	if p["DSPZOFFSET"].(float64) != 2.5 {
		t.Fatal("DSPZOFFSET for PEG must be 2.5")
	}
	if p["POST_MIX_Z"].(float64) != 3.5 {
		t.Fatal("POST_MIX_Z for PEG must be 3.5")
	}

	for i := 0; i < 100; i++ {
		q := pft.GetPolicyFor(ins1)

		if q["ASPZOFFSET"] != p["ASPZOFFSET"] {
			t.Fatal("Inconsistent Z offsets returned")
		}
	}

}
