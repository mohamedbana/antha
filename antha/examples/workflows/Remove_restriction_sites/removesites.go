// This protocol is intended to check sequences for restriction sites and remove according to
// specified conditions

package main

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)
var (
	Sequencekey string = "agcgcgttaaaaattcttatttgagcaagagtatagaatgaggttgtagcttatacatgtgattttaggacattgacatggttgaaaatgccaaagactttgacataggttttttggacgataacgactgtttttaatatataagatggatgattatgagcagtatgaatagaggatagaaagaggagaaatactagatggcttcctccgaagacgttatcaaagagttcatgcgtttcaaagttcgtatggaaggttccgttaacggtcacgagttcgaaatcgaaggtgaaggtgaaggtcgtccgtacgaaggtacccagaccgctaaactgaaagttaccaaaggtggtccgctgccgttcgcttgggacatcctgtccccgcagttccagtacggttccaaagcttacgttaaacacccggctgacatcccggactacctgaaactgtccttcccggaaggtttcaaatgggaacgtgttatgaacttcgaagacggtggtgttgttaccgttacccaggactcctccctgcaagacggtgagttcatctacaaagttaaactgcgtggtaccaacttcccgtccgacggtccggttatgcagaaaaaaaccatgggttgggaagcttccaccgaacgtatgtacccggaagacggtgctctgaaaggtgaaatcaaaatgcgtctgaaactgaaagacggtggtcactacgacgctgaagttaaaaccacctacatggctaaaaaaccggttcagctgccgggtgcttacaaaaccgacatcaaactggacatcacctcccacaacgaagactacaccatcgttgaacagtacgaacgtgctgaaggtcgtcactccaccggtgcttaataaactagaccaggcatcaaataaaacgaaaggctcagtcgaaagactgggcctttcgttttatctgttgtttgtcggtgaacgctctctactagagtcacactggctcaccttcgggtgggcctttctgcgtttatatactagagctgatagctagctcagtcctagggattatgctagctactagagaaagaggagaaatactagatggatgcacttagtaaaatttttgatgatattcatttaaataaatctgaatatctctatgtaaaaactcaaggagaatgggcatttcatatgcttgcccaagatgctttgattgctcacattattttaatgggttctgctcatttcacattacaggatggtacgacacaaaccgcttattcgggagatattgtgcttattccttctgggcaagctcattttgtctcaaatgatgctcagaacaaactaattgattgttccaacattcagtctttttttgatggacaccgcaatgatgccattgagttaggcacaacatcttcggatcatggtttaatttttactgtgcggagccatatcgatatgcatatgggtagccctttactcaatgctttaccatcactgatacatattcatcacgcaatgagctcaactggcccagagtggttacgcgtcggtttatattttgtggcacttgaaacccagcgtattcaaccaggacgagataaaatatttgaccacttgatgagcatcctatttatagaatgtgtccgcgatcatattgcccagcttgatgaccccaaaagctggttaactgcactcatgcatcccgaattatctaatgcgctcgcagctattcatgcttaccccgaaaaaccttggactgtagagagtttggcagatcaatgttgcatgtcacggtctaaatttgccacactttttcaaagtattgttcatgagacccctctcgcttatcttcaacaacatcgccttcgtctggccgtacagttattaaaaaccagtcagttaaatattcagcaaattgccaataaagtcggatactcctcagaaacagcttttagtcaggcatttaaacgtcaatttgaacaatctcctaaacactatcgccagcaatcatgataataatactagagccaggcatcaaataaaacgaaaggctcagtcgaaagactgggcctttcgttttatctgttgtttgtcggtgaacgctctctactagagtcacactggctcaccttcgggtgggcctttctgcgtttata"
	//	Sequencekey                 string = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGGGGGGGGGGGGGGGGGGGGGGGGGAGCTCGGGGGGGGGGGGGGGGGGGGGGGGG"
	RestrictionsitetoAvoid             = []string{"SacI"}
	RemoveifnotinORF            bool   = true
	PreserveTranslatedseq       bool   = true
	EnzymeforRestrictionmapping string = "BsaI"
)

// Data which is returned from this protocol, and data types
var (
	Warnings                            string
	Status                              string
	SiteFreeSequence                    wtype.DNASequence // i.e. parts to order
	Sitesfoundinoriginal                []enzymes.Restrictionsites
	FragmentSizesfromRestrictionmapping []int
)

func main() {

	Sequence := wtype.MakeLinearDNASequence("Test", Sequencekey)

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// first lookup enzyme properties for all enzymes and make a new array
	enzlist := make([]wtype.RestrictionEnzyme, 0)
	for _, site := range RestrictionsitetoAvoid {
		enzsite := lookup.EnzymeLookup(site)
		enzlist = append(enzlist, enzsite)
	}

	// check for sites in the sequence
	sitesfound := enzymes.Restrictionsitefinder(Sequence, enzlist)

	// if no sites found skip to restriction map stage
	if len(sitesfound) == 0 {
		Warnings = "none"
		Status = "No sites found in sequence to remove so same sequence returned"
		SiteFreeSequence = Sequence
		Sitesfoundinoriginal = sitesfound

	} else {

		// make a list of sequences to avoid before modifying the sequence
		allsitestoavoid := make([]string, 0)

		// add all restriction sites (we need this step since the functions coming up require strings)
		for _, enzy := range enzlist {
			allsitestoavoid = append(allsitestoavoid, enzy.RecognitionSequence)
		}

		for _, site := range sitesfound {
			if site.Sitefound {

				var tempseq wtype.DNASequence
				var err error

				orfs := sequences.FindallNonOverlappingORFS(Sequence.Seq)
				warnings = append(warnings, text.Print("orfs: ", orfs))
				features := sequences.ORFs2Features(orfs)

				//set up a boolean to change to true if a sequence is found in an ORF
				foundinorf := false
				//set up an index for each orf found with site within it (need enzyme name too but will recheck all anyway!)
				orfswithsites := make([]int, 0)

				if len(orfs) > 0 {
					for i, orf := range orfs {

						// change func to handle this step of making dnaseq first

						dnaseq := wtype.MakeLinearDNASequence("orf"+strconv.Itoa(i), orf.DNASeq)

						foundinorfs := enzymes.Restrictionsitefinder(dnaseq, enzlist) // won't work yet orf is actually type features

						for _, siteinorf := range foundinorfs {
							if siteinorf.Sitefound == true {
								foundinorf = true
							}
						}

						if foundinorf == true {

							warning := text.Print("sites found in orf"+dnaseq.Nm, orf)
							warnings = append(warnings, warning)
						}
					}
				}
				if RemoveifnotinORF {
					if foundinorf == false {
						tempseq, err = sequences.RemoveSite(Sequence, site.Enzyme, allsitestoavoid)
						if err != nil {
							warning := text.Print("removal of site failed! improve your algorithm!", err.Error())
							warnings = append(warnings, warning)

						}
						fmt.Println("in 1")
						SiteFreeSequence = tempseq

						// all done if all sites are not in orfs!
						// make proper remove allsites func
					}
					if foundinorf == true {
						fmt.Println("in 2")
						fmt.Println("Sequence", Sequence)
						fmt.Println("enz", site.Enzyme.RecognitionSequence)
						//fmt.Println("features", features)
						SiteFreeSequence, err = sequences.RemoveSitesOutsideofFeatures(Sequence, site.Enzyme.RecognitionSequence, sequences.ReplaceBycomplement, features)
						fmt.Println("sitefreesequence", SiteFreeSequence)
						if err != nil {
							warnings = append(warnings, err.Error())
						}
					}
				} else {
					if PreserveTranslatedseq {
						// make func to check codon and swap site to preserve aa sequence product
						for _, orfnumber := range orfswithsites {
							fmt.Println("in 3")
							for _, position := range site.Positions("ALL") {
								orfcoordinates := sequences.MakeStartendPair(orfs[orfnumber].StartPosition, orfs[orfnumber].EndPosition)
								tempseq, _, _, err = sequences.ReplaceCodoninORF(tempseq, orfcoordinates, position, allsitestoavoid)
								if err != nil {
									warning := text.Print("removal of site from orf "+strconv.Itoa(orfnumber), " failed! improve your algorithm! "+err.Error())
									warnings = append(warnings, warning)
								}
							}

						}
					}
					SiteFreeSequence = tempseq
				}

			}
		}
	}
	fmt.Println("18")
	// Now let's find out the size of fragments we would get if digested with a common site cutter
	mapenz := lookup.EnzymeLookup(EnzymeforRestrictionmapping)

	FragmentSizesfromRestrictionmapping = enzymes.RestrictionMapper(Sequence, mapenz)

	// allow the data to be exported by capitalising the first letter of the variable
	Sitesfoundinoriginal = sitesfound

	Warnings = strings.Join(warnings, ";")
	fmt.Println("19")
	// Print status
	if Status == "all parts available" {
		/*Status =*/ fmt.Println(Status)
	} else {
		///*Status =*/ fmt.Println(
		fmt.Println("Warnings:", Warnings)
		fmt.Println("Sequence", Sequence)
		fmt.Println("SiteFreeSequence", SiteFreeSequence)
		fmt.Println("Sitesfound", Sitesfoundinoriginal)
		fmt.Println("Test digestion sizes with"+EnzymeforRestrictionmapping, FragmentSizesfromRestrictionmapping)

		fmt.Println("20")
	}

	fmt.Println("20")

}
