// meltingtemp.go
//Tm= 64.9 +41*(yG+zC-16.4)/(wA+xT+yG+zC)
//Wallace,R.B., Shaffer,J., Murphy,R.F., Bonner,J., Hirose,T., and Itakura,K. (1979) Nucleic Acids Res 6:3543-3557 (Abstract) and Sambrook,J., and Russell,D.W. (2001) Molecular Cloning: A Laboratory Manual. Cold Spring Harbor Laboratory Press; Cold Spring Harbor, NY. (CHSL Press)

// Package to calculate melting temp of dna sequences for primer design
package sequences

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

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
