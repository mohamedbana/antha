// antha/AnthaStandardLibrary/Packages/REBASE/rebase.go: Part of the Antha language
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

package rebase

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/internal/github.com/jlaffaye/ftp"
	"io"
	"io/ioutil"
	"log"
	//"net/http"
	"os"
	//"strings"
)

func EnzymeLookup(name string) (enzyme wtype.LogicalRestrictionEnzyme) {
	if Exists("REBASETypeII.txt") == false {
		err := UpdateRebasefile()
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	enzymes, err := ioutil.ReadFile("REBASETypeII.txt")
	if err != nil {
		fmt.Println("error:", err)
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range parser.RebaseParse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if record.Name == name {
			fmt.Println(record)
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzyme = record
		}

	}
	return enzyme
}

/*
func LookupEnzymes(names []string) (enzymelist []wtype.LogicalRestrictionEnzyme) {
	if Exists("REBASETypeII.txt") == false {
		err := UpdateRebasefile()
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	enzymelist = make([]wtype.LogicalRestrictionEnzyme, 0)

	enzymes, err := ioutil.ReadFile("REBASETypeII.txt")
	if err != nil {
		fmt.Println("error:", err)
	}

	rebaseFh := bytes.NewReader(enzymes)

	for record := range parser.RebaseParse(rebaseFh) {


		for _, name := range names {
			if strings.Contains(record.Name, name) == true {
				fmt.Println(record)
				enzymelist = append(enzymelist, record)
				//enzyme = record
			}
		}

	}
	return enzymelist
}
*/

func UpdateRebasefile() (err error) {
	fmt.Println("Slurping...", "ftp://ftp.neb.com/pub/rebase/type2.txt")

	c, err := ftp.Dial("ftp.neb.com:21")
	if err != nil {
		log.Fatal(err)
	}
	err = c.Login("anonymous", "anonymous")
	if err != nil {
		panic(err)
	}

	err = c.ChangeDir("/pub/rebase")
	if err != nil {
		panic(err)
	}

	r, err := c.Retr("type2.txt")
	if err != nil {
		panic(err)
	} else {
		f, _ := os.Create("REBASETypeII.txt")
		fmt.Println("step 2: copying registry")
		_, err = io.Copy(f, r)
		if err != nil {
			panic(err)
		}

		r.Close()
		c.Quit()
	}

	return err
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
