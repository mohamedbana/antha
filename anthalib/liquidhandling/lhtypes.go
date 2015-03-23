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
// 1 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package liquidhandling

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wutil"
	"strconv"
	"strings"
)

// describes a liquid handler, its capabilities and current state
type LHProperties struct {
	ID                 string
	Nposns             int
	Positions          []*LHPosition
	Model              string
	Manfr              string
	LHType             string
	TPType             string
	Formats            []string
	Cnfvol             []*LHParameter
	CurrConf           *LHParameter
	Tip_preferences    []int
	Input_preferences  []int
	Output_preferences []int
}

// constructor for the above
func NewLHProperties(num_positions int, model, manufacturer, lhtype, tptype string, formats []string) *LHProperties {
	var lhp LHProperties

	lhp.Nposns = num_positions

	lhp.Model = model
	lhp.Manfr = manufacturer

	positions := make([]*LHPosition, num_positions)

	for i := 0; i < num_positions; i++ {
		// not overriding these defaults seems like a
		// bad idea --- TODO: Fix, e.g., MAXH here
		positions[i] = NewLHPosition(i+1, "position_"+strconv.Itoa(i+1), 50.0)
	}

	lhp.Positions = positions

	lhp.Cnfvol = make([]*LHParameter, 2)

	// lhp.Curcnf, lhp.Cmnvol etc. intentionally left blank

	return &lhp
}

// describes sets of parameters which can be used to create a configuration
type LHParameter struct {
	ID      string
	Name    string
	Minvol  float64
	Maxvol  float64
	Volunit string
	Policy  LHPolicy
}

func NewLHParameter(name string, minvol, maxvol float64, volunit string) *LHParameter {
	var lhp LHParameter
	lhp.ID = wtype.GetUUID()
	lhp.Name = name
	lhp.Minvol = minvol
	lhp.Maxvol = maxvol
	lhp.Volunit = volunit
	lhp.Policy = make(LHPolicy, 1)
	return &lhp
}

// map for anything else

type LHPolicy map[string]interface{}

// defines an addendum to a liquid handler
// not much to say yet

type LHDevice struct {
	ID   string
	Name string
	Mnfr string
}

func NewLHDevice(name, mfr string) *LHDevice {
	var dev LHDevice
	dev.ID = wtype.GetUUID()
	dev.Name = name
	dev.Mnfr = mfr
	return &dev
}

// describes a position on the liquid handling deck and its current state
type LHPosition struct {
	ID    string
	Name  string
	Num   int
	Extra []LHDevice
	Maxh  float64
}

func NewLHPosition(position_number int, name string, maxh float64) *LHPosition {
	var lhp LHPosition
	lhp.ID = wtype.GetUUID()
	lhp.Name = name
	lhp.Num = position_number
	lhp.Extra = make([]LHDevice, 0, 2)
	lhp.Maxh = maxh
	return &lhp
}

/*
// question over whether this is necessary
//@implement wtype.SolidContainer
func (lhp *LHPosition) Contents() []Solid {
	return nil
}
func (lhp *LHPosition) ContainerType() string {
	return lhp.Name
}
func Empty() bool {

}
func PartOf() Entity {

}
*/
// structure for defining a request to the liquid handler
type LHRequest struct {
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
	Output_assignments         []string
	Input_plates               map[string]*LHPlate
	Output_plates              map[string]*LHPlate
	Input_platetype            *LHPlate
	Input_major_group_layouts  map[int][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[int]string
	Output_platetype           *LHPlate
	Output_major_group_layouts map[int][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[int]string
	Plate_lookup               map[string]int
	Stockconcs                 map[string]float64
}

func NewLHRequest() *LHRequest {
	var lhr LHRequest
	lhr.ID = wtype.GetUUID()
	lhr.Output_solutions = make(map[string]*LHSolution)
	lhr.Input_solutions = make(map[string][]*LHComponent)
	lhr.Plates = make(map[string]*LHPlate)
	lhr.Tips = make([]*LHTipbox, 1)
	lhr.Locats = make([]string, 1)
	lhr.Instructions = make([]RobotInstruction, 1)
	lhr.Input_plates = make(map[string]*LHPlate)
	lhr.Output_plates = make(map[string]*LHPlate)
	lhr.Input_major_group_layouts = make(map[int][]string)
	lhr.Input_minor_group_layouts = make([][]string, 1)
	lhr.Output_major_group_layouts = make(map[int][]string)
	lhr.Output_minor_group_layouts = make([][]string, 1)
	lhr.Output_plate_layout = make(map[int]string)
	lhr.Plate_lookup = make(map[string]int)
	lhr.Stockconcs = make(map[string]float64)
	return &lhr
}

// structure describing a solution: a combination of liquid components
type LHSolution struct {
	*wtype.GenericPhysical
	ID               string
	Inst             string
	SName            string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	Platetype        string
	Vol              float64
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Minorlayoutgroup int
}

func NewLHSolution() *LHSolution {
	var lhs LHSolution
	lhs.ID = wtype.GetUUID()
	var gp wtype.GenericPhysical
	lhs.GenericPhysical = &gp
	return &lhs
}

func (sol LHSolution) GetComponentVolume(key string) float64 {
	vol := 0.0

	for _, v := range sol.Components {
		if v.CName == key {
			vol += v.Vol
		}
	}

	return vol
}

// @implement wtype.LHSolution

// structure describing a liquid component and its desired properties
type LHComponent struct {
	*wtype.GenericPhysical
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
	LContainer  *LHWell
	Destination string
}

// @implement wtype.Liquid

func (lhc *LHComponent) Viscosity() float64 {
	return lhc.Visc
}

func (lhc *LHComponent) Name() string {
	return lhc.CName
}

func (lhc *LHComponent) Container() wtype.LiquidContainer {
	return lhc.LContainer
}

func (lhc *LHComponent) Sample(v wunit.Volume) wtype.Liquid {
	// need to jig around with units a bit here
	// Should probably just make Vunit, Cunit etc. wunits anyway
	meas := wunit.ConcreteMeasurement{lhc.Vol, wunit.ParsePrefixedUnit(lhc.Vunit)}

	// we need some logic potentially

	if v.SIValue() > meas.SIValue() {
		wutil.Error(errors.New(fmt.Sprintf("LHComponent ID: %s Not enough volume for sample", lhc.ID)))
	} else if v.SIValue() == meas.SIValue() {
		return lhc
	}
	smp := CopyLHComponent(lhc)
	// need a convention here

	smp.Vol = v.RawValue()
	smp.Vunit = v.Unit().PrefixedSymbol()
	meas.Subtract(&v.ConcreteMeasurement)
	lhc.Vol = meas.RawValue()
	return smp
}

func NewLHComponent() *LHComponent {
	var lhc LHComponent
	var gp wtype.GenericPhysical
	lhc.GenericPhysical = &gp
	lhc.ID = wtype.GetUUID()
	return &lhc
}

func CopyLHComponent(lhc *LHComponent) *LHComponent {
	tmp, _ := json.Marshal(lhc)
	var lhc2 LHComponent
	json.Unmarshal(tmp, &lhc2)
	lhc2.ID = wtype.GetUUID()
	if lhc2.Inst != "" {
		lhc2.Inst = wtype.GetUUID()
		// this needs some thought
	}
	return &lhc2
}

// structure defining a liquid handler setup

type LHSetup map[string]interface{}

func NewLHSetup() LHSetup {
	return make(LHSetup, 10)
}

// structure describing a microplate
// this needs to be harmonised with the wtype version
type LHPlate struct {
	*wtype.GenericEntity
	ID         string
	Inst       string
	Loc        string
	PlateName  string
	Type       string
	Mnfr       string
	WlsX       int
	WlsY       int
	Nwells     int
	HWells     map[string]*LHWell
	Height     float64
	Hunit      string
	Rows       [][]*LHWell
	Cols       [][]*LHWell
	Welltype   *LHWell
	Wellcoords map[string]*LHWell
}

// @implement wtype.Location

func (lhp *LHPlate) Location_ID() string {
	return lhp.ID
}

func (lhp *LHPlate) Location_Name() string {
	return lhp.PlateName
}

func (lhp *LHPlate) Positions() []wtype.Location {
	ret := make([]wtype.Location, lhp.Nwells)
	x := 0
	for _, v := range lhp.Cols {
		for _, w := range v {
			ret[x] = wtype.Location(w)
			x += 1
		}
	}
	return ret
}

func (lhp *LHPlate) Container() wtype.Location {
	return lhp
}

// Shape() deferred to GenericPhysical

// @implement wtype.Labware
func (lhp *LHPlate) Wells() [][]wtype.Well {
	ret := make([][]wtype.Well, len(lhp.Rows))
	for i := 0; i < len(lhp.Rows); i++ {
		ret[i] = make([]wtype.Well, len(lhp.Rows[i]))
		for j := 0; j < len(lhp.Rows[i]); j++ {
			ret[i][j] = lhp.Rows[i][j]
		}
	}
	return ret
}

func (lhp *LHPlate) WellAt(crds wtype.WellCoords) wtype.Well {
	return lhp.Cols[crds.X][crds.Y]
}

func (lhp *LHPlate) WellsX() int {
	return lhp.WlsX
}

func (lhp *LHPlate) WellsY() int {
	return lhp.WlsY

}

func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell) *LHPlate {
	var lhp LHPlate
	lhp.Type = platetype
	lhp.ID = wtype.GetUUID()
	lhp.Mnfr = mfr
	lhp.WlsX = ncols
	lhp.WlsY = nrows
	lhp.Nwells = ncols * nrows
	lhp.Height = height
	lhp.Hunit = hunit
	lhp.Welltype = welltype

	wellcoords := make(map[string]*LHWell, ncols*nrows)

	// make wells
	rowarr := make([][]*LHWell, nrows)
	colarr := make([][]*LHWell, ncols)
	arr := make([][]*LHWell, ncols)
	wellmap := make(map[string]*LHWell, ncols*nrows)

	for i := 0; i < ncols; i++ {
		arr[i] = make([]*LHWell, nrows)
		colarr[i] = make([]*LHWell, nrows)
		for j := 0; j < nrows; j++ {
			if rowarr[j] == nil {
				rowarr[j] = make([]*LHWell, ncols)
			}
			arr[i][j] = NewLHWellCopy(welltype)

			crds := wutil.NumToAlpha(j+1) + ":" + strconv.Itoa(i+1)
			wellcoords[crds] = arr[i][j]
			arr[i][j].Crds = crds
			colarr[i][j] = arr[i][j]
			rowarr[j][i] = arr[i][j]
			wellmap[arr[i][j].ID] = arr[i][j]
			// fill in necessary bits of callback info

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

// structure representing a well on a microplate - description of a destination
type LHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Platetype string
	Crds      string
	Vol       float64
	Vunit     string
	WContents []*LHComponent
	Rvol      float64
	Currvol   float64
	WShape    wtype.Shape
	Bottom    int
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
	Plate     *LHPlate
}

func (w *LHWell) updateVolume() {
	w.Vol = 0.0
	for _, val := range w.WContents {
		w.Vol += val.Vol
	}
}

//@implement wtype.Location

func (lhw *LHWell) Location_ID() string {
	return lhw.ID
}

func (lhw *LHWell) Location_Name() string {
	return lhw.Platetype
}

func (lhw *LHWell) Positions() []wtype.Location {
	return nil
}

func (lhw *LHWell) Container() wtype.Location {
	return lhw.Plate
}

func (lhw *LHWell) Shape() wtype.Shape {
	return lhw.WShape
}

// @implement wtype.Well
func (w *LHWell) WellTypeName() string {
	return w.Platetype
}

func (w *LHWell) ResidualVolume() wunit.Volume {
	return wunit.Volume{wunit.NewMeasurement(w.Rvol, "", w.Vunit)}
}

func (w *LHWell) Coords() wtype.WellCoords {
	return wtype.MakeWellCoordsXY(w.Crds)
}

func (w *LHWell) ContainerVolume() wunit.Volume {
	return wunit.Volume{wunit.NewMeasurement(w.Vol, "", w.Vunit)}
}

func (w *LHWell) Contents() []wtype.Physical {
	ret := make([]wtype.Physical, len(w.WContents))
	for i := 0; i < len(w.WContents); i++ {
		ret[i] = wtype.Physical(w.WContents[i])
	}
	return ret
}

func (w *LHWell) Add(p wtype.Physical) {
	switch t := p.(type) {
	default:
		wutil.Error(errors.New(fmt.Sprintf("LHWell: Cannot add type %T", t)))
	case *LHSolution:
		// do something
	case *LHComponent:
		w.WContents = append(w.WContents, p.(*LHComponent))
	}

}

// this is pretty dodgy... we will have to be quite careful here
// the core problem is how to maintain a list of components and volumes
// but respect the physical fact that we can't actually unmix things
func (w *LHWell) Remove(v wunit.Volume) wtype.Physical {
	defer w.updateVolume()
	ret := w.WContents[0]

	if ret.Vol > v.SIValue() {
		ret.Vol = v.SIValue()
		w.WContents[0].Vol -= v.SIValue()
	} else {
		w.WContents = w.WContents[1:len(w.WContents)]
	}
	return ret
}

func (w *LHWell) ContainerType() string {
	return w.Platetype
}

func (w *LHWell) PartOf() wtype.Entity {
	return w.Plate
}

func (w *LHWell) Empty() bool {
	if w.Vol <= 0.000001 {
		return true
	} else {
		return false
	}
}

func NewLHWellCopy(template *LHWell) *LHWell {
	cp := NewLHWell(template.Platetype, template.Plateid, template.Crds, template.Vol, template.Rvol, template.WShape, template.Bottom, template.Xdim, template.Ydim, template.Zdim, template.Bottomh, template.Dunit)

	return cp
}

// make a new well structure
func NewLHWell(platetype, plateid, crds string, vol, rvol float64, shape wtype.Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	var well LHWell
	well.ID = wtype.GetUUID()
	well.Platetype = platetype
	well.Plateid = plateid
	well.Crds = crds
	well.Vol = vol
	well.Rvol = rvol
	well.Currvol = 0.0
	well.WShape = shape
	well.Bottom = bott
	well.Xdim = xdim
	well.Ydim = ydim
	well.Zdim = zdim
	well.Bottomh = bottomh
	well.Dunit = dunit
	return &well
}

func get_next_well(plate *LHPlate, component *LHComponent, curwell *LHWell) (*LHWell, bool) {
	nrow, ncol := 0, 1

	vol := component.Vol

	if curwell != nil {
		// quick check to see if we have room
		vol_left := get_vol_left(curwell)

		if vol_left >= vol {
			// fine we can just return this one
			return curwell, true
		}

		// we need a defined traversal of the wells

		crds := curwell.Crds

		tx := strings.Split(crds, ":")

		nrow = wutil.AlphaToNum(tx[0])
		ncol = wutil.ParseInt(tx[1])
	}

	wellsx := plate.WlsX
	wellsy := plate.WlsY

	var new_well *LHWell

	for {
		nrow, ncol = next_well_to_try(nrow, ncol, wellsy, wellsx)

		if nrow == -1 {
			return nil, false
		}
		crds := wutil.NumToAlpha(nrow) + ":" + strconv.Itoa(ncol)

		new_well = plate.Wellcoords[crds]

		cnts := new_well.WContents

		if len(cnts) == 0 {
			break
		}

		cont := cnts[0].Name()
		if cont != component.Name() {
			continue
		}

		vol_left := get_vol_left(new_well)
		if vol < vol_left {
			break
		}
	}

	return new_well, true
}

func get_vol_left(well *LHWell) float64 {
	cnts := well.WContents

	// in the first instance we have a fixed constant times the number of
	// transfers... volumes are in microlitres as always

	carry_vol := 10.0 // microlitres
	total_carry_vol := float64(len(cnts)) * carry_vol
	currvol := well.Currvol
	rvol := well.Rvol
	vol := well.Vol
	return vol - (currvol + total_carry_vol + rvol)
}

func next_well_to_try(row, col, nrows, ncols int) (int, int) {
	// this needs to be refactored into an iterator

	nrow := -1
	ncol := -1

	// iterate down columns

	if row+1 > nrows {
		if col+1 <= ncols {
			nrow = 1
			ncol = col + 1
		}
	} else {
		ncol = col
		nrow = row + 1
	}

	// note that the default should be to leave ncol/nrow unchanged
	// and return -1 -1

	return nrow, ncol
}

func new_component(name, ctype string, vol float64) *LHComponent {
	var component LHComponent
	component.ID = wtype.GetUUID()
	component.CName = name
	component.Type = ctype
	component.Vol = vol
	return &component
}

func new_solution() *LHSolution {
	var solution LHSolution
	solution.ID = wtype.GetUUID()
	solution.Components = make([]*LHComponent, 0, 4)
	return &solution
}

func new_plate(platetype *LHPlate) *LHPlate {
	new_plate := NewLHPlate(platetype.Type, platetype.Mnfr, platetype.WlsY, platetype.WlsX, platetype.Height, platetype.Hunit, platetype.Welltype)
	initialize_wells(new_plate)
	return new_plate
}

func initialize_wells(plate *LHPlate) {
	id := (*plate).ID
	wells := (*plate).HWells
	newwells := make(map[string]*LHWell, len(wells))
	wellcrds := (*plate).Wellcoords
	for _, well := range wells {
		well.ID = wtype.GetUUID()
		well.Plateid = id
		well.Currvol = 0.0
		newwells[well.ID] = well
		wellcrds[well.Crds] = well
	}
	(*plate).HWells = newwells
	(*plate).Wellcoords = wellcrds
}

/* tip box */

type LHTipbox struct {
	*wtype.GenericSolid
	ID    string
	Type  string
	Mnfr  string
	Nrows int
	Ncols int
	Tips  map[string]*LHTipholder
	Loc   wtype.Location
}

// @implement wtype.Labware
// most methods deferred to GenericSolid
func (lhtb *LHTipbox) Location() wtype.Location {
	return lhtb.Loc
}

func (lhtb *LHTipbox) Manufacturer() string {
	return lhtb.Mnfr
}

func (lhtb *LHTipbox) LabwareType() string {
	return lhtb.Type
}

// @implement wtype.Location
func (lhtb *LHTipbox) Location_ID() string {
	return lhtb.ID
}

func (lhtb *LHTipbox) Location_Name() string {
	return lhtb.Name()
}

func (lhtb *LHTipbox) Positions() []wtype.Location {
	ret := make([]wtype.Location, len(lhtb.Tips))
	i := 0
	for _, t := range lhtb.Tips {
		ret[i] = wtype.Location(t)
		i += 1
	}

	return ret
}

func (lhtb *LHTipbox) Container() wtype.Location {
	return lhtb.Loc
}

// Shape() deferred to wtype.GenericPhysical

func NewLHTipbox(nrows, ncols int, manufacturer string, tiptype *LHTip) *LHTipbox {
	var tipbox LHTipbox
	tipbox.ID = wtype.GetUUID()
	tipbox.Mnfr = manufacturer
	tipbox.Nrows = nrows
	tipbox.Ncols = ncols
	return initialize_tips(&tipbox, tiptype)
}

type LHTipholder struct {
	ID       string
	ParentID string
	Cnts     []*LHTip
	Parent   *LHTipbox
}

// @implement SolidContainer

func (lht *LHTipholder) Contents() []wtype.Solid {
	ret := make([]wtype.Solid, len(lht.Cnts))

	for i, t := range lht.Cnts {
		ret[i] = wtype.Solid(t)
	}

	return ret
}

func (lht *LHTipholder) ContainerType() string {
	return "TipHolder"
}

func (lht *LHTipholder) Empty() bool {
	if lht.Cnts[0] != nil {
		return false
	}
	return true
}

func (lht *LHTipholder) PartOf() wtype.Entity {
	return lht.Parent
}

// @implement Location

func (lht *LHTipholder) Location_ID() string {
	return lht.ID
}

func (lht *LHTipholder) Location_Name() string {
	// TODO -- add the proper location name here
	return "LHTipHolder"
}

func (lht *LHTipholder) Positions() []wtype.Location {
	return nil
}

func (lht *LHTipholder) Container() wtype.Location {
	return lht.Parent
}

func (lht *LHTipholder) Shape() wtype.Shape {
	// TODO this should return the right answer
	return nil
}

func NewLHTipholder(parentid string) *LHTipholder {
	var holder LHTipholder
	holder.ID = wtype.GetUUID()
	holder.ParentID = parentid
	holder.Cnts = make([]*LHTip, 1, 1)
	return &holder
}

func initialize_tips(tipbox *LHTipbox, tiptype *LHTip) *LHTipbox {
	nr := tipbox.Nrows
	nc := tipbox.Ncols
	wells := make(map[string]*LHTipholder, nr*nc)
	id := tipbox.ID
	wellcoords := make(map[string]string, nr*nc)

	for i := 0; i < nr; i++ {
		row := wutil.NumToAlpha(i + 1)
		for j := 0; j < nc; j++ {
			col := strconv.Itoa(j + 1)
			coords := row + ":" + col
			holder := NewLHTipholder(id)
			cnts := holder.Cnts
			cnts[0] = new_tip_copy(tiptype)
			holder.Cnts = cnts
			wells[holder.ID] = holder
			wellcoords[holder.ID] = coords
		}
	}
	tipbox.Tips = wells
	return tipbox
}

type LHTip struct {
	*wtype.GenericSolid
	ID       string
	Mnfr     string
	Type     string
	Minvol   float64
	Maxvol   float64
	Curvol   float64
	Contents string
	Dirty    bool
	Loc      wtype.Location
}

// @implement wtype.Labware

func (lht *LHTip) Manufacturer() string {
	return lht.Mnfr
}

func (lht *LHTip) LabwareType() string {
	return lht.Type
}

// @implement wtype.Entity
func (lht *LHTip) Location() wtype.Location {
	return lht.Loc
}

func NewLHTip(manufacturer, tiptype string, minvol, maxvol float64) *LHTip {
	var tip LHTip
	tip.ID = wtype.GetUUID()
	tip.Mnfr = manufacturer
	tip.Type = tiptype
	tip.Minvol = minvol
	tip.Maxvol = maxvol
	tip.Curvol = 0.0
	tip.Contents = ""
	tip.Dirty = false
	return &tip
}

func new_tip_copy(parent *LHTip) *LHTip {
	tip := NewLHTip(parent.Mnfr, parent.Type, parent.Minvol, parent.Maxvol)
	return tip
}

type RobotInstruction interface {
	InstructionType() int
	GetParameter(name string) interface{}
}
