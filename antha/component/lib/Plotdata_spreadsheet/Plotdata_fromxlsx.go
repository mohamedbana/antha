package Plotdata_spreadsheet

import (
	"fmt"
	//"math/rand"
	//"github.com/antha-lang/antha/internal/github.com/montanaflynn/stats"
	graph "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//                                                                         = "plotinumdata.xlsx"
//                                                                        = 0
/*datarange*/ //  = []string{"a4", "a16"}                                                           // row in A1 format i.e string{A,E} would indicate all data between those points
/*datarange*/ //= [][]string{[]string{"b4", "b16"}, []string{"c4", "c16"}, []string{"d4", "d16"}} // column in A1 format i.e string{1,12} would indicate all data between those points
//= "Excelfile.jpg"
//	HeaderRange []string // if Bycolumn == true, format would be e.g. string{A1,E1} else e.g. string{A1,A12}

// Data which is returned from this protocol, and data types

//	OutputData       []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	// Get some data.

	file, err := spreadsheet.OpenFile(_input.Filename)

	sheet := file.Sheets[_input.Sheet]

	Xdatarange, err := spreadsheet.ConvertMinMaxtoArray(_input.Xminmax)
	if err != nil {
		fmt.Println(_input.Xminmax, Xdatarange)
		panic(err)
	}
	fmt.Println(Xdatarange)

	Ydatarangearray := make([][]string, 0)
	for i, Yminmax := range _input.Yminmaxarray {
		Ydatarange, err := spreadsheet.ConvertMinMaxtoArray(Yminmax)
		if err != nil {
			panic(err)
		}
		if len(Xdatarange) != len(Ydatarange) {
			panicmessage := fmt.Errorf("for index", i, "of array", "len(Xdatarange) != len(Ydatarange)")
			panic(panicmessage.Error())
		}
		Ydatarangearray = append(Ydatarangearray, Ydatarange)
		fmt.Println(Ydatarange)
	}

	// now plot the graph

	// the data points

	graph.PlotfromMinMaxpairs(sheet, _input.Xminmax, _input.Yminmaxarray, _input.Exportedfilename)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _validation(_ctx context.Context, _input *Input_, _output *Output_) {

}

func _run(_ctx context.Context, value inject.Value) (inject.Value, error) {
	input := &Input_{}
	output := &Output_{}
	if err := inject.Assign(value, input); err != nil {
		return nil, err
	}
	_setup(_ctx, input)
	_steps(_ctx, input, output)
	_analysis(_ctx, input, output)
	_validation(_ctx, input, output)
	return inject.MakeValue(output), nil
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

func New() interface{} {
	return &Element_{
		inject.CheckedRunner{
			RunFunc: _run,
			In:      &Input_{},
			Out:     &Output_{},
		},
	}
}

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	Exportedfilename string
	Filename         string
	Sheet            int
	Xminmax          []string
	Yminmaxarray     [][]string
}

type Output_ struct {
}
