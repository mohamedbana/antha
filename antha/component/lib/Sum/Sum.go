package Sum

import

// Input parameters for this protocol
(
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/execution"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"sync"
)

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *Sum) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *Sum) setup(p SumParamBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

// Core process of the protocol: steps to be performed for each input
func (e *Sum) steps(p SumParamBlock, r *SumResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

	r.Sum = p.A + p.B
}

// Actions to perform after steps block to analyze data
func (e *Sum) analysis(p SumParamBlock, r *SumResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

func (e *Sum) validation(p SumParamBlock, r *SumResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

// AsyncBag functions
func (e *Sum) Complete(params interface{}) {
	p := params.(SumParamBlock)
	if p.Error {
		e.Sum <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(SumResultBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)
	if r.Error {
		e.Sum <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

	e.Sum <- execute.ThreadParam{Value: r.Sum, ID: p.ID, Error: false}

	e.analysis(p, r)
	if r.Error {
		e.Sum <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

	e.validation(p, r)
	if r.Error {
		e.Sum <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Sum) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Sum) NewConfig() interface{} {
	return &SumConfig{}
}

func (e *Sum) NewParamBlock() interface{} {
	return &SumParamBlock{}
}

func NewSum() interface{} { //*Sum {
	e := new(Sum)
	e.init()
	return e
}

// Mapper function
func (e *Sum) Map(m map[string]interface{}) interface{} {
	var res SumParamBlock
	res.Error = false || m["A"].(execute.ThreadParam).Error || m["B"].(execute.ThreadParam).Error

	vA, is := m["A"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SumJSONBlock
		json.Unmarshal([]byte(vA.JSONString), &temp)
		res.A = *temp.A
	} else {
		res.A = m["A"].(execute.ThreadParam).Value.(int)
	}

	vB, is := m["B"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SumJSONBlock
		json.Unmarshal([]byte(vB.JSONString), &temp)
		res.B = *temp.B
	} else {
		res.B = m["B"].(execute.ThreadParam).Value.(int)
	}

	res.ID = m["A"].(execute.ThreadParam).ID

	return res
}

func (e *Sum) OnA(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("A", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Sum) OnB(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("B", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Sum struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	A              <-chan execute.ThreadParam
	B              <-chan execute.ThreadParam
	Sum            chan<- execute.ThreadParam
}

type SumParamBlock struct {
	ID    execute.ThreadID
	Error bool
	A     int
	B     int
}

type SumConfig struct {
	ID    execute.ThreadID
	Error bool
	A     int
	B     int
}

type SumResultBlock struct {
	ID    execute.ThreadID
	Error bool
	Sum   int
}

type SumJSONBlock struct {
	ID    *execute.ThreadID
	Error *bool
	A     *int
	B     *int
	Sum   *int
}

func (c *Sum) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("A", "int", "A", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("B", "int", "B", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sum", "int", "Sum", true, true, nil, nil))

	ci := execute.NewComponentInfo("Sum", "Sum", "", false, inp, outp)

	return ci
}
