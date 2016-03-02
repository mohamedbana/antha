// PrimerDesign

// Get DNA sequence
// Find region within that sequence
// Find primer sequence within that region that fits criteria
//
package oligos

import (
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"strings"
)

// calculates the basic melting temperature of a DNA sequence
func BasicMeltingTemp(primersequence wtype.DNASequence) (meltingtemp wunit.Temperature) {

	primerseq := primersequence.Sequence()

	primerseq = strings.ToUpper(primerseq)

	g := strings.Count(primerseq, "G")
	c := strings.Count(primerseq, "C")
	a := strings.Count(primerseq, "A")
	t := strings.Count(primerseq, "T")

	var mt float64

	if len(primerseq) < 14 {
		//err = fmt.Errorf("cannot use this algorithm for sequences less than 13 nucleotides")

		mt = float64((a+t)*2 + (g+c)*4)

	} else {
		mt = 64.9 + 41.0*(float64(g+c)-16.4)/float64(a+t+c+g)

		fmt.Println(mt)

	}

	meltingtemp = wunit.NewTemperature(mt, "℃")
	return
}

//define region in DNA sequence
func DNAregion(sequence string, startposition int, endposition int) (region string) {

	dnaseq := sequence

	//define region in sequence to create primer. NB: Sequence position will start from 0 not 1.

	region = dnaseq[startposition-1 : endposition]

	return

}

// Takes defined region and makes an oligosequence between a defined minimum and maximum length
// with a melting temperature between a defined minimum and maximum and a maximum GC content ( between 0 and 1).
// function finds oligo by starting at position 0 and making sequence of the minimum length, calculating parameters
// and if they do not match then adds one basepair to end of sequence until the maximum length is reached.
// if still unsuccessful, the function begins again at position 1 and cycles through until a matching oligo sequence is found.

func FWDOligoSeq(region string, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature) (oligoseq string, GCpercentage float64) {

	//var start int
	//var end int

	for start := 0; start < maxlength; start++ {

		for end := minlength + start; end < start+maxlength; end++ {
			tempoligoseq := region[start:end]

			ssoligo := wtype.MakeSingleStrandedDNASequence("oligo", tempoligoseq)

			temppercentage := sequences.GCcontent(tempoligoseq)

			meltingtemp := BasicMeltingTemp(ssoligo)

			fmt.Println(ssoligo.Seq, temppercentage, meltingtemp.ToString())

			if GCpercentage <= maxGCcontent && minmeltingtemp.SIValue() < meltingtemp.SIValue() && maxmeltingtemp.SIValue() > meltingtemp.SIValue() {
				fmt.Println(tempoligoseq, temppercentage)
				oligoseq = tempoligoseq
				GCpercentage = temppercentage
				return

			}
		}
	}

	//}else {
	//	fmt.Println("no oligos")
	//	}

	return
}
