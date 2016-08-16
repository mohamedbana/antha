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
	"fmt"
	"math"
)

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

	bounds BBox
	parent LHObject
}

func NewLHTipbox(nrows, ncols int, size Coordinates, manufacturer, boxtype string, tiptype *LHTip, well *LHWell, tipxoffset, tipyoffset, tipxstart, tipystart, tipzstart float64) *LHTipbox {
	var tipbox LHTipbox
	tipbox.ID = GetUUID()
	tipbox.Type = boxtype
	tipbox.Boxname = fmt.Sprintf("%s_%s", boxtype, tipbox.ID[1:len(tipbox.ID)-2])
	tipbox.Mnfr = manufacturer
	tipbox.Nrows = nrows
	tipbox.Ncols = ncols
	tipbox.Tips = make([][]*LHTip, ncols)
	tipbox.NTips = tipbox.Nrows * tipbox.Ncols
	tipbox.bounds.SetSize(size)
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
Width     : %f,
Length    : %f,
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
		tb.bounds.GetSize().X,
		tb.bounds.GetSize().Y,
		tb.bounds.GetSize().Z,
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
	return NewLHTipbox(tb.Nrows, tb.Ncols, tb.bounds.GetSize(), tb.Mnfr, tb.Type, tb.Tiptype, tb.AsWell, tb.TipXOffset, tb.TipYOffset, tb.TipXStart, tb.TipYStart, tb.TipZStart)
}

// @implement named

func (tb *LHTipbox) GetName() string {
	return tb.Boxname
}

func (tb *LHTipbox) GetType() string {
	return tb.Type
}

func (self *LHTipbox) GetClass() string {
	return "tipbox"
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

//##############################################
//@implement LHObject
//##############################################

func (self *LHTipbox) GetPosition() Coordinates {
	return self.bounds.GetPosition()
}

func (self *LHTipbox) GetSize() Coordinates {
	return self.bounds.GetSize()
}

func (self *LHTipbox) GetBoxIntersections(box BBox) []LHObject {
	ret := []LHObject{}
	if self.bounds.IntersectsBox(box) {
		ret = append(ret, self)
	}
	//todo, scan tips
	return ret
}

func (self *LHTipbox) GetPointIntersections(point Coordinates) []LHObject {
	ret := []LHObject{}
	if self.bounds.IntersectsPoint(point) {
		ret = append(ret, self)
	}
	//todo, scan tips
	return ret
}

func (self *LHTipbox) SetOffset(o Coordinates) error {
	if self.parent != nil {
		o = o.Add(self.parent.GetPosition())
	}
	self.bounds.SetPosition(o)
	return nil
}

func (self *LHTipbox) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

func (self *LHTipbox) GetParent() LHObject {
	return self.parent
}

//##############################################
//@implement Addressable
//##############################################

func (tb *LHTipbox) AddressExists(c WellCoords) bool {
	return c.X >= 0 &&
		c.Y >= 0 &&
		c.X < tb.Ncols &&
		c.Y < tb.Nrows
}

func (self *LHTipbox) NRows() int {
	return self.Nrows
}

func (self *LHTipbox) NCols() int {
	return self.Ncols
}

func (tb *LHTipbox) GetChildByAddress(c WellCoords) LHObject {
	if !tb.AddressExists(c) {
		return nil
	}
	//Tips aren't yet an LHObject...
	//return tb.Tips[c.X][c.Y]
	return nil
}

func (tb *LHTipbox) CoordsToWellCoords(r Coordinates) (WellCoords, Coordinates) {
	wc := WellCoords{
		int(math.Floor(((r.X - tb.TipXStart) / tb.TipXOffset))), // + 0.5)), Don't have to add .5 because
		int(math.Floor(((r.Y - tb.TipYStart) / tb.TipYOffset))), // + 0.5)), TipX/YStart is to TL corner, not center
	}
	if wc.X < 0 {
		wc.X = 0
	} else if wc.X >= tb.Ncols {
		wc.X = tb.Ncols - 1
	}
	if wc.Y < 0 {
		wc.Y = 0
	} else if wc.Y >= tb.Nrows {
		wc.Y = tb.Nrows - 1
	}

	r2, _ := tb.WellCoordsToCoords(wc, TopReference)

	return wc, r.Subtract(r2)
}

func (tb *LHTipbox) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
	if !tb.AddressExists(wc) {
		return Coordinates{}, false
	}

	var z float64
	if r == BottomReference {
		z = tb.TipZStart
	} else if r == TopReference {
		z = tb.Height
	} else {
		return Coordinates{}, false
	}

	return Coordinates{
		tb.TipXStart + (float64(wc.X)+0.5)*tb.TipXOffset,
		tb.TipYStart + (float64(wc.Y)+0.5)*tb.TipYOffset,
		z}, true
}

//HasTipAt
func (tb *LHTipbox) HasTipAt(c WellCoords) bool {
	return tb.AddressExists(c) && tb.Tips[c.X][c.Y] != nil
}

//RemoveTip
func (tb *LHTipbox) RemoveTip(c WellCoords) *LHTip {
	if !tb.AddressExists(c) {
		return nil
	}
	tip := tb.Tips[c.X][c.Y]
	tb.Tips[c.X][c.Y] = nil
	return tip
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
					//fmt.Println(j, "Getting TIP from ", wc.FormatA1())
					ret = append(ret, wc.FormatA1())
					n += 1
					if n >= multi {
						break
					}
				}

				//fmt.Println("RET: ", ret)
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
