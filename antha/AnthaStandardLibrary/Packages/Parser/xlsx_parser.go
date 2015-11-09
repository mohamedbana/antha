package parser

import (
	"errors"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"strings"
)

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
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {
				var vals []string
				if row != nil {
					for _, cell := range row.Cells {
						vals = append(vals, fmt.Sprintf("%q", cell.String()))
					}
					outputf(strings.Join(vals, delimiter) + "\n")
				}
			}
		}

	}
	return nil
}

func Xlsxparser(filename string, sheetIndex int, outputprefix string) (f *os.File, status string) {
	status = "complete"
	f, err := ioutil.TempFile(os.TempDir(), outputprefix)

	defer f.Close()

	if err != nil {
		panic(err)
		status = "error"
	}
	printer := func(s string) {
		_, _ = f.WriteString(s)

	}
	if err := generateCSVFromXLSXsheet(filename, sheetIndex, printer); err != nil {
		status = "error"

	}
	return f, status
}

func Assemblyfromxlsx(designandpartsfile string) (Assemblies []enzymes.Assemblyparameters) {
	parts, _ := Xlsxparser(designandpartsfile, 0, "PartsTemp")
	design, _ := Xlsxparser(designandpartsfile, 1, "DesignTemp")

	Assemblies = Assemblyfromcsv(design.Name(), parts.Name())

	return Assemblies
}
