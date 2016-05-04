// multipleassemblies.go
package doe

import (
	"bufio"
	"fmt"
	//"io/ioutil"
	"strconv"
	"strings"

	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func ReadGeneFiles(genefiles []string) {

	geneseqs := make([]interface{}, 0)

	for _, genefile := range genefiles {
		genes, err := parser.DNAFiletoDNASequence(genefile, false)

		file, err := parser.Fastatocsv(genefile, "output.csv")
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		fmt.Println(lines)

		for _, gene := range genes {
			fmt.Println("genefile:", genefile, "GENE name:", gene.Nm, "plasmid?", gene.Plasmid)
		}
		if err != nil {
			panic(err.Error())
		}
		vals := make([]interface{}, len(genes))
		for i, v := range genes {
			vals[i] = v
		}
		for _, gene := range vals {
			geneseqs = append(geneseqs, gene)
		}

	}

	fmt.Println("number of sequences", len(geneseqs))

}

func AssemblyparametersfromRuns(runs []Run, enzymename string) (assemblyparameters []enzymes.Assemblyparameters) {

	assemblyparameters = make([]enzymes.Assemblyparameters, 0)
	parts := make([]wtype.DNASequence, 0)
	var parameters enzymes.Assemblyparameters

	enzyme, _ := lookup.TypeIIsLookup(enzymename)

	for j, run := range runs {

		parameters.Constructname = "run" + strconv.Itoa(j)
		parameters.Enzyme = enzyme

		for i, _ := range run.Setpoints {

			if strings.Contains(run.Factordescriptors[i], "Vector") {

				parameters.Vector = run.Setpoints[i].(wtype.DNASequence)
			} else if i < len(runs)-1 {
				parts = append(parts, run.Setpoints[i].(wtype.DNASequence))
			}

		}
		parameters.Partsinorder = parts
		assemblyparameters = append(assemblyparameters, parameters)
	}

	return
}

func DNASequencetoInterface(genes []wtype.DNASequence) (geneseqs []interface{}) {

	vals := make([]interface{}, len(genes))
	for i, v := range genes {
		vals[i] = v
	}
	for _, gene := range vals {
		geneseqs = append(geneseqs, gene)
	}
	return
}

/*
// inputs
var (
	//vectors       []string = []string{"Rewiring_DaughterVector_I.gb"} // put vector sequence files here
	//promoterfiles []string = []string{"All_proms.fasta"}              // put promoter sequence files here
	genefiles []string = []string{"genes.fasta", "notworking.fasta", "working.fasta"} //genes    []string
	//part3files    []string = []string("")                             // put other part sequence files here if it's not in the right order you'll have to check the code and swap the order
//	enzyme  string = "bbsI"
//	dirname string = "assembly_export"
)
*/
