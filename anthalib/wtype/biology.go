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
// 1 Royal College St, London NW1 0NH UK

package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wutil"
	"math/rand"
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
	RecognitionSequence string
	CutDist             int
	EndLength           int
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
	Nm  string
	Seq string
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

func makeABunchaRandomSeqs(n_seq_sets, seqs_per_set, min_len, len_var int) [][]DNASequence {
	var seqs [][]DNASequence

	seqs = make([][]DNASequence, n_seq_sets)

	for i := 0; i < n_seq_sets; i++ {
		seqs[i] = make([]DNASequence, seqs_per_set)
		for j := 0; j < seqs_per_set; j++ {
			seqs[i][j] = DNASequence{fmt.Sprintf("SEQ%04d", i*seqs_per_set+j+1), wutil.Random_dna_seq(rand.Intn(len_var) + min_len)}
		}
	}
	return seqs
}
