package liquidhandling

import (
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// it would probably make more sense for this to be a method of the robot
// in general the instruction generator might well be moved there wholesale
// so that drivers can have specific versions of it... this could lead to some
// very interesting situations though

type ChannelScore float64

type ChannelScoreFunc interface {
	ScoreCombinedChannel(wunit.Volume, *wtype.LHHead, *wtype.LHAdaptor, *wtype.LHTip) ChannelScore
}

type DefaultChannelScoreFunc struct {
}

func (sc DefaultChannelScoreFunc) ScoreCombinedChannel(vol wunit.Volume, head *wtype.LHHead, adaptor *wtype.LHAdaptor, tip *wtype.LHTip) ChannelScore {
	// something pretty simple
	// higher is better
	// 0 == don't bother

	// first we merge the parameters together and see if we can do this at all

	lhcp := head.Params.MergeWithTip(tip)

	// we should always make sure we do not send a volume which is too low

	if lhcp.Minvol.LessThan(vol) {
		return 0
	}

	// clearly now vol >= Minvol

	// the main idea is to estimate the error from each source: head, adaptor, tip
	// and make the choice on that basis

	// a big head with a tiny tip will make pretty big errors... a big tip on a tiny
	// head likewise

	// we therefore make our choice as Min(1/tip_error, 1/adaptor_error, 1/head_error)

	err := 999999999.0

	chans := []*wtype.LHChannelParameter{head.GetParams(), tip.GetParams()}

	for _, ch := range chans {
		myerr := sc.ScoreChannel(vol, ch)
		if myerr < err {
			err = myerr
		}
	}

	return ChannelScore(err)
}

func (sc DefaultChannelScoreFunc) ScoreChannel(vol wunit.Volume, lhcp *wtype.LHChannelParameter) float64 {
	// cannot have 0 error
	extra := 0.1
	// we try to estimate the error of using a channel outside its limits
	// first of all how many movements do we need?

	v := vol.RawValue()
	mx := lhcp.Maxvol.ConvertTo(vol.Unit())
	mn := lhcp.Minvol.ConvertTo(vol.Unit())

	n, _ := math.Modf(vol.RawValue() / lhcp.Maxvol.ConvertTo(vol.Unit()))

	// we assume errors scale linearly
	// and that the error is generally greatest at the lowest levels

	tv := v
	if n >= 1 {
		tv = mx
	}

	err := (mx-tv)/(mx-mn) + extra

	if n > 1 {
		err *= (n + 1)
	}

	score := 1.0 / err

	return score
}

func ChooseChannel(vol wunit.Volume, prms *LHProperties) (*wtype.LHChannelParameter, string) {
	var headchosen *wtype.LHHead = nil
	var tipchosen *wtype.LHTip = nil
	var bestscore ChannelScore = ChannelScore(0.0)

	scorer := prms.GetChannelScoreFunc()

	// just choose the best... need to iterate on this sometime though
	// we don't consider head or adaptor changes now

	for _, head := range prms.HeadsLoaded {
		for _, tip := range prms.Tips {
			sc := scorer.ScoreCombinedChannel(vol, head, head.Adaptor, tip)
			if sc > bestscore {
				headchosen = head
				tipchosen = tip
				bestscore = sc
			}
		}

	}

	if headchosen == nil {
		return nil, ""
	}

	// shouldn't we also return the adaptor?

	return headchosen.GetParams(), tipchosen.Type

	/*
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
			//if tip.Minvol.LessThan(vol) || tip.Minvol.EqualTo(vol) {
			//	tiptype = tip.Type
			//}

			dif := vol.ConvertTo(tip.Minvol.Unit()) - tip.Minvol.RawValue()
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
	*/
}
