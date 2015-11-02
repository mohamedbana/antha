// wtype/biology.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math/rand"
	"strings"
)

// the following are all physical things; we need a way to separate
// out just the logical part

// structure which defines an enzyme -- solutions containing
// enzymes need careful handling as they can be quite delicate
type Enzyme struct {
	Properties map[string]wunit.Measurement
}

type LogicalRestrictionEnzyme struct {
	// other fields required but for now the main things are...
	RecognitionSequence               string
	EndLength                         int
	Name                              string
	Prototype                         string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
	MethylationSite                   string   //"attr, <4>"
	CommercialSource                  []string //string "attr, <5>"
	References                        []int
	Class                             string
}

// structure which defines an organism. These need specific handling
// -- some detail is derived using the TOL structure
type Organism struct {
	Species *TOL // position on the TOL
}

// a set of organisms, can be mixed or homogeneous
type Population struct {
}

// defines a plasmid
type Plasmid struct {
}

// defines things which have biosequences... useful for operations
// valid on biosequences such as BLASTing / other alignment methods
type BioSequence interface {
	Name() string
	Sequence() string
	Append(string)
	Prepend(string)
}

// defines something as physical DNA
// hence it is physical and has a DNASequence
type DNA struct {
	GenericPhysical
	Seq DNASequence
}

// DNAsequence is a type of Biosequence
type DNASequence struct {
	Nm             string
	Seq            string
	Plasmid        bool
	Singlestranded bool
	Overhang5prime Overhang
	Overhang3prime Overhang
	Methylation    string // add histones etc?
}

func MakeDNASequence(name string, seqstring string, properties []string) (seq DNASequence, err error) {
	seq.Nm = name
	seq.Seq = seqstring
	for _, property := range properties {
		property = strings.ToUpper(property)

		if strings.Contains(property, "DCM") || strings.Contains(property, "DAM") || strings.Contains(property, "CPG") {
			seq.Methylation = property
		}

		if strings.Contains(property, "PLASMID") || strings.Contains(property, "CIRCULAR") || strings.Contains(property, "VECTOR") {
			seq.Plasmid = true
			break
		}
		if strings.Contains(property, "SS") || strings.Contains(property, "SINGLE STRANDED") {
			seq.Singlestranded = true
			break
		}
		/*
		   // deal with overhangs separately
		   if strings.Contains(property,"5'") {
		   	seq.Overhang5prime.End = 5
		   	seq.Overhang5prime.Type =
		   }
		*/
	}
	return
}
func MakeLinearDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring

	return
}
func MakePlasmidDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring
	seq.Plasmid = true
	return
}
func MakeSingleStrandedDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring
	seq.Singlestranded = true
	return
}

func MakeOverhang(sequence DNASequence, end int, toporbottom int, length int, phosphorylated bool) (overhang Overhang, err error) {

	if sequence.Singlestranded {
		err = fmt.Errorf("Can't have overhang on single stranded dna")
		return
	}
	if sequence.Plasmid {
		err = fmt.Errorf("Can't have overhang on Plasmid(circular) dna")
		return
	}
	if end == 0 {
		err = fmt.Errorf("if end = 0, all fields are returned empty")
		return
	}

	if end == 5 || end == 3 || end == 0 {
		overhang.End = end
	} else {
		err = fmt.Errorf("invalid entry for end: 5PRIME = 5, 3PRIME = 3, NA = 0")
		return
	}
	if toporbottom == 0 && length == 0 {
		overhang.Type = 1
		return
	}
	if toporbottom == 0 && length != 0 {
		err = fmt.Errorf("If length of overhang is not 0, toporbottom must be 0")
		return
	}
	if toporbottom != 0 && length == 0 {
		err = fmt.Errorf("If length of overhang is not 0, toporbottom must be 0")
		return
	}
	if toporbottom > 2 {
		err = fmt.Errorf("invalid entry for toporbottom: NEITHER = 0, TOP    = 1, BOTTOM = 2")
		return
	}
	if toporbottom == 1 {
		overhang.Type = 2
		overhang.Sequence = sequences.Prefix(sequence.Seq, length)
	}
	if toporbottom == 2 {
		overhang.Type = -1
		overhang.Sequence = sequences.Suffix(sequences.RevComp(sequence.Seq), length)
	}
	overhang.Phosphorylation = phosphorylated
	return
}

func Phosphorylate(dnaseq DNASequence) (phosphorylateddna DNASequence, err error) {
	if dnaseq.Plasmid == true {
		err = fmt.Errorf("Can't phosphorylate circular dna")
		phosphorylateddna = dnaseq
		return
	}
	if dnaseq.Overhang5prime.Type != 0 {
		dnaseq.Overhang5prime.Phosphorylation = true
	}
	if dnaseq.Overhang3prime.Type != 0 {
		dnaseq.Overhang3prime.Phosphorylation = true
	}
	if dnaseq.Overhang3prime.Type == 0 && dnaseq.Overhang5prime.Type == 0 {
		err = fmt.Errorf("No ends available, but not plasmid! This doesn't seem possible!")
		phosphorylateddna = dnaseq
	}
	return
}

const (
	FALSE     = 0
	BLUNT     = 1
	OVERHANG  = 2
	UNDERHANG = -1
)

const (
	NEITHER = 0
	TOP     = 1
	BOTTOM  = 2
)

/*const (
	5PRIME = 5
	3PRIME = 3
	NA = 0
)*/

type Overhang struct {
	//Strand          int // i.e. 1 or 2 (top or bottom
	End             int // i.e. 5 or 3 or 0
	Type            int //as contants above
	Length          int
	Sequence        string
	Phosphorylation bool
}

func (dna *DNASequence) Sequence() string {
	return dna.Seq
}
func (dna *DNASequence) Name() string {
	return dna.Nm
}

func (dna *DNASequence) Append(s string) {
	dna.Seq = dna.Seq + s
}

func (dna *DNASequence) Prepend(s string) {
	dna.Seq = s + dna.Seq
}

// RNA sample: physical RNA, has an RNASequence object
type RNA struct {
	GenericPhysical
	Seq RNASequence
}

// RNASequence object is a type of Biosequence
type RNASequence struct {
	Nm  string
	Seq string
}

func (rna *RNASequence) Sequence() string {
	return rna.Seq
}

func (rna *RNASequence) Name() string {
	return rna.Nm
}

func (rna *RNASequence) Append(s string) {
	rna.Seq = rna.Seq + s
}

func (rna *RNASequence) Prepend(s string) {
	rna.Seq = s + rna.Seq
}

// physical protein sample
// has a ProteinSequence
type Protein struct {
	GenericPhysical
	Seq ProteinSequence
}

// ProteinSequence object is a type of Biosequence
type ProteinSequence struct {
	Nm  string
	Seq string
}

func (prot *ProteinSequence) Sequence() string {
	return prot.Seq
}

func (prot *ProteinSequence) Name() string {
	return prot.Nm
}

func (prot *ProteinSequence) Append(s string) {
	prot.Seq = prot.Seq + s
}

func (prot *ProteinSequence) Prepend(s string) {
	prot.Seq = s + prot.Seq
}

func random_dna_seq(leng int) string {
	s := ""
	for i := 0; i < leng; i++ {
		s += random_char("ACTG")
	}
	return s
}

func random_char(chars string) string {
	return string(chars[rand.Intn(len(chars))])
}

func makeABunchaRandomSeqs(n_seq_sets, seqs_per_set, min_len, len_var int) [][]DNASequence {
	var seqs [][]DNASequence

	seqs = make([][]DNASequence, n_seq_sets)

	for i := 0; i < n_seq_sets; i++ {
		seqs[i] = make([]DNASequence, seqs_per_set)
		for j := 0; j < seqs_per_set; j++ {
			seqs[i][j] = DNASequence{fmt.Sprintf("SEQ%04d", i*seqs_per_set+j+1), random_dna_seq(rand.Intn(len_var) + min_len), false, false, Overhang{0, 0, 0, "", false}, Overhang{0, 0, 0, "", false}, ""}
		}
	}
	return seqs
}
