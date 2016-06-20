// Part of the Antha language
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

// Package for looking up restriction enzyme properties
package lookup

import (
	"bytes"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/REBASE"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func TypeIIsLookup(name string) (enzyme wtype.TypeIIs, err error) {
	enz := EnzymeLookup(name)

	enzyme, err = wtype.ToTypeIIs(enz)
	return
}

func EnzymeLookup(name string) (enzyme wtype.RestrictionEnzyme) {
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

func FindEnzymesofClass(class string) (enzymelist []wtype.RestrictionEnzyme) {

	var enzyme wtype.RestrictionEnzyme

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
