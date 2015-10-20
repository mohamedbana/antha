// antha/AnthaStandardLibrary/Packages/Parser/csv_parser.go: Part of the Antha language
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

// Needs to generate array of part names per construct for sockets
// Needs to generate an array of part sequences for validation.
//

package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"os"
	"strings"
)

/*type Part struct {
	Name     string
	Sequence string
	Plasmid  string
}
*/
type ConstructAssemblyParams struct {
	//parts is an array of array of strings.
	// each string is the name of a part, they are grouped into arrays for the ones that form a construct
	// each construct (which is an array of part names) is an entry of the uppermost array. Ex:
	// { {"part1", "part2", "part3"}, //construct 1 instance 1
	//   {"part1", "part2", "part3"}, //construct 1 instance 2
	//   {"part3, "part4"}, //construct 2 instance 1
	//   {"part3, "part4"}, //construct 2 instance 2
	//   {"part3, "part4"}, //construct 2 instance 3
	//   {"part3, "part4"}, //construct 2 instance 4
	// }
	// each array entry is translated into 1 construct, for a construct to be executed more than once, it must
	// be included as many times as desired inside the construct array
	Parts             [][]string
	Vector            string
	RestrictionEnzyme string
	Buffer            string
	Ligase            string
	Atp               string
	Outplate          string
	TipType           string
	ReactionVolume    string
	PartConc          string
	VectorConc        string
	AtpVol            string
	ReVol             string
	LigVol            string
	ReactionTemp      string
	ReactionTime      string
	InactivationTemp  string
	InactivationTime  string
}

func SetupParams(cnsts [][]string) ConstructAssemblyParams {

	//lets create an assembly params object
	cap := ConstructAssemblyParams{}
	cap.Vector = "component:standard_cloning_vector_mark_1"
	cap.RestrictionEnzyme = "component:SapI"
	cap.Buffer = "component:CutsmartBuffer"
	cap.Ligase = "component:T4Ligase"
	cap.Atp = "component:ATP"
	cap.Outplate = "plate:pcrplate_with_cooler"
	cap.TipType = "tipbox:Gilson50"
	cap.ReactionVolume = "20ul"
	cap.PartConc = "0.0001g/l"
	cap.VectorConc = "0.001g/l"
	cap.AtpVol = "1ul"
	cap.ReVol = "1ul"
	cap.LigVol = "1ul"
	cap.ReactionTemp = "25C"
	cap.ReactionTime = "1800s"
	cap.InactivationTemp = "40C"
	cap.InactivationTime = "60s"
	//Now the PARTS!
	cap.Parts = cnsts

	return cap

}

func ReadDesign(filename string) [][]string {

	var constructs [][]string

	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return constructs
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {
		var parts []string

		//fmt.Println("test", string(strings.TrimSpace(each[0])[0]))
		//fmt.Printf("name : %s and restriction enzyme : %s\n", each[0], each[1])

		if string(strings.TrimSpace(each[0])[0]) != "#" {
			for _, p := range each {
				if p != "" {
					parts = append(parts, strings.TrimSpace(p))
				}
			}
			constructs = append(constructs, parts)
		}
	}

	return constructs
}

func ReadParts(filename string) map[string]wtype.DNASequence {

	m := make(map[string]wtype.DNASequence)

	var parts []wtype.DNASequence

	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return m
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {

		var part wtype.DNASequence

		if string(strings.TrimSpace(each[0])[0]) != "#" && len(each) > 1 {

			part.Nm = each[0]

			if strings.ToUpper(each[3]) == "AA" || strings.ToUpper(each[3]) == "Protein" || strings.ToUpper(each[3]) == "Amino Acid" {

				part.Seq = enzymes.RevTranslatetoNstring(each[1])
			} else if each[3] == "DNA" || each[3] == "RNA" {
				part.Seq = strings.ToUpper(each[1])
			}

			part.Plasmid = false

			if len(each) > 2 {
				if strings.ToUpper(each[2]) == "TRUE" {
					part.Plasmid = true
				}
				if strings.ToUpper(each[2]) == "1" {
					part.Plasmid = true
				}
				/*	if each[2] == 1 {
					part.Plasmid = true
				}*/
				if strings.ToUpper(each[2]) == "PLASMID" {
					part.Plasmid = true
				}
				if strings.ToUpper(each[2]) == "YES" {
					part.Plasmid = true
				}
				if strings.ToUpper(each[2]) == "FALSE" {
					part.Plasmid = false
				}
				if strings.ToUpper(each[2]) == "LINEAR" {
					part.Plasmid = false
				}
				if strings.ToUpper(each[2]) == "NO" {
					part.Plasmid = false
				}
				if strings.ToUpper(each[2]) == "0" {
					part.Plasmid = false
				}
				/*	if each[2] == 0 {
					part.Plasmid = false
				}*/
			}
			parts = append(parts, part)
			m[part.Nm] = part
		}
	}

	return m

}

func Assemblyfromcsv(designfile string, partsfile string) (assemblyparameters []enzymes.Assemblyparameters) {

	var designedconstructs [][]string
	//var partsforvalidation [][]

	designedconstructs = ReadDesign(designfile)

	var definedparts map[string]wtype.DNASequence

	definedparts = ReadParts(partsfile)

	fmt.Println("Number of defined parts: ")
	fmt.Println(len(definedparts))

	for i, p := range definedparts {
		fmt.Println(i, p)
	}

	// Number of constructs:
	fmt.Println("Number of designed constructs: ")

	fmt.Println(len(designedconstructs))

	assemblyparameters = make([]enzymes.Assemblyparameters, 0)

	for i, c := range designedconstructs {
		var newassemblyparameters enzymes.Assemblyparameters
		fmt.Println("i=", i, "c=", c)
		newassemblyparameters.Constructname = c[0]
		newassemblyparameters.Enzymename = c[1]
		newassemblyparameters.Vector = definedparts[c[2]]
		var nextpart wtype.DNASequence

		partsinorder := make([]wtype.DNASequence, 0)
		for k := 3; k < len(c); k++ {
			//if definedparts[c[k]].Nm != "" {
			nextpart = definedparts[c[k]]
			if nextpart.Nm != "" {
				partsinorder = append(partsinorder, nextpart)
				//	}
			}
		}
		newassemblyparameters.Partsinorder = partsinorder

		assemblyparameters = append(assemblyparameters, newassemblyparameters)
		/*for j, p := range c {
		if j != 0 {

				fmt.Println("p=", p, "definedparts[p]", definedparts[p])
		}*/

	}
	params := SetupParams(designedconstructs)

	print := 0

	if print == 1 {
		jsn, err := json.Marshal(params)
		if err != nil {
			panic(err)
		}
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, jsn, "", "\t")
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(json))
		fmt.Println(prettyJSON.String())
	}

	fmt.Println("assemblyparameters=", assemblyparameters)
	return assemblyparameters
}
