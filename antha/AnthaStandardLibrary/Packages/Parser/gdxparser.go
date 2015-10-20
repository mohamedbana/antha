// antha/AnthaStandardLibrary/Packages/Parser/gdxparser.go: Part of the Antha language
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

package parser

import (
	"encoding/xml"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Project struct {
	DesignConstruct []DesignConstruct `xml:"DesignConstruct"`
}

type DesignConstruct struct {
	Label       string       `xml:"label,attr"`
	Plasmid     string       `xml:"circular,attr"`
	Rev         string       `xml:"reverseComplement,attr"`
	Notes       string       `xml:"notes"`
	DNAElements []DNAElement `xml:"DNAElement"`
	AAElements  []AAElement  `xml:"AAElement"`
}

type DNAElement struct {
	Label    string `xml:"label,attr"`
	Sequence string `xml:"sequence"`
	Notes    string `xml:"notes"`
}

type AAElement struct {
	Label    string `xml:"label,attr"`
	Sequence string `xml:"sequence"`
	Notes    string `xml:"notes"`
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func Parse(filename string) []string {

	fmt.Println("parsing ", filename)
	str, _ := ioutil.ReadFile(filename)

	var gdx Project

	err := xml.Unmarshal(str, &gdx)
	checkError(err)

	fmt.Println(gdx)

	parts_list := make([]string, len(gdx.DesignConstruct))

	nconstructs := 0
	for _, c := range gdx.DesignConstruct {

		//fmt.Println(c.Label)
		//f, err:= os.Create("outfiles/"+c.Label+".csv")
		//checkError(err)
		//defer f.Close()
		//nparts := 0
		//for _,e := range c.DNAElement {
		//    nparts +=1
		//    row := e.Label  + "," + e.Sequence + "\n"
		//    fmt.Println(row)
		//   // _,err = f.Write([]byte(row))
		//}
		//fmt.Println("Construct: " , nconstructs, " n parts: ", nparts)
		//fmt.Println(len(gdx.DesignConstruct))
		parts_list[nconstructs] = "Construct: " + strconv.Itoa(nconstructs) + " n parts: " + strconv.Itoa(len(c.DNAElements)+len(c.AAElements))
		nconstructs += 1

	}

	return parts_list
}

/*
type AnnotatedSeq struct {
	Nm                string
	Seq               string
	Plasmid           bool
	Reversecomplement bool
	Elements          []Element
}

type Element struct {
	DNASequence       wtype.DNASequence
	ORF               bool
	Promoter          bool
	RBS               bool
	Terminator        bool
	Reversecomplement bool
}*/
/*type Feature struct {
	Name          string
	Class         string //int // defined by constants above
	Reverse       bool
	StartPosition int
	EndPosition   int
	DNASeq        string
	Protseq       string
}

type AnnotatedSeq struct {
	Nm       string
	Seq      string
	features []Feature
}

func MakeFeature(name string, seq string, sequencetype string, class string, reverse string) (feature Feature) {

	feature.Name = name
	feature.DNASeq = strings.ToUpper(seq)

	//features := make([]Feature,0)
	//feature := Feature
	//fmt.Println("len seq =", len(seq))
	feature.Class = class
	if sequencetype == "aa" {
		feature.DNASeq = RevTranslatetoNstring(seq)
		feature.Protseq = seq
		feature.StartPosition = 1
		feature.EndPosition = len(feature.DNASeq)
		fmt.Println("len seq =", len(feature.DNASeq))
	} else {
		feature.DNASeq = seq
		feature.StartPosition = 1
		feature.EndPosition = len(seq)
		fmt.Println("len seq =", len(seq))
	}

	if reverse == "Reverse" {
		feature.Reverse = true
	}
	if feature.Class == "orf" {
		orf, orftrue := FindORF(seq)
		if orftrue == true {
			fmt.Println("orftrue!)")
			feature.Protseq = orf.Protseq
			feature.StartPosition = orf.StartPosition
			feature.EndPosition = orf.EndPosition
		}
	}
	return feature
}

func ConcatenateFeatures(name string, featuresinorder []Feature) (annotated AnnotatedSeq) {

	annotated.Nm = name
	annotated.Seq = featuresinorder[0].DNASeq
	annotated.features = make([]Feature, 0)
	annotated.features = append(annotated.features, featuresinorder[0])
	for i := 1; i < len(featuresinorder); i++ {
		nextfeature := featuresinorder[i]
		nextfeature.StartPosition = nextfeature.StartPosition + annotated.features[i-1].EndPosition
		nextfeature.EndPosition = nextfeature.EndPosition + annotated.features[i-1].EndPosition
		annotated.Seq = annotated.Seq + featuresinorder[i].DNASeq
		annotated.features = append(annotated.features, nextfeature)
	}
	return annotated
}

func ParsetoAnnotatedSeq(filename string) []AnnotatedSeq {

	fmt.Println("parsing ", filename)
	str, _ := ioutil.ReadFile(filename)

	var gdx Project

	err := xml.Unmarshal(str, &gdx)
	checkError(err)

	fmt.Println(gdx)

	parts_list := make([]wtype.DNASequence, 0)

	for _, a := range gdx.DesignConstruct {
		for _, b := range a.DNAElement {
			var newseq wtype.DNASequence
			for i := 0; i < len(a.DNAElement); i++ {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if a.Plasmid == "true" {
					newseq.Plasmid = true
				}
				parts_list = append(parts_list, newseq)
			}
		}
	}

	return parts_list
}

/*
/*func ParsetoAnnotatedSeq(filename string) []AnnotatedSeq {/*

	fmt.Println("parsing ", filename)
	str, _ := ioutil.ReadFile(filename)

	var gdx Project

	err := xml.Unmarshal(str, &gdx)
	checkError(err)

	fmt.Println(gdx)

	elements_list := make([]wtype.DNASequence, 0)

	for _, a := range gdx.DesignConstruct {
		for _, b := range a.DNAElement {
			var newseq wtype.DNASequence
			for i := 0; i < len(a.DNAElement); i++ {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if a.Plasmid == "true" {
					newseq.Plasmid = true
				}
				elements_list = append(elements_list, newseq)
			}
		}
	}

	return parts_list
}
*/
/*type Assemblyparameters struct {
	Constructname string
	Enzymename    string
	Vector        wtype.DNASequence
	Partsinorder  []wtype.DNASequence
}*/

func ParsetoAssemblyParameters(filename string) []enzymes.Assemblyparameters {

	fmt.Println("parsing ", filename)
	str, _ := ioutil.ReadFile(filename)

	var gdx Project

	err := xml.Unmarshal(str, &gdx)
	checkError(err)

	fmt.Println(gdx)

	construct_list := make([]enzymes.Assemblyparameters, 0)
	//parts_list := make([]wtype.DNASequence, 0)

	for _, a := range gdx.DesignConstruct {
		fmt.Println("len(designconstructs)", len(gdx.DesignConstruct))
		var newconstruct enzymes.Assemblyparameters
		newconstruct.Constructname = a.Label
		fmt.Println("name=", newconstruct.Constructname)
		if strings.Contains(a.Notes, "Enzyme:") == true {
			newconstruct.Enzymename = strings.TrimSpace(strings.TrimPrefix(a.Notes, "Enzyme:")) // add trim function to trim after space
		}
		fmt.Println("enzyme=", newconstruct.Enzymename)
		for _, b := range a.DNAElements {
			var newseq wtype.DNASequence
			//for i := 0; i < len(a.DNAElements); i++ {
			if strings.Contains(strings.ToUpper(b.Notes), "VECTOR") == true {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if strings.Contains(strings.ToUpper(a.Notes), "PLASMID") == true || strings.Contains(strings.ToUpper(a.Notes), "CIRCULAR") == true {
					newseq.Plasmid = true
				}
				newconstruct.Vector = newseq
			} else {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if strings.Contains(a.Notes, "Plasmid") == true {
					newseq.Plasmid = true
				}
				newconstruct.Partsinorder = append(newconstruct.Partsinorder, newseq)

			}
			fmt.Println("vector", newconstruct.Vector)
			fmt.Println("len(newconstruct.Partsinorder)=", len(newconstruct.Partsinorder))

		}

		construct_list = append(construct_list, newconstruct)
	}
	fmt.Println("len(constructlist)=", len(construct_list))
	return construct_list
}

func ParsetoDNASequence(filename string) []wtype.DNASequence {

	fmt.Println("parsing ", filename)
	str, _ := ioutil.ReadFile(filename)

	var gdx Project

	err := xml.Unmarshal(str, &gdx)
	checkError(err)

	fmt.Println(gdx)

	parts_list := make([]wtype.DNASequence, 0)

	for _, a := range gdx.DesignConstruct {
		for _, b := range a.DNAElements {
			var newseq wtype.DNASequence
			for i := 0; i < len(a.DNAElements); i++ {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if a.Plasmid == "true" {
					newseq.Plasmid = true
				}
				parts_list = append(parts_list, newseq)
			}
		}
	}

	return parts_list
}
