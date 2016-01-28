package wtype_testing

import (
	"testing"

	. "github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"
)

func TestColumniterator_96(t *testing.T) {
	p := factory.GetPlateByType("pcrplate_skirted")
	it := NewColumnIterator(p)

	for i := 0; i < 96; i++ {

		if i == 12 {
			wt := WellCoords{0, 1}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}
		it.Next()
	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}
func TestColumniterator_384(t *testing.T) {
	p := factory.GetPlateByType("greiner384")
	it := NewColumnIterator(p)
	for i := 0; i < 384; i++ {

		if i == 24 {
			wt := WellCoords{0, 1}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}

		it.Next()
	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}
func TestColumniterator_1536(t *testing.T) {
	p := factory.GetPlateByType("nunc1536_riser")
	it := NewColumnIterator(p)
	for i := 0; i < 1536; i++ {

		if i == 48 {
			wt := WellCoords{0, 1}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}

		it.Next()
	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}

func TestRowiterator_96(t *testing.T) {
	p := factory.GetPlateByType("pcrplate_skirted")
	it := NewRowIterator(p)

	for i := 0; i < 96; i++ {

		if i == 8 {
			wt := WellCoords{1, 0}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}
		it.Next()

	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}
func TestRowiterator_384(t *testing.T) {
	p := factory.GetPlateByType("greiner384")
	it := NewRowIterator(p)
	for i := 0; i < 384; i++ {

		if i == 16 {
			wt := WellCoords{1, 0}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}

		it.Next()
	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}
func TestRowiterator_1536(t *testing.T) {
	p := factory.GetPlateByType("nunc1536_riser")
	it := NewRowIterator(p)
	for i := 0; i < 1536; i++ {

		if i == 32 {
			wt := WellCoords{1, 0}
			if !wt.Equals(it.Curr()) {
				t.Error()
			}
		}

		it.Next()
	}

	wc := it.Curr()

	wt := WellCoords{0, 0}

	if !wc.Equals(wt) {
		t.Error()
	}
}
