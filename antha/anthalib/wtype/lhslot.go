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

type SlotListType int
const (
    WhiteList SlotListType = iota
    BlackList
)

type SlotListItem int
type SlotListItems []SlotListItem
const (
    TipboxItem SlotListItem = iota
    PlateItem
    TipwasteItem
)

//DeckSlot A slot on an lh robot
type DeckSlot struct {
    name            string
    bounds          BBox
    listType        SlotListType
    listItems       SlotListItems
    child           LHObject
    parent          LHObject
}

func NewDeckSlot(name string, position, size Coordinates, list_type SlotListType) *DeckSlot {
    r := DeckSlot{name, *NewBBox(position, size), list_type, SlotListItems{}, nil, nil}
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

func (self *DeckSlot) AddItem(i SlotListItem) {
    self.listItems = append(self.listItems, i)
}

//@implements LHSlot
func (self *DeckSlot) GetContents() LHObject {
    if self == nil {
        return nil
    }
    return self.child
}

func (self *DeckSlot) GetContentsPosition() Coordinates {
    return self.bounds.GetPosition()
}

func (self *DeckSlot) SetContents(o LHObject) error {
    if err := self.Accepts(o); err != nil {
        return err
    }
    o.SetParent(self)
    o.SetOffset(Coordinates{})
    self.child = o
    return nil
}

func (self *DeckSlot) Accepts(o LHObject) error {
    ss := self.bounds.GetSize()
    os := o.GetBounds().GetSize()
    if ss.X != os.X && ss.Y != os.Y {
        var b bool
        switch o.(type) {
        case *LHPlate:
            b = self.acceptsType(PlateItem)
        case *LHTipbox:
            b = self.acceptsType(TipboxItem)
        case *LHTipwaste:
            b = self.acceptsType(TipwasteItem)
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

func (self *DeckSlot) acceptsType(t SlotListItem) bool {
    for _,v := range self.listItems {
        if v == t {
            return self.listType == WhiteList
        }
    }
    return self.listType == BlackList
}

