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
    "errors"
    "fmt"
)

//defining this here for now
type LHSlot interface {
    HasChild() bool
    GetChild() wtype.LHObject
    SetChild(wtype.LHObject) error
    Accepts(wtype.LHObject) error
    GetChildPosition() wtype.Coordinates
}

// -------------------------------------------------------------------------------
//                            GeneralSlot
// -------------------------------------------------------------------------------

//GeneralSlot this slot will accept anything that fits
type GeneralSlot struct {
    wtype.BBox
    name            string
    child           wtype.LHObject
}

func NewGeneralSlot(name string, position, size wtype.Coordinates) *GeneralSlot {
    r := GeneralSlot{*wtype.NewBBox(position, size), name, nil}
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

func (self *GeneralSlot) SetChild(o wtype.LHObject) error {
    if err := self.Accepts(o); err != nil {
        return err
    }
    self.child = o
    return nil
}

func (self *GeneralSlot) Accepts(o wtype.LHObject) error {
    if self.GetSize().X >= o.GetSize().X && self.GetSize().Y >= o.GetSize().Y {
        return nil
    }
    if n, ok := o.(wtype.Named); ok {
        return fmt.Errorf("Footprint of object \"%s\"[%s] too large for slot \"%s\"[%s]", 
                            n.GetName(), o.GetSize(), self.GetName(), self.GetSize())
    }
    return fmt.Errorf("Footprint of unnamed object[%s] too large for slot \"%s\"[%s]", 
                            o.GetSize(), self.GetName(), self.GetSize())
}

func (self *GeneralSlot) GetChildPosition() wtype.Coordinates {
    return self.GetPosition()
}

// -------------------------------------------------------------------------------
//                            TipwasteSlot
// -------------------------------------------------------------------------------

//TipwasteSlot only accepts TipWaste
type TipwasteSlot struct {
    GeneralSlot
}

func NewTipwasteSlot(name string, position, size wtype.Coordinates) *TipwasteSlot {
    r := TipwasteSlot{}
    r.GeneralSlot = *NewGeneralSlot(name, position, size)
    return &r
}

//GetType
func (self *TipwasteSlot) GetType() string {
    return "tipwaste_slot"
}

func (self *TipwasteSlot) Accepts(o wtype.LHObject) error {
    if err := self.GeneralSlot.Accepts(o); err != nil {
        return err
    } else if _, ok := o.(*wtype.LHTipwaste); !ok {
        if n, ok := o.(wtype.Named); ok {
            return errors.New(fmt.Sprintf("Slot \"%s\" cannot accept non-tipwaste object \"%s\"",
            self.GetName(), n.GetName()))
        } else {
            return errors.New(fmt.Sprintf("Slot \"%s\" cannot accept unnamed object as it is not a TipWaste",
            self.GetName()))
        }
    }
    return nil
}

func (self *TipwasteSlot) SetChild(o wtype.LHObject) error {
    if err := self.Accepts(o); err != nil {
        return err
    }
    self.child = o
    return nil
}

// -------------------------------------------------------------------------------
//                            NonTipwasteSlot
// -------------------------------------------------------------------------------

//NonTipwasteSlot only accepts TipWaste
type NonTipwasteSlot struct {
    GeneralSlot
}

func NewNonTipwasteSlot(name string, position, size wtype.Coordinates) *NonTipwasteSlot {
    r := NonTipwasteSlot{}
    r.GeneralSlot = *NewGeneralSlot(name, position, size)
    return &r
}

//GetType
func (self *NonTipwasteSlot) GetType() string {
    return "nontipwaste_slot"
}

func (self *NonTipwasteSlot) Accepts(o wtype.LHObject) error {
    if err := self.GeneralSlot.Accepts(o); err != nil {
        return err
    } else if _, ok := o.(*wtype.LHTipwaste); ok {
        if n, ok := o.(wtype.Named); ok {
            return errors.New(fmt.Sprintf("Slot \"%s\" cannot accept tipwaste object \"%s\"",
            self.GetName(), n.GetName()))
        } else {
            return errors.New(fmt.Sprintf("Slot \"%s\" cannot accept unnamed tipwaste object",
            self.GetName()))
        }
    }
    return nil
}

func (self *NonTipwasteSlot) SetChild(o wtype.LHObject) error {
    if err := self.Accepts(o); err != nil {
        return err
    }
    self.child = o
    return nil
}
