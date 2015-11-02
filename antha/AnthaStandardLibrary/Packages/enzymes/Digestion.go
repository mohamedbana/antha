// antha/AnthaStandardLibrary/Packages/enzymes/Digestion.go: Part of the Antha language
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
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"sort"
	"strconv"
	"strings"
)

//should expand to be more general, i.e. 3prime overhangs
type DoublestrandedDNA struct {
	Fwdsequence           wtype.DNASequence
	Reversesequence       wtype.DNASequence
	TopStickyend5prime    string
	Bottomstickyend5prime string
	Phosphorylated        bool
}

func MakedoublestrandedDNA(sequence wtype.DNASequence) (Doublestrandedpair []wtype.DNASequence) {
	fwdsequence := strings.TrimSpace(strings.ToUpper(sequence.Seq))
	revcomp := RevComp(fwdsequence)
	reversesequence := strings.TrimSpace(strings.ToUpper(revcomp))

	var Fwdsequence = wtype.DNASequence{Nm: "Fwdsequence", Seq: fwdsequence}
	var Reversesequence = wtype.DNASequence{Nm: "Reversecomplement", Seq: reversesequence}
	if sequence.Plasmid == true {
		Fwdsequence.Plasmid = true
		Reversesequence.Plasmid = true
	}
	Doublestrandedpair = []wtype.DNASequence{Fwdsequence, Reversesequence}
	return Doublestrandedpair
}

type Restrictionsites struct {
	Enzyme              wtype.LogicalRestrictionEnzyme
	Recognitionsequence string
	Sitefound           bool
	Numberofsites       int
	Forwardpositions    []int
	Reversepositions    []int
}

func SitepositionString(sitesperpart Restrictionsites) (sitepositions string) {
	Num := make([]string, 0)

	for _, site := range sitesperpart.Forwardpositions {
		Num = append(Num, strconv.Itoa(site))
	}
	for _, site := range sitesperpart.Reversepositions {
		Num = append(Num, strconv.Itoa(site))
	}
	sitepositions = strings.Join(Num, ", ")
	return
}

func Restrictionsitefinder(sequence wtype.DNASequence, enzymelist []wtype.LogicalRestrictionEnzyme) (sites []Restrictionsites) {

	sites = make([]Restrictionsites, 0)

	for _, enzyme := range enzymelist {
		var enzymesite Restrictionsites
		//var siteafterwobble Restrictionsites
		enzymesite.Enzyme = enzyme
		enzymesite.Recognitionsequence = strings.ToUpper(enzyme.RecognitionSequence)
		sequence.Seq = strings.ToUpper(sequence.Seq)

		wobbleproofrecognitionoptions := Wobble(enzymesite.Recognitionsequence)

		for _, wobbleoption := range wobbleproofrecognitionoptions {

			options := Findall(sequence.Seq, wobbleoption)
			for _, option := range options {
				if option != 0 {
					enzymesite.Forwardpositions = append(enzymesite.Forwardpositions, option)
				}
			}
			if enzyme.RecognitionSequence != strings.ToUpper(RevComp(wobbleoption)) {
				revoptions := Findall(sequence.Seq, RevComp(wobbleoption))
				for _, option := range revoptions {
					if option != 0 {
						enzymesite.Reversepositions = append(enzymesite.Reversepositions, option)
					}
				}

			}
			enzymesite.Numberofsites = len(enzymesite.Forwardpositions) + len(enzymesite.Reversepositions)
			if enzymesite.Numberofsites > 0 {
				enzymesite.Sitefound = true
			}

		}

		sites = append(sites, enzymesite)
	}

	return sites
}

/*
func CutatSite(startingdnaseq wtype.DNASequence, typeIIenzyme wtype.LogicalRestrictionEnzyme) (Digestproducts []wtype.DNASequence) {
	// not tested and not finished

	Digestproducts = make([]wtype.DNASequence, 0)
	originalfwdsequence := strings.ToUpper(startingdnaseq.Seq)

	recogseq := strings.ToUpper(typeIIenzyme.RecognitionSequence)
	sites := Restrictionsitefinder(startingdnaseq, []wtype.LogicalRestrictionEnzyme{typeIIenzyme})

	if len(sites) == 0 {
		Digestproducts = append(Digestproducts, startingdnaseq)
	} else {
		for _, site := range sites {

			fragments := make([]string, 0)
			fragment := ""
			for i, position := range site.forwardpositions {
				if i == 0 {
					fragment = originalfwdsequence[0:position]
				} else {
					fragment = originalfwdsequence[site.forwardpositions[i-1]:site.forwardpositions[i]]
					fragments = append(fragments, fragment)
				}
			}
			for i, fragment := range fragments {
				//not tested
				cutup := ""
				cutdown := ""
				if typeIIenzyme.Class == "TypeII" {
					if i != 0 {
						cutup = Prefix(fragment, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))
					}
					if i != len(fragments) {
						cutdown = Prefix(fragment, (len(recogseq)))
						cutdown = Suffix(cutdown, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
					}
					fragment = cutup + fragment + cutdown
				} else if typeIIenzyme.Class == "TypeIIs" {
					if i != 0 {
						fragment = Suffix(fragment, len(fragment)-(len(recogseq)+typeIIenzyme.Topstrand3primedistancefromend)) //cutdown = suffix(cutdown,(-1 * typeIIenzyme.Topstrand3primedistancefromend))
					}
					if i != len(fragments) {
						cutdown = Prefix(fragment, (len(recogseq) + typeIIenzyme.Topstrand3primedistancefromend))
					}
					fragment = fragment + cutdown

				}

			}


			var digestproduct wtype.DNASequence
			for i, frag := range fragments {
				digestproduct.Nm = startingdnaseq.Nm + "fragment" + strconv.Itoa(i)
				digestproduct.Seq = frag

				//digestproduct.Overhang5prime = Overhang{5, 2}
				//digestproduct.Overhang3prime = Overhang{3, -1}
				Digestproducts = append(Digestproducts, digestproduct)
			}
		}
	}
	return
}
*/

type Digestedfragment struct {
	Topstrand              string
	Bottomstrand           string
	TopStickyend_5prime    string
	TopStickyend_3prime    string
	BottomStickyend_5prime string
	BottomStickyend_3prime string
}

func Pair(digestedtopstrand []string, digestedbottomstrand []string, topstickyend5prime []string, topstickyend3prime []string, bottomstickyend5prime []string, bottomstickyend3prime []string) (pairs []Digestedfragment) {

	pairs = make([]Digestedfragment, 0)

	var pair Digestedfragment

	if len(digestedtopstrand) == len(digestedbottomstrand) { //}|| len(topstickyend5prime) || len(topstickyend3prime) {
		for i := 0; i < len(digestedtopstrand); i++ {
			pair.Topstrand = digestedtopstrand[i]
			pair.Bottomstrand = digestedbottomstrand[i]
			pair.TopStickyend_5prime = topstickyend5prime[i]
			pair.TopStickyend_3prime = topstickyend3prime[i]
			pair.BottomStickyend_5prime = bottomstickyend5prime[i]
			pair.BottomStickyend_3prime = bottomstickyend3prime[i]
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

func DigestionPairs(Doublestrandedpair []wtype.DNASequence, typeIIsenzyme TypeIIs) (digestionproducts []Digestedfragment) {
	topstrands, topstickyends5, topstickyends3 := TypeIIsdigest(Doublestrandedpair[0], typeIIsenzyme)
	bottomstrands, bottomstickyends5, bottomstickyends3 := TypeIIsdigest(Doublestrandedpair[1], typeIIsenzyme)
	if len(topstrands) == len(bottomstrands) {
		if len(topstrands) == 2 {
			digestionproducts = Pair(topstrands, bottomstrands, topstickyends5, topstickyends3, bottomstickyends5, bottomstickyends3)
		}
		if len(topstrands) == 3 {
			digestionproducts = Pair(topstrands, Revarrayorder(bottomstrands), topstickyends5, topstickyends3, Revarrayorder(bottomstickyends5), Revarrayorder(bottomstickyends3))
		}
	}
	return digestionproducts
}

func Digest(sequence wtype.DNASequence, typeIIenzyme wtype.LogicalRestrictionEnzyme) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {
	if typeIIenzyme.Class == "TypeII" {
		Finalfragments, Stickyends_5prime, Stickyends_3prime = TypeIIDigest(sequence, typeIIenzyme)
	}
	if typeIIenzyme.Class == "TypeIIs" {

		var isoschizomers = make([]string, 0)
		/*for _, lookup := range ...
		add code to lookup isoschizers from rebase
		*/
		var typeIIsenz = TypeIIs{typeIIenzyme, typeIIenzyme.Name, isoschizomers, typeIIenzyme.Topstrand3primedistancefromend, typeIIenzyme.Bottomstrand5primedistancefromend}

		Finalfragments, Stickyends_5prime, Stickyends_3prime = TypeIIsdigest(sequence, typeIIsenz)
	}
	return
}

func RestrictionMapper(seq wtype.DNASequence, enzyme wtype.LogicalRestrictionEnzyme) (fraglengths []int) {
	enzlist := []wtype.LogicalRestrictionEnzyme{enzyme}
	frags, _, _ := Digest(seq, enzlist[0]) // doesn't handle non cutters well - returns 1 seq string, blunt, blunt therefore inaccurate representation
	fraglengths = make([]int, 0)
	for _, frag := range frags {
		fraglengths = append(fraglengths, len(frag))
	}
	fragslice := sort.IntSlice(fraglengths)
	fragslice.Sort()

	return fraglengths
}
func SearchandCut(typeIIenzyme wtype.LogicalRestrictionEnzyme, topstranddigestproducts []string, topstrandstickyends_5prime []string, topstrandstickyends_3prime []string) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {
	finalfragments, topstrandstickyends_5primeFW, topstrandstickyends_3primeFW :=
		SearchandCutFWD(typeIIenzyme, topstranddigestproducts, topstrandstickyends_5prime, topstrandstickyends_3prime)

	Finalfragments, Stickyends_5prime, Stickyends_3prime = SearchandCutRev(typeIIenzyme, finalfragments, topstrandstickyends_5primeFW, topstrandstickyends_3primeFW)
	return
}

func SearchandCutFWD(typeIIenzyme wtype.LogicalRestrictionEnzyme, topstranddigestproducts []string, topstrandstickyends_5prime []string, topstrandstickyends_3prime []string) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {

	Finalfragments = make([]string, 0)

	originalfwdsequence := strings.ToUpper(strings.Join(topstranddigestproducts, ""))
	recogseq := strings.ToUpper(typeIIenzyme.RecognitionSequence)
	sites := Findall(originalfwdsequence, recogseq)
	// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
	if len(sites) == 0 {
		Finalfragments = topstranddigestproducts
		Stickyends_5prime = topstrandstickyends_5prime
		Stickyends_3prime = topstrandstickyends_3prime
	} else {
		finaldigestproducts := make([]string, 0)
		finaltopstrandstickyends_5prime := make([]string, 0)
		finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, "blunt")
		finaltopstrandstickyends_3prime := make([]string, 0)
		for _, fragment := range topstranddigestproducts {
			cuttopstrand := strings.Split(fragment, recogseq)
			// reversed
			recognitionsiteup := Prefix(recogseq, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))
			recognitionsitedown := Suffix(recogseq, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
			firstfrag := strings.Join([]string{cuttopstrand[0], recognitionsiteup}, "")
			finaldigestproducts = append(finaldigestproducts, firstfrag)

			for i := 1; i < len(cuttopstrand); i++ {
				joineddownstream := strings.Join([]string{recognitionsitedown, cuttopstrand[i]}, "")
				if i != len(cuttopstrand)-1 {
					joineddownstream = strings.Join([]string{joineddownstream, recognitionsiteup}, "")
				}
				finaldigestproducts = append(finaldigestproducts, joineddownstream)
			}
			frag2topStickyend5prime := ""
			frag2topStickyend3prime := ""
			// cut with 5prime overhang
			if len(recognitionsitedown) > len(recognitionsiteup) {

				for i := 1; i < len(cuttopstrand); i++ {
					frag2topStickyend5prime = recognitionsitedown[:typeIIenzyme.EndLength]
					finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
					frag2topStickyend3prime = ""
					finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)

				}

			}
			// blunt cut
			if len(recognitionsitedown) == len(recognitionsiteup) {
				for i := 1; i < len(cuttopstrand); i++ {
					frag2topStickyend5prime = "blunt"
					finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
					frag2topStickyend3prime = "blunt"
					finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)
				}
			}
			// cut with 3prime overhang
			if len(recognitionsitedown) < len(recognitionsiteup) {
				for i := 1; i < len(cuttopstrand); i++ {
					frag2topStickyend5prime = ""
					finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
					frag2topStickyend3prime = recognitionsiteup[typeIIenzyme.EndLength:]
					finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)
				}
			}
		}
		finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, "blunt")
		Finalfragments = finaldigestproducts
		Stickyends_5prime = finaltopstrandstickyends_5prime
		Stickyends_3prime = finaltopstrandstickyends_3prime
	}
	return
}

func SearchandCutRev(typeIIenzyme wtype.LogicalRestrictionEnzyme, topstranddigestproducts []string, topstrandstickyends_5prime []string, topstrandstickyends_3prime []string) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {
	Finalfragments = make([]string, 0)
	reverseenzymeseq := RevComp(strings.ToUpper(typeIIenzyme.RecognitionSequence))

	if reverseenzymeseq == strings.ToUpper(typeIIenzyme.RecognitionSequence) {
		Finalfragments = topstranddigestproducts
		Stickyends_5prime = topstrandstickyends_5prime
		Stickyends_3prime = topstrandstickyends_3prime
	} else {
		originalfwdsequence := strings.Join(topstranddigestproducts, "")
		sites := Findall(originalfwdsequence, reverseenzymeseq)
		// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
		if len(sites) == 0 {
			Finalfragments = topstranddigestproducts
		} else {
			finaldigestproducts := make([]string, 0)
			finaltopstrandstickyends_5prime := make([]string, 0)
			finaltopstrandstickyends_3prime := make([]string, 0)
			for _, fragment := range topstranddigestproducts {
				cuttopstrand := strings.Split(fragment, reverseenzymeseq)
				// reversed
				recognitionsiteup := Prefix(reverseenzymeseq, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))
				recognitionsitedown := Suffix(reverseenzymeseq, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
				firstfrag := strings.Join([]string{cuttopstrand[0], recognitionsiteup}, "")
				finaldigestproducts = append(finaldigestproducts, firstfrag)
				for i := 1; i < len(cuttopstrand); i++ {
					joineddownstream := strings.Join([]string{recognitionsitedown, cuttopstrand[i]}, "")
					if i != len(cuttopstrand)-1 {
						joineddownstream = strings.Join([]string{joineddownstream, recognitionsiteup}, "")
					}
					finaldigestproducts = append(finaldigestproducts, joineddownstream)
				}
				frag2topStickyend5prime := ""
				frag2topStickyend3prime := ""
				// cut with 5prime overhang
				if len(recognitionsitedown) > len(recognitionsiteup) {
					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = recognitionsitedown[:typeIIenzyme.EndLength]
						finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
						if i != len(cuttopstrand)-1 {
							frag2topStickyend3prime = ""
						} else {
							frag2topStickyend3prime = "blunt"
						}
						finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)
					}
				}
				// blunt cut
				if len(recognitionsitedown) == len(recognitionsiteup) {
					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = "blunt"
						finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
						frag2topStickyend3prime = "blunt"
						finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)
					}
				}
				// cut with 3prime overhang
				if len(recognitionsitedown) < len(recognitionsiteup) {

					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = ""
						finaltopstrandstickyends_5prime = append(finaltopstrandstickyends_5prime, frag2topStickyend5prime)
						if i != len(cuttopstrand)-1 {
							frag2topStickyend3prime = recognitionsiteup[typeIIenzyme.EndLength:]
						} else {
							frag2topStickyend3prime = "blunt"
						}
						finaltopstrandstickyends_3prime = append(finaltopstrandstickyends_3prime, frag2topStickyend3prime)
					}
				}
				for _, strand5 := range finaltopstrandstickyends_5prime {
					topstrandstickyends_5prime = append(topstrandstickyends_5prime, strand5)
				}
				for _, strand3 := range finaltopstrandstickyends_3prime {
					topstrandstickyends_3prime = append(topstrandstickyends_3prime, strand3)
				}
				Finalfragments = finaldigestproducts
				Stickyends_5prime = topstrandstickyends_5prime
				Stickyends_3prime = topstrandstickyends_3prime
			}
		}
	}
	return
}

func LineartoPlasmid(fragmentsiflinearstart []string) (fragmentsifplasmidstart []string) {

	// make linear plasmid part by joining last part to first part
	plasmidcutproducts := make([]string, 0)
	plasmidcutproducts = append(plasmidcutproducts, fragmentsiflinearstart[len(fragmentsiflinearstart)-1])
	plasmidcutproducts = append(plasmidcutproducts, fragmentsiflinearstart[0])
	linearpartfromplasmid := strings.Join(plasmidcutproducts, "")

	// fix order of final fragments
	fragmentsifplasmidstart = make([]string, 0)
	fragmentsifplasmidstart = append(fragmentsifplasmidstart, linearpartfromplasmid)
	for i := 1; i < (len(fragmentsiflinearstart) - 1); i++ {
		fragmentsifplasmidstart = append(fragmentsifplasmidstart, fragmentsiflinearstart[i])
	}

	return
}

func LineartoPlasmidEnds(endsiflinearstart []string) (endsifplasmidstart []string) {

	endsifplasmidstart = make([]string, 0)

	endsifplasmidstart = append(endsifplasmidstart, endsiflinearstart[len(endsiflinearstart)-1])

	for i := 1; i < (len(endsiflinearstart)); i++ {
		endsifplasmidstart = append(endsifplasmidstart, endsiflinearstart[i])

	}

	return
}

func TypeIIDigest(sequence wtype.DNASequence, typeIIenzyme wtype.LogicalRestrictionEnzyme) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {
	// step 1. get sequence in string format from DNASequence, make sure all spaces are removed and all upper case

	if typeIIenzyme.Class != "TypeII" {
		panic("This is not the function you are looking for! Wrong enzyme class for this function")
	}

	originalfwdsequence := strings.TrimSpace(strings.ToUpper(sequence.Seq))
	//originalreversesequence := strings.TrimSpace(strings.ToUpper(RevComp(sequence.Seq)))
	sites := Findall(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence))

	// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
	topstranddigestproducts := make([]string, 0)
	topstrandstickyends_5prime := make([]string, 0)
	topstrandstickyends_3prime := make([]string, 0)

	if len(sites) != 0 {

		cuttopstrand := strings.Split(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence))
		recognitionsitedown := Suffix(typeIIenzyme.RecognitionSequence, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
		recognitionsiteup := Prefix(typeIIenzyme.RecognitionSequence, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))

		//repairedfrag := ""
		//repairedfrags := make([]string,0)

		//if sequence.Plasmid != true{

		firstfrag := strings.Join([]string{cuttopstrand[0], recognitionsiteup}, "")
		topstranddigestproducts = append(topstranddigestproducts, firstfrag)

		for i := 1; i < len(cuttopstrand); i++ {
			joineddownstream := strings.Join([]string{recognitionsitedown, cuttopstrand[i]}, "")
			if i != len(cuttopstrand)-1 {
				joineddownstream = strings.Join([]string{joineddownstream, recognitionsiteup}, "")
			}
			topstranddigestproducts = append(topstranddigestproducts, joineddownstream)

		}

		frag2topStickyend5prime := ""
		frag2topStickyend3prime := ""
		// cut with 5prime overhang
		if len(recognitionsitedown) > len(recognitionsiteup) {
			frag2topStickyend5prime = "blunt"
			topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
			frag2topStickyend3prime := ""
			topstrandstickyends_3prime = append(topstrandstickyends_3prime, frag2topStickyend3prime)
			for i := 1; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = recognitionsitedown[:typeIIenzyme.EndLength]
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
				if i != len(cuttopstrand)-1 {
					frag2topStickyend3prime = ""
				} else {
					frag2topStickyend3prime = "blunt"
				}
				topstrandstickyends_3prime = append(topstrandstickyends_3prime, frag2topStickyend3prime)

			}

		}
		// blunt cut
		if len(recognitionsitedown) == len(recognitionsiteup) {
			for i := 0; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = "blunt"
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
				frag2topStickyend3prime = "blunt"
				topstrandstickyends_3prime = append(topstrandstickyends_3prime, frag2topStickyend3prime)
			}
		}
		// cut with 3prime overhang
		if len(recognitionsitedown) < len(recognitionsiteup) {
			frag2topStickyend5prime = "blunt"
			topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)

			frag2topStickyend3prime = Suffix(recognitionsiteup, typeIIenzyme.EndLength)
			topstrandstickyends_3prime = append(topstrandstickyends_3prime, frag2topStickyend3prime)

			for i := 1; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = ""
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
				if i != len(cuttopstrand)-1 {
					frag2topStickyend3prime = recognitionsiteup[typeIIenzyme.EndLength:]
				} else {
					frag2topStickyend3prime = "blunt"
				}
				topstrandstickyends_3prime = append(topstrandstickyends_3prime, frag2topStickyend3prime)

			}
		}
	} else {
		topstranddigestproducts = []string{originalfwdsequence}
		topstrandstickyends_5prime = []string{"blunt"}
		topstrandstickyends_3prime = []string{"blunt"}
	}

	Finalfragments, topstrandstickyends_5prime, topstrandstickyends_3prime = SearchandCutRev(typeIIenzyme, topstranddigestproducts, topstrandstickyends_5prime, topstrandstickyends_3prime)

	if len(Finalfragments) == 1 && sequence.Plasmid == true {
		// TODO
		// need to really return an uncut plasmid, maybe an error?
		//	fmt.Println("uncut plasmid returned with no sticky ends!")

	}
	if len(Finalfragments) > 1 && sequence.Plasmid == true {
		ifplasmidfinalfragments := LineartoPlasmid(Finalfragments)
		Finalfragments = ifplasmidfinalfragments
		// now change order of sticky ends
		//5'
		ifplasmidsticky5prime := make([]string, 0)
		ifplasmidsticky5prime = append(ifplasmidsticky5prime, topstrandstickyends_5prime[len(topstrandstickyends_5prime)-1])
		for i := 1; i < (len(Finalfragments)); i++ {
			ifplasmidsticky5prime = append(ifplasmidsticky5prime, topstrandstickyends_5prime[i])
		}
		topstrandstickyends_5prime = ifplasmidsticky5prime
		//hack to fix wrong sticky end assignment in certain cases
		reverseenzymeseq := RevComp(typeIIenzyme.RecognitionSequence)
		if strings.Index(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence)) > strings.Index(originalfwdsequence, reverseenzymeseq) {
			topstrandstickyends_5prime = Revarrayorder(topstrandstickyends_5prime)
		}
		//3'
		ifplasmidsticky3prime := make([]string, 0)
		ifplasmidsticky3prime = append(ifplasmidsticky3prime, topstrandstickyends_3prime[0])
		for i := 1; i < (len(Finalfragments)); i++ {
			ifplasmidsticky3prime = append(ifplasmidsticky3prime, topstrandstickyends_3prime[i])
		}
		topstrandstickyends_3prime = ifplasmidsticky3prime
	}
	Stickyends_5prime = topstrandstickyends_5prime
	// deal with this later
	Stickyends_3prime = topstrandstickyends_3prime
	return Finalfragments, Stickyends_5prime, Stickyends_3prime
}

// A function is called by the first word (note the capital letter!); it takes in the input variables in the first parenthesis and returns the contents of the second parenthesis
// currently this doesn't work well for plasmids which are cut on reverse strand or cut twice
func TypeIIsdigest(sequence wtype.DNASequence, typeIIsenzyme TypeIIs) (Finalfragments []string, Stickyends_5prime []string, Stickyends_3prime []string) {
	if typeIIsenzyme.Class != "TypeIIs" {
		return Finalfragments, Stickyends_5prime, Stickyends_3prime
	}
	// step 1. get sequence in string format from DNASequence, make sure all spaces are removed and all upper case
	originalfwdsequence := strings.TrimSpace(strings.ToUpper(sequence.Seq))

	// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
	topstranddigestproducts := make([]string, 0)
	topstrandstickyends_5prime := make([]string, 0)
	topstrandstickyends_3prime := make([]string, 0)
	if strings.Contains(originalfwdsequence, strings.ToUpper(typeIIsenzyme.LogicalRestrictionEnzyme.RecognitionSequence)) == false {
		topstranddigestproducts = append(topstranddigestproducts, originalfwdsequence)
		topstrandstickyends_5prime = append(topstrandstickyends_5prime, "blunt")
	} else {
		// step 3. split the sequence (into an array of daughter seqs) after the recognition site! Note! this is a preliminary step, we'll fix the sequence to reflect reality in subsequent steps
		cuttopstrand := strings.SplitAfter(originalfwdsequence, strings.ToUpper(typeIIsenzyme.LogicalRestrictionEnzyme.RecognitionSequence))
		// step 4. If this results in only 2 fragments (i.e. only one site in upper strand) it means we can continue. We can add the ability to handle multiple sites later!
		// add boolean for direction of cut (i.e. need to use different strategy for 3' or 5')
		if len(cuttopstrand) == 2 {
			// step 5. name the two fragments
			frag1 := cuttopstrand[0]
			frag2 := cuttopstrand[1]
			// step 6. find the length of the downstream fragment
			sz := len(frag2)
			// step 7. remove extra base pairs from downstream fragment according to typeIIs enzyme properties (i.e. N bp downstream (or 3') of recognition site e.g. in the case of SapI it cuts 1bp 3' to the recognittion site on the top strand
			Cuttop2 := frag2[sz-(sz-typeIIsenzyme.Topstrand3primedistancefromend):]
			// step 8. then add these extra base pairs to the 3' end of upstream fragment; first we find the base pairs
			bittoaddtopriorsequence := frag2[:sz-(sz-typeIIsenzyme.Topstrand3primedistancefromend)]
			// step 9. Now we join back together
			firstsequenceparts := make([]string, 0)
			firstsequenceparts = append(firstsequenceparts, frag1)
			firstsequenceparts = append(firstsequenceparts, bittoaddtopriorsequence)
			joinedfirstpart := strings.Join(firstsequenceparts, "")
			// for use in sticky end caclulation later (added here to be before if statements
			frag2topStickyend5prime := Cuttop2[:sz-(sz-(typeIIsenzyme.Bottomstrand5primedistancefromend-typeIIsenzyme.Topstrand3primedistancefromend))]
			// step 10. Now we bundle them back up again into an array to access later
			topstranddigestproducts = append(topstranddigestproducts, joinedfirstpart)
			topstranddigestproducts = append(topstranddigestproducts, Cuttop2)
			// now for sticky ends
			// add nothing as sticky end for fragment 1
			topstrandstickyends_5prime = append(topstrandstickyends_5prime, "blunt")
			if len(frag2topStickyend5prime) == 0 {
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, "blunt")
			}
			if len(frag2topStickyend5prime) > 0 {
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
			}
			//no 3' sticky ends in this case...
		}
	}
	Finalfragments = make([]string, 0)

	reverseenzymeseq := RevComp(strings.ToUpper(typeIIsenzyme.LogicalRestrictionEnzyme.RecognitionSequence))

	for _, digestedfragment := range topstranddigestproducts {
		if strings.Contains(digestedfragment, reverseenzymeseq) == false {
			Finalfragments = append(Finalfragments, digestedfragment)
			topstrandstickyends_3prime = append(topstrandstickyends_3prime, "")
		} else {
			cuttopstrandat3prime := strings.Split(digestedfragment, reverseenzymeseq)
			if len(cuttopstrandat3prime) < 3 || len(cuttopstrandat3prime) > 0 {
				// step 5. name the two fragments
				frag1 := cuttopstrandat3prime[0]
				frag2 := cuttopstrandat3prime[1]
				// step 6. find the length of the upstream fragment
				new_sz := len(frag1)
				// step 7. remove extra base pairs 3 from upstream fragment according to typeIIs enzyme properties (i.e. N bp upstream (or 5') of reverserecognition site e.g. in the case of SapI it cuts 4bp 5' to the recognittion site on the top strand (since reverse comp)
				//s = s[:sz-1]
				Cuttop3 := frag1[:new_sz-(typeIIsenzyme.Bottomstrand5primedistancefromend)]
				// step 8. then add these extra base pairs to the 3' end of upstream fragment; first we find the base pairs
				bittoaddtopostsequence := frag1[new_sz-(typeIIsenzyme.Bottomstrand5primedistancefromend):]
				// step 9. Now we join back together
				step2sequenceparts := make([]string, 0)
				step2sequenceparts = append(step2sequenceparts, bittoaddtopostsequence)
				step2sequenceparts = append(step2sequenceparts, reverseenzymeseq)
				step2sequenceparts = append(step2sequenceparts, frag2)
				joinedsecondpart := strings.Join(step2sequenceparts, "")
				//bitlength := len(bittoaddtopostsequence)
				frag2topStickyend5prime := bittoaddtopostsequence[:(typeIIsenzyme.Bottomstrand5primedistancefromend - typeIIsenzyme.Topstrand3primedistancefromend)]

				// step 10. Now we bundle them back up again into an array to access later
				Finalfragments = append(Finalfragments, Cuttop3)
				Finalfragments = append(Finalfragments, joinedsecondpart)
				topstrandstickyends_5prime = append(topstrandstickyends_5prime, frag2topStickyend5prime)
				topstrandstickyends_3prime = append(topstrandstickyends_3prime, "")

				topstrandstickyends_3prime = append(topstrandstickyends_3prime, "")
				// step 12. we then return this!
			}
		}
	}

	if len(Finalfragments) == 1 && sequence.Plasmid == true {
		// TODO
		// need to really return an uncut plasmid, maybe an error?
		//	fmt.Println("uncut plasmid returned with no sticky ends!")

	}
	if len(Finalfragments) > 1 && sequence.Plasmid == true {

		// make linear plasmid part
		plasmidcutproducts := make([]string, 0)
		plasmidcutproducts = append(plasmidcutproducts, Finalfragments[len(Finalfragments)-1])
		plasmidcutproducts = append(plasmidcutproducts, Finalfragments[0])
		linearpartfromplasmid := strings.Join(plasmidcutproducts, "")

		// fix order of final fragments
		ifplasmidfinalfragments := make([]string, 0)
		ifplasmidfinalfragments = append(ifplasmidfinalfragments, linearpartfromplasmid)
		for i := 1; i < (len(Finalfragments) - 1); i++ {
			ifplasmidfinalfragments = append(ifplasmidfinalfragments, Finalfragments[i])
		}

		Finalfragments = ifplasmidfinalfragments

		// now change order of sticky ends
		//5'
		ifplasmidsticky5prime := make([]string, 0)

		ifplasmidsticky5prime = append(ifplasmidsticky5prime, topstrandstickyends_5prime[len(topstrandstickyends_5prime)-1])

		for i := 1; i < (len(Finalfragments)); i++ {
			ifplasmidsticky5prime = append(ifplasmidsticky5prime, topstrandstickyends_5prime[i])
		}
		topstrandstickyends_5prime = ifplasmidsticky5prime

		//hack to fix wrong sticky end assignment in certain cases
		if strings.Index(originalfwdsequence, strings.ToUpper(typeIIsenzyme.LogicalRestrictionEnzyme.RecognitionSequence)) > strings.Index(originalfwdsequence, reverseenzymeseq) {
			topstrandstickyends_5prime = Revarrayorder(topstrandstickyends_5prime)
		}
		//3'
		ifplasmidsticky3prime := make([]string, 0)
		ifplasmidsticky3prime = append(ifplasmidsticky3prime, topstrandstickyends_3prime[0])
		for i := 1; i < (len(Finalfragments)); i++ {
			ifplasmidsticky3prime = append(ifplasmidsticky3prime, topstrandstickyends_3prime[i])
		}
		topstrandstickyends_3prime = ifplasmidsticky3prime
	}

	Stickyends_5prime = topstrandstickyends_5prime

	// deal with this later
	Stickyends_3prime = topstrandstickyends_3prime

	return Finalfragments, Stickyends_5prime, Stickyends_3prime
}

func Digestionsimulator(assemblyparameters Assemblyparameters) (digestedfragementarray [][]Digestedfragment) {
	// fetch enzyme properties from map (this is basically a look up table for those who don't know)
	digestedfragementarray = make([][]Digestedfragment, 0)
	enzymename := strings.ToUpper(assemblyparameters.Enzymename)
	enzyme := TypeIIsEnzymeproperties[enzymename]
	//assemble (note that sapIenz is found in package enzymes)
	doublestrandedvector := MakedoublestrandedDNA(assemblyparameters.Vector)
	digestedvector := DigestionPairs(doublestrandedvector, enzyme)
	digestedfragementarray = append(digestedfragementarray, digestedvector)
	for _, part := range assemblyparameters.Partsinorder {
		doublestrandedpart := MakedoublestrandedDNA(part)
		digestedpart := DigestionPairs(doublestrandedpart, enzyme)
		digestedfragementarray = append(digestedfragementarray, digestedpart)
	}
	return digestedfragementarray
}
