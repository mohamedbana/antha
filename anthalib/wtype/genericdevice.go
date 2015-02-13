// wtype/genericdevice.go: Part of the Antha language
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
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package wtype

// Generic device structure defining the manufacturer
// device type and current state (true/false indicating "ready/not ready")
type GenericDevice struct{
	Mf string
	Tp string
	State bool
}

func (gd *GenericDevice) Manufacturer() string {
	return gd.Mf
}

func (gd *GenericDevice) Type() string{
	return gd.Tp
}

func (gd *GenericDevice) Ready() bool{
	return gd.State
}


// structure to define interface to a physical device which can
// perform more than one operation
type CompositeDevice struct{
	Mf string
	Tp string
	Components []*Device
}

func (cd *CompositeDevice)Manufacturer() string{
	return cd.Mf
}

func (cd *CompositeDevice)Type()string{
	return cd.Tp
}

// the default for the whole device is to AND the Ready()s for
// the whole set of devices
func (cd *CompositeDevice)Ready()bool{
	// this could be one of a few options
	// for the device as a whole it's bascally AND

	b:=true

	for _,t:=range cd.Components{
		if(!(*t).Ready()){
			b=false
			break
		}
	}
	return b
}



