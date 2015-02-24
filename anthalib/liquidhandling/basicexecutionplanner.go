// liquidhandling/executionplanner.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"os"
)

// a default execution planner which relies on a call to code external to the
// Antha project.
func BasicExecutionPlanner(request *LHRequest, params *LHProperties) *LHRequest {
	// essentially we defer everything to the existing liquid handling planner
	// which is NOT included as part of the language

	cnf := "config.csv"
	plf := "plan.tsv"

	MakeConfigFile(cnf, *request)
	MakePlanFile(plf, *request)

	// need to run the software and get the instructions etc. from it

	return request
}

func MakeConfigFile(fn string, request LHRequest) {
	f, _ := os.Create(fn)
	defer f.Close()

	// simple enough

	fmt.Fprintf(f, "stage,liquidhandling\n")

	inputs := request.Input_solutions

	for _, inputarr := range inputs {
		for _, input := range inputarr {
			// if both concentration and volume are set for this then
			// volume has priority

			if input.Vol != 0.0 {
				fmt.Fprintf(f, "%s,F,%s,\n", input.Name, input.Type)
			} else if input.Conc != 0.0 {
				fmt.Fprintf(f, "[%s],F,%s,\n", input.Name, input.Type)
			} else if input.Tvol != 0.0 {
				fmt.Fprintf(f, "%s,V,%s,\n", input.Name, input.Type)
			}
		}
	}

	plates := request.Plates
	inplat := plates["input"]
	robotfn := request.Robotfn

	if robotfn == "" {
		robotfn = "robots/defaultFelix.rbt"
	}

	fmt.Fprintln(f, "groups,\n")
	fmt.Fprintln(f, "parameters,\n")
	fmt.Fprintln(f, "volumeunit,1000000,# specifies microlitres\n")
	fmt.Fprintf(f, "platefile,plates/%s.txt,\n", inplat.Type)
	fmt.Fprintf(f, "platedir,plates/\n")
	fmt.Fprintf(f, "robotfile,%s,\n", robotfn)
	fmt.Fprintf(f, "tipdir,plates/,\n")
	fmt.Fprintf(f, "test_execution,1,\n")

	// TODO many more parameters here:
	// e.g. if component order is specified
}

func MakePlanFile(fn string, request LHRequest) {
	f, _ := os.Create(fn)
	defer f.Close()

	inputs := request.Input_solutions
	inputnames := make([]string, 0, len(inputs))
	c := 0
	for _, inputarr := range inputs {
		for _, input := range inputarr {
			inputnames[c] = input.CName
			fmt.Fprintf(f, "%s\t", inputnames[c])
			c += 1
		}
	}
	fmt.Fprintf(f, "Y\n")

	solns := request.Output_solutions

	for _, solution := range solns {
		for _, n := range inputnames {
			val := solution.GetComponentVolume(n)
			fmt.Fprintf(f, "%f\t", val)
		}
		fmt.Fprintf(f, " \n")
	}
}
