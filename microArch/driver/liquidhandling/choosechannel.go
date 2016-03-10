package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/logger"
)

func ChooseChannel(vol wunit.Volume, prms *LHProperties) (*wtype.LHChannelParameter, string) {
	// very quick and dirty, basically just makes sure the minimum volume is
	// below the required volume

	var headchosen *wtype.LHHead = nil

	minvol := 99999.0

	v := vol.RawValue()

	for _, head := range prms.HeadsLoaded {
		//fmt.Println("Trying head ", head.Name, " Which has minimum volume ", head.Adaptor.Params.Minvol.ConvertTo(vol.Unit()))

		minv := head.Params.Minvol.ConvertTo(vol.Unit())
		maxv := head.Params.Maxvol.ConvertTo(vol.Unit())

		d := v - minv

		if d >= 0.0 && minv < minvol {

			if headchosen == nil {
				headchosen = head
				minvol = minv
			}

			if v <= maxv {

				if head.GetParams().Minvol.SIValue() < headchosen.GetParams().Minvol.SIValue() {
					headchosen = head
					minvol = minv
				}
			} else {
				if head.GetParams().Maxvol.SIValue() > headchosen.GetParams().Maxvol.SIValue() {
					headchosen = head
					minvol = minv
				}
				//minvol = minv
				//headchosen = head
			}
		}
		//headchosen = prms.Heads[0]
	}

	if headchosen == nil {
		logger.Fatal(fmt.Sprintf("Cannot find a head with suitable capacity to handle volume %s", vol.ToString()))
		panic("NO HEAD CHOSEN")
	}

	// check if we need to change adaptor

	//logger.Debug(fmt.Sprintf("want vol %s min vol %s", vol.ToString(), headchosen.Adaptor.Params.Minvol.ToString()))

	if headchosen.Adaptor.Params.Minvol.GreaterThan(vol) {
		logger.Fatal(fmt.Sprintf("Handling volume %s is possible but an adaptor change is required first. This is not presently implemented. Sorry.", vol.ToString()))
		panic("ADAPTOR CHANGE NEEDED BUT NOT IMPLEMENTED")
	}

	// now get the requisite tip
	// this is just a big bowl of wrong... </JeffGreene>
	// need to make this more dependent on what's actually there
	tiptype := ""
	// get the closest to the min vol
	d := 99999.0
	for _, tip := range prms.Tips {
		//if tip.MinVol.LessThan(vol) || tip.MinVol.EqualTo(vol) {
		//	tiptype = tip.Type
		//}

		dif := vol.ConvertTo(tip.MinVol.Unit()) - tip.MinVol.RawValue()
		if dif >= 0.0 && dif < d {
			tiptype = tip.Type
			d = dif
		}

	}

	if tiptype == "" {
		logger.Fatal(fmt.Sprintf("No tips are available for servicing a volume of %s.", vol.ToString()))
		panic("NO TIP TYPE FOUND")
	}

	return headchosen.GetParams(), tiptype
}
