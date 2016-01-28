package trace

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

// An instruction
type Inst interface {
	Args() []Value
}

type NoopInst struct{}

func (a *NoopInst) Args() []Value {
	return nil
}

func MakeNoop() *NoopInst {
	return &NoopInst{}
}

type DebugInst struct {
	Message string
	Values  []Value
}

func (a *DebugInst) Args() []Value {
	return a.Values
}

func MakeDebug(msg string, values []Value) *DebugInst {
	return &DebugInst{Message: msg, Values: values}
}

type MixInst struct {
	Opt    MixOpt
	Values []Value
}

func (a *MixInst) Args() []Value {
	return a.Values
}

type MixOpt struct {
	OutputSol *wtype.LHSolution
	OutPlate  *wtype.LHPlate
	PlateType string
	Address   string
	PlateNum  int
	Caller    string
}

func MakeMix(opt MixOpt, values []Value) *MixInst {
	return &MixInst{Opt: opt, Values: values}
}

type IncubateInst struct {
	Opt   IncubateOpt
	Value Value
}

func (a *IncubateInst) Args() []Value {
	return []Value{a.Value}
}

type IncubateOpt struct {
	BlockID      string
	OutputSol    *wtype.LHSolution
	Temp         wunit.Temperature
	Time         wunit.Time
	ShakingForce interface{}
}

func MakeIncubate(opt IncubateOpt, value Value) *IncubateInst {
	return &IncubateInst{Opt: opt, Value: value}
}

// Pair to track instructions with their promised return values
type instp struct {
	inst    Inst
	promise *Promise
}

// Issue an instruction to execute in the trace context; return promise for
// return value.
func Issue(ctx context.Context, inst Inst) *Promise {
	result := getScope(ctx).MakeName("")
	out := make(chan interface{})

	p := &Promise{
		construct: func(v interface{}) Value {
			return &promiseValue{
				name: result,
				inst: inst,
				v:    v,
			}
		},
		out: out,
	}

	t := getTrace(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()

	t.issue = append(t.issue, instp{inst: inst, promise: p})

	return p
}

// Force code generation and execution for current set of issued commands
func Flush(ctx context.Context) {
	p := Issue(ctx, MakeNoop())
	Read(ctx, p)
}
