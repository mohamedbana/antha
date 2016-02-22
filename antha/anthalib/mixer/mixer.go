// mixer/mixer.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// mix needs to define the interface with liquid handling
// in order to do this it has to make the appropriate liquid handling
// request structure

// take all of this liquid
func SampleAll(l *wtype.LHComponent) *wtype.LHComponent {
	return Sample(l, l.Volume())
}

// below need to account for having locations for liquids specified...

// take a sample of volume v from this liquid
func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	ret.ID = l.ID

	l.AddDaughter(ret.ID)
	if l.HasAnyParent() {
		ret.ParentID = l.ParentID
	}
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()

	//logger.Track(fmt.Sprintf("SAMPLE V %s %s %s", l.ID, ret.ID, v.ToString()))

	return ret
}

// take an array of samples and array of corresponding volumes and sample them all
func MultiSample(l []*wtype.LHComponent, v []wunit.Volume) []*wtype.LHComponent {
	reta := make([]*wtype.LHComponent, 0)

	for i, j := range l {
		ret := wtype.NewLHComponent()
		vi := v[i]
		ret.ID = j.ID
		j.AddDaughter(ret.ID)
		if j.HasAnyParent() {
			ret.ParentID = j.ParentID
		}
		ret.CName = j.Name()
		ret.Type = j.Type
		ret.Vol = vi.RawValue()
		ret.Vunit = vi.Unit().PrefixedSymbol()
		ret.Extra = j.GetExtra()
		ret.Smax = j.GetSmax()
		ret.Visc = j.GetVisc()
		//	logger.Track(fmt.Sprintf("SAMPLE V %s %s %s", j.ID, ret.ID, vi.ToString()))
		reta = append(reta, ret)
	}

	return reta
}

// take a sample of this liquid and aim for a particular concentration
func SampleForConcentration(l *wtype.LHComponent, c wunit.Concentration) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	ret.ID = l.ID
	l.AddDaughter(ret.ID)
	if l.HasAnyParent() {
		ret.ParentID = l.ParentID
	}
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Conc = c.RawValue()
	ret.Cunit = c.Unit().PrefixedSymbol()
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	//logger.Track(fmt.Sprintf("SAMPLE C %s %s %s", l.ID, ret.ID, c.ToString()))
	return ret
}

func SampleMass(s *wtype.LHComponent, m wunit.Mass, d wunit.Density) *wtype.LHComponent {

	// calculate volume to add from density
	v := wunit.MasstoVolume(m, d)

	ret := wtype.NewLHComponent()
	ret.ID = s.ID
	s.AddDaughter(ret.ID)
	if s.HasAnyParent() {
		ret.ParentID = s.ParentID
	}
	ret.CName = s.Name()
	ret.Type = s.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = s.GetExtra()
	ret.Smax = s.GetSmax()
	ret.Visc = s.GetVisc()
	//logger.Track(fmt.Sprintf("SAMPLE M %s %s %s %s", s.ID, ret.ID, m.ToString(), d.ToString()))
	return ret
}

// take a sample ofs this liquid to be used to make the solution up to
// a particular total volume
// edited to take into account the volume of the other solution components
func SampleForTotalVolume(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	ret.ID = l.ID
	l.AddDaughter(ret.ID)
	if l.HasAnyParent() {
		ret.ParentID = l.ParentID
	}
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Tvol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	//logger.Track(fmt.Sprintf("SAMPLE T %s %s %s", l.ID, ret.ID, v.ToString()))

	return ret
}

/*
func SampleSolidtoLiquid(s wtype.Powder, m wunit.Mass, d wunit.Density) *wtype.LHComponent {

	// calculate volume to add from density
	v := wunit.MasstoVolume(m, d)

	ret := wtype.NewLHComponent()
	ret.CName = s.Name()
	ret.Type = s.GetType()
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = s.GetExtra()
	ret.Smax = s.GetSmax()
	ret.Visc = s.GetVisc()

	return ret
}
*/

type MixOptions struct {
	Components  []*wtype.LHComponent // Components to mix (required)
	Instruction *wtype.LHInstruction // used to be LHSolution
	Result      *wtype.LHComponent   // the resultant component
	Destination *wtype.LHPlate       // Destination plate; if nil, select one later
	PlateType   string               // type of destination plate
	Address     string               // Well in destination to place result; if nil, select one later
	PlateNum    int                  // which plate to stick these on
}

func GenericMix(opt MixOptions) *wtype.LHInstruction {
	r := opt.Instruction
	if r == nil {
		r = wtype.NewLHInstruction()
	}
	r.Components = opt.Components

	if opt.Destination != nil {
		r.ContainerType = opt.Destination.Type
		r.Platetype = opt.Destination.Type
		r.PlateID = opt.Destination.ID
	}

	if opt.PlateType != "" {
		r.ContainerType = opt.PlateType
		r.Platetype = opt.PlateType
	}

	if len(opt.Address) > 0 {
		r.Welladdress = opt.Address
	}

	if opt.PlateNum > 0 {
		r.Majorlayoutgroup = opt.PlateNum - 1
	}

	// We must respect the order in which things are mixed. The convention is
	// that mix(X,Y) corresponds to "Add Y to X".
	/*
		-- this has already happened in execute.mix
		for idx, comp := range r.Components {
			comp.Order = idx
			r.Result.Mix(comp)
			r.Result.AddParent(comp.ID)
			comp.AddDaughter(r.Result.ID)
		}
	*/

	r.Result = opt.Result

	return r
}

// Mix the specified wtype.LHComponents together and leave the destination TBD
func Mix(components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components: components,
	})
	return r.Result
}

// Mix the specified wtype.LHComponents together into a specific plate
func MixInto(destination *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components:  components,
		Destination: destination,
		Address:     address,
	})

	return r.Result
}

// Mix the specified wtype.LHComponents together into a plate of a particular type
func MixTo(platetype string, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	r := GenericMix(MixOptions{
		Components: components,
		PlateType:  platetype,
		Address:    address,
		PlateNum:   platenum,
	})
	return r.Result
}
