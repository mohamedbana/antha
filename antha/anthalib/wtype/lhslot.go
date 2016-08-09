// /anthalib/simulator/liquidhandling/robotstate.go: Part of the Antha language
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
	"fmt"
)

// -------------------------------------------------------------------------------
//                            DeckSlot
// -------------------------------------------------------------------------------

//DeckSlot A slot on an lh robot
type DeckSlot struct {
	name            string
	acceptsTip      bool
	acceptsPlate    bool
	acceptsTipwaste bool
	bounds          BBox
	child           LHObject
	parent          LHObject
}

func NewDeckSlot(name string, position, size Coordinates) *DeckSlot {
	r := DeckSlot{name, false, false, false, *NewBBox(position, size), nil, nil}
	return &r
}

//@implement Named
//GetName
func (self *DeckSlot) GetName() string {
	return self.name
}

//@implement Typed
//GetType
func (self *DeckSlot) GetType() string {
	return "deck_slot"
}

func (self *DeckSlot) SetAcceptsTip(b bool) {
	self.acceptsTip = b
}

func (self *DeckSlot) SetAcceptsPlate(b bool) {
	self.acceptsPlate = b
}

func (self *DeckSlot) SetAcceptsTipwaste(b bool) {
	self.acceptsTipwaste = b
}

//@Implement LHObject
func (self *DeckSlot) GetBounds() BBox {
	return self.bounds
}

func (self *DeckSlot) SetOffset(o Coordinates) {
	if self.parent != nil {
		o = o.Add(self.parent.GetBounds().GetPosition())
	}
	self.bounds.SetPosition(o)
}

func (self *DeckSlot) GetParent() LHObject {
	return self.parent
}

func (self *DeckSlot) SetParent(o LHObject) {
	self.parent = o
}

//@implements LHSlot
func (self *DeckSlot) GetChild() LHObject {
	return self.child
}

func (self *DeckSlot) GetChildPosition() Coordinates {
	//Deck slots put their child at the same location as themselves
	return self.bounds.GetPosition()
}

func (self *DeckSlot) SetChild(o LHObject) error {
	if err := self.accepts(o); err != nil {
		return err
	}
	if self.child != nil {
		o_name := "unnamed"
		c_name := "unnamed"
		if on, ok := o.(Named); ok {
			o_name = on.GetName()
		}
		if cn, ok := self.child.(Named); ok {
			c_name = cn.GetName()
		}
		return fmt.Errorf("Cannot add object \"%s\" to slot \"%s\" which already contains \"%s\"",
			o_name, self.GetName(), c_name)
	}
	o.SetParent(self)
	o.SetOffset(Coordinates{})
	self.child = o
	return nil
}

func (self *DeckSlot) Accepts(o LHObject) bool {
	return self.accepts(o) == nil
}

func (self *DeckSlot) accepts(o LHObject) error {
	ss := self.bounds.GetSize()
	os := o.GetBounds().GetSize()
	if ss.X == os.X && ss.Y == os.Y {
		var b bool
		switch o.(type) {
		case *LHPlate:
			b = self.acceptsPlate
		case *LHTipbox:
			b = self.acceptsTip
		case *LHTipwaste:
			b = self.acceptsTipwaste
		}

		if b {
			return nil
		} else {
			return fmt.Errorf("Cannot accept object type %T to slot \"%s\"", o, self.GetName())
		}
	}
	if n, ok := o.(Named); ok {
		return fmt.Errorf("Footprint of object \"%s\"[%vmm x %vmm] does not fit slot \"%s\"[%vmm x %vmm]",
			n.GetName(), os.X, os.Y,
			self.GetName(), ss.X, ss.Y)
	}
	return fmt.Errorf("Footprint of unnamed object[%s] does not fit slot \"%s\"[%s]",
		os, self.GetName(), ss)
}
