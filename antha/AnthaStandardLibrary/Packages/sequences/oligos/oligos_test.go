// oligos_test.go
package oligos

import (
	//"fmt"
	//"strconv"
	//	"strings"

	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	//"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

// simple reverse complement check to test testing methodology initially
type testpair struct {
	//BasicMeltingTemp test
	sequence    wtype.DNASequence
	meltingtemp wunit.Temperature
	// FWDoligoTest
	mintemp             wunit.Temperature
	maxtemp             wunit.Temperature
	maxGCcontent        float64
	minlength           int
	maxlength           int
	outputoligoseq      string
	calculatedGCcontent float64
}

var meltingtemptests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("whatever", "AAAAAAAAAAAAAAAAAAA"),
		meltingtemp: wunit.NewTemperature(29.511, "C")},
}

var oligotests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("sfGFP", "ATGAGCAAAGGAGAAGAACTTTTCACTGGAGTTGTCCCAATTCTTGTTGAATTAGATGGTGATGTTAATGGGCACAAATTTTCTGTCCGTGGAGAGGGTGAAGGTGATGCTACAAACGGAAAACTCACCCTTAAATTTATTTGCACTACTGGAAAACTACCTGTTCCATGGCCAACACTTGTCACTACTCTGACCTATGGTGTTCAATGCTTTTCCCGTTATCCGGATCACATGAAACGGCATGACTTTTTCAAGAGTGCCATGCCCGAAGGTTATGTACAGGAACGCACTATATCTTTCAAAGATGACGGGACCTACAAGACGCGTGCTGAAGTCAAGTTTGAAGGTGATACCCTTGTTAATCGTATCGAGTTAAAAGGTATTGATTTTAAAGAAGATGGAAACATTCTCGGACACAAACTCGAGTACAACTTTAACTCACACAATGTATACATCACGGCAGACAAACAAAAGAATGGAATCAAAGCTAACTTCAAAATTCGCCACAACGTTGAAGATGGTTCCGTTCAACTAGCAGACCATTATCAACAAAATACTCCAATTGGCGATGGCCCTGTCCTTTTACCAGACAACCATTACCTGTCGACACAATCTGTCCTTTCGAAAGATCCCAACGAAAAGCGTGACCACATGGTCCTTCTTGAGTTTGTAACTGCTGCTGGGATTACACATGGCATGGATGAGCTCTACAAATAA"),
		meltingtemp:         wunit.NewTemperature(52.764, "C"),
		mintemp:             wunit.NewTemperature(50, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           25,
		maxlength:           45,
		outputoligoseq:      "ATGAGCAAAGGAGAAGAACTTTTCA",
		calculatedGCcontent: 0.36},
}

func TestBasicMeltingTemp(t *testing.T) {
	for _, oligo := range meltingtemptests {
		result := BasicMeltingTemp(oligo.sequence)
		if result != oligo.meltingtemp {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.meltingtemp.ToString(), "\n",
				"got", result.ToString(), "\n",
			)
		}
	}

}

func TestFWDOligoSeq(t *testing.T) {
	for _, oligo := range oligotests {
		oligoseq, gc := FWDOligoSeq(oligo.sequence.Sequence(), oligo.maxGCcontent, oligo.minlength, oligo.maxlength, oligo.mintemp, oligo.maxtemp)
		if oligoseq != oligo.outputoligoseq {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.outputoligoseq, "\n",
				"got", oligoseq, "\n",
			)
		}
		if gc != oligo.calculatedGCcontent {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.calculatedGCcontent, "\n",
				"got", gc, "\n",
			)
		}
	}

}
