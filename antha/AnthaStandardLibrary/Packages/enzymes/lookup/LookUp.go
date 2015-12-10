// LookUp
package lookup

import (
	"bytes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/REBASE"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Package for looking up restriction enzyme properties
func TypeIIsLookup(name string) (enzyme wtype.TypeIIs, err error) {
	enz := EnzymeLookup(name)

	enzyme, err = wtype.ToTypeIIs(enz)
	return
}

func EnzymeLookup(name string) (enzyme wtype.LogicalRestrictionEnzyme) {
	if anthapath.Anthafileexists("REBASETypeII.txt") == false {
		err := rebase.UpdateRebasefile()
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	enzymes, err := ioutil.ReadFile(filepath.Join(anthapath.Dirpath(), "REBASETypeII.txt"))
	if err != nil {
		fmt.Println("error:", err)
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

func FindEnzymesofClass(class string) (enzymelist []wtype.LogicalRestrictionEnzyme) {

	var enzyme wtype.LogicalRestrictionEnzyme

	if anthapath.Anthafileexists("REBASETypeII.txt") == false {
		err := rebase.UpdateRebasefile()
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	enzymes, err := ioutil.ReadFile(filepath.Join(anthapath.Dirpath(), "REBASETypeII.txt"))
	if err != nil {
		fmt.Println("error:", err)
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Class) == strings.ToUpper(class) {
			fmt.Println(record)
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzyme = record
			enzymelist = append(enzymelist, enzyme)
		}

	}
	return enzymelist
}

func FindEnzymeNamesofClass(class string) (enzymelist []string) {

	var enzyme string

	if anthapath.Anthafileexists("REBASETypeII.txt") == false {
		err := rebase.UpdateRebasefile()
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	enzymes, err := ioutil.ReadFile(filepath.Join(anthapath.Dirpath(), "REBASETypeII.txt"))
	if err != nil {
		fmt.Println("error:", err)
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range rebase.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.ToUpper(record.Class) == strings.ToUpper(class) {
			fmt.Println(record)
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzyme = record.Name
			enzymelist = append(enzymelist, enzyme)
		}

	}
	return enzymelist
}
