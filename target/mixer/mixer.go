package mixer

import (
	"errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	cannotGetCap = errors.New("cannot get capabilities")
)

var (
	_ target.Device = &Mixer{}
	_ target.Mixer  = &Mixer{}
	_ target.Shaper = &Mixer{}
)

type Mixer struct {
	driver     driver.ExtendedLiquidhandlingDriver
	properties driver.LHProperties
	opt        Opt
	Out        *lh.LHRequest // XXX: remove me
}

func (a *Mixer) Can(...target.Request) bool {
	return true // XXX: implement me
}

func (a *Mixer) MoveCost(from target.Device) int {
	if from == a {
		return 0
	}
	return human.HumanByXCost - 1
}

func (a *Mixer) Shape() interface{} {
	return &a.properties
}

func (a *Mixer) makeReq() (*lh.LHRequest, *lh.Liquidhandler) {
	req := lh.NewLHRequest()
	req.Policies = driver.GetLHPolicyForTest()
	planner := lh.Init(&a.properties)

	if p := a.opt.MaxPlates; p != nil {
		req.Input_setup_weights["MAX_N_PLATES"] = *p
	}

	if p := a.opt.MaxWells; p != nil {
		req.Input_setup_weights["MAX_N_WELLS"] = *p
	}

	if p := a.opt.ResidualVolumeWeight; p != nil {
		req.Input_setup_weights["RESIDUAL_VOLUME_WEIGHT"] = *p
	}

	if p := a.opt.InputPlateType; len(p) != 0 {
		for _, v := range p {
			req.Input_platetypes = append(req.Input_platetypes, factory.GetPlateByType(v))
		}
	}

	if p := a.opt.OutputPlateType; len(p) != 0 {
		for _, v := range p {
			req.Output_platetypes = append(req.Output_platetypes, factory.GetPlateByType(v))
		}
	}

	return req, planner
}

func (a *Mixer) PrepareMix(mixes []*wtype.LHInstruction) (*target.MixResult, error) {
	hasPlate := func(plates []*wtype.LHPlate, typ, id string) bool {
		for _, p := range plates {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	getId := func(mixes []*wtype.LHInstruction) (r wtype.BlockID) {
		m := make(map[wtype.BlockID]bool)
		for _, mix := range mixes {
			m[mix.BlockID] = true
		}
		if len(m) > 1 {
			panic("aa") // XXX
		}
		for k := range m {
			r = k
			break
		}
		return
	}

	req, planner := a.makeReq()
	req.BlockID = getId(mixes)

	for _, mix := range mixes {
		if len(mix.Platetype) != 0 && !hasPlate(req.Output_platetypes, mix.Platetype, mix.PlateID) {
			p := factory.GetPlateByType(mix.Platetype)
			p.ID = mix.PlateID
			req.Output_platetypes = append(req.Output_platetypes, p)
		}
		req.Add_instruction(mix)
	}

	planner.MakeSolutions(req)

	a.Out = req

	return &target.MixResult{
		Request: req,
	}, nil
}

func New(opt Opt, d driver.ExtendedLiquidhandlingDriver) (*Mixer, error) {
	p, status := d.GetCapabilities()
	if !status.OK {
		return nil, cannotGetCap
	}
	p.Driver = d
	return &Mixer{driver: d, properties: p, opt: opt}, nil
}
