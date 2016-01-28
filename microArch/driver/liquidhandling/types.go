// /anthalib/driver/liquidhandling/types.go: Part of the Antha language
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
	"fmt"
	"strconv"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/material"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

// describes a liquid handler, its capabilities and current state
// probably needs splitting up to separate out the state information
// from the properties information
type LHProperties struct {
	ID                   string
	Nposns               int
	Positions            map[string]*wtype.LHPosition
	PlateLookup          map[string]interface{}
	PosLookup            map[string]string
	PlateIDLookup        map[string]string
	Plates               map[string]*wtype.LHPlate
	Tipboxes             map[string]*wtype.LHTipbox
	Tipwastes            map[string]*wtype.LHTipwaste
	Wastes               map[string]*wtype.LHPlate
	Washes               map[string]*wtype.LHPlate
	Devices              map[string]string
	Model                string
	Mnfr                 string
	LHType               string
	TipType              string
	Heads                []*wtype.LHHead
	HeadsLoaded          []*wtype.LHHead
	Adaptors             []*wtype.LHAdaptor
	Tips                 []*wtype.LHTip
	Tip_preferences      []string
	Input_preferences    []string
	Output_preferences   []string
	Tipwaste_preferences []string
	Waste_preferences    []string
	Wash_preferences     []string
	Driver               LiquidhandlingDriver        `gotopb:"-"`
	CurrConf             *wtype.LHChannelParameter   // TODO: initialise
	Cnfvol               []*wtype.LHChannelParameter // TODO: initialise
	Layout               map[string]wtype.Coordinates
	MaterialType         material.MaterialType
}

// validator for LHProperties structure

func ValidateLHProperties(props *LHProperties) (bool, string) {
	bo := true
	so := "OK"

	be := false
	se := "LHProperties Error: No position"

	if props.Positions == nil || len(props.Positions) == 0 {
		return be, se + "s"
	}

	for k, p := range props.Positions {
		if p == nil || p.ID == "" {
			return be, se + " " + k + " not set"
		}
	}

	se = "LHProperties error: No position lookup"

	if props.PosLookup == nil || len(props.PosLookup) == 0 {
		return be, se
	}

	se = "LHProperties Error: No tip preference information"

	if props.Tip_preferences == nil || len(props.Tip_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No input preference information"

	if props.Input_preferences == nil || len(props.Input_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No output preference information"

	if props.Output_preferences == nil || len(props.Output_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No waste preference information"

	if props.Waste_preferences == nil || len(props.Waste_preferences) == 0 {
		return be, se
	}

	se = "LHProperties Error: No tipwaste preference information"

	if props.Tipwaste_preferences == nil || len(props.Tipwaste_preferences) == 0 {
		return be, se
	}
	se = "LHProperties Error: No wash preference information"

	if props.Wash_preferences == nil || len(props.Wash_preferences) == 0 {
		return be, se
	}
	se = "LHProperties Error: No Plate ID lookup information"

	if props.PlateIDLookup == nil {
		return be, se
	}

	se = "LHProperties Error: No tip defined"

	if props.Tips == nil {
		return be, se
	}

	se = "LHProperties Error: No headsloaded array"

	if props.HeadsLoaded == nil {
		return be, se
	}

	return bo, so
}

// copy constructor

func (lhp *LHProperties) Dup() *LHProperties {
	lo := make(map[string]wtype.Coordinates, len(lhp.Layout))
	for k, v := range lhp.Layout {
		lo[k] = v
	}
	r := NewLHProperties(lhp.Nposns, lhp.Model, lhp.Mnfr, lhp.LHType, lhp.TipType, lo)

	for _, a := range lhp.Adaptors {
		r.Adaptors = append(r.Adaptors, a)
	}

	for _, h := range lhp.Heads {
		r.Heads = append(r.Heads, h)
	}

	for _, hl := range lhp.HeadsLoaded {
		r.HeadsLoaded = append(r.HeadsLoaded, hl)
	}

	for name, pl := range lhp.PosLookup {
		r.PosLookup[name] = pl
	}

	for name, pl := range lhp.PlateIDLookup {
		r.PlateIDLookup[name] = pl
	}

	for name, pt := range lhp.PlateLookup {
		r.PlateLookup[name] = pt
	}

	for name, tb := range lhp.Tipboxes {
		r.Tipboxes[name] = tb.Dup()
	}

	for name, tw := range lhp.Tipwastes {
		r.Tipwastes[name] = tw.Dup()
	}

	for name, waste := range lhp.Wastes {
		r.Wastes[name] = waste.Dup()
	}

	for name, wash := range lhp.Washes {
		r.Washes[name] = wash.Dup()
	}

	for name, dev := range lhp.Devices {
		r.Devices[name] = dev
	}

	for name, head := range lhp.Heads {
		r.Heads[name] = head.Dup()
	}

	for i, hl := range lhp.HeadsLoaded {
		r.HeadsLoaded[i] = hl.Dup()
	}

	for i, ad := range lhp.Adaptors {
		r.Adaptors[i] = ad.Dup()
	}

	for _, tip := range lhp.Tips {
		r.Tips = append(r.Tips, tip.Dup())
	}

	for _, pref := range lhp.Tip_preferences {
		r.Tip_preferences = append(r.Tip_preferences, pref)
	}
	for _, pref := range lhp.Input_preferences {
		r.Input_preferences = append(r.Input_preferences, pref)
	}
	for _, pref := range lhp.Output_preferences {
		r.Output_preferences = append(r.Output_preferences, pref)
	}
	for _, pref := range lhp.Waste_preferences {
		r.Waste_preferences = append(r.Waste_preferences, pref)
	}
	for _, pref := range lhp.Tipwaste_preferences {
		r.Tipwaste_preferences = append(r.Tipwaste_preferences, pref)
	}
	for _, pref := range lhp.Wash_preferences {
		r.Wash_preferences = append(r.Wash_preferences, pref)
	}

	if lhp.CurrConf != nil {
		r.CurrConf = lhp.CurrConf.Dup()
	}

	for i, v := range lhp.Cnfvol {
		r.Cnfvol[i] = v
	}

	for i, v := range lhp.Layout {
		r.Layout[i] = v
	}

	r.MaterialType = lhp.MaterialType

	return r
}

// constructor for the above
func NewLHProperties(num_positions int, model, manufacturer, lhtype, tiptype string, layout map[string]wtype.Coordinates) *LHProperties {
	var lhp LHProperties

	lhp.Nposns = num_positions

	lhp.Model = model
	lhp.Mnfr = manufacturer
	lhp.LHType = lhtype
	lhp.TipType = tiptype

	lhp.Adaptors = make([]*wtype.LHAdaptor, 0, 2)
	lhp.Heads = make([]*wtype.LHHead, 0, 2)
	lhp.HeadsLoaded = make([]*wtype.LHHead, 0, 2)

	positions := make(map[string]*wtype.LHPosition, num_positions)

	for i := 0; i < num_positions; i++ {
		// not overriding these defaults seems like a
		// bad idea --- TODO: Fix, e.g., MAXH here
		posname := fmt.Sprintf("position_%d", i+1)
		positions[posname] = wtype.NewLHPosition(i+1, "position_"+strconv.Itoa(i+1), 80.0)
	}

	lhp.Positions = positions
	lhp.PosLookup = make(map[string]string, lhp.Nposns)
	lhp.PlateLookup = make(map[string]interface{}, lhp.Nposns)
	lhp.PlateIDLookup = make(map[string]string, lhp.Nposns)
	lhp.Plates = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Tipboxes = make(map[string]*wtype.LHTipbox, lhp.Nposns)
	lhp.Tipwastes = make(map[string]*wtype.LHTipwaste, lhp.Nposns)
	lhp.Wastes = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Washes = make(map[string]*wtype.LHPlate, lhp.Nposns)
	lhp.Devices = make(map[string]string, lhp.Nposns)
	lhp.Heads = make([]*wtype.LHHead, 0, 2)
	lhp.Tips = make([]*wtype.LHTip, 0, 3)

	lhp.Layout = layout

	// lhp.Curcnf, lhp.Cmnvol etc. intentionally left blank

	lhp.MaterialType = material.DEVICE

	return &lhp
}

func (lhp *LHProperties) AddTipBox(tipbox *wtype.LHTipbox) bool {
	for _, pref := range lhp.Tip_preferences {
		if lhp.PosLookup[pref] != "" {
			//fmt.Println(pref, " ", lhp.PlateLookup[lhp.PosLookup[pref]])
			continue
		}

		lhp.AddTipBoxTo(pref, tipbox)
		return true
	}

	logger.Debug("NO TIP SPACES LEFT")
	return false
}
func (lhp *LHProperties) AddTipBoxTo(pos string, tipbox *wtype.LHTipbox) bool {
	/*
		fmt.Println("Adding tip box of type, ", tipbox.Type, " To position ", pos)
		if lhp.PosLookup[pos] != "" {
			logger.Fatal("CAN'T ADD TIPBOX TO FULL POSITION")
			panic("CAN'T ADD TIPBOX TO FULL POSITION")
		}
	*/

	if lhp.PosLookup[pos] != "" {
		logger.Debug(fmt.Sprintf("Tried to add tipbox to full position: %s", pos))
		return false
	}
	lhp.Tipboxes[pos] = tipbox
	lhp.PlateLookup[tipbox.ID] = tipbox
	lhp.PosLookup[pos] = tipbox.ID
	lhp.PlateIDLookup[tipbox.ID] = pos

	return true
}

func (lhp *LHProperties) RemoveTipBoxes() {
	for pos, tbx := range lhp.Tipboxes {
		lhp.PlateLookup[tbx.ID] = nil
		lhp.PosLookup[pos] = ""
		lhp.PlateIDLookup[tbx.ID] = ""
	}

	lhp.Tipboxes = make(map[string]*wtype.LHTipbox)
}

func (lhp *LHProperties) AddTipWaste(tipwaste *wtype.LHTipwaste) bool {
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {

			//fmt.Println(pref, " ", lhp.PlateLookup[lhp.PosLookup[pref]])

			continue
		}

		lhp.AddTipWasteTo(pref, tipwaste)
		return true
	}
	/*
		logger.Fatal("NO TIPWASTE SPACES LEFT")
		panic("NO TIPWASTE SPACES LEFT")
	*/

	logger.Debug("NO TIPWASTE SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddTipWasteTo(pos string, tipwaste *wtype.LHTipwaste) bool {
	if lhp.PosLookup[pos] != "" {
		logger.Debug("CAN'T ADD TIPWASTE TO FULL POSITION")
		//panic("CAN'T ADD TIPWASTE TO FULL POSITION")
		return false
	}
	lhp.Tipwastes[pos] = tipwaste
	lhp.PlateLookup[tipwaste.ID] = tipwaste
	lhp.PosLookup[pos] = tipwaste.ID
	lhp.PlateIDLookup[tipwaste.ID] = pos
	return true
}

func (lhp *LHProperties) AddPlate(pos string, plate *wtype.LHPlate) bool {
	if lhp.PosLookup[pos] != "" {
		logger.Debug("CAN'T ADD PLATE TO FULL POSITION")
		return false
	}
	lhp.Plates[pos] = plate
	lhp.PlateLookup[plate.ID] = plate
	lhp.PosLookup[pos] = plate.ID
	lhp.PlateIDLookup[plate.ID] = pos
	return true
}

func (lhp *LHProperties) addWaste(waste *wtype.LHPlate) bool {
	for _, pref := range lhp.Waste_preferences {
		if lhp.PosLookup[pref] != "" {

			//fmt.Println(pref, " ", lhp.PlateLookup[lhp.PosLookup[pref]])

			continue
		}

		lhp.AddWasteTo(pref, waste)
		return true
	}

	logger.Debug("NO WASTE SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddWasteTo(pos string, waste *wtype.LHPlate) bool {
	if lhp.PosLookup[pos] != "" {
		logger.Debug("CAN'T ADD WASTE TO FULL POSITION")
		return false
	}
	lhp.Wastes[pos] = waste
	lhp.PlateLookup[waste.ID] = waste
	lhp.PosLookup[pos] = waste.ID
	lhp.PlateIDLookup[waste.ID] = pos
	return true
}

func (lhp *LHProperties) AddWash(wash *wtype.LHPlate) bool {
	for _, pref := range lhp.Wash_preferences {
		if lhp.PosLookup[pref] != "" {
			//fmt.Println(pref, " ", lhp.PlateLookup[lhp.PosLookup[pref]])
			continue
		}

		lhp.AddWashTo(pref, wash)
		return true
	}

	logger.Debug("NO WASH SPACES LEFT")
	return false
}

func (lhp *LHProperties) AddWashTo(pos string, wash *wtype.LHPlate) bool {
	if lhp.PosLookup[pos] != "" {

		logger.Debug("CAN'T ADD WASH TO FULL POSITION")
		return false
	}
	lhp.Washes[pos] = wash
	lhp.PlateLookup[wash.ID] = wash
	lhp.PosLookup[pos] = wash.ID
	lhp.PlateIDLookup[wash.ID] = pos
	return true
}

func (lhp *LHProperties) GetCleanTips(tiptype string, channel *wtype.LHChannelParameter, mirror bool, multi int) (wells, positions, boxtypes []string) {
	positions = make([]string, multi)
	boxtypes = make([]string, multi)

	// hack -- count tips left
	n_tips_left := 0
	for _, pos := range lhp.Tip_preferences {
		bx, ok := lhp.Tipboxes[pos]

		if !ok || bx.Tiptype.Type != tiptype {
			continue
		}

		n_tips_left += bx.N_clean_tips()
	}

	//	logger.Debug(fmt.Sprintf("There are %d clean tips of type %s left", n_tips_left, tiptype))

	foundit := false

	for _, pos := range lhp.Tip_preferences {
		//	for i := len(lhp.Tip_preferences) - 1; i >= 0; i-- {
		//		pos := lhp.Tip_preferences[i]
		bx, ok := lhp.Tipboxes[pos]
		if !ok || bx.Tiptype.Type != tiptype {
			continue
		}
		wells = bx.GetTips(mirror, multi, channel.Orientation)
		if wells != nil {
			foundit = true
			for i := 0; i < multi; i++ {
				positions[i] = pos
				boxtypes[i] = bx.Boxname
			}
			break
		}
	}

	// if you don't find any suitable tips, why just make a
	// new box full of them!
	// nothing can possibly go wrong
	// surely

	if !foundit {

		// try adding a new tip box
		bx := factory.GetTipboxByType(tiptype)

		if bx == nil {
			return nil, nil, nil
		}

		r := lhp.AddTipBox(bx)

		if !r {
			return nil, nil, nil
		}

		return lhp.GetCleanTips(tiptype, channel, mirror, multi)
		//		return nil, nil, nil
	}

	return
}

func (lhp *LHProperties) DropDirtyTips(channel *wtype.LHChannelParameter, multi int) (wells, positions, boxtypes []string) {
	wells = make([]string, multi)
	positions = make([]string, multi)
	boxtypes = make([]string, multi)

	foundit := false

	for pos, bx := range lhp.Tipwastes {
		yes := bx.Dispose(multi)
		if yes {
			foundit = true
			for i := 0; i < multi; i++ {
				wells[i] = "A1"
				positions[i] = pos
				boxtypes[i] = bx.Type
			}

			break
		}
	}

	if !foundit {
		return nil, nil, nil
	}

	return
}

//GetMaterialType implement stockableMaterial
func (lhp *LHProperties) GetMaterialType() material.MaterialType {
	return lhp.MaterialType
}
func (lhp *LHProperties) GetTimer() *LHTimer {
	return GetTimerFor(lhp.Mnfr, lhp.Model)
}

// records timing info
// preliminary implementation assumes all instructions of a given
// type have the same timing, TimeFor is expressed in terms of the instruction
// however so it will be possible to modify this behaviour in future

type LHTimer struct {
	Times []time.Duration
}

func NewTimer() *LHTimer {
	var t LHTimer
	t.Times = make([]time.Duration, 50)
	return &t
}

func (t *LHTimer) TimeFor(r RobotInstruction) time.Duration {
	var d time.Duration
	if r.InstructionType() > 0 && r.InstructionType() < len(t.Times) {
		d = t.Times[r.InstructionType()]
	} else {
	}
	return d
}
