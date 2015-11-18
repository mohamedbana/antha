// antha/AnthaStandardLibrary/Packages/enzymes/Translation.go: Part of the Antha language
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

package sequences

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"strings"
)

/*
func Siteinorfs(Features features, site string) bool {

}
*/

var Algorithmlookuptable = map[string]ReplacementAlgorithm{
	"ReplacebyComplement": ReplaceBycomplement,
}

// will potentially be generalisable for codon optimisation
type ReplacementAlgorithm func(sequence, thingtoreplace string, otherseqstoavoid []string) (replacement string, err error)

func ReplaceBycomplement(sequence, thingtoreplace string, otherseqstoavoid []string) (replacement string, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, thingtoreplace)

	for i, _ := range thingtoreplace {

		replacementnucleotide := Comp(string(thingtoreplace[i]))
		replacement := strings.Replace(thingtoreplace, string(thingtoreplace[i]), replacementnucleotide, 1)
		newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return replacement, err
		}
	}

	for i, _ := range thingtoreplace {

		replacementnucleotide := Comp(thingtoreplace[i : i+1])
		replacement := strings.Replace(thingtoreplace, thingtoreplace[i:i+1], replacementnucleotide, 1)
		newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return replacement, err
		}
		if i+2 == len(thingtoreplace) {
			specificseqs := text.Print("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
			break
		}
	}
	return
}

// iterates through each position of a restriction site and replaces with the complementary base and then removes these from the main sequence
// if that fails the algorithm will attempt to find the complements of two adjacent positions. The algorithm needs improvement
func RemoveSite(sequence wtype.DNASequence, enzyme wtype.LogicalRestrictionEnzyme, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzyme.RecognitionSequence)

	for i, _ := range enzyme.RecognitionSequence {

		replacementnucleotide := Comp(string(enzyme.RecognitionSequence[i]))
		replacement := strings.Replace(enzyme.RecognitionSequence, string(enzyme.RecognitionSequence[i]), replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzyme.RecognitionSequence, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
	}

	for i, _ := range enzyme.RecognitionSequence {

		replacementnucleotide := Comp(enzyme.RecognitionSequence[i : i+1])
		replacement := strings.Replace(enzyme.RecognitionSequence, enzyme.RecognitionSequence[i:i+1], replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzyme.RecognitionSequence, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
		if i+2 == len(enzyme.RecognitionSequence) {
			specificseqs := text.Print("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
			break
		}
	}
	return
}

// this replaces all instances but this is not what we want
func ReplaceString(sequence string, seq string, otherseqstoavoid []string) (newseq string, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, seq)

	for i, _ := range seq {

		replacementnucleotide := Comp(string(seq[i]))
		replacement := strings.Replace(seq, string(seq[i]), replacementnucleotide, 1)
		newseq = strings.Replace(sequence, seq, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
	}

	for i, _ := range seq {

		replacementnucleotide := Comp(seq[i : i+1])
		replacement := strings.Replace(seq, seq[i:i+1], replacementnucleotide, 1)
		newseq = strings.Replace(sequence, seq, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
		if i+2 == len(seq) {
			specificseqs := text.Print("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
			break
		}
	}
	return
}

/*
// working on this
func RemoveSiteFromORF(orf ORF, enzyme wtype.LogicalRestrictionEnzyme, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzyme.RecognitionSequence)

	sitesfound := search.Findallthings(orf.DNASeq, allthingstoavoid)

	if len(sitesfound) == 0 {
		err = fmt.Errorf("no sites found in this Orf!", orf, enzyme, otherseqstoavoid)
		return
	}

	allpositions := search.Findall(orf.DNASeq, enzyme.RecognitionSequence)

	if len(allpositions) != 0 {
		for _, position := range allpositions {
			if orf.Direction != "Reverse" {

				_, _ = Codonfromposition(orf.DNASeq, position)

			}
		}
	}
	return
}
*/
/*
func RemoveSiteFromSeq(annotated AnnotatedSeq, enzyme wtype.LogicalRestrictionEnzyme, otherseqstoavoid []string) (newseq AnnotatedSeq, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzyme.RecognitionSequence)

	sitesfound := search.Findallthings(orf.DNASeq, allthingstoavoid)

	if len(sitesfound) == 0 {
		err = fmt.Errorf("no sites found in this Orf!", orf, enzyme, otherseqstoavoid)
		return
	}

	allpositions := search.Findall(orf.DNASeq, enzyme.RecognitionSequence)

	if len(allpositions) != 0 {
		for _, position := range allpositions {
			if orf.Direction != "Reverse" {

				_ = Codonfromposition(orf.DNASeq, position)

			}
		}
	}
	return
}
*/
func RemoveSitesOutsideofFeatures(dnaseq wtype.DNASequence, site string, algorithm ReplacementAlgorithm, featurelisttoavoid []Feature) (newseq wtype.DNASequence, err error) {

	newseq = dnaseq

	pairs := make([]StartEndPair, 0)
	var pair StartEndPair

	for _, feature := range featurelisttoavoid {
		pair[0] = feature.StartPosition
		pair[1] = feature.EndPosition
		pairs = append(pairs, pair)
	}

	var otherseqstoavoid = []string{}

	replacement, err := algorithm(dnaseq.Seq, site, otherseqstoavoid)
	if err != nil {
		panic("choose different replacement choice func or change parameters")
	}

	newseq.Seq = ReplaceAvoidingPositionPairs(dnaseq.Seq, pairs, site, replacement)

	return
}

func ReplacefrombetweenPositions(seq string, start int, end int, original string, replacement string) (newseq string) {

	newseq = strings.Replace(seq[start-1:end-1], original, replacement, -1)
	return
}

func ReplaceAvoidingPositionPairs(seq string, positionpairs []StartEndPair, original string, replacement string) (newseq string) {

	temp := "£££££££££££"
	newseq = ""
	for _, pair := range positionpairs {
		if pair[1] < pair[0] {
			newseq = strings.Replace(seq[pair[0]-1:pair[1]-1], original, temp, -1)
		}
	}

	newseq = strings.Replace(newseq, original, replacement, -1)

	newseq = strings.Replace(newseq, temp, original, -1)

	// now look for reverse
	for _, pair := range positionpairs {
		if pair[1] > pair[0] {
			newseq = strings.Replace(seq[pair[1]-1:pair[0]-1], RevComp(original), temp, -1)
		}
	}

	newseq = strings.Replace(newseq, RevComp(original), RevComp(replacement), -1)

	newseq = strings.Replace(newseq, temp, RevComp(original), -1)
	return
}

type StartEndPair []int

func MakeStartendPair(start, end int) (pair StartEndPair) {

	pair[0] = start
	pair[1] = end
	return
}

func AAPosition(dnaposition int) (aaposition int) {

	remainder := dnaposition % 3
	aaposition = wutil.RoundInt(float64(dnaposition/3) + float64(remainder/3))

	return
}

func CodonOptions(codon string) (replacementoptions []string) {
	replacementoptions = RevCodonTable[codon]
	return
}

func SwapCodon(codon string, position int) (replacement string) {

	replacementarray := RevCodonTable[codon]

	replacement = replacementarray[position]

	return
}

func ReplaceCodoninORF(sequence wtype.DNASequence, startandendoforf StartEndPair, position int, seqstoavoid []string) (newseq wtype.DNASequence, err error) {

	// only handling cases where orf is not in reverse strand currently
	if startandendoforf[0] < startandendoforf[1] {
		seqslice := sequence.Seq[startandendoforf[0]-1 : startandendoforf[1]-1]
		orf, orftrue := FindORF(seqslice)
		if orftrue && len(orf.DNASeq) == len(seqslice) {
			codontochange, pair := Codonfromposition(orf.DNASeq, (position - startandendoforf[0] - 1))

			options := CodonOptions(codontochange)
			for _, option := range options {
				tempseq := ReplacePosition(sequence.Seq, pair, option)
				seqslice := tempseq[startandendoforf[0]-1 : startandendoforf[1]-1]
				temporf, _ := FindORF(seqslice)

				sitesfound := search.Findallthings(tempseq, seqstoavoid)

				if temporf.ProtSeq == orf.ProtSeq && len(sitesfound) == 0 {
					newseq := sequence
					newseq.Seq = tempseq
					return newseq, err
				}

			}
		}
	}
	return
}

func ReplacePosition(sequence string, position StartEndPair, replacement string) (newseq string) {

	if position[0] < position[1] {
		one := sequence[0 : position[0]-1]
		_ = sequence[position[0] : position[1]-1]
		three := sequence[position[1]:]

		newseq = one + replacement + three
	}
	return
}

func Codonfromposition(sequence string, dnaposition int) (codontoreturn string, position StartEndPair) {

	nucleotides := []rune(sequence)
	//fmt.Println("codons=", string(codons))
	res := ""
	aas := make([]string, 0)
	codon := ""
	for i, r := range nucleotides {
		res = res + string(r)
		//fmt.Printf("i%d r %c\n", i, r)

		if i > 0 && (i+1)%3 == 0 {
			//fmt.Printf("=>(%d) '%v'\n", i, res)
			codon = res
			aas = append(aas, res)
			res = ""
		}
		if i+1 > dnaposition && i > 0 && (i+1)%3 == 0 {
			codontoreturn = codon
			position[1] = i + 1
			position[0] = i - 2
			return
		}
	}
	return
}
