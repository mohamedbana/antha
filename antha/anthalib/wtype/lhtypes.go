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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/logger"
)

const (
	LHVChannel = iota // vertical orientation
	LHHChannel        // horizontal orientation
)

// describes sets of parameters which can be used to create a configuration
type LHChannelParameter struct {
	ID          string
	Name        string
	Minvol      wunit.Volume
	Maxvol      wunit.Volume
	Minspd      wunit.FlowRate
	Maxspd      wunit.FlowRate
	Multi       int
	Independent bool
	Orientation int
	Head        int
}

func (lhcp LHChannelParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          string
		Name        string
		Minvol      wunit.Volume
		Maxvol      wunit.Volume
		Minspd      wunit.FlowRate
		Maxspd      wunit.FlowRate
		Multi       int
		Independent bool
		Orientation int
		Head        int
	}{
		lhcp.ID,
		lhcp.Name,
		lhcp.Minvol,
		lhcp.Maxvol,
		lhcp.Minspd,
		lhcp.Maxspd,
		lhcp.Multi,
		lhcp.Independent,
		lhcp.Orientation,
		lhcp.Head,
	})
}

func (lhcp *LHChannelParameter) Dup() *LHChannelParameter {
	r := NewLHChannelParameter(lhcp.Name, lhcp.Minvol, lhcp.Maxvol, lhcp.Minspd, lhcp.Maxspd, lhcp.Multi, lhcp.Independent, lhcp.Orientation, lhcp.Head)

	return r
}

func NewLHChannelParameter(name string, minvol, maxvol wunit.Volume, minspd, maxspd wunit.FlowRate, multi int, independent bool, orientation int, head int) *LHChannelParameter {
	var lhp LHChannelParameter
	lhp.ID = GetUUID()
	lhp.Name = name
	lhp.Minvol = minvol
	lhp.Maxvol = maxvol
	lhp.Minspd = minspd
	lhp.Maxspd = maxspd
	lhp.Multi = multi
	lhp.Independent = independent
	lhp.Orientation = orientation
	lhp.Head = head
	return &lhp
}

func (lhcp *LHChannelParameter) MergeWithTip(tip *LHTip) *LHChannelParameter {
	lhcp2 := *lhcp
	if tip.MinVol.GreaterThan(lhcp2.Minvol) {
		lhcp2.Minvol = wunit.CopyVolume(tip.MinVol)
	}

	if tip.MaxVol.LessThan(lhcp2.Maxvol) {
		lhcp2.Maxvol = wunit.CopyVolume(tip.MaxVol)
	}

	return &lhcp2
}

// defines an addendum to a liquid handler
// not much to say yet

type LHDevice struct {
	ID   string
	Name string
	Mnfr string
}

func NewLHDevice(name, mfr string) *LHDevice {
	var dev LHDevice
	dev.ID = GetUUID()
	dev.Name = name
	dev.Mnfr = mfr
	return &dev
}

func (lhd *LHDevice) Dup() *LHDevice {
	d := NewLHDevice(lhd.Name, lhd.Mnfr)
	return d
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
	lhp.ID = GetUUID()
	lhp.Name = name
	lhp.Num = position_number
	lhp.Extra = make([]LHDevice, 0, 2)
	lhp.Maxh = maxh
	return &lhp
}

// @implement Location
// -- this is clearly somewhere that something can be
// need to implement the liquid handler as a location as well

func (lhp *LHPosition) Location_ID() string {
	return lhp.ID
}

func (lhp *LHPosition) Location_Name() string {
	return lhp.Name
}

func (lhp *LHPosition) Container() Location {
	return lhp
}

func (lhp *LHPosition) Positions() []Location {
	return nil
}

func (lhp *LHPosition) Shape() *Shape {
	return NewShape("box", "mm", 0.08548, 0.12776, 0.0)
}

//  instruction to a liquid handler
type LHInstruction struct {
	ID               string
	ProductID        string
	BlockID          BlockID
	SName            string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	Plateaddress     string
	PlateID          string
	Platetype        string
	Vol              float64
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Result           *LHComponent
}

func (ins *LHInstruction) HasAnyParent() bool {
	for _, v := range ins.Components {
		if v.HasAnyParent() {
			return true
		}
	}

	return false
}

func (ins *LHInstruction) HasParent(id string) bool {
	for _, v := range ins.Components {
		if v.HasParent(id) {
			return true
		}
	}
	return false
}

func (ins *LHInstruction) ParentString() string {
	if ins == nil {
		return ""
	}

	tx := make([]string, 0, 1)

	for _, v := range ins.Components {
		//s += v.ParentID + "_"

		pid := v.ParentID

		if pid != "" {
			tx = append(tx, pid)
		}
	}

	if len(tx) == 0 {
		return ""
	} else {
		return strings.Join(tx, "_")
	}

}

// structure describing a solution: a combination of liquid components
type LHSolution struct {
	ID               string
	BlockID          BlockID
	Inst             string
	SName            string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	Plateaddress     string
	PlateID          string
	Platetype        string
	Vol              float64 // in S.I units only for now
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Minorlayoutgroup int
}

func NewLHInstruction() *LHInstruction {
	var lhi LHInstruction
	lhi.ID = GetUUID()
	lhi.Majorlayoutgroup = -1
	return &lhi
}

func NewLHSolution() *LHSolution {
	var lhs LHSolution
	lhs.ID = GetUUID()
	lhs.Majorlayoutgroup = -1
	lhs.Minorlayoutgroup = -1
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

func (sol LHSolution) String() string {
	one := fmt.Sprintf(
		"%s, %s, %s, %s, %d",
		sol.ID,
		sol.BlockID,
		sol.Inst,
		sol.SName,
		sol.Order,
	)
	for _, c := range sol.Components {
		one = one + fmt.Sprintf("[%s], ", c.CName)
	}
	two := fmt.Sprintf("%s, %s, %s, %g, %s, %g, %d, %d",
		sol.ContainerType,
		sol.Welladdress,
		sol.Platetype,
		sol.Vol,
		sol.Type,
		sol.Conc,
		sol.Tvol,
		sol.Majorlayoutgroup,
		sol.Minorlayoutgroup,
	)
	return one + two
}

func (lhs *LHSolution) GetAssignment() string {
	return lhs.Plateaddress + ":" + lhs.Welladdress
}

// structure describing a liquid component and its desired properties
type LHComponent struct {
	ID                 string
	BlockID            BlockID
	DaughterID         string
	ParentID           string
	Inst               string
	Order              int
	CName              string
	Type               LiquidType
	Vol                float64
	Conc               float64
	Vunit              string
	Cunit              string
	Tvol               float64
	Smax               float64
	Visc               float64
	StockConcentration float64
	Extra              map[string]interface{}
	Loc                string
	Destination        string
}

func (lhc *LHComponent) IsZero() bool {
	if lhc == nil || lhc.Type == LTNIL || lhc.CName == "" || lhc.Vol < 0.0000001 {
		return true
	}
	return false
}

func (lhc *LHComponent) SetVolume(v wunit.Volume) {
	lhc.Vol = v.RawValue()
	lhc.Vunit = v.Unit().PrefixedSymbol()
}

func (lhc *LHComponent) HasParent(s string) bool {
	return strings.Contains(lhc.ParentID, s)
}

func (lhc *LHComponent) HasDaughter(s string) bool {
	return strings.Contains(lhc.DaughterID, s)
}

func (lhc *LHComponent) Name() string {
	return lhc.CName
}

func (lhc *LHComponent) TypeName() string {
	return LiquidTypeName(lhc.Type)
}

func (lhc *LHComponent) Volume() wunit.Volume {
	return wunit.NewVolume(lhc.Vol, lhc.Vunit)
}

func (lhc *LHComponent) Remove(v wunit.Volume) {
	///TODO -- catch errors
	lhc.Vol -= v.ConvertToString(lhc.Vunit)
}

func (lhc *LHComponent) Dup() *LHComponent {
	c := NewLHComponent()
	c.ID = lhc.ID
	c.Order = lhc.Order
	c.CName = lhc.CName
	c.Type = lhc.Type
	c.Vol = lhc.Vol
	c.Conc = lhc.Conc
	c.Vunit = lhc.Vunit
	c.Tvol = lhc.Vol
	c.Smax = lhc.Smax
	c.Visc = lhc.Visc
	c.StockConcentration = lhc.StockConcentration
	c.Extra = make(map[string]interface{}, len(lhc.Extra))
	for k, v := range lhc.Extra {
		c.Extra[k] = v
	}
	c.Loc = lhc.Loc
	c.Destination = lhc.Destination
	return c
}

func (cmp *LHComponent) HasAnyParent() bool {
	if cmp.ParentID != "" {
		return true
	}

	return false
}

func (cmp *LHComponent) AddParent(parentID string) {
	cmp.ParentID += parentID + "_"
}

func (cmp *LHComponent) AddDaughter(daughterID string) {
	cmp.DaughterID += daughterID + "_"
}

func (cmp *LHComponent) Mix(cmp2 *LHComponent) {
	// if this component is zero we inherit the id of the other one
	if cmp.IsZero() {
		cmp.ID = cmp2.ID
	}
	cmp.Smax = mergeSolubilities(cmp, cmp2)
	// determine type of final
	cmp.Type = mergeTypes(cmp, cmp2)
	// add cmp2 to cmp
	vcmp := wunit.NewVolume(cmp.Vol, cmp.Vunit)
	vcmp2 := wunit.NewVolume(cmp2.Vol, cmp2.Vunit)
	vcmp.Add(vcmp2)
	cmp.Vol = vcmp.RawValue() // same units
	cmp.CName = mergeNames(cmp.CName, cmp2.CName)
	// allow trace back
	logger.Track(fmt.Sprintf("MIX %s %s %s", cmp.ID, cmp2.ID, vcmp.ToString()))
}

// @implement Liquid
// @deprecate Liquid

func (lhc *LHComponent) GetSmax() float64 {
	return lhc.Smax
}

func (lhc *LHComponent) GetVisc() float64 {
	return lhc.Visc
}

func (lhc *LHComponent) GetExtra() map[string]interface{} {
	return lhc.Extra
}

func (lhc *LHComponent) GetConc() float64 {
	return lhc.Conc
}

func (lhc *LHComponent) GetCunit() string {
	return lhc.Cunit
}

// new
func (lhc *LHComponent) Concentration() (conc wunit.Concentration) {
	conc = wunit.NewConcentration(lhc.Conc, lhc.Cunit)
	return conc
}

func (lhc *LHComponent) GetVunit() string {
	return lhc.Vunit
}

func (lhc *LHComponent) GetType() string {
	return LiquidTypeName(lhc.Type)
}

func NewLHComponent() *LHComponent {
	var lhc LHComponent
	lhc.ID = GetUUID()
	lhc.Vunit = "ul"
	return &lhc
}

// XXX -- why is this different from Dup?
func CopyLHComponent(lhc *LHComponent) *LHComponent {
	tmp, _ := json.Marshal(lhc)
	var lhc2 LHComponent
	json.Unmarshal(tmp, &lhc2)
	lhc2.ID = GetUUID()
	if lhc2.Inst != "" {
		lhc2.Inst = GetUUID()
		// this needs some thought
	}
	return &lhc2
}

// structure describing a sample

type LHSample LHComponent

// structure describing a microplate
// this needs to be harmonised with the version
type LHPlate struct {
	ID          string
	Inst        string
	Loc         string
	PlateName   string
	Type        string
	Mnfr        string
	WlsX        int
	WlsY        int
	Nwells      int
	HWells      map[string]*LHWell
	Height      float64
	Hunit       string
	Rows        [][]*LHWell
	Cols        [][]*LHWell
	Welltype    *LHWell
	Wellcoords  map[string]*LHWell
	WellXOffset float64
	WellYOffset float64
	WellXStart  float64
	WellYStart  float64
	WellZStart  float64
}

func (lhp LHPlate) Name() string {
	return lhp.PlateName
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
	Height      : %f,
	Hunit       : %s,
	Rows        : %p,
	Cols        : %p,
	Welltype    : %p,
	Wellcoords  : %p,
	WellXOffset : %f,
	WellYOffset : %f,
	WellXStart  : %f,
	WellYStart  : %f,
	WellZStart  : %f,
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
		lhp.Height,
		lhp.Hunit,
		lhp.Rows,
		lhp.Cols,
		lhp.Welltype,
		lhp.Wellcoords,
		lhp.WellXOffset,
		lhp.WellYOffset,
		lhp.WellXStart,
		lhp.WellYStart,
		lhp.WellZStart,
	)
}

// convenience method

func (lhp *LHPlate) GetComponent(cmp *LHComponent, exact bool) ([]WellCoords, bool) {
	ret := make([]WellCoords, 0, 1)

	it := NewOneTimeColumnWiseIterator(lhp)
	volGot := wunit.NewVolume(0.0, "ul")

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]
		//	logger.Debug(fmt.Sprint("WANT$$$: ", cmp.CName, " :: ", wc.FormatA1(), " ", w.Contents().CName))
		if w.Contents().CName == cmp.CName {
			if exact && w.Contents().ID != cmp.ID {
				continue
			}

			v := w.WorkingVolume()
			volGot.Add(v)
			ret = append(ret, wc)

			if volGot.GreaterThan(cmp.Volume()) {
				break
			}
		}
	}

	if !volGot.GreaterThan(cmp.Volume()) {
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

// @implement named

func (lhp *LHPlate) GetName() string {
	return lhp.PlateName
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

func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	var lhp LHPlate
	lhp.Type = platetype
	lhp.ID = GetUUID()
	lhp.Mnfr = mfr
	lhp.WlsX = ncols
	lhp.WlsY = nrows
	lhp.Nwells = ncols * nrows
	lhp.Height = height
	lhp.Hunit = hunit
	lhp.Welltype = welltype
	lhp.WellXOffset = wellXOffset
	lhp.WellYOffset = wellYOffset
	lhp.WellXStart = wellXStart
	lhp.WellYStart = wellYStart
	lhp.WellZStart = wellZStart

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
	ret := NewLHPlate(lhp.Type, lhp.Mnfr, lhp.WlsY, lhp.WlsX, lhp.Height, lhp.Hunit, lhp.Welltype, lhp.WellXOffset, lhp.WellYOffset, lhp.WellXStart, lhp.WellYStart, lhp.WellZStart)

	ret.PlateName = lhp.PlateName

	for i, row := range lhp.Rows {
		for j, well := range row {
			d := well.Dup()
			ret.Rows[i][j] = d
			ret.Cols[j][i] = d
			ret.Wellcoords[d.Crds] = d
		}
	}

	return ret
}

// structure representing a well on a microplate - description of a destination
type LHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Platetype string
	Crds      string
	MaxVol    float64
	Vunit     string
	WContents *LHComponent
	Rvol      float64
	WShape    *Shape
	Bottom    int
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
	Extra     map[string]interface{}
	Plate     *LHPlate `gotopb:"-" json:"-"`
}

func (w LHWell) String() string {
	return fmt.Sprintf(
		`LHWELL{
ID        : %s,
Inst      : %s,
Plateinst : %s,
Plateid   : %s,
Platetype : %s,
Crds      : %s,
MaxVol    : %g,
Vunit     : %s,
WContents : %v,
Rvol      : %g,
WShape    : %v,
Bottom    : %d,
Xdim      : %g,
Ydim      : %g,
Zdim      : %g,
Bottomh   : %g,
Dunit     : %s,
Extra     : %v,
Plate     : %v,
}`,
		w.ID,
		w.Inst,
		w.Plateinst,
		w.Plateid,
		w.Platetype,
		w.Crds,
		w.MaxVol,
		w.Vunit,
		w.WContents,
		w.Rvol,
		w.WShape,
		w.Bottom,
		w.Xdim,
		w.Ydim,
		w.Zdim,
		w.Bottomh,
		w.Dunit,
		w.Extra,
		w.Plate,
	)
}

func (w *LHWell) Contents() *LHComponent {
	// be careful
	if w == nil {
		logger.Debug("CONTENTS OF NIL WELL REQUESTED")
		return NewLHComponent()
	}
	if w.WContents == nil {
		return NewLHComponent()
	}

	return w.WContents
}

func (w *LHWell) Currvol() float64 {
	return w.Contents().Vol
}

func (w *LHWell) CurrVolume() wunit.Volume {
	return w.Contents().Volume()
}

func (w *LHWell) MaxVolume() wunit.Volume {
	return wunit.NewVolume(w.MaxVol, w.Vunit)
}
func (w *LHWell) Add(c *LHComponent) {
	mv := wunit.NewVolume(w.MaxVol, w.Vunit)
	cv := wunit.NewVolume(c.Vol, c.Vunit)
	wv := w.CurrentVolume()
	cv.Add(wv)
	if cv.GreaterThan(mv) {
		// could make this fatal but we don't track state well enough
		// for that to be worthwhile
		logger.Debug("WARNING: OVERFULL WELL AT ", w.Crds)
	}
	w.Contents().Mix(c)
}

func (w *LHWell) Remove(v wunit.Volume) *LHComponent {
	// if the volume is too high we complain

	if v.GreaterThan(w.CurrentVolume()) {
		logger.Debug("You ask too much: ", w.Crds, v.ToString())
		return nil
	}

	ret := w.Contents().Dup()
	ret.Vol = v.ConvertToString(w.Vunit)

	return ret
}

func (w *LHWell) WorkingVolume() wunit.Volume {
	v := wunit.NewVolume(w.Currvol(), w.Vunit)
	v2 := wunit.NewVolume(w.Rvol, w.Vunit)
	v.Subtract(v2)
	return v
}

func (w *LHWell) ResidualVolume() wunit.Volume {
	v := wunit.NewVolume(w.Rvol, w.Vunit)
	return v
}

func (w *LHWell) CurrentVolume() wunit.Volume {
	return w.Contents().Volume()
}

//@implement Location

func (lhw *LHWell) Location_ID() string {
	return lhw.ID
}

func (lhw *LHWell) Location_Name() string {
	return lhw.Platetype
}

func (lhw *LHWell) Shape() *Shape {
	if lhw.WShape == nil {
		// return the non-shape
		return NewNilShape()
	}
	return lhw.WShape
}

// @implement Well
// @deprecate Well

func (w *LHWell) ContainerType() string {
	return w.Platetype
}

func (w *LHWell) Empty() bool {
	if w.Currvol() <= 0.000001 {
		return true
	} else {
		return false
	}
}

func (lhw *LHWell) Dup() *LHWell {
	cp := NewLHWell(lhw.Platetype, lhw.Plateid, lhw.Crds, lhw.Vunit, lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.Xdim, lhw.Ydim, lhw.Zdim, lhw.Bottomh, lhw.Dunit)

	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	cp.WContents = lhw.Contents().Dup()

	return cp
}

func (lhw *LHWell) CalculateCrossSectionArea() (ca wunit.Area, err error) {

	ca, err = lhw.Shape().CrossSectionalArea()

	return
}

func (lhw *LHWell) CalculateMaxVolume() (vol wunit.Volume, err error) {

	if lhw.Bottom == 0 { // flat
		vol, err = lhw.Shape().Volume()
	} /*else if lhw.Bottom == 1 { // round
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	} else if lhw.Bottom == 2 { // Pointed / v-shaped /pyramid
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	}
	*/
	return
}

// make a new well structure
func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	var well LHWell

	well.WContents = NewLHComponent()
	well.ID = GetUUID()
	well.Platetype = platetype
	well.Plateid = plateid
	well.Crds = crds
	well.MaxVol = vol
	well.Rvol = rvol
	well.Vunit = vunit
	well.WShape = shape.Dup()
	well.Bottom = bott
	well.Xdim = xdim
	well.Ydim = ydim
	well.Zdim = zdim
	well.Bottomh = bottomh
	well.Dunit = dunit
	well.Extra = make(map[string]interface{})
	return &well
}

// this function tries to find somewhere to put something... it was written before
// i had an iterator. fml
func Get_Next_Well(plate *LHPlate, component *LHComponent, curwell *LHWell) (*LHWell, bool) {
	vol := component.Vol

	it := NewOneTimeColumnWiseIterator(plate)

	if curwell != nil {
		// quick check to see if we have room
		vol_left := get_vol_left(curwell)

		if vol_left >= vol {
			// fine we can just return this one
			return curwell, true
		}

		startcoords := MakeWellCoords(curwell.Crds)
		it.SetStartTo(startcoords)
		it.Rewind()
		it.Next()
	}

	var new_well *LHWell

	for wc := it.Curr(); it.Valid(); wc = it.Next() {

		crds := wc.FormatA1()

		new_well = plate.Wellcoords[crds]

		if new_well.Empty() {
			break
		}
		cnts := new_well.Contents()

		cont := cnts.Name()
		if cont != component.Name() {
			continue
		}

		vol_left := get_vol_left(new_well)

		if vol < vol_left {
			break
		}
	}

	if new_well == nil {
		return nil, false
	}

	return new_well, true
}

//XXX sloboda? This makes no sense now; need to revise
func get_vol_left(well *LHWell) float64 {
	//cnts := well.WContents
	// this is very odd... I can see how this works as a heuristic
	// but it doesn't make much sense to me
	carry_vol := 10.0 // microlitres
	//	total_carry_vol := float64(len(cnts)) * carry_vol
	total_carry_vol := carry_vol // yeah right
	Currvol := well.Currvol
	rvol := well.Rvol
	vol := well.MaxVol
	return vol - (Currvol() + total_carry_vol + rvol)
}

func New_Solution() *LHSolution {
	var solution LHSolution
	solution.ID = GetUUID()
	solution.Components = make([]*LHComponent, 0, 4)
	return &solution
}

func New_Plate(platetype *LHPlate) *LHPlate {
	new_plate := NewLHPlate(platetype.Type, platetype.Mnfr, platetype.WlsY, platetype.WlsX, platetype.Height, platetype.Hunit, platetype.Welltype, platetype.WellXOffset, platetype.WellYOffset, platetype.WellXStart, platetype.WellYStart, platetype.WellZStart)
	Initialize_Wells(new_plate)
	return new_plate
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

/* tip box */

type LHTipbox struct {
	ID         string
	Boxname    string
	Type       string
	Mnfr       string
	Nrows      int
	Ncols      int
	Height     float64
	Tiptype    *LHTip
	AsWell     *LHWell
	NTips      int
	Tips       [][]*LHTip
	TipXOffset float64
	TipYOffset float64
	TipXStart  float64
	TipYStart  float64
	TipZStart  float64
}

func NewLHTipbox(nrows, ncols int, height float64, manufacturer, boxtype string, tiptype *LHTip, well *LHWell, tipxoffset, tipyoffset, tipxstart, tipystart, tipzstart float64) *LHTipbox {
	var tipbox LHTipbox
	tipbox.ID = GetUUID()
	tipbox.Type = boxtype
	tipbox.Boxname = fmt.Sprintf("%s_%s", boxtype, tipbox.ID[1:len(tipbox.ID)-2])
	tipbox.Mnfr = manufacturer
	tipbox.Nrows = nrows
	tipbox.Ncols = ncols
	tipbox.Tips = make([][]*LHTip, ncols)
	tipbox.NTips = tipbox.Nrows * tipbox.Ncols
	tipbox.Height = height
	tipbox.Tiptype = tiptype
	tipbox.AsWell = well
	for i := 0; i < ncols; i++ {
		tipbox.Tips[i] = make([]*LHTip, nrows)
	}
	tipbox.TipXOffset = tipxoffset
	tipbox.TipYOffset = tipyoffset
	tipbox.TipXStart = tipxstart
	tipbox.TipYStart = tipystart
	tipbox.TipZStart = tipzstart
	return initialize_tips(&tipbox, tiptype)
}

func (tb LHTipbox) String() string {
	return fmt.Sprintf(
		`LHTipbox {
ID        : %s,
Boxname   : %s,
Type      : %s,
Mnfr      : %s,
Nrows     : %d,
Ncols     : %d,
Height    : %f,
Tiptype   : %p,
AsWell    : %v,
NTips     : %d,
Tips      : %p,
TipXOffset: %f,
TipYOffset: %f,
TipXStart : %f,
TipYStart : %f,
TipZStart : %f,
}`,
		tb.ID,
		tb.Boxname,
		tb.Type,
		tb.Mnfr,
		tb.Nrows,
		tb.Ncols,
		tb.Height,
		tb.Tiptype,
		tb.AsWell,
		tb.NTips,
		tb.Tips,
		tb.TipXOffset,
		tb.TipYOffset,
		tb.TipXStart,
		tb.TipYStart,
		tb.TipZStart,
	)
}

func (tb *LHTipbox) Dup() *LHTipbox {
	return NewLHTipbox(tb.Nrows, tb.Ncols, tb.Height, tb.Mnfr, tb.Type, tb.Tiptype, tb.AsWell, tb.TipXOffset, tb.TipYOffset, tb.TipXStart, tb.TipYStart, tb.TipZStart)
}

// @implement named

func (tb *LHTipbox) GetName() string {
	return tb.Boxname
}

func (tb *LHTipbox) N_clean_tips() int {
	c := 0
	for j := 0; j < tb.Nrows; j++ {
		for i := 0; i < tb.Ncols; i++ {
			if tb.Tips[i][j] != nil && !tb.Tips[i][j].Dirty {
				c += 1
			}
		}
	}
	return c
}

// actually useful functions
// TODO implement Mirror

func (tb *LHTipbox) GetTips(mirror bool, multi, orient int) []string {
	// this removes the tips as well
	var ret []string = nil
	if orient == LHHChannel {
		for j := 0; j < tb.Nrows; j++ {
			c := 0
			s := -1
			for i := 0; i < tb.Ncols; i++ {
				if tb.Tips[i][j] != nil && !tb.Tips[i][j].Dirty {
					c += 1
					if s == -1 {
						s = i
					}
				}
			}

			if c >= multi {
				ret = make([]string, multi)
				for i := 0; i < multi; i++ {
					tb.Tips[i+s][j] = nil
					wc := WellCoords{i + s, j}
					ret[i] = wc.FormatA1()
				}
				break
			}
		}

	} else if orient == LHVChannel {
		// find the first column with a contiguous set of at least multi
		for i := 0; i < tb.Ncols; i++ {
			c := 0
			s := -1
			// if we're picking up < the maxium number of tips we need to be careful
			// that there are no tips beneath the ones we're picking up

			for j := tb.Nrows - 1; j >= 0; j-- {
				if tb.Tips[i][j] != nil { // && !tb.Tips[i][j].Dirty
					c += 1
					if s == -1 {
						s = j
					}
				} else {
					if s != -1 {
						break // we've reached a gap
					}
				}
			}

			if c >= multi {
				ret = make([]string, 0, multi)
				n := 0
				for j := s; j >= 0; j-- {
					tb.Tips[i][j] = nil
					wc := WellCoords{i, j}
					fmt.Println(j, "Getting TIP from ", wc.FormatA1())
					ret = append(ret, wc.FormatA1())
					n += 1
					if n >= multi {
						break
					}
				}

				fmt.Println("RET: ", ret)
				break
			}
		}

	}

	tb.NTips -= multi
	return ret
}

func initialize_tips(tipbox *LHTipbox, tiptype *LHTip) *LHTipbox {
	nr := tipbox.Nrows
	nc := tipbox.Ncols
	for i := 0; i < nc; i++ {
		for j := 0; j < nr; j++ {
			tipbox.Tips[i][j] = CopyTip(*tiptype)
		}
	}
	tipbox.NTips = tipbox.Nrows * tipbox.Ncols
	return tipbox
}

type LHTip struct {
	ID     string
	Type   string
	Mnfr   string
	Dirty  bool
	MaxVol wunit.Volume
	MinVol wunit.Volume
}

func (tip *LHTip) IsNil() bool {
	if tip == nil || tip.Type == "" || tip.MaxVol.IsZero() || tip.MinVol.IsZero() {
		return true
	}
	return false
}

func (tip *LHTip) Dup() *LHTip {
	t := NewLHTip(tip.Mnfr, tip.Type, tip.MinVol.RawValue(), tip.MaxVol.RawValue(), tip.MinVol.Unit().PrefixedSymbol())
	t.Dirty = tip.Dirty
	return t
}

func NewLHTip(mfr, ttype string, minvol, maxvol float64, volunit string) *LHTip {
	var lht LHTip
	lht.ID = GetUUID()
	lht.Mnfr = mfr
	lht.Type = ttype
	lht.MaxVol = wunit.NewVolume(maxvol, volunit)
	lht.MinVol = wunit.NewVolume(minvol, volunit)
	return &lht
}

func CopyTip(tt LHTip) *LHTip {
	return &tt
}

// tip waste

type LHTipwaste struct {
	ID         string
	Type       string
	Mnfr       string
	Capacity   int
	Contents   int
	Height     float64
	WellXStart float64
	WellYStart float64
	WellZStart float64
	AsWell     *LHWell
}

func (te LHTipwaste) String() string {
	return fmt.Sprintf(
		`LHTipwaste {
	ID: %s,
	Type: %s,
	Mnfr: %s,
	Capacity: %d,
	Contents: %d,
	Height: %f,
	WellXStart: %f,
	WellYStart: %f,
	WellZStart: %f,
	AsWell: %p,
}
`,
		te.ID,
		te.Type,
		te.Mnfr,
		te.Capacity,
		te.Contents,
		te.Height,
		te.WellXStart,
		te.WellYStart,
		te.WellZStart,
		te.AsWell, //AsWell is printed as pointer to kepp things short
	)
}

func (tw *LHTipwaste) Dup() *LHTipwaste {
	return NewLHTipwaste(tw.Capacity, tw.Type, tw.Mnfr, tw.Height, tw.AsWell, tw.WellXStart, tw.WellYStart, tw.WellZStart)
}

func (tw *LHTipwaste) GetName() string {
	return tw.Type
}

func NewLHTipwaste(capacity int, typ, mfr string, height float64, w *LHWell, wellxstart, wellystart, wellzstart float64) *LHTipwaste {
	var lht LHTipwaste
	lht.ID = GetUUID()
	lht.Type = typ
	lht.Mnfr = mfr
	lht.Capacity = capacity
	lht.Height = height
	lht.AsWell = w
	lht.WellXStart = wellxstart
	lht.WellYStart = wellystart
	lht.WellZStart = wellzstart
	return &lht
}

func (lht *LHTipwaste) Empty() {
	lht.Contents = 0
}

func (lht *LHTipwaste) Dispose(n int) bool {
	if lht.Capacity-lht.Contents < n {
		return false
	}

	lht.Contents += n
	return true
}

// head
type LHHead struct {
	Name         string
	Manufacturer string
	ID           string
	Adaptor      *LHAdaptor
	Params       *LHChannelParameter
}

func NewLHHead(name, mf string, params *LHChannelParameter) *LHHead {
	var lhh LHHead
	lhh.Manufacturer = mf
	lhh.Name = name
	lhh.Params = params
	return &lhh
}

func (head *LHHead) Dup() *LHHead {
	h := NewLHHead(head.Name, head.Manufacturer, head.Params.Dup())
	if head.Adaptor != nil {
		h.Adaptor = head.Adaptor.Dup()
	}

	return h
}

func (lhh *LHHead) GetParams() *LHChannelParameter {
	if lhh.Adaptor == nil {
		return lhh.Params
	} else {
		return lhh.Adaptor.GetParams()
	}
}

// adaptor
// TODO -- should be an array of loaded tips
type LHAdaptor struct {
	Name          string
	ID            string
	Manufacturer  string
	Params        *LHChannelParameter
	Ntipsloaded   int
	Tiptypeloaded *LHTip
}

func NewLHAdaptor(name, mf string, params *LHChannelParameter) *LHAdaptor {
	var lha LHAdaptor
	lha.Name = name
	lha.Manufacturer = mf
	lha.Params = params
	return &lha
}

func (lha *LHAdaptor) Dup() *LHAdaptor {
	ad := NewLHAdaptor(lha.Name, lha.Manufacturer, lha.Params.Dup())
	ad.Ntipsloaded = lha.Ntipsloaded
	if !ad.Tiptypeloaded.IsNil() {
		ad.Tiptypeloaded = lha.Tiptypeloaded.Dup()
	}

	return ad
}

func (lha *LHAdaptor) LoadTips(n int, tiptype *LHTip) bool {
	if lha.Ntipsloaded > 0 {
		return false
	}

	lha.Ntipsloaded = n
	lha.Tiptypeloaded = tiptype
	return true
}

func (lha *LHAdaptor) UnloadTips() bool {
	if lha.Ntipsloaded == 0 {
		return false
	}

	lha.Ntipsloaded = 0
	lha.Tiptypeloaded = nil

	return true
}

func (lha *LHAdaptor) GetParams() *LHChannelParameter {
	if lha.Ntipsloaded == 0 {
		return lha.Params
	} else {
		return lha.Params.MergeWithTip(lha.Tiptypeloaded)
	}
}
