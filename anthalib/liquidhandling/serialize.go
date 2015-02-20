// anthalib//liquidhandling/serialize.go: Part of the Antha language
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

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/anthalib/wutil"
	"strconv"
)

// functions to deal with how to serialize / deserialize the relevant objects.
// Despite the availability of good JSON serialization in Go it is necessary
// to include this to allow object structure to be sensibly defined for
// runtime purposes without making the network traffic too heavy

// serializable, stripped-down version of the LHPlate
type SLHPlate struct {
	ID             string
	Inst           string
	Loc            string
	Name           string
	Type           string
	Mnfr           string
	WellsX         int
	WellsY         int
	Nwells         int
	Height         float64
	Hunit          string
	Welltype       *LHWell
	Wellcoords     map[string]*LHWell
	Welldimensions *LHWellType
}

func (slhp SLHPlate) FillPlate(plate *LHPlate) {
	plate.ID = slhp.ID
	plate.Inst = slhp.Inst
	plate.Loc = slhp.Loc
	plate.PlateName = slhp.Name
	plate.Type = slhp.Type
	plate.Mnfr = slhp.Mnfr
	plate.WlsX = slhp.WellsX
	plate.WlsY = slhp.WellsY
	plate.Nwells = slhp.Nwells
	plate.Height = slhp.Height
	plate.Hunit = slhp.Hunit
	plate.Welltype = slhp.Welltype
	plate.Wellcoords = slhp.Wellcoords
}

// this is for keeping track of the well type

type LHWellType struct {
	Vol     float64
	Vunit   string
	Rvol    float64
	Shape   int
	Bottom  int
	Xdim    float64
	Ydim    float64
	Zdim    float64
	Bottomh float64
	Dunit   string
}

func (w *LHWell) AddDimensions(lhwt *LHWellType) {
	w.Vol = lhwt.Vol
	w.Vunit = lhwt.Vunit
	w.Rvol = lhwt.Rvol
	w.Shape = lhwt.Shape
	w.Bottom = lhwt.Bottom
	w.Xdim = lhwt.Xdim
	w.Ydim = lhwt.Ydim
	w.Zdim = lhwt.Zdim
	w.Bottomh = lhwt.Bottomh
	w.Dunit = lhwt.Dunit
}

func (plate *LHPlate) Welldimensions() *LHWellType {
	t := plate.Welltype
	lhwt := LHWellType{t.Vol, t.Vunit, t.Rvol, t.Shape, t.Bottom, t.Xdim, t.Ydim, t.Zdim, t.Bottomh, t.Dunit}
	return &lhwt
}

func (plate *LHPlate) MarshalJSON() ([]byte, error) {
	slp := SLHPlate{plate.ID, plate.Inst, plate.Loc, plate.PlateName, plate.Type, plate.Mnfr, plate.WlsX, plate.WlsY, plate.Nwells, plate.Height, plate.Hunit, plate.Welltype, plate.Wellcoords, plate.Welldimensions()}

	return json.Marshal(slp)
}

func (plate *LHPlate) UnmarshalJSON(b []byte) error {
	var slp SLHPlate

	e := json.Unmarshal(b, &slp)

	if e != nil {
		return e
	}

	// push the info into the plate

	slp.FillPlate(plate)

	// allocate and fill the other structures

	plate.HWells = make(map[string]*LHWell, len(plate.Wellcoords))
	plate.Rows = make([][]*LHWell, plate.WlsY)
	plate.Cols = make([][]*LHWell, plate.WlsX)

	wt := slp.Welldimensions

	for s, w := range plate.Wellcoords {
		plate.HWells[w.ID] = w

		// give w its properties back

		w.AddDimensions(wt)
		x, y := wutil.DecodeCoords(s)

		if len(plate.Rows[x]) == 0 {
			plate.Rows[x] = make([]*LHWell, plate.WlsX)
		}
		plate.Rows[x][y] = w

		if len(plate.Cols[y]) == 0 {
			plate.Cols[y] = make([]*LHWell, plate.WlsY)
		}
		plate.Cols[y][x] = w
	}

	// don't forget to add them back to the welltype!

	plate.Welltype.AddDimensions(wt)

	return e
}

type SLHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Coords    string
	Contents  []*LHComponent
	Currvol   float64
}

func (slw SLHWell) FillWell(lw *LHWell) {
	lw.ID = slw.ID
	lw.Inst = slw.Inst
	lw.Plateinst = slw.Plateinst
	lw.Plateid = slw.Plateid
	lw.Crds = slw.Coords
	lw.WContents = slw.Contents
	lw.Currvol = slw.Currvol
}

func (well *LHWell) MarshalJSON() ([]byte, error) {
	slw := SLHWell{well.ID, well.Inst, well.Plateinst, well.Plateid, well.Crds, well.WContents, well.Currvol}
	return json.Marshal(slw)
}

func (well *LHWell) UnmarshalJSON(ar []byte) error {
	var slw SLHWell
	err := json.Unmarshal(ar, &slw)

	slw.FillWell(well)

	return err
}

// marshal / unmarshal methods for the top-level lhrequest class

type SLHRequest struct {
	ID                         string
	Output_solutions           map[string]*LHSolution
	Input_solutions            map[string][]*LHComponent
	Plates                     map[string]*LHPlate
	Tips                       []*LHTipbox
	Locats                     []string
	Setup                      LHSetup
	Instructions               []RobotInstruction
	Robotfn                    string
	Input_assignments          map[string][]string
	Output_plates              map[string]*LHPlate
	Input_platetype            *LHPlate
	Input_major_group_layouts  map[string][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[string]string
	Output_platetype           *LHPlate
	Output_major_group_layouts map[string][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[string]string
	Plate_lookup               map[string]int
	Stockconcs                 map[string]float64
}

func (req *LHRequest) MarshalJSON() ([]byte, error) {
	new_input_major_layouts := make(map[string][]string, len(req.Input_major_group_layouts))

	for k, v := range req.Input_major_group_layouts {
		new_input_major_layouts[strconv.Itoa(k)] = v
	}

	new_input_plate_layout := make(map[string]string, len(req.Input_plate_layout))

	for k, v := range req.Input_plate_layout {
		new_input_plate_layout[strconv.Itoa(k)] = v
	}
	new_output_major_layouts := make(map[string][]string, len(req.Output_major_group_layouts))

	for k, v := range req.Output_major_group_layouts {
		new_output_major_layouts[strconv.Itoa(k)] = v
	}
	new_output_plate_layout := make(map[string]string, len(req.Output_plate_layout))

	for k, v := range req.Output_plate_layout {
		new_output_plate_layout[strconv.Itoa(k)] = v
	}

	slhr := SLHRequest{req.ID, req.Output_solutions, req.Input_solutions, req.Plates, req.Tips, req.Locats, req.Setup, req.Instructions, req.Robotfn, req.Input_assignments, req.Output_plates, req.Input_platetype, new_input_major_layouts, req.Input_minor_group_layouts, new_input_plate_layout, req.Output_platetype, new_output_major_layouts, req.Output_minor_group_layouts, new_output_plate_layout, req.Plate_lookup, req.Stockconcs}

	return json.Marshal(slhr)
}

func (req *LHRequest) UnmarshalJSON(ar []byte) error {
	var slhr SLHRequest
	e := json.Unmarshal(ar, req)

	fmt.Println("ERR: ", e)

	e = json.Unmarshal(ar, slhr)

	req.Input_major_group_layouts = make(map[int][]string, len(slhr.Input_major_group_layouts))

	for k, v := range slhr.Input_major_group_layouts {
		req.Input_major_group_layouts[wutil.ParseInt(k)] = v
	}

	req.Input_plate_layout = make(map[int]string, len(slhr.Input_plate_layout))

	for k, v := range slhr.Input_plate_layout {
		req.Input_plate_layout[wutil.ParseInt(k)] = v
	}

	req.Output_major_group_layouts = make(map[int][]string, len(slhr.Output_major_group_layouts))

	for k, v := range slhr.Output_major_group_layouts {
		req.Output_major_group_layouts[wutil.ParseInt(k)] = v
	}

	req.Output_plate_layout = make(map[int]string, len(slhr.Output_plate_layout))

	for k, v := range slhr.Output_plate_layout {
		req.Output_plate_layout[wutil.ParseInt(k)] = v
	}

	return e
}
