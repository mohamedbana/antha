// LookUp
package lookup

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/REBASE"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/asset"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Package for looking up restriction enzyme properties
func TypeIIsLookup(name string) (enzyme wtype.TypeIIs, err error) {
	enz := EnzymeLookup(name)

	enzyme, err = wtype.ToTypeIIs(enz)
	return
}

func EnzymeLookup(name string) (enzyme wtype.RestrictionEnzyme) {
	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Name) == strings.ToUpper(name) {
			fmt.Println(record)
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzyme = record
		}

	}
	return enzyme
}

func FindEnzymesofClass(class string) (enzymelist []wtype.RestrictionEnzyme) {
	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Class) == strings.ToUpper(class) {
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzymelist = append(enzymelist, record)
		}
	}
	return enzymelist
}

func FindEnzymeNamesofClass(class string) (enzymelist []string) {
	for _, enzyme := range FindEnzymesofClass(class) {
		enzymelist = append(enzymelist, enzyme.Name)
	}
	return
}
