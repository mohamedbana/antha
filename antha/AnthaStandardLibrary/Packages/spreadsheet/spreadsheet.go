package spreadsheet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

var (
	alphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Parses an a1 style excel cell coordinate into ints for row and column for use by plotinum library
// note that 1 is subtracted from the column number in accordance with the go convention of counting from 0
func A1formattorowcolumn(a1 string) (row, column int, err error) {

	a1 = strings.ToUpper(a1)

	column, err = strconv.Atoi(a1[1:])
	column = column - 1
	if err == nil {
		row := strings.Index(alphabet, string(a1[0]))

		return row, column, err
	}

	column, err = strconv.Atoi(a1[2:])
	column = column - 1
	if err == nil {
		row := strings.Index(alphabet, a1[0:1])
		return row, column, err
	}

	column, err = strconv.Atoi(a1[3:])
	column = column - 1
	if err == nil {
		row := strings.Index(alphabet, a1[0:1])
		return row, column, err
	}

	newerr := fmt.Errorf(err.Error() + "more than first three letters of coordinate not int! seems unlikely")
	err = newerr
	return

}

// from a pair of cell coordinates an aray of all entrires between the pair will be returned (e.g. a1:a12 or a1:e1)
func ConvertMinMaxtoArray(minmax []string) (array []string, err error) {

	if len(minmax) != 2 {
		err = fmt.Errorf("can only make array from a pair of values")
		return
	}

	minrow, mincol, err := A1formattorowcolumn(minmax[0])
	if err != nil {
		return
	}
	maxrow, maxcol, err := A1formattorowcolumn(minmax[1])
	if err != nil {
		return
	}

	if minrow == maxrow {
		// fill by column
		array = make([]string, 0)
		for i := mincol; i < maxcol+1; i++ {
			fmt.Println(minrow)
			fmt.Println(alphabet)
			fmt.Println(alphabet[1])
			rowstring := alphabet[minrow]
			colstring := strconv.Itoa(i + 1)
			fmt.Println("bycolumn: ", rowstring, colstring)
			array = append(array, string(rowstring)+colstring)
		}

	} else if mincol == maxcol {
		// fill by row
		array = make([]string, 0)
		for i := minrow; i < maxrow+1; i++ {
			colstring := strconv.Itoa(mincol)
			rowstring := alphabet[i]
			fmt.Println(rowstring, colstring)
			array = append(array, string(rowstring)+colstring)
		}
	} else {
		err = fmt.Errorf("either column or row needs to be the same to make an array from two cordinates")
	}
	return

}

func OpenFile(filename string) (file *xlsx.File, err error) {
	file, err = xlsx.OpenFile(filename)
	return
}

func Sheet(file *xlsx.File, sheetnum int) (sheet *xlsx.Sheet) {
	sheet = file.Sheets[sheetnum]
	return
}
