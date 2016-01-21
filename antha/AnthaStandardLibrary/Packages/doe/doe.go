// doe.go
package doe

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
)

type DOEPair struct {
	Factor string
	Levels []interface{}
}

func (pair DOEPair) LevelCount() (numberoflevels int) {
	numberoflevels = len(pair.Levels)
	return
}
func Pair(factordescription string, levels []interface{}) (doepair DOEPair) {
	doepair.Factor = factordescription
	doepair.Levels = levels
	return
}

type Run struct {
	RunNumber           int
	StdNumber           int
	Factordescriptors   []string
	Setpoints           []interface{}
	Responsedescriptors []string
	ResponseValues      []interface{}
}

func AllComboCount(pairs []DOEPair) (numberofuniquecombos int) {
	fmt.Println("In AllComboCount", "len(pairs)", len(pairs))
	var movingcount int
	movingcount = (pairs[0]).LevelCount()
	fmt.Println("levelcount", movingcount)
	fmt.Println("len(levels)", len(pairs[0].Levels))
	for i := 1; i < len(pairs); i++ {
		fmt.Println("levelcount", movingcount)
		movingcount = movingcount * (pairs[i]).LevelCount()
	}
	numberofuniquecombos = movingcount
	return
}

func AllCombinations(factors []DOEPair) (runs []Run) {
	fmt.Println(factors)
	descriptors := make([]string, 0)
	setpoints := make([]interface{}, 0)
	runs = make([]Run, AllComboCount(factors))
	var run Run
	for i, factor := range factors {
		fmt.Println(factor, i, "of", AllComboCount(factors))
		for j, level := range factor.Levels {
			fmt.Println(factor, level, i, j, i+j)
			descriptors = append(descriptors, factor.Factor)
			setpoints = append(setpoints, level)
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints
			runs[i+j] = run
		}
	}
	return
}

func RunsFromDXDesign(xlsx string) (runs []Run, err error) {
	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return runs, err
	}
	sheet := spreadsheet.Sheet(file, 0)

	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 3; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		factordescriptors := make([]string, 0)
		responsedescriptors := make([]string, 0)
		setpoints := make([]interface{}, 0)
		responsevalues := make([]interface{}, 0)

		run.RunNumber, err = sheet.Cell(i, 1).Int()
		if err != nil {
			return runs, err
		}
		run.StdNumber, err = sheet.Cell(i, 0).Int()
		if err != nil {
			return runs, err
		}

		for j := 2; j < sheet.MaxCol; j++ {
			factororresponse := sheet.Cell(0, j).String()
			fmt.Println(i, j, factororresponse)
			if strings.Contains(factororresponse, "Factor") {
				//	maxfactorcol = j
				descriptor = strings.Split(sheet.Cell(1, j).String(), ":")[1]
				factrodescriptor := descriptor
				fmt.Println(i, j, descriptor)

				cell := sheet.Cell(i, j)

				celltype := cell.Type()

				if celltype == 1 {
					setpoint, err = cell.Float()
				} else if celltype == 3 {
					setpoint = cell.Bool()
				} else {
					setpoint = cell.String()
				}
				factordescriptors = append(factordescriptors, factrodescriptor)
				setpoints = append(setpoints, setpoint)

			} else {
				//run.Factordescriptors = factordescriptors

				//for k := maxfactorcol; k < sheet.MaxCol; k++ {
				//	factororresponse := sheet.Cell(0, maxfactorcol).String()
				if strings.Contains(factororresponse, "Response") {
					descriptor = sheet.Cell(1, j).String()
					responsedescriptor := descriptor
					fmt.Println("response", i, j, descriptor)
					responsedescriptors = append(responsedescriptors, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						responsevalues = append(responsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						responsevalues = append(responsevalues, responsevalue)
					}

				}

			}
		}
		run.Factordescriptors = factordescriptors
		run.Responsedescriptors = responsedescriptors
		run.Setpoints = setpoints
		run.ResponseValues = responsevalues

		runs = append(runs, run)
		factordescriptors = make([]string, 0)
		responsedescriptors = make([]string, 0)
	}

	return
}
