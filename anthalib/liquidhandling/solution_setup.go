// anthalib//liquidhandling/solution_setup.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import "errors"
import "github.com/antha-lang/antha/anthalib/wutil"
import "fmt"

// determines how to
// fulfil the requirements for making
// solutions to specifications
func solution_setup(request *LHRequest, prms *LHProperties) (map[string]*LHSolution, map[string]float64) {
	solutions := request.Output_solutions

	// need a component-wise index of total volumes and components

	mtvols := make(map[string][]float64, 10)
	mconcs := make(map[string][]float64, 10)

	// we need to get info on the maximum solubilities
	Smax := make(map[string]float64, 10)

	for _, solution := range solutions {
		components := solution.Components

		// we need to identify the concentration components
		// and the total volume components, if we have
		// concentrations but not the tvols we have to return
		// an error

		arrCncs := make([]*LHComponent, 0, len(components))
		arrTvol := make([]*LHComponent, 0, len(components))
		cmpvol := 0.0
		totalvol := 0.0

		for _, component := range components {
			// what sort of component is it?
			// what is the total volume ?
			conc := component.Conc
			tvol := component.Tvol

			if conc != 0.0 {
				arrCncs = append(arrCncs, component)
			} else if tvol != 0.0 {
				arrTvol = append(arrTvol, component)
				tv := component.Tvol
				if totalvol == 0.0 || totalvol == tv {
					totalvol = tv
				} else {
					// error
					wutil.Error(errors.New(fmt.Sprintf("Inconsistent total volumes %-6.4f and %-6.4f at component %s", totalvol, tv, component.Name)))
				}
			} else {
				cmpvol += component.Vol
			}
		}

		// add everything to the maps

		for _, cmp := range arrCncs {
			nm := cmp.CName
			cnc := cmp.Conc

			_, ok := Smax[nm]

			if !ok {
				Smax[nm] = cmp.Smax
			}

			var cncslc []float64

			cncslc, ok = mconcs[nm]

			if !ok {
				cncslc = make([]float64, 0, 10)
			}

			cncslc = append(cncslc, cnc)

			mconcs[nm] = cncslc
		}

		// now the total volumes

		for _, cmp := range arrTvol {
			nm := cmp.CName
			tvol := cmp.Tvol

			var tvslc []float64

			tvslc, ok := mtvols[nm]

			if !ok {
				tvslc = make([]float64, 0, 10)
			}

			tvslc = append(tvslc, tvol)

			mtvols[nm] = tvslc
		}

	}
	// so now we should be able to make stock concentrations
	// first we need the min and max for each

	minrequired := make(map[string]float64, len(mtvols))
	maxrequired := make(map[string]float64, len(mtvols))
	Tvols := make(map[string]float64, len(mtvols))
	vmin := prms.CurrConf.Minvol

	for cmp, arr := range mtvols {
		min := wutil.FMin(arr)
		max := wutil.FMax(arr)
		tvol := wutil.FMin(mtvols[cmp])
		minrequired[cmp] = min
		maxrequired[cmp] = max

		// if smax undefined we need to deal  - we assume infinite solubility!!

		_, ok := Smax[cmp]

		if !ok {
			Smax[cmp] = 9999999
			wutil.Warn(fmt.Sprintf("Max solubility undefined for component %s -- assuming infinite solubility!", cmp))
		}

		Tvols[cmp] = tvol
	}

	stockconcs := choose_stock_concentrations(minrequired, maxrequired, Smax, vmin, Tvols)

	// nearly there now! Need to turn all the components into volumes, then we're done

	// make an array for the new solutions

	newSolutions := make(map[string]*LHSolution, len(solutions))

	for _, solution := range solutions {
		components := solution.Components
		arrCncs := make([]*LHComponent, 0, len(components))
		arrTvol := make([]*LHComponent, 0, len(components))
		arrSvol := make([]*LHComponent, 0, len(components))
		cmpvol := 0.0
		totalvol := 0.0

		for _, component := range components {
			// what sort of component is it?
			// what is the total volume ?
			if component.Conc != 0.0 {
				arrCncs = append(arrCncs, component)
			} else if component.Tvol != 0.0 {
				arrTvol = append(arrTvol, component)
				tv := component.Tvol
				if totalvol == 0.0 || totalvol == tv {
					totalvol = tv
				} else {
					// error
					wutil.Error(errors.New(fmt.Sprintf("Inconsistent total volumes %-6.4f and %-6.4f at component %s", totalvol, tv, component.Name)))
				}
			} else {
				cmpvol += component.Vol
				arrSvol = append(arrSvol, component)
			}
		}

		// first we add the volumes to the concentration components

		arrFinalComponents := make([]*LHComponent, 0, len(components))

		for _, component := range arrCncs {
			name := component.CName
			cnc := component.Conc
			vol := totalvol * cnc / stockconcs[name]
			cmpvol += vol
			component.Vol = vol
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// next we get the final volume for total volume components

		for _, component := range arrTvol {
			vol := totalvol - cmpvol
			component.Vol = vol
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// then we add the rest

		arrFinalComponents = append(arrFinalComponents, arrSvol...)

		// finally we replace the components in this solution

		solution.Components = arrFinalComponents

		// and put the new solution in the array

		newSolutions[solution.ID] = solution
	}

	return newSolutions, stockconcs
}
