package enzymes

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"testing"
)

func TestNoRotationNeeded(t *testing.T) {
	enzyme, _ := lookup.TypeIIsLookup("SAPI")
	seq := "GCTCTTCxxxxx"
	rseq := "GCTCTTCxxxxx"

	s := wtype.DNASequence{Nm: "nevermind", Seq: seq}

	rs, err := rotate_vector(s, enzyme)

	if err != nil {
		t.Fatal(err)
	}

	if rs.Seq != rseq {
		t.Fatal(fmt.Sprintf("Error with vector rotation: got %s expected %s", s, rs))
	}
}
func TestSomeRotationNeeded(t *testing.T) {
	enzyme, _ := lookup.TypeIIsLookup("SAPI")
	seq := "xxxxxGCTCTTCn"
	rseq := "GCTCTTCnxxxxx"

	s := wtype.DNASequence{Nm: "nevermind", Seq: seq}

	rs, err := rotate_vector(s, enzyme)

	if err != nil {
		t.Fatal(err)
	}

	if rs.Seq != rseq {
		t.Fatal(fmt.Sprintf("Error with vector rotation: got %s expected %s", rs.Seq, rseq))
	}
}