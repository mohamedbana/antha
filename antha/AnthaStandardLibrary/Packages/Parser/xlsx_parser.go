// antha/AnthaStandardLibrary/Packages/Parser/xlsx_parser.go: Part of the Antha language
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

package parser

import (
	"errors"
	//"flag"
	"fmt"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
	//"time"
)

//var xlsxPath = flag.String("f", "", "Path to an XLSX file")
//var sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
//var delimiter = flag.String("d", ",", "Delimiter to use between fields")

//var sheetIndex = 0
var delimiter = ","

type outputer func(s string)

func generateCSVFromXLSXsheet(excelFileName string, sheetIndex int, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				vals = append(vals, fmt.Sprintf("%q", cell.String()))
			}
			outputf(strings.Join(vals, delimiter) + "\n")
			fmt.Println("output", outputf)
		}
	}
	return nil
}

func generateCSVFromspecificXLSXsheet(excelFileName string, sheetname string, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
		/*case sheetIndex >= sheetLen:
			return fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
		}*/
		for _, sheet := range xlFile.Sheets {
			/*if strings.Contains(sheet.Name, sheetname) == false {
				return errors.New("no sheet of that name")
				break
			}*/
			for _, row := range sheet.Rows {
				var vals []string
				if row != nil {
					for _, cell := range row.Cells {
						vals = append(vals, fmt.Sprintf("%q", cell.String()))
					}
					outputf(strings.Join(vals, delimiter) + "\n")
					fmt.Println("output", outputf)
				}
			}
		}

	}
	return nil
}

func Xlsxparser(filename string, sheetIndex int, outputfilename string) (status string) {
	//flag.Parse()
	/*	if len(os.Args) < 3 {
			flag.PrintDefaults()
			return
		}
		flag.Parse()*/
	status = "complete"

	f, err := os.Create(outputfilename)
	fmt.Println(outputfilename)
	if err != nil {
		panic(err)
		status = "error"
	}
	defer f.Close()
	printer := func(s string) {
		_, _ = f.WriteString(s)

	}
	/*_, err = io.Copy(out, partsfile) // file not files[i] !

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}*/
	if err := generateCSVFromXLSXsheet(filename, sheetIndex, printer); err != nil {
		fmt.Println(err)
		status = "error"

	}
	fmt.Println("err", err)

	return status
}

func Assemblyfromxlsx(designandpartsfile string) (Assemblies []enzymes.Assemblyparameters) {

	Xlsxparser(designandpartsfile, 0, "tmp/Parts1.csv")
	Xlsxparser(designandpartsfile, 1, "tmp/Design1.csv")

	//time.Sleep(10000)
	Assemblies = Assemblyfromcsv("tmp/Design1.csv", "tmp/Parts1.csv")

	return Assemblies
}
