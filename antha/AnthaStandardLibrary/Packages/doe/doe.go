// Part of the Antha language
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

// Package for facilitating DOE methodology in antha
package doe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/tealeg/xlsx"
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
	AdditionalHeaders    []string // could represent a class e.g. Environment variable, processed, raw, location
	AdditionalSubheaders []string // e.g. well ID, Ambient Temp, order,
	AdditionalValues     []interface{}
}

func (run Run) AddResponseValue(responsedescriptor string, responsevalue interface{}) {

	for i, descriptor := range run.Responsedescriptors {
		if strings.ToUpper(descriptor) == strings.ToUpper(responsedescriptor) {
			run.ResponseValues[i] = responsevalue
		}
	}

}

func (run Run) AllResponses() (headers []string, values []interface{}) {
	headers = make([]string, 0)
	values = make([]interface{}, 0)

	for _, header := range run.Responsedescriptors {
		headers = append(headers, header)
	}
	for _, value := range run.ResponseValues {
		values = append(values, value)
	}
	return
}

func (run Run) AllFactors() (headers []string, values []interface{}) {
	headers = make([]string, 0)
	values = make([]interface{}, 0)

	for _, header := range run.Factordescriptors {
		fmt.Println(header)
		headers = append(headers, header)
	}
	for _, value := range run.Setpoints {
		values = append(values, value)
	}
	fmt.Println(headers, values)
	return
}

func (run Run) GetResponseValue(responsedescriptor string) (responsevalue interface{}, err error) {

	var tempresponsevalue interface{}
	headers, _ := run.AllResponses()

	errstr := fmt.Sprint("response descriptor", responsedescriptor, "not found in ", headers)
	err = fmt.Errorf(errstr)
	for i, descriptor := range run.Responsedescriptors {
		if strings.TrimSpace(strings.ToUpper(descriptor)) == strings.TrimSpace(strings.ToUpper(responsedescriptor)) {
			responsevalue = run.ResponseValues[i]
			return responsevalue, nil
		} else if strings.Contains(strings.TrimSpace(strings.ToUpper(descriptor)), strings.TrimSpace(strings.ToUpper(responsedescriptor))) {

			errstr := fmt.Sprint("response descriptor", responsedescriptor, "found within ", descriptor, "but no exact match")
			err = fmt.Errorf(errstr)
			tempresponsevalue = run.ResponseValues[i]
			return tempresponsevalue, err
		} else if strings.Contains(strings.TrimSpace(strings.ToUpper(responsedescriptor)), strings.TrimSpace(strings.ToUpper(descriptor))) {

			errstr := fmt.Sprint("response descriptors of ", descriptor, "found within ", responsedescriptor, "but not exact match")
			err = fmt.Errorf(errstr)
			tempresponsevalue = run.ResponseValues[i]
			return tempresponsevalue, err
		}
	}

	responsevalue = tempresponsevalue
	return
}

func (run Run) GetFactorValue(factordescriptor string) (factorvalue interface{}, err error) {

	var tempresponsevalue interface{}
	headers, values := run.AllFactors()

	errstr := fmt.Sprint("factor descriptor ", factordescriptor, " not found in ", headers, values)
	err = fmt.Errorf(errstr)
	for i, descriptor := range run.Factordescriptors {
		if strings.TrimSpace(strings.ToUpper(descriptor)) == strings.TrimSpace(strings.ToUpper(factordescriptor)) {
			factorvalue = run.Setpoints[i]
			return factorvalue, nil
		} else if strings.Contains(strings.TrimSpace(strings.ToUpper(descriptor)), strings.TrimSpace(strings.ToUpper(factordescriptor))) {

			errstr := fmt.Sprint("factor descriptor", factordescriptor, "found within ", descriptor, "but no exact match")
			err = fmt.Errorf(errstr)
			tempresponsevalue = run.Setpoints[i]
			return tempresponsevalue, err
		} else if strings.Contains(strings.TrimSpace(strings.ToUpper(factordescriptor)), strings.TrimSpace(strings.ToUpper(descriptor))) {

			errstr := fmt.Sprint("factor descriptors of ", descriptor, "found within ", factordescriptor, "but not exact match")
			err = fmt.Errorf(errstr)
			tempresponsevalue = run.Setpoints[i]
			return tempresponsevalue, err
		}
	}

	factorvalue = tempresponsevalue
	return
}

func AddNewResponseFieldandValue(run Run, responsedescriptor string, responsevalue interface{}) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, len(run.Responsedescriptors))
	responsevalues := make([]interface{}, len(run.ResponseValues))

	responsedescriptors = run.Responsedescriptors
	responsevalues = run.ResponseValues

	responsedescriptors = append(responsedescriptors, responsedescriptor)
	responsevalues = append(responsevalues, responsevalue)

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func AddNewResponseField(run Run, responsedescriptor string) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, len(run.Responsedescriptors))
	responsevalues := make([]interface{}, len(run.ResponseValues)+1)

	responsedescriptors = run.Responsedescriptors

	for i := range run.ResponseValues {
		responsevalues[i] = run.ResponseValues[i]
	}
	responsedescriptors = append(responsedescriptors, responsedescriptor)

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func DeleteResponseField(run Run, responsedescriptor string) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	for i, descriptor := range run.Responsedescriptors {
		if strings.ToUpper(descriptor) != strings.ToUpper(responsedescriptor) {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, run.ResponseValues[i])
		}
	}

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func ReplaceResponseValue(run Run, responsedescriptor string, responsevalue interface{}) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	for i, descriptor := range run.Responsedescriptors {
		if strings.ToUpper(descriptor) != strings.ToUpper(responsedescriptor) {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, run.ResponseValues[i])
		} else if strings.ToUpper(descriptor) == strings.ToUpper(responsedescriptor) {
			responsedescriptors = append(responsedescriptors, descriptor)
			responsevalues = append(responsevalues, responsevalue)
		}
	}

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func DeleteAllResponses(run Run) (newrun Run) {

	newrun = run

	responsedescriptors := make([]string, 0)
	responsevalues := make([]interface{}, 0)

	newrun.Responsedescriptors = responsedescriptors
	newrun.ResponseValues = responsevalues

	return
}

func AddNewFactorFieldandValue(run Run, factordescriptor string, factorvalue interface{}) (newrun Run) {

	newrun = run

	factordescriptors := make([]string, len(run.Factordescriptors))
	factorvalues := make([]interface{}, len(run.Setpoints))

	factordescriptors = run.Factordescriptors
	factorvalues = run.Setpoints

	factordescriptors = append(factordescriptors, factordescriptor)
	factorvalues = append(factorvalues, factorvalue)

	newrun.Factordescriptors = factordescriptors
	newrun.Setpoints = factorvalues

	return
}

func AddAdditionalValue(run Run, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	newrun = run

	values := make([]interface{}, 0)

	for _, value := range run.AdditionalValues {
		values = append(values, value)
	}

	for _, descriptor := range run.AdditionalSubheaders {
		if strings.ToUpper(descriptor) == strings.ToUpper(additionalsubheader) {
			values = append(values, additionalvalue)
		}
	}

	newrun.AdditionalValues = values

	return
}

func ReplaceAdditionalValue(run Run, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	newrun = run

	values := make([]interface{}, len(run.AdditionalSubheaders))

	for i, descriptor := range run.AdditionalSubheaders {
		if strings.ToUpper(descriptor) == strings.ToUpper(additionalsubheader) {
			values[i] = additionalvalue
		} else {
			values[i] = run.AdditionalValues[i]
		}
	}

	newrun.AdditionalValues = values

	return
}

func AddAdditionalHeaders(run Run, additionalheader string, additionalsubheader string) (newrun Run) {

	newrun = run

	headers := make([]string, 0)

	for _, header := range run.AdditionalHeaders {
		headers = append(headers, header)
	}

	headers = append(headers, additionalheader)

	subheaders := make([]string, 0)

	for _, subheader := range run.AdditionalSubheaders {
		subheaders = append(subheaders, subheader)
	}

	subheaders = append(subheaders, additionalsubheader)

	newrun.AdditionalHeaders = headers
	newrun.AdditionalSubheaders = subheaders

	// fmt.Println("newrun: ", newrun)
	return

}

func AddAdditionalHeaderandValue(run Run, additionalheader string, additionalsubheader string, additionalvalue interface{}) (newrun Run) {

	// only add column if no column with header exists
	if search.InSlice(additionalsubheader, run.AdditionalSubheaders) == false {

		midrun := AddAdditionalHeaders(run, additionalheader, additionalsubheader)
		// fmt.Println("midrun: ", midrun)
		newrun = AddAdditionalValue(midrun, additionalsubheader, additionalvalue)
	} else {
		newrun = ReplaceAdditionalValue(run, additionalsubheader, additionalvalue)
	}
	return
}

func (run Run) CheckAdditionalInfo(subheader string, value interface{}) bool {

	for i, header := range run.AdditionalSubheaders {
		if strings.ToUpper(header) == strings.ToUpper(subheader) && run.AdditionalValues[i] == value {
			return true
		}
	}
	return false
}

func (run Run) GetAdditionalInfo(subheader string) (value interface{}, err error) {

	for i, header := range run.AdditionalSubheaders {
		if strings.ToUpper(header) == strings.ToUpper(subheader) {
			value = run.AdditionalValues[i]
			return value, err

		}
	}
	// fmt.Println("Header: ", subheader)
	return value, fmt.Errorf("header, ", subheader, " not found in ", run.AdditionalSubheaders)
}

func AddFixedFactors(runs []Run, fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	if len(runs) > 0 {
		for _, run := range runs {
			descriptors := make([]string, len(run.Factordescriptors))
			setpoints := make([]interface{}, len(run.Setpoints))

			for i, descriptor := range run.Factordescriptors {
				descriptors[i] = descriptor
			}

			for i, setpoint := range run.Setpoints {
				setpoints[i] = setpoint
			}

			for _, fixed := range fixedfactors {

				descriptors = append(descriptors, fixed.Factor)
				setpoints = append(setpoints, fixed.Levels[0])

			}
			run.Factordescriptors = descriptors
			run.Setpoints = setpoints

		}

	} else {
		runs = RunsFromFixedFactors(fixedfactors)
	}

	runswithfixedfactors = runs

	return
}

func RunsFromFixedFactors(fixedfactors []DOEPair) (runswithfixedfactors []Run) {

	var run Run
	var descriptors = make([]string, 0)
	var setpoints = make([]interface{}, 0)

	for _, factor := range fixedfactors {

		descriptors = append(descriptors, factor.Factor)
		setpoints = append(setpoints, factor.Levels[0])

	}

	run.Factordescriptors = descriptors
	run.Setpoints = setpoints

	runswithfixedfactors = make([]Run, 1)

	runswithfixedfactors[0] = run

	return

}

func AllComboCount(pairs []DOEPair) (numberofuniquecombos int) {
	// fmt.Println("In AllComboCount", "len(pairs)", len(pairs))
	var movingcount int
	movingcount = (pairs[0]).LevelCount()
	// fmt.Println("levelcount", movingcount)
	// fmt.Println("len(levels)", len(pairs[0].Levels))
	for i := 1; i < len(pairs); i++ {
		// fmt.Println("levelcount", movingcount)
		movingcount = movingcount * (pairs[i]).LevelCount()
	}
	numberofuniquecombos = movingcount
	return
}

func FixedAndNonFixed(factors []DOEPair) (fixedfactors []DOEPair, nonfixed []DOEPair) {

	fixedfactors = make([]DOEPair, 0)
	nonfixed = make([]DOEPair, 0)

	for _, factor := range factors {
		if len(factor.Levels) == 1 {
			fixedfactors = append(fixedfactors, factor)
		} else if len(factor.Levels) > 1 {
			nonfixed = append(nonfixed, factor)
		}
	}
	return
}

func IsFixedFactor(factor DOEPair) (yesorno bool) {
	if len(factor.Levels) == 1 {
		yesorno = true
	}
	return
}

func AllCombinations(factors []DOEPair) (runs []Run) {

	//fixed, nonfixed := FixedAndNonFixed(factors)

	numberofruns := AllComboCount(factors)

	runs = make([]Run, numberofruns)

	var swapevery int
	var numberofswaps int
	for i, factor := range factors {

		counter := 0
		runswitheachlevelforthisfactor := numberofruns / factor.LevelCount()

		if i == 0 {
			swapevery = runswitheachlevelforthisfactor
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		} else {
			swapevery = swapevery / factor.LevelCount()
			numberofswaps = runswitheachlevelforthisfactor / swapevery
		}

		for j := 0; j < numberofswaps; j++ {
			for _, level := range factor.Levels {
				for k := 0; k < swapevery; k++ {

					runs[counter] = AddNewFactorFieldandValue(runs[counter], factor.Factor, level)
					runs[counter].RunNumber = counter + 1
					runs[counter].StdNumber = counter + 1
					counter++
				}
			}
		}

	}
	//runs = AddFixedFactors(runs, fixed)
	return
}

func ParseRunWellPair(pair string, nameappendage string) (runnumber int, well string, err error) {
	split := strings.Split(pair, ":")

	numberstring := strings.SplitAfter(split[0], nameappendage)

	runnumber, err = strconv.Atoi(string(numberstring[1]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "+ Failed at", pair, nameappendage)
	}
	well = split[1]
	return
}

func AddWelllocations(DXORJMP string, xlsxfile string, oldsheet int, runnumbertowellcombos []string, nameappendage string, pathtosave string, extracolumnheaders []string, extracolumnvalues []interface{}) error {

	var xlsxcell *xlsx.Cell

	file, err := spreadsheet.OpenFile(xlsxfile)
	if err != nil {
		return err
	}

	sheet := spreadsheet.Sheet(file, oldsheet)

	_ = file.AddSheet("hello")

	//extracolumn := sheet.MaxCol + 1

	// add extra column headers first
	for _, extracolumnheader := range extracolumnheaders {
		xlsxcell = sheet.Rows[0].AddCell()

		xlsxcell.Value = "Extra column added"
		// fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
		xlsxcell = sheet.Rows[1].AddCell()
		xlsxcell.Value = extracolumnheader
	}

	// now add well position column
	xlsxcell = sheet.Rows[0].AddCell()

	xlsxcell.Value = "Location"
	// fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
	xlsxcell = sheet.Rows[1].AddCell()
	xlsxcell.Value = "Well ID"

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
					xlsxcell.SetValue(extracolumnvalue)
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

		factordescriptors := make([]string, 0)
		responsedescriptors := make([]string, 0)
		setpoints := make([]interface{}, 0)
		responsevalues := make([]interface{}, 0)
		otherheaders := make([]string, 0)
		othersubheaders := make([]string, 0)
		otherresponsevalues := make([]interface{}, 0)

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

			if strings.Contains(factororresponse, "Factor") {

				descriptor = strings.Split(sheet.Cell(1, j).String(), ":")[1]
				factrodescriptor := descriptor
				//fmt.Println(i, j, descriptor)

				cell := sheet.Cell(i, j)

				celltype := cell.Type()

				_, err := cell.Float()

				if strings.ToUpper(cell.Value) == "TRUE" {
					setpoint = true //cell.SetBool(true)
				} else if strings.ToUpper(cell.Value) == "FALSE" {
					setpoint = false //cell.SetBool(false)
				} else if celltype == 3 {
					setpoint = cell.Bool()
				} else if err == nil || celltype == 1 {
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

				factordescriptors = append(factordescriptors, factrodescriptor)
				setpoints = append(setpoints, setpoint)

			} else if strings.Contains(factororresponse, "Response") {
				descriptor = sheet.Cell(1, j).String()
				responsedescriptor := descriptor
				//// fmt.Println("response", i, j, descriptor)
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

			} else {
				descriptor = sheet.Cell(1, j).String()
				responsedescriptor := descriptor

				otherheaders = append(otherheaders, factororresponse)
				othersubheaders = append(othersubheaders, responsedescriptor)

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
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				}

			}
		}
		run.Factordescriptors = factordescriptors
		run.Responsedescriptors = responsedescriptors
		run.Setpoints = setpoints
		run.ResponseValues = responsevalues
		run.AdditionalHeaders = otherheaders
		run.AdditionalSubheaders = othersubheaders
		run.AdditionalValues = otherresponsevalues

		runs = append(runs, run)
		factordescriptors = make([]string, 0)
		responsedescriptors = make([]string, 0)

		// assuming this is necessary too
		otherheaders = make([]string, 0)
		othersubheaders = make([]string, 0)
	}

	return
}

func DXXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet = xlsxfile.AddSheet("Sheet1")

	// add headers
	row = sheet.AddRow()

	// 2 blank cells
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = ""

	// take factor and run descriptors from first run (assuming they're all the same)
	for i, _ := range runs[0].Factordescriptors {
		cell = row.AddCell()
		cell.Value = "Factor " + strconv.Itoa(i+1)

	}
	for i, _ := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = "Response " + strconv.Itoa(i+1)

	}
	for _, additionalheader := range runs[0].AdditionalHeaders {
		cell = row.AddCell()
		cell.Value = additionalheader

	}
	// new row
	row = sheet.AddRow()

	// add Std and Run number headers
	cell = row.AddCell()
	cell.Value = "Std"
	cell = row.AddCell()
	cell.Value = "Run"

	// then add subheadings and descriptors
	for i, descriptor := range runs[0].Factordescriptors {
		letter := wutil.NumToAlpha(i + 1)
		cell = row.AddCell()
		cell.Value = letter + ":" + descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}

	// add blank row

	row = sheet.AddRow()

	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()
		// Std
		cell = row.AddCell()
		cell.SetValue(run.StdNumber)

		// Run
		cell = row.AddCell()
		cell.SetValue(run.RunNumber)

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}

// jmp

func RunsFromJMPDesign(xlsx string, factorcolumns []int, responsecolumns []int, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return runs, err
	}
	sheet := spreadsheet.Sheet(file, 0)

	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 1; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		factordescriptors := make([]string, 0)
		responsedescriptors := make([]string, 0)
		setpoints := make([]interface{}, 0)
		responsevalues := make([]interface{}, 0)
		otherheaders := make([]string, 0)
		othersubheaders := make([]string, 0)
		otherresponsevalues := make([]interface{}, 0)

		run.RunNumber = i //sheet.Cell(i, 1).Int()

		run.StdNumber = i //sheet.Cell(i, 0).Int()

		for j := 0; j < sheet.MaxCol; j++ {

			var factororresponse string

			if search.Contains(factorcolumns, j) {
				factororresponse = "Factor"
			} else if search.Contains(responsecolumns, j) {
				factororresponse = "Response"
			}

			if strings.Contains(factororresponse, "Factor") {

				descriptor = sheet.Cell(0, j).String()
				factrodescriptor := descriptor
				fmt.Println(i, j, descriptor)

				cell := sheet.Cell(i, j)

				celltype := cell.Type()

				_, err := cell.Float()

				if strings.ToUpper(cell.Value) == "TRUE" {
					setpoint = true //cell.SetBool(true)
				} else if strings.ToUpper(cell.Value) == "FALSE" {
					setpoint = false //cell.SetBool(false)
				} else if celltype == 3 {
					setpoint = cell.Bool()
				} else if err == nil || celltype == 1 {
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
				factordescriptors = append(factordescriptors, factrodescriptor)
				setpoints = append(setpoints, setpoint)

			} else if strings.Contains(factororresponse, "Response") {
				descriptor = sheet.Cell(0, j).String()
				responsedescriptor := descriptor

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

			} else /*if j != patterncolumn*/ {
				descriptor = sheet.Cell(0, j).String()
				responsedescriptor := descriptor

				otherheaders = append(otherheaders, factororresponse)
				othersubheaders = append(othersubheaders, responsedescriptor)

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
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				} else {
					responsevalue := cell.String()
					otherresponsevalues = append(otherresponsevalues, responsevalue)
				}

			}
		}
		run.Factordescriptors = factordescriptors
		run.Responsedescriptors = responsedescriptors
		run.Setpoints = setpoints
		run.ResponseValues = responsevalues
		run.AdditionalHeaders = otherheaders
		run.AdditionalSubheaders = othersubheaders
		run.AdditionalValues = otherresponsevalues

		runs = append(runs, run)
		factordescriptors = make([]string, 0)
		responsedescriptors = make([]string, 0)

		// assuming this is necessary too
		otherheaders = make([]string, 0)
		othersubheaders = make([]string, 0)
	}

	return
}

func JMPXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	//var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet = xlsxfile.AddSheet("Sheet1")

	// new row
	row = sheet.AddRow()

	// then add subheadings and descriptors
	for _, descriptor := range runs[0].Factordescriptors {

		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}

func XLSXFileFromRuns(runs []Run, outputfilename string, dxorjmp string) (xlsxfile *xlsx.File) {
	if dxorjmp == "DX" {
		xlsxfile = DXXLSXFilefromRuns(runs, outputfilename)
	}
	if dxorjmp == "JMP" {
		xlsxfile = JMPXLSXFilefromRuns(runs, outputfilename)
	}
	return
}

func RunsFromDesign(designfile string, intfactors []string, responsecolumns []int, dxorjmp string) (runs []Run, err error) {

	if dxorjmp == "DX" {

		runs, err = RunsFromDXDesign(designfile, intfactors)
		if err != nil {
			return runs, err
		}

	} else if dxorjmp == "JMP" {

		factorcolumns := findFactorColumns(designfile, responsecolumns)

		runs, err = RunsFromJMPDesign(designfile, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return runs, err
		}
	}
	return
}

func RunsFromDesignPreResponses(designfile string, intfactors []string, dxorjmp string) (runs []Run, err error) {

	if dxorjmp == "DX" {

		runs, err = RunsFromDXDesign(designfile, intfactors)
		if err != nil {
			return runs, err
		}

	} else if dxorjmp == "JMP" {

		factorcolumns, responsecolumns, _ := findJMPFactorandResponseColumnsinEmptyDesign(designfile)

		runs, err = RunsFromJMPDesign(designfile, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return runs, err
		}
	}
	return

}

func findFactorColumns(xlsx string, responsefactors []int) (factorcolumns []int) {

	factorcolumns = make([]int, 0)

	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return factorcolumns
	}
	sheet := spreadsheet.Sheet(file, 0)

	for i := 0; i < sheet.MaxCol; i++ {
		if search.BinarySearch(responsefactors, i) == false && strings.ToUpper(sheet.Cell(0, i).String()) != "PATTERN" {
			factorcolumns = append(factorcolumns, i)
		}
	}

	return
}

// add func to auto check for Response and factor status based on empty entries implying Response column
func findJMPFactorandResponseColumnsinEmptyDesign(xlsx string) (factorcolumns []int, responsecolumns []int, PatternColumn int) {

	factorcolumns = make([]int, 0)
	responsecolumns = make([]int, 0)

	file, err := spreadsheet.OpenFile(xlsx)
	if err != nil {
		return
	}
	sheet := spreadsheet.Sheet(file, 0)

	//descriptors := make([]string, 0)

	for j := 0; j < sheet.MaxCol; j++ {

		descriptor := sheet.Cell(0, j).String()
		//	descriptors = append(descriptors,descriptor)
		if strings.ToUpper(descriptor) == "PATTERN" {
			PatternColumn = j
		}
	}
	// iterate through every run of the design sheet (row) and if all values for that row == "", the column is interpreted as a response
	for i := 1; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		for j := 0; j < sheet.MaxCol; j++ {

			if j != PatternColumn && sheet.Cell(i, j).String() != "" {
				factorcolumns = append(factorcolumns, j)
			} else if sheet.Cell(i, j).String() == "" {

				responsecolumns = append(responsecolumns, j)
			}

		}

	}
	return
}
