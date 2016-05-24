package wtype

import (
	"fmt"
	"testing"
)

func TestLHError(t *testing.T) {
	e := LHError(LH_OK, "")

	if e.Error() != "0 (LH_OK): no problem -- " {
		t.Fatal("LH Error format changed... this may break stuff")
	}

	e = LHError(LH_ERR_NO_DECK_SPACE, "can't fit tip box in")

	if e.Error() != "1 (LH_ERR_NO_DECK_SPACE): not sufficient deck space to fit all required items; this may be due to constraints -- can't fit tip box in" {
		t.Fatal("LH Error format changed... this may break stuff")
	}
}
