// doe.go
package doe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
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
	RunNumber            int
	StdNumber            int
	Factordescriptors    []string
	Setpoints            []interface{}
	Responsedescriptors  []string
	ResponseValues       []interface{}
	AdditionalHeaders    []string
	AdditionalSubheaders []string
	AdditionalValues     []interface{}
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

func ParseRunWellPair(pair string, nameappendage string) (runnumber int, well string, err error) {
	split := strings.Split(pair, ":")

	numberstring := strings.SplitAfter(split[0], nameappendage)

	//numberstring := split[0]
	fmt.Println("Pair", pair, "SPLIT", split /*string(numberstring[0])*/)
	fmt.Println("NUMBERSTRING!!", numberstring /*string(numberstring[0])*/)
	runnumber, err = strconv.Atoi(string(numberstring[1]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "+ Failed at", pair, nameappendage)
	}
	well = split[1]
	return
}

func AddWelllocations(xlsxfile string, oldsheet int, runnumbertowellcombos []string, nameappendage string, pathtosave string, extracolumnheaders []string, extracolumnvalues []interface{}) error {

	/*

		    var row *xlsx.Row
		    var cell *xlsx.Cell
		    var err error

		    file = xlsx.NewFile()
		    //sheet, err = file.AddSheet("Sheet1")
		    if err != nil {
		        fmt.Printf(err.Error())
		    }
		    row = sheet.AddRow()
		    cell = row.AddCell()
		    cell.Value = "I am a cell!"
		    err = file.Save("MyXLSXFile.xlsx")
		    if err != nil {
		        fmt.Printf(err.Error())
		    }
		}



	*/

	var xlsxcell *xlsx.Cell

	file, err := spreadsheet.OpenFile(xlsxfile)
	if err != nil {
		return err
	}

	sheet := spreadsheet.Sheet(file, oldsheet)

	_ = file.AddSheet("hello")

	/*
		for _, row := range sheet.Rows {
			//newrow := newsheet.AddRow()

			newsheet.Rows = append(newsheet.Rows, row)
		}
	*/
	extracolumn := sheet.MaxCol + 1

	// add extra column headers first
	for _, extracolumnheader := range extracolumnheaders {
		xlsxcell = sheet.Rows[0].AddCell()

		xlsxcell.Value = "Extra column added"
		fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
		xlsxcell = sheet.Rows[1].AddCell()
		xlsxcell.Value = extracolumnheader
	}

	// now add well position column
	xlsxcell = sheet.Rows[0].AddCell()

	xlsxcell.Value = "Location"
	fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
	xlsxcell = sheet.Rows[1].AddCell()
	xlsxcell.Value = "Well ID"
	//	sheet.Cell(0, extracolumn).SetString("Location")
	//	sheet.Cell(1, extracolumn).SetString("Well")

	for i := 3; i < sheet.MaxRow; i++ {
		for _, pair := range runnumbertowellcombos {
			runnumber, well, err := ParseRunWellPair(pair, nameappendage)
			if err != nil {
				return err
			}
			xlsxrunmumber, err := sheet.Cell(i, 1).Int()
			if err != nil {
				return err
			}
			if xlsxrunmumber == runnumber {
				for _, extracolumnvalue := range extracolumnvalues {
					xlsxcell = sheet.Rows[i].AddCell()
					xlsxcell.Value = extracolumnvalue.(string)
				}
				xlsxcell = sheet.Rows[i].AddCell()
				xlsxcell.Value = well

			}
		}
	}

	err = file.Save(pathtosave)

	return err
}

func RunsFromDXDesign(xlsx string, intfactors []string) (runs []Run, err error) {
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

				_, err := cell.Float()

				if err == nil || celltype == 1 {
					setpoint, _ = cell.Float()
					if search.InSlice(descriptor, intfactors) {
						setpoint, err = cell.Int()
						if err != nil {
							return runs, err
						}
					}
				} else {
					setpoint = cell.String()
				}

				if celltype == 3 {
					setpoint = cell.Bool()
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
