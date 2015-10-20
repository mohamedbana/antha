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
// 2 Royal College St, London NW1 0NH UK

package wtype

import (
	"encoding/json"
//	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// functions to deal with how to serialize / deserialize the relevant objects.
// Despite the availability of good JSON serialization in Go it is necessary
// to include this to allow object structure to be sensibly defined for
// runtime purposes without making the network traffic too heavy

// serializable version of LHComponent
/*
type SLHComponent struct {
	*GenericPhysical
	ID          string
	Inst        string
	Order       int
	CName       string
	Type        string
	Vol         float64
	Conc        float64
	Vunit       string
	Cunit       string
	Tvol        float64
	Loc         string
	Smax        float64
	Visc        float64
	ContainerID string
	Destination string
}

func (lhc *LHComponent) MarshalJSON() ([]byte, error) {
	id := ""
	if lhc.LContainer != nil {
		id = lhc.LContainer.ID
	}

		slhc := SLHComponent{lhc.GenericPhysical, lhc.ID, lhc.Inst, lhc.Order, lhc.CName, lhc.Type, lhc.Vol, lhc.Conc, lhc.Vunit, lhc.Cunit, lhc.Tvol, lhc.Loc, lhc.Smax, lhc.Visc, id, lhc.Destination}
		return json.Marshal(slhc)


}

func (lhc *LHComponent) UnmarshalJSON(b []byte) error {
	var slhc SLHComponent
	err := json.Unmarshal(b, &slhc)

	if err != nil {
		return err
	}
	// fill in the component

	return err
}
*/
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
	Vol       float64
	Vunit     string
	Rvol      float64
	ShapeName string
	Bottom    int
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
}

func (w *LHWell) AddDimensions(lhwt *LHWellType) {
	w.Vol = lhwt.Vol
	w.Vunit = lhwt.Vunit
	w.Rvol = lhwt.Rvol
	w.WShape = NewShape(lhwt.ShapeName, lhwt.Dunit, lhwt.Xdim, lhwt.Ydim, lhwt.Zdim)
	w.Bottom = lhwt.Bottom
	w.Xdim = lhwt.Xdim
	w.Ydim = lhwt.Ydim
	w.Zdim = lhwt.Zdim
	w.Bottomh = lhwt.Bottomh
	w.Dunit = lhwt.Dunit
}

func (plate *LHPlate) Welldimensions() *LHWellType {
	t := plate.Welltype
	lhwt := LHWellType{t.Vol, t.Vunit, t.Rvol, t.WShape.ShapeName, t.Bottom, t.Xdim, t.Ydim, t.Zdim, t.Bottomh, t.Dunit}
	return &lhwt
}
//
//func (plate *LHPlate) MarshalJSON() ([]byte, error) {
//	slp := SLHPlate{plate.ID, plate.Inst, plate.Loc, plate.PlateName, plate.Type, plate.Mnfr, plate.WlsX, plate.WlsY, plate.Nwells, plate.Height, plate.Hunit, plate.Welltype, plate.Wellcoords, plate.Welldimensions()}
//
//	return json.Marshal(slp)
//}
//
//func (plate *LHPlate) UnmarshalJSON(b []byte) error {
//	var slp SLHPlate
//
//	e := json.Unmarshal(b, &slp)
//
//	if e != nil {
//		return e
//	}
//
//	// push the info into the plate
//
//	slp.FillPlate(plate)
//
//	// allocate and fill the other structures
//
//	plate.HWells = make(map[string]*LHWell, len(plate.Wellcoords))
//	plate.Rows = make([][]*LHWell, plate.WlsY)
//	plate.Cols = make([][]*LHWell, plate.WlsX)
//
//	wt := slp.Welldimensions
//
//	for s, w := range plate.Wellcoords {
//		// give w's contents their proper references
//
//		for _, contents := range w.WContents {
//			contents.LContainer = w
//		}
//
//		plate.HWells[w.ID] = w
//
//		// give w its properties back
//
//		w.AddDimensions(wt)
//		x, y := wutil.DecodeCoords(s)
//
//		if len(plate.Rows[x]) == 0 {
//			plate.Rows[x] = make([]*LHWell, plate.WlsX)
//		}
//		plate.Rows[x][y] = w
//
//		if len(plate.Cols[y]) == 0 {
//			plate.Cols[y] = make([]*LHWell, plate.WlsY)
//		}
//		plate.Cols[y][x] = w
//	}
//
//	// don't forget to add them back to the welltype!
//
//	plate.Welltype.AddDimensions(wt)
//
//	return e
//}

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
	for _, c := range lw.WContents {
		c.LContainer = lw
	}
}

//func (well LHWell) MarshalJSON() ([]byte, error) {
//	// make sure we don't cause an infinite loop
//	for _, c := range well.WContents {
//		c.LContainer = nil
//	}
//	slw := SLHWell{well.ID, well.Inst, well.Plateinst, well.Plateid, well.Crds, well.WContents, well.Currvol}
//	return json.Marshal(slw)
//}

//func (well *LHWell) UnmarshalJSON(ar []byte) error {
//	var slw SLHWell
//	err := json.Unmarshal(ar, &slw)
//
//	slw.FillWell(well)
//
//	return err
//}

type FromFactory struct {
	String string
}

func (f *FromFactory) MarshalJSON() ([]byte, error) {
	v, e := json.Marshal(f.String)
	return v, e
}

func (f *FromFactory) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	f.String = s
	return nil
}
