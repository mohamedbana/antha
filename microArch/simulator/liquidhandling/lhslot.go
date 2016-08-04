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

package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
    "fmt"
)

//defining this here for now
type LHSlot interface {
    HasChild() bool
    GetChild() wtype.LHObject
    SetChild(wtype.LHObject) bool
    Accepts(wtype.LHObject) bool
    GetChildPosition() wtype.Coordinates
}

// -------------------------------------------------------------------------------
//                            GeneralSlot
// -------------------------------------------------------------------------------

//GeneralSlot this slot will accept anything that fits
type GeneralSlot struct {
    name            string
    position        wtype.Coordinates
    size            wtype.Coordinates
    child           wtype.LHObject
}

func NewGeneralSlot(name string, position, size wtype.Coordinates) *GeneralSlot {
    r := GeneralSlot{name, position, size, nil}
    return &r
}

//@implement Named
//GetName
func (self *GeneralSlot) GetName() string {
    return self.name
}

//@implement Typed
//GetType
func (self *GeneralSlot) GetType() string {
    return "general_slot"
}

//@implements LHObject
func (self *GeneralSlot) GetSize() wtype.Coordinates {
    return self.size
}
func (self *GeneralSlot) GetPosition() wtype.Coordinates {
    return self.position
}
func (self *GeneralSlot) SetPosition(p wtype.Coordinates) {
    self.position = p
}

//@implements LHSlot
func (self *GeneralSlot) HasChild() bool {
    return self.child != nil
}

func (self *GeneralSlot) GetChild() wtype.LHObject {
    if self == nil {
        return nil
    }
    return self.child
}

func (self *GeneralSlot) SetChild(o wtype.LHObject) bool {
    if self.Accepts(o) {
        self.child = o
        return true
    }
    return false
}

func (self *GeneralSlot) Accepts(o wtype.LHObject) bool {
    return self.size.X >= o.GetSize().X && self.size.Y >= o.GetSize().Y
}

func (self *GeneralSlot) GetChildPosition() wtype.Coordinates {
    return self.position
}

// -------------------------------------------------------------------------------
//                            TipwasteSlot
// -------------------------------------------------------------------------------

//TipwasteSlot only accepts TipWaste
type TipwasteSlot struct {
    GeneralSlot
}

func NewTipwasteSlot(name string, position, size wtype.Coordinates) *TipwasteSlot {
    r := TipwasteSlot{GeneralSlot{name, position, size, nil}}
    return &r
}

//GetType
func (self *TipwasteSlot) GetType() string {
    return "tipwaste_slot"
}

func (self *TipwasteSlot) Accepts(o wtype.LHObject) bool {
    if self.GeneralSlot.Accepts(o) {
        //check that this is a tipwaste
        _, ok := o.(*wtype.LHTipwaste)
        return ok
    }
    fmt.Println("GeneralFail")
    fmt.Println("object size: ", o.GetSize())
    return false
}

// -------------------------------------------------------------------------------
//                            NonTipwasteSlot
// -------------------------------------------------------------------------------

//NonTipwasteSlot only accepts TipWaste
type NonTipwasteSlot struct {
    GeneralSlot
}

func NewNonTipwasteSlot(name string, position, size wtype.Coordinates) *NonTipwasteSlot {
    r := NonTipwasteSlot{GeneralSlot{name, position, size, nil}}
    return &r
}

//GetType
func (self *NonTipwasteSlot) GetType() string {
    return "nontipwaste_slot"
}

func (self *NonTipwasteSlot) Accepts(o wtype.LHObject) bool {
    if self.GeneralSlot.Accepts(o) {
        //check that this isn't a tipwaste
        _, ok := o.(*wtype.LHTipwaste)
        return !ok
    }
    return false
}
