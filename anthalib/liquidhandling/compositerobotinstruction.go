// anthalib//compositerobotinstruction.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/anthalib/wutil"
)

// a set of instructions which are higher-level than the
// basic kind, along with some default implementations

type TransferOutputFunc func(TransferInstruction) []RobotInstruction

type TransferInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []float64 // this could be a Measurement
	VolumeUnit []string
	Prms       *LHParameter
}

func (ti TransferInstruction) InstructionType() int {
	return ti.Type
}

func (ins TransferInstruction) GetParameter(s string) interface{} {
	switch s {
	case "TYPES":
		return ins.What
	case "VOLUMES":
		return ins.Volume
	case "VOLUNTS":
		return ins.VolumeUnit
	case "POSFROM":
		return ins.PltFrom
	case "POSTO":
		return ins.PltTo
	case "WELLSFROM":
		return ins.WellFrom
	case "WELLSTO":
		return ins.WellTo
	case "PARAMS":
		return ins.Prms
	default:
		RaiseError(fmt.Sprintf("Illegal parameter: %s", s))
	}
	return nil
}

func Transfer(what []string, pfrom, pto []string, wfrom, wto []string, v []float64, vu []string, prms *LHParameter) TransferInstruction {
	ti := TransferInstruction{TFR, what, pfrom, pto, wfrom, wto, v, vu, prms}
	return ti
}

// placeholder function: the intention is to have flexible rewriting of transfers
// this should happen in the device driver with input from other policies, liquid classes etc.
func SimpleOutput(ti TransferInstruction, rq LHRequest) []RobotInstruction {
	posFrom := make([]int, len(ti.PltFrom))
	posTo := make([]int, len(ti.PltTo))

	for i, _ := range ti.PltFrom {
		posFrom[i] = PlateLookup(rq, ti.PltFrom[i])
		posTo[i] = PlateLookup(rq, ti.PltTo[i])
	}

	return (SimpleTransfer(posFrom, posTo, ti.WellFrom, ti.WellTo, ti.Volume, ti.VolumeUnit, ti.What, ti.Prms))
}

func SimpleTransfer(posfrom, posto []int, wellfrom, wellto []string, amount []float64, unit []string, what []string, prms *LHParameter) []RobotInstruction {
	ret := make([]RobotInstruction, 0, 4)

	// the usages below are OK but need to account for how we specify null values

	aspXOffset := wutil.GetFloat64FromMap(prms.Policy, "AspirateXOffset")
	aspYOffset := wutil.GetFloat64FromMap(prms.Policy, "AspirateYOffset")
	aspZOffset := wutil.GetFloat64FromMap(prms.Policy, "AspirateZOffset")

	aspSpeed := wutil.GetFloat64FromMap(prms.Policy, "AspirateSpeed")
	dspSpeed := wutil.GetFloat64FromMap(prms.Policy, "DispenseSpeed")

	pipSpeedUnit := wutil.GetStringFromMap(prms.Policy, "PipetteSpeedUnit")

	aspHeight := wutil.GetIntFromMap(prms.Policy, "AspirateHeight")
	dspHeight := wutil.GetIntFromMap(prms.Policy, "DispenseHeight")

	dspXOffset := wutil.GetFloat64FromMap(prms.Policy, "DispenseXOffset")
	dspYOffset := wutil.GetFloat64FromMap(prms.Policy, "DispenseYOffset")
	dspZOffset := wutil.GetFloat64FromMap(prms.Policy, "DispenseZOffset")

	// how many transfers?
	// defer that question to the transferVolumes function

	for i := 0; i < len(what); i++ {
		vols := transferVolumes(what[i], amount[i], prms)
		fmt.Println(what[i], " ", amount[i])
		for _, transferVolume := range vols {
			ret = append(ret, Move(posfrom[i], wellfrom[i], aspHeight, aspXOffset, aspYOffset, aspZOffset, what[i]))
			ret = append(ret, Aspirate(transferVolume, unit[i], aspSpeed, pipSpeedUnit, what[i]))
			ret = append(ret, Move(posto[i], wellto[i], dspHeight, dspXOffset, dspYOffset, dspZOffset, what[i]))
			ret = append(ret, Dispense(transferVolume, unit[i], dspSpeed, pipSpeedUnit, what[i]))
		}
	}

	return ret
}

func transferVolumes(what string, amount float64, prms *LHParameter) []float64 {
	vols := make([]float64, 2)

	// this needs added flexibility, for the moment we need to just make this work

	min := prms.Minvol
	max := prms.Maxvol

	// if the volume is OK we just leave it as-is

	if amount <= max && amount >= min {
		vols[0] = amount
	} else if amount > min {
		// we must be > max by definition

		// in future we need to add more liquid handling policies here

		t := amount / max

		n_max := int(t)
		remainder := t - float64(int(t))

		v := max * remainder

		n_extra := make([]float64, 2)

		if v < min {
			// need to make up the difference
			n_max -= 1
			v2 := v + max
			v2 /= 2
			n_extra = append(n_extra, v2)
			n_extra = append(n_extra, v2)
		}

		for x := 0; x < n_max; x++ {
			vols = append(vols, max)
		}
		for x := 0; x < len(n_extra); x++ {
			vols = append(vols, n_extra[x])
		}
	} else {
		// nothing we can do here
		panic("Min volume exceeded")
	}

	return vols
}
