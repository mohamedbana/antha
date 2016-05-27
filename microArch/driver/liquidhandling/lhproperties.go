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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/material"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/microArch/sampletracker"
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
	Driver               LiquidhandlingDriver `gotopb:"-"`
	CurrConf             *wtype.LHChannelParameter
	Cnfvol               []*wtype.LHChannelParameter
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
		r.Adaptors = append(r.Adaptors, a.Dup())
	}

	for _, h := range lhp.Heads {
		r.Heads = append(r.Heads, h.Dup())
	}

	for _, hl := range lhp.HeadsLoaded {
		r.HeadsLoaded = append(r.HeadsLoaded, hl.Dup())
	}

	for name, pl := range lhp.PosLookup {
		r.PosLookup[name] = pl
	}

	for name, pl := range lhp.PlateIDLookup {
		r.PlateIDLookup[name] = pl
	}

	// plate lookup can contain anything

	m := make(map[string]interface{})
	for name, pt := range lhp.PlateLookup {
		var pt2 interface{}
		switch pt.(type) {
		case *wtype.LHTipwaste:
			tmp := pt.(*wtype.LHTipwaste).Dup()
			tmp.ID = pt.(*wtype.LHTipwaste).ID
			pt2 = tmp
		case *wtype.LHPlate:
			tmp := pt.(*wtype.LHPlate).Dup()
			tmp.ID = pt.(*wtype.LHPlate).ID
			pt2 = tmp
		case *wtype.LHTipbox:
			tmp := pt.(*wtype.LHTipbox).Dup()
			tmp.ID = pt.(*wtype.LHTipbox).ID
			pt2 = tmp
		}
		m[name] = pt2
		r.PlateLookup[name] = pt2
	}

	for name, plate := range lhp.Plates {
		p2 := m[plate.ID]
		r.Plates[name] = p2.(*wtype.LHPlate)
	}

	for name, tb := range lhp.Tipboxes {
		tb2 := tb.Dup()
		tb2.ID = tb.ID
		r.Tipboxes[name] = tb2
	}

	for name, tw := range lhp.Tipwastes {
		tw2 := tw.Dup()
		tw2.ID = tw.ID
		r.Tipwastes[name] = tw2
	}

	for name, waste := range lhp.Wastes {
		w2 := waste.Dup()
		w2.ID = waste.ID
		r.Wastes[name] = w2
	}

	for name, wash := range lhp.Washes {
		w2 := wash.Dup()
		w2.ID = wash.ID
		r.Washes[name] = w2
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

func (lhp *LHProperties) TipsLeftOfType(tiptype string) int {
	n := 0

	for _, pref := range lhp.Tip_preferences {
		tb := lhp.Tipboxes[pref]
		if tb != nil {
			n += tb.N_clean_tips()
		}
	}

	return n
}

func (lhp *LHProperties) AddTipBox(tipbox *wtype.LHTipbox) error {
	/// TODO --- get rid of this!
	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip box")
	for _, pref := range lhp.Tip_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		lhp.AddTipBoxTo(pref, tipbox)
		return nil
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip box")
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

func (lhp *LHProperties) TipWastesMounted() int {
	r := 0
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			_, ok := lhp.Tipwastes[lhp.PosLookup[pref]]

			if !ok {
				logger.Debug(fmt.Sprintf("Position %s claims to have a tipbox but is empty", pref))
				continue
			}

			r += 1
		}
	}

	return r

}

func (lhp *LHProperties) TipSpacesLeft() int {
	r := 0
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			bx, ok := lhp.Tipwastes[lhp.PosLookup[pref]]

			if !ok {
				logger.Debug(fmt.Sprintf("Position %s claims to have a tipbox but is empty", pref))
				continue
			}

			r += bx.SpaceLeft()
		}
	}

	return r
}

func (lhp *LHProperties) AddTipWaste(tipwaste *wtype.LHTipwaste) error {
	for _, pref := range lhp.Tipwaste_preferences {
		if lhp.PosLookup[pref] != "" {
			continue
		}

		err := lhp.AddTipWasteTo(pref, tipwaste)
		return err
	}

	return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "Trying to add tip waste")
}

func (lhp *LHProperties) AddTipWasteTo(pos string, tipwaste *wtype.LHTipwaste) error {
	if lhp.PosLookup[pos] != "" {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add tip waste to full position %s", pos))
	}
	lhp.Tipwastes[pos] = tipwaste
	lhp.PlateLookup[tipwaste.ID] = tipwaste
	lhp.PosLookup[pos] = tipwaste.ID
	lhp.PlateIDLookup[tipwaste.ID] = pos
	return nil
}

func (lhp *LHProperties) AddPlate(pos string, plate *wtype.LHPlate) error {
	if lhp.PosLookup[pos] != "" {
		return wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, fmt.Sprintf("Trying to add plate to full position %s", pos))
	}
	lhp.Plates[pos] = plate
	lhp.PlateLookup[plate.ID] = plate
	lhp.PosLookup[pos] = plate.ID
	lhp.PlateIDLookup[plate.ID] = pos
	return nil
}

// reverse the above

func (lhp *LHProperties) RemovePlateWithID(id string) {
	pos := lhp.PlateIDLookup[id]
	delete(lhp.PosLookup, pos)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, pos)
}

func (lhp *LHProperties) RemovePlateAtPosition(pos string) {
	id := lhp.PosLookup[pos]
	delete(lhp.PosLookup, pos)
	delete(lhp.PlateIDLookup, id)
	delete(lhp.PlateLookup, id)
	delete(lhp.Plates, pos)
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

// of necessity, this must be destructive of state so we have to work on a copy
// NB this is not properly specified yet
func (lhp *LHProperties) GetComponents(cmps []*wtype.LHComponent) ([]string, []string, error) {
	r1 := make([]string, len(cmps))
	r2 := make([]string, len(cmps))

	fudgevol := wunit.NewVolume(0.5, "ul")

	for i, v := range cmps {
		foundIt := false

		if v.HasAnyParent() {
			//fmt.Println("Trying to get component ", v.CName, v.ParentID)
			// this means it was already made with a previous call
			tx := strings.Split(v.Loc, ":")

			// maybe we can look it up?

			if len(tx) < 2 || len(v.Loc) == 0 {
				st := sampletracker.GetSampleTracker()
				loc, _ := st.GetLocationOf(v.ID)
				tx = strings.Split(loc, ":")
			}

			r1[i] = tx[0]
			r2[i] = tx[1]

			vol := v.Volume().Dup()
			vol.Add(fudgevol)
			lhp.RemoveComponent(tx[0], tx[1], vol)

			foundIt = true

		} else {
			for _, ipref := range lhp.Input_preferences {
				// check if the plate at position ipref has the
				// component we seek

				p, ok := lhp.Plates[ipref]

				if ok {
					// whaddya got?
					// nb this won't work if we need to split a volume across several plates
					wcarr, ok := p.GetComponent(v, false)

					if ok {
						foundIt = true
						// update r1 and r2
						r1[i] = p.ID
						// XXX XXX XXX FFS this should aggregate
						// TODO -- fix this
						r2[i] = wcarr[0].FormatA1()

						vol := v.Volume().Dup()
						vol.Add(fudgevol)

						lhp.RemoveComponent(r1[i], r2[i], vol)

						break
					}
				}
			}

			if !foundIt {
				//logger.Fatal("NOSOURCE FOR ", v.CName, " at volume ", v.Volume().ToString())
				err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("NO SOURCE FOR ", v.CName, " at volume ", v.Volume().ToString()))
				return r1, r2, err
			}

		}
	}

	return r1, r2, nil
}

func (lhp *LHProperties) GetCleanTips(tiptype string, channel *wtype.LHChannelParameter, mirror bool, multi int) (wells, positions, boxtypes []string, err error) {
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
			err = wtype.LHError(wtype.LH_ERR_NO_TIPS, fmt.Sprint("No tipbox of type ", tiptype, " is known"))
			return nil, nil, nil, err
		}

		r := lhp.AddTipBox(bx)

		if r != nil {
			err = r
			return nil, nil, nil, err
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

func (lhp *LHProperties) GetChannelScoreFunc() ChannelScoreFunc {
	// this is to permit us to make this flexible

	sc := DefaultChannelScoreFunc{}

	return sc
}

// convenience method

func (lhp *LHProperties) RemoveComponent(plateID string, well string, volume wunit.Volume) bool {
	p := lhp.Plates[lhp.PlateIDLookup[plateID]]

	if p == nil {
		logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find plate"))
		return false
	}

	r := p.RemoveComponent(well, volume)

	if r == nil {
		logger.Info(fmt.Sprint("CAN'T REMOVE COMPONENT ", plateID, " ", well, " ", volume.ToString()))
		return false
	}

	/*
		w := p.Wellcoords[well]

		if w == nil {
			logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find well"))
			return false
		}

		c:=w.Remove(volume)

		if c==nil{
			logger.Info(fmt.Sprint("RemoveComponent ", plateID, " ", well, " ", volume.ToString(), " can't find well"))
			return false
		}
	*/

	return true
}

func (lhp *LHProperties) RemoveTemporaryComponents() {
	ids := make([]string, 0, 1)
	for _, p := range lhp.Plates {
		// if the whole plate is temporary we can just delete the whole thing
		if p.IsTemporary() {
			ids = append(ids, p.ID)
			continue
		}

		// now remove any components in wells still marked temporary

		for _, w := range p.Wellcoords {
			if w.IsTemporary() {
				w.Clear()
			}
		}
	}

	for _, id := range ids {
		lhp.RemovePlateWithID(id)
	}

	// good
}
