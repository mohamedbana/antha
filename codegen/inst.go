package codegen

type Inst interface {
	Id() string
}

type LHInst struct{}

type IncubateInst struct{}

type ManualMoveInst struct{}

type ForkInst struct{}

type JoinInst struct{}
