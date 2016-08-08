// liquidhandling/lhtypes.Go: Part of the Antha language
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
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package wtype

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
    "math"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/logger"
)

// structure describing a microplate
type LHPlate struct {
	ID          string
	Inst        string
	Loc         string             // location of plate
	PlateName   string             // user-definable plate name
	Type        string             // plate type
	Mnfr        string             // manufacturer
	WlsX        int                // wells along long axis
	WlsY        int                // wells along short axis
	Nwells      int                // total number of wells
	HWells      map[string]*LHWell // map of well IDs to well
	Rows        [][]*LHWell
	Cols        [][]*LHWell
	Welltype    *LHWell
	Wellcoords  map[string]*LHWell // map of coords in A1 format to wells
	WellXOffset float64            // distance (mm) between well centres in X direction
	WellYOffset float64            // distance (mm) between well centres in Y direction
	WellXStart  float64            // offset (mm) to first well in X direction
	WellYStart  float64            // offset (mm) to first well in Y direction
	WellZStart  float64            // offset (mm) to bottom of well in Z direction
    size        Coordinates        // size of the plate (mm)
    offset      Coordinates        // (relative) position of the plate (mm), set by parent
    parent      LHObject
} 

func (lhp LHPlate) String() string {
	return fmt.Sprintf(
		`LHPlate {
	ID          : %s,
	Inst        : %s,
	Loc         : %s,
	PlateName   : %s,
	Type        : %s,
	Mnfr        : %s,
	WlsX        : %d,
	WlsY        : %d,
	Nwells      : %d,
	HWells      : %p,
	Rows        : %p,
	Cols        : %p,
	Welltype    : %p,
	Wellcoords  : %p,
	WellXOffset : %f,
	WellYOffset : %f,
	WellXStart  : %f,
	WellYStart  : %f,
	WellZStart  : %f,
	Size  : %f x %f x %f,
}`,
		lhp.ID,
		lhp.Inst,
		lhp.Loc,
		lhp.PlateName,
		lhp.Type,
		lhp.Mnfr,
		lhp.WlsX,
		lhp.WlsY,
		lhp.Nwells,
		lhp.HWells,
		lhp.Rows,
		lhp.Cols,
		lhp.Welltype,
		lhp.Wellcoords,
		lhp.WellXOffset,
		lhp.WellYOffset,
		lhp.WellXStart,
		lhp.WellYStart,
		lhp.WellZStart,
        lhp.size.X,
        lhp.size.Y,
        lhp.size.Z,
	)
}

// convenience method

func (lhp *LHPlate) GetComponent(cmp *LHComponent, exact bool) ([]WellCoords, bool) {
	ret := make([]WellCoords, 0, 1)

	it := NewOneTimeColumnWiseIterator(lhp)

	var volGot wunit.Volume
	volGot = wunit.NewVolume(0.0, "ul")

	x := 0

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]

		/*
			if !w.Empty() {
				logger.Debug(fmt.Sprint("WANT: ", cmp.CName, " :: ", wc.FormatA1(), " ", w.Contents().CName, " ", w.CurrVolume().ToString()))
			}
		*/
		if w.Contents().CName == cmp.CName {
			if exact && w.Contents().ID != cmp.ID {
				continue
			}
			x += 1

			v := w.WorkingVolume()
			if v.LessThan(cmp.Volume()) {
				continue
			}
			volGot.Add(v)
			ret = append(ret, wc)

			if volGot.GreaterThan(cmp.Volume()) || volGot.EqualTo(cmp.Volume()) {
				break
			}
		}
	}

	//	fmt.Println("FOUND: ", cmp.CName, " WANT ", cmp.Volume().ToString(), " GOT ", volGot.ToString(), "  ", ret)

	if !(volGot.GreaterThan(cmp.Volume()) || volGot.EqualTo(cmp.Volume())) {
		return ret, false
	}

	return ret, true
}

func (lhp *LHPlate) Wells() [][]*LHWell {
	return lhp.Rows
}
func (lhp *LHPlate) WellMap() map[string]*LHWell {
	return lhp.Wellcoords
}

func (lhp *LHPlate) AllWellPositions() (wellpositionarray []string) {

	wellpositionarray = make([]string, 0)

	// range through well coordinates
	for j := 0; j < lhp.WlsX; j++ {
		for i := 0; i < lhp.WlsY; i++ {
			wellposition := wutil.NumToAlpha(i+1) + strconv.Itoa(j+1)
			wellpositionarray = append(wellpositionarray, wellposition)
		}
	}
	return
}

// @implement named

func (lhp *LHPlate) GetName() string {
	return lhp.PlateName
}

// @implement Typed
func (lhp *LHPlate) GetType() string {
    return lhp.Type
} 

func (lhp *LHPlate) WellAt(wc WellCoords) *LHWell {
	return lhp.Wellcoords[wc.FormatA1()]
}

func (lhp *LHPlate) WellAtString(s string) (*LHWell, bool) {
	// improve later, start by assuming these are in FormatA1()
	w, ok := lhp.Wellcoords[s]

	return w, ok
}

func (lhp *LHPlate) WellsX() int {
	return lhp.WlsX
}

func (lhp *LHPlate) WellsY() int {
	return lhp.WlsY
}

func (lhp *LHPlate) NextEmptyWell(it PlateIterator) WellCoords {
	c := 0
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		if c == lhp.Nwells {
			// prevent iterators from ever making this loop infinitely
			break
		}

		if lhp.Cols[wc.X][wc.Y].Empty() {
			return wc
		}
	}

	return ZeroWellCoords()
}

func NewLHPlate(platetype, mfr string, nrows, ncols int, size Coordinates, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	var lhp LHPlate
	lhp.Type = platetype
	lhp.ID = GetUUID()
    lhp.PlateName = fmt.Sprintf("%s_%s", platetype, lhp.ID[1:len(lhp.ID)-2])
	lhp.Mnfr = mfr
	lhp.WlsX = ncols
	lhp.WlsY = nrows
	lhp.Nwells = ncols * nrows
	lhp.Welltype = welltype
	lhp.WellXOffset = wellXOffset
	lhp.WellYOffset = wellYOffset
	lhp.WellXStart = wellXStart
	lhp.WellYStart = wellYStart
	lhp.WellZStart = wellZStart
    lhp.size = size

	wellcoords := make(map[string]*LHWell, ncols*nrows)

	// make wells
	rowarr := make([][]*LHWell, nrows)
	colarr := make([][]*LHWell, ncols)
	arr := make([][]*LHWell, nrows)
	wellmap := make(map[string]*LHWell, ncols*nrows)

	for i := 0; i < nrows; i++ {
		arr[i] = make([]*LHWell, ncols)
		rowarr[i] = make([]*LHWell, ncols)
		for j := 0; j < ncols; j++ {
			if colarr[j] == nil {
				colarr[j] = make([]*LHWell, nrows)
			}
			arr[i][j] = welltype.Dup()

			//crds := wutil.NumToAlpha(i+1) + ":" + strconv.Itoa(j+1)
			crds := WellCoords{j, i}.FormatA1()
			wellcoords[crds] = arr[i][j]
			arr[i][j].Crds = crds
			colarr[j][i] = arr[i][j]
			rowarr[i][j] = arr[i][j]
			wellmap[arr[i][j].ID] = arr[i][j]
			arr[i][j].Plate = &lhp
			arr[i][j].Plateinst = lhp.Inst
			arr[i][j].Plateid = lhp.ID
			arr[i][j].Platetype = lhp.Type
			arr[i][j].Crds = crds
		}
	}

	lhp.Wellcoords = wellcoords
	lhp.HWells = wellmap
	lhp.Cols = colarr
	lhp.Rows = rowarr

	return &lhp
}

func (lhp *LHPlate) Dup() *LHPlate {
	ret := NewLHPlate(lhp.Type, lhp.Mnfr, lhp.WlsY, lhp.WlsX, lhp.GetSize(), lhp.Welltype, lhp.WellXOffset, lhp.WellYOffset, lhp.WellXStart, lhp.WellYStart, lhp.WellZStart)
    
	ret.PlateName = lhp.PlateName

	ret.HWells = make(map[string]*LHWell, len(ret.HWells))

	for i, row := range lhp.Rows {
		for j, well := range row {
			d := well.Dup()
			ret.Rows[i][j] = d
			ret.Cols[j][i] = d
			ret.Wellcoords[d.Crds] = d
			ret.HWells[d.ID] = d
		}
	}

	return ret
}

func (p *LHPlate) ProtectAllWells() {
	for _, v := range p.Wellcoords {
		v.Protect()
	}
}

func (p *LHPlate) UnProtectAllWells() {
	for _, v := range p.Wellcoords {
		v.UnProtect()
	}
}

func Initialize_Wells(plate *LHPlate) {
	id := (*plate).ID
	wells := (*plate).HWells
	newwells := make(map[string]*LHWell, len(wells))
	wellcrds := (*plate).Wellcoords
	for _, well := range wells {
		well.ID = GetUUID()
		well.Plateid = id
		newwells[well.ID] = well
		wellcrds[well.Crds] = well
	}
	(*plate).HWells = newwells
	(*plate).Wellcoords = wellcrds
}

func (p *LHPlate) RemoveComponent(well string, vol wunit.Volume) *LHComponent {
	w := p.Wellcoords[well]

	if w == nil {
		logger.Debug(fmt.Sprint("RemoveComponent (plate) ERROR: ", well, " ", vol.ToString(), " Can't find well"))
		return nil
	}

	err := w.Remove(vol)

	return err
}

func (p *LHPlate) DeclareTemporary() {
	for _, w := range p.Wellcoords {
		w.DeclareTemporary()
	}
}

func (p *LHPlate) IsTemporary() bool {
	for _, w := range p.Wellcoords {
		if !w.IsTemporary() {
			return false
		}
	}

	return true
}

func (p *LHPlate) DeclareAutoallocated() {
	for _, w := range p.Wellcoords {
		w.DeclareAutoallocated()
	}
}

func (p *LHPlate) IsAutoallocated() bool {
	for _, w := range p.Wellcoords {
		if !w.IsAutoallocated() {
			return false
		}
	}

	return true
}

func ExportPlateCSV(outputpilename string, plate *LHPlate, platename string, wells []string, liquids []*LHComponent, Volumes []wunit.Volume) error {

	csvfile, err := os.Create(outputpilename)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer csvfile.Close()

	records := make([][]string, 0)

	//record := make([]string, 0)

	headerrecord := []string{plate.Type, platename, "", "", ""}

	records = append(records, headerrecord)

	for i, well := range wells {

		volfloat := Volumes[i].RawValue()

		volstr := strconv.FormatFloat(volfloat, 'G', -1, 64)

		/*
			fmt.Println("len(wells)", len(wells))
			fmt.Println("len(liquids)", len(liquids))
			fmt.Println("len(Volumes)", len(Volumes))
		*/

		record := []string{well, liquids[i].CName, liquids[i].TypeName(), volstr, Volumes[i].Unit().PrefixedSymbol()}
		records = append(records, record)
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, record := range records {

		err = csvwriter.Write(record)

		if err != nil {
			return err
		}
	}
	csvwriter.Flush()

	return err
}
func (p *LHPlate) SetConstrained(platform string, positions []string) {
	p.Welltype.Extra[platform] = positions
}

func (p *LHPlate) IsConstrainedOn(platform string) ([]string, bool) {
	var pos []string

	par, ok := p.Welltype.Extra[platform]

	if ok {
		pos = par.([]string)
		return pos, true
	}

	return pos, false

}

//##############################################
//@implement LHObject
//##############################################

func (self *LHPlate) GetOffset() Coordinates {
    if self.parent != nil {
        return self.offset.Add(self.parent.GetOffset())
    }
    return self.offset
}

func (self *LHPlate) SetOffset(o Coordinates) {
    self.offset = o
}

func (self *LHPlate) GetSize() Coordinates {
    return self.size
}

func (self *LHPlate) GetBounds() *BBox {
    r := BBox{self.GetOffset(), self.GetSize()}
    return &r
}

func (self *LHPlate) SetParent(p LHObject) {
    self.parent = p
}

func (self *LHPlate) GetParent() LHObject {
    return self.parent
}

//##############################################
//@implement Addressable
//##############################################

func (self *LHPlate) HasCoords(c WellCoords) bool {
    return c.X >= 0 &&
           c.Y >= 0 &&
           c.X < self.WlsX &&
           c.Y < self.WlsY
}

func (self *LHPlate) GetCoords(c WellCoords) (interface{}, bool) {
    if !self.HasCoords(c) {
        return nil, false
    }
    return self.Cols[c.X][c.Y], true
}

func (self *LHPlate) CoordsToWellCoords(r Coordinates) (WellCoords, Coordinates) {
    wc := WellCoords{
        int(math.Floor(((r.X-self.WellXStart) / self.WellXOffset))),// + 0.5), Don't need to add .5 because
        int(math.Floor(((r.Y-self.WellYStart) / self.WellYOffset))),// + 0.5), WellXStart is to edge, not center
    }
    if wc.X < 0 {
        wc.X = 0
    } else if wc.X >= self.WlsX {
        wc.X = self.WlsX - 1
    }
    if wc.Y < 0 {
        wc.Y = 0
    } else if wc.Y >= self.WlsY {
        wc.Y = self.WlsY - 1
    }

    r2, _ := self.WellCoordsToCoords(wc, TopReference)

    return wc, r.Subtract(r2)
}

func (self *LHPlate) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
    if !self.HasCoords(wc) {
        return Coordinates{}, false
    }

    var z float64
    if r == BottomReference {
        z = self.WellZStart
    } else if r == TopReference {
        z = self.size.Z
    } else if r == LiquidReference {
        panic("Haven't implemented liquid level yet")
    }

    return Coordinates{
        self.WellXStart + (float64(wc.X)+0.5) * self.WellXOffset,
        self.WellYStart + (float64(wc.Y)+0.5) * self.WellYOffset,
        z}, true
}
