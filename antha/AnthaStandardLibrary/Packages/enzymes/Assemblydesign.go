// antha/AnthaStandardLibrary/Packages/enzymes/Assemblydesign.go: Part of the Antha language
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

package enzymes

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"strings"
)

/*
not finished
func LengthofPrefixOverlap(seq string, subseq string) (number int, end string) { // add end string


	i:=0;i<len(subseq);i++{
	truncated := subseq[i:]
	fmt.Println("truncated", truncated)
	if strings.HasPrefix(part.Seq, truncated) == true {
		number = i
		end = "end"
	}
	/*start := subseq[:i]
	fmt.Println("start", start)
	if strings.HasPrefix(part.Seq, start) == true {
		number = i
		end = "start"
	}
	return number
}
*/
func AddStandardStickyEnds(part wtype.DNASequence, assemblystandard string, level string, numberofparts int, position int) (Partwithends wtype.DNASequence) {

	if part.Seq != "" {
		if position > numberofparts {
			panic("Cannot have position greater than number of parts")
		}
		/*if position == 0 {
			("1st position = 1, not 0")
		}
		*/
		vectorends := Vectorends[assemblystandard][level] // this could also look up Endlinks[assemblystandard][level][numberofparts][0]

		enzyme := Enzymelookup[assemblystandard][level]

		bittoadd := Endlinks[assemblystandard][level][numberofparts][position]
		if strings.HasPrefix(part.Seq, bittoadd) == true {
			bittoadd = ""
		}

		// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

		truncated := bittoadd[1:]
		fmt.Println("truncated", truncated)
		if strings.HasPrefix(part.Seq, truncated) == true {
			//truncated = bittoadd[:1]
			//fmt.Println("truncated", truncated)
			bittoadd = bittoadd[:1]
		}

		bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
		fmt.Println("bittoadd5prime", bittoadd5prime)
		partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

		bittoadd3 := Endlinks[assemblystandard][level][numberofparts][position+1]
		fmt.Println("bittoadd3", bittoadd3)

		if numberofparts == position {
			bittoadd3 = RevComp(vectorends[0])
		}
		if strings.HasSuffix(part.Seq, bittoadd3) == true {
			bittoadd3 = ""
		}
		//fmt.Println("bittoadd3", bittoadd3)
		bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
		fmt.Println("bittoadd3prime", bittoadd3prime)
		partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

		Partwithends.Nm = part.Nm
		Partwithends.Plasmid = part.Plasmid
		Partwithends.Seq = partwithends
	}
	return Partwithends
}

func AddStandardStickyEndsfromClass(part wtype.DNASequence, assemblystandard string, level string, class string) (Partwithends wtype.DNASequence) {

	//vectorends := Vectorends[assemblystandard][level] // this could also look up Endlinks[assemblystandard][level][numberofparts][0]

	enzyme := Enzymelookup[assemblystandard][level]

	bitstoadd := EndlinksString[assemblystandard][level][class]
	bittoadd := bitstoadd[0]
	if strings.HasPrefix(part.Seq, bittoadd) == true {
		bittoadd = ""
	}

	// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

	truncated := bittoadd[1:]
	fmt.Println("truncated", truncated)
	if strings.HasPrefix(part.Seq, truncated) == true {
		//truncated = bittoadd[:1]
		//fmt.Println("truncated", truncated)
		bittoadd = bittoadd[:1]
	}

	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	fmt.Println("bittoadd5prime", bittoadd5prime)
	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := bitstoadd[1]
	fmt.Println("bittoadd3", bittoadd3)

	if strings.HasSuffix(part.Seq, bittoadd3) == true {
		bittoadd3 = ""
	}
	//fmt.Println("bittoadd3", bittoadd3)
	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	fmt.Println("bittoadd3prime", bittoadd3prime)
	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends

	return Partwithends
}

func AddCustomEnds(part wtype.DNASequence, enzyme wtype.TypeIIs, desiredstickyend5prime string, desiredstickyend3prime string) (Partwithends wtype.DNASequence) {

	bittoadd := desiredstickyend5prime
	if strings.HasPrefix(part.Seq, bittoadd) == true {
		bittoadd = ""
	}
	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := desiredstickyend3prime

	if strings.HasSuffix(part.Seq, bittoadd3) == true {
		bittoadd3 = ""
	}

	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends
	return Partwithends
}

func MakeStandardTypeIIsassemblyParts(parts []wtype.DNASequence, assemblystandard string, level string, optionalpartclasses []string) (partswithends []wtype.DNASequence) {

	partswithends = make([]wtype.DNASequence, 0)
	var partwithends wtype.DNASequence

	if len(optionalpartclasses) != 0 {
		if len(optionalpartclasses) == len(parts) {
			for i := 0; i < len(parts); i++ {
				partwithends = AddStandardStickyEndsfromClass(parts[i], assemblystandard, level, optionalpartclasses[i])
				partswithends = append(partswithends, partwithends)
			}
		}
	} else {

		for i := 0; i < len(parts); i++ {
			partwithends = AddStandardStickyEnds(parts[i], assemblystandard, level, (len(parts)), i+1)
			partswithends = append(partswithends, partwithends)
		}
	}
	return partswithends
}

func CheckForExistingTypeIISEnds(part wtype.DNASequence, enzyme wtype.TypeIIs) (numberofsitesfound int, stickyends5 []string, stickyends3 []string) {

	enz := lookup.EnzymeLookup(enzyme.Name)

	sites := Restrictionsitefinder(part, []wtype.LogicalRestrictionEnzyme{enz})

	numberofsitesfound = sites[0].Numberofsites
	_, stickyends5, stickyends3 = TypeIIsdigest(part, enzyme)

	return
}

/*

func HandleExistingEnds (parts []wtype.DNASequence, enzymewtype.LogicalRestrictionEnzyme)(partswithoverhangs []wtype.DNASequence {
	partswithexistingsites := make([]RestrictionSites, 0)

	for _, part := range parts {
		sites := Restrictionsitefinder(part, wtype.LogicalRestrictionEnzyme{enzyme})
		if len(sites) != 0 {
			partswithexistingsites = append(partswithexistingsites, sites)
		}

	}
	return
}

func AddStandardVectorEnds (vector wtype.DNASequence, standard, level string) (vectrowithends wtype.DNASequence) {

	}
*/

func MakeScarfreeCustomTypeIIsassemblyParts(parts []wtype.DNASequence, vector wtype.DNASequence, enzyme wtype.TypeIIs) (partswithends []wtype.DNASequence) {

	partswithends = make([]wtype.DNASequence, 0)
	var partwithends wtype.DNASequence

	// find sticky ends from cutting vector with enzyme

	_, stickyends5, _ := TypeIIsdigest(vector, enzyme)

	//initialise

	desiredstickyend5prime := ""

	vector3primestickyend := ""

	// add better logic for the scenarios where the vector is cut more than twice or we want to add fragment in either direction
	for i := 0; i < len(stickyends5); i++ {
		if stickyends5[i] != "" {

			// RevComp() // fill in later
			vector3primestickyend = stickyends5[i]
			desiredstickyend5prime = stickyends5[i+1]
			{
				break
			}
		}
	} // fill in later

	// declare as blank so no end added
	desiredstickyend3prime := ""

	for i := 0; i < len(parts); i++ {
		if i == (len(parts) - 1) {
			desiredstickyend3prime = vector3primestickyend
		}
		partwithends = AddCustomEnds(parts[i], enzyme, desiredstickyend5prime, desiredstickyend3prime)
		partswithends = append(partswithends, partwithends)

		desiredstickyend5prime = Suffix(parts[i].Seq, enzyme.LogicalRestrictionEnzyme.EndLength)

	}

	return partswithends
}

func Addoverhang(seq string, bittoadd string, end string) (seqwithoverhang string) {

	if end == "5prime" {
		seqwithoverhang = strings.Join([]string{bittoadd, seq}, "")
	}
	if end == "3prime" {
		seqwithoverhang = strings.Join([]string{seq, bittoadd}, "")
	}
	return seqwithoverhang
}

func Makeallspaceroptions(spacerlength int) (finalarray []string) {
	// only works for spacer length 1 or 2

	// new better code, but untested! test and replace code below
	newarray := make([][]string, 0)
	for i := 0; i < spacerlength; i++ {
		newarray = append(newarray, nucleotides)
	}

	finalarray = AllCombinations(newarray)

	return finalarray
}

func ChooseSpacer(spacerlength int, seq string, seqstoavoid []string) (spacer string) {
	// very simple case to start with

	possibilities := Makeallspaceroptions(spacerlength)

	if len(seqstoavoid) == 0 {
		spacer = possibilities[0]
	} else {
		for _, possibility := range possibilities {
			if len(Findallthings(strings.Join([]string{seq, possibility}, ""), seqstoavoid)) == 0 &&
				len(Findallthings(strings.Join([]string{possibility, seq}, ""), seqstoavoid)) == 0 &&
				len(Findallthings(RevComp(strings.Join([]string{possibility, seq}, "")), seqstoavoid)) == 0 &&
				len(Findallthings(RevComp(strings.Join([]string{seq, possibility}, "")), seqstoavoid)) == 0 {
				spacer = possibility
			}
		}
	}
	return spacer
}

var nucleotides = []string{"A", "T", "C", "G"}

func Addnucleotide(s string) (splus1array []string) {

	splus1 := s
	splus1array = make([]string, 0)
	for _, nucleotide := range nucleotides {
		splus1 = strings.Join([]string{s, nucleotide}, "")
		splus1array = append(splus1array, splus1)
	}
	return splus1array
}

func Makeoverhang(enzyme wtype.TypeIIs, end string, stickyendseq string, spacer string) (seqwithoverhang string) {
	if end == "5prime" {
		if enzyme.Topstrand3primedistancefromend < 0 {
			panic("Unlikely to work with this enzyme in making a 5'prime spacer")
		}

		if len(spacer) != enzyme.Topstrand3primedistancefromend {
			panic("length of spacer will lead to cutting at run position! change length to match enzyme NN region length")
		}
		seqwithoverhang = strings.Join([]string{enzyme.LogicalRestrictionEnzyme.RecognitionSequence, spacer, stickyendseq}, "")
	}

	// This case needs work, but may not appear in reality so is a place holder for now until a real scenario becomes apparent
	if end == "3prime" {
		/*if enzyme.Topstrand3primedistancefromend < 0 && len(spacer) == enzyme.Bottomstrand5primedistancefromend {
			seqwithoverhang = strings.Join([]string{stickyendseq, spacer, enzyme.LogicalRestrictionEnzyme.RecognitionSequence}, "")
		}*/
		seqwithoverhang = strings.Join([]string{stickyendseq, spacer, RevComp(enzyme.LogicalRestrictionEnzyme.RecognitionSequence)}, "")
	}
	return seqwithoverhang

}

var Endlinks = map[string]map[string]map[int]map[int]string{
	//map["assembly strategy"]map[number of parts]map[part number]"sticky end to add"
	"MoClo_Raven": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			4: map[int]string{
				0: "AAGC",
				1: "GAGG",
				2: "TACT",
				3: "AATG",
				4: "AGGT",
			},
		},
	},
	"MoClo": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			4: map[int]string{
				0: "AAGC",
				1: "GGTA",
				2: "",
				3: "",
				4: "",
			},
		},
	},
	"Electra": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			1: map[int]string{
				0: "ATG",
			},
		},
	},
}
var EndlinksString = map[string]map[string]map[string][]string{
	"MoClo": map[string]map[string][]string{
		"Level0": map[string][]string{
			"Pro":         []string{"GGAG", "TACT"},
			"5U":          []string{"TACT", "CCAT"},
			"5U(f)":       []string{"TACT", "CCAT"},
			"Pro + 5U(f)": []string{"GGAG", "CCAT"},
			"Pro + 5U":    []string{"GGAG", "AATG"},
			"NT1":         []string{"CCAT", "AATG"},
			"5U + NT1":    []string{"TACT", "AATG"},
			"CDS1":        []string{"AATG", "GCTT"},
			"CDS1 ns":     []string{"AATG", "TTCG"},
			"NT2":         []string{"AATG", "AGGT"},
			"SP":          []string{"AATG", "AGGT"},
			"CDS2 ns":     []string{"AGGT", "TTCG"},
			"CDS2":        []string{"AGGT", "GCTT"},
			"CT":          []string{"TTCG", "GCTT"},
			"3U":          []string{"GCTT", "GGTA"},
			"Ter":         []string{"GGTA", "CGCT"},
			"3U + Ter":    []string{"GCTT", "CGCT"},
		},
	},
	"MoClo_Raven": map[string]map[string][]string{
		"Level0": map[string][]string{
			"Pro":         []string{"GAGG", "TACT"},
			"5U":          []string{"TACT", "CCAT"},
			"5U(f)":       []string{"TACT", "CCAT"},
			"Pro + 5U(f)": []string{"GGAG", "CCAT"},
			"Pro + 5U":    []string{"GGAG", "AATG"},
			"NT1":         []string{"CCAT", "AATG"},
			"5U + NT1":    []string{"TACT", "AATG"},
			"CDS1":        []string{"AATG", "GCTT"},
			"CDS1 ns":     []string{"AATG", "TTCG"},
			"NT2":         []string{"AATG", "AGGT"},
			"SP":          []string{"AATG", "AGGT"},
			"CDS2 ns":     []string{"AGGT", "TTCG"},
			"CDS2":        []string{"AGGT", "GCTT"},
			"CT":          []string{"TTCG", "GCTT"},
			"3U":          []string{"GCTT", "GGTA"},
			"Ter":         []string{"GGTA", "CGCT"},
			"3U + Ter":    []string{"GCTT", "GCTT"},
		},
	},
}

const (
	// for indexing part position based on part name/class
	VECTOR = iota
	PROMOTER
	RBS
	CDS
	TERMINATOR
)

var Vectorends = map[string]map[string][]string{
	// array of strings returned correspond to [3'overhang and 5'overhang]
	"MoClo_Raven": map[string][]string{
		"Level0": []string{"AAGC", "CCTC"}, //
		"Level1": []string{"", ""},
	},
	"MoClo": map[string][]string{
		"Level0": []string{"CGCT", "GGAG"}, //
		"Level1": []string{"", ""},
	},
	"Synthace": map[string][]string{
		"Level0": []string{"GGT", "GAA"},
		"Level1": []string{"", ""},
	},
	"Electra": map[string][]string{
		"Level0": []string{"GGT", "ATG"},
		"Level1": []string{"", ""},
	},
}

var Enzymelookup = map[string]map[string]wtype.TypeIIs{
	// array of strings returned correspond to 5'overhang and 3'overhang
	"MoClo_Raven": map[string]wtype.TypeIIs{
		"Level0": BsaIenz,
		"Level1": BpiIenz,
	},
	"MoClo": map[string]wtype.TypeIIs{
		"Level0": BsaIenz,
		"Level1": BpiIenz,
	},
	"Electra": map[string]wtype.TypeIIs{
		"Level0": SapIenz,
	},
}

/*
var MoClo AssemblyStandard{
	[]AssemblyStandardLevel{{BsaIenz,"Level0"},{BpiIenz,"Level1"}}"Moclo"}

*/
type AssemblyStandardLevel struct {
	Enzyme    wtype.TypeIIs
	Levelname string
}

type AssemblyStandard struct {
	Endstable       map[string]map[string]map[int]map[int]string
	EnzymeTable     map[string]map[string]wtype.TypeIIs
	VectorEndstable map[string]map[string][]string // Vector 5prime can also be found in Endstable position 0
}
