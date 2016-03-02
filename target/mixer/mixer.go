package mixer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	planner "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	cannotGetCap = errors.New("cannot get capabilities")
)

var (
	_ target.Device = &Mixer{}
)

type Mixer struct {
	driver     driver.ExtendedLiquidhandlingDriver
	properties driver.LHProperties
	opt        Opt
}

func (a *Mixer) Can(req ast.Request) bool {
	// TODO: remove when mixers have wait instruction
	if req.Time != nil {
		return false
	}
	if req.Temp != nil {
		return false
	}
	if req.Move != nil {
		return false
	}
	// TODO: Add specific volume constraints
	return req.MixVol != nil
}

func (a *Mixer) MoveCost(from target.Device) int {
	if from == a {
		return 0
	}
	return human.HumanByXCost + 1
}

func (a *Mixer) makeReq() (*planner.LHRequest, *planner.Liquidhandler) {
	req := planner.NewLHRequest()
	req.Policies = driver.GetLHPolicyForTest()
	p := planner.Init(&a.properties)

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

	return req, p
}

func (a *Mixer) Compile(cmds []ast.Command) ([]target.Inst, error) {
	var mixes []*wtype.LHInstruction
	for _, c := range cmds {
		if m, ok := c.(*ast.Mix); !ok {
			return nil, fmt.Errorf("cannot compile %T", c)
		} else {
			mixes = append(mixes, m.Inst)
		}
	}
	if inst, err := a.makeMix(mixes); err != nil {
		return nil, err
	} else {
		return []target.Inst{inst}, nil
	}
}

func (a *Mixer) saveFile(name string) ([]byte, error) {
	data, status := a.driver.GetOutputFile()
	if !status.OK {
		return nil, fmt.Errorf("%d: %s", status.Errorcode, status.Msg)
	} else if len(data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	w := tar.NewWriter(gzip.NewWriter(&buf))
	bs := []byte(data)

	if err := w.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    0644,
		Size:    int64(len(bs)),
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	} else if _, err := w.Write(bs); err != nil {
		return nil, err
	} else if err := w.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (a *Mixer) makeMix(mixes []*wtype.LHInstruction) (target.Inst, error) {
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
		for k := range m {
			r = k
			break
		}
		return
	}

	req, p := a.makeReq()
	req.BlockID = getId(mixes)

	for _, mix := range mixes {
		if len(mix.Platetype) != 0 && !hasPlate(req.Output_platetypes, mix.Platetype, mix.PlateID) {
			p := factory.GetPlateByType(mix.Platetype)
			p.ID = mix.PlateID
			req.Output_platetypes = append(req.Output_platetypes, p)
		}
		req.Add_instruction(mix)
	}

	p.MakeSolutions(req)

	tarball, err := a.saveFile("input")
	if err != nil {
		return nil, err
	}

	return &target.Mix{
		Request:    req,
		Properties: a.properties,
		Files:      tarball,
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
