---
layout: default
type: api
navgroup: docs
shortname: anthalib/liquidhandling
title: anthalib/liquidhandling
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: anthalib/liquidhandling
---
# liquidhandling
--
    import "."

defines types for dealing with liquid handling requests

## Usage

#### func  DefineOrderOrFail

```go
func DefineOrderOrFail(mapin map[string]map[string]int) []string
```

#### func  MakeConfigFile

```go
func MakeConfigFile(fn string, request LHRequest)
```

#### func  MakePlanFile

```go
func MakePlanFile(fn string, request LHRequest)
```

#### func  PlateLookup

```go
func PlateLookup(rq LHRequest, id string) string
```
looks up where a plate is mounted on a liquid handler as expressed in a request

#### func  RaiseError

```go
func RaiseError(err string)
```

#### type LHPolicyManager

```go
type LHPolicyManager struct {
	SystemPolicies *liquidhandling.LHPolicyRuleSet
	UserPolicies   *liquidhandling.LHPolicyRuleSet
}
```


#### func (*LHPolicyManager) MergePolicies

```go
func (mgr *LHPolicyManager) MergePolicies(protocolpolicies *liquidhandling.LHPolicyRuleSet) *liquidhandling.LHPolicyRuleSet
```

#### type LHRequest

```go
type LHRequest struct {
	ID                         string
	BlockID                    string
	BlockName                  string
	Output_solutions           map[string]*wtype.LHSolution
	Input_solutions            map[string][]*wtype.LHComponent
	Plates                     map[string]*wtype.LHPlate
	Tips                       []*wtype.LHTipbox
	Tip_Type                   *wtype.LHTipbox
	Locats                     []string
	Setup                      wtype.LHSetup
	InstructionSet             *liquidhandling.RobotInstructionSet
	Instructions               []liquidhandling.TerminalRobotInstruction
	Robotfn                    string
	Outputfn                   string
	Input_assignments          map[string][]string
	Output_assignments         []string
	Input_plates               map[string]*wtype.LHPlate
	Output_plates              map[string]*wtype.LHPlate
	Input_platetypes           []*wtype.LHPlate
	Input_major_group_layouts  map[int][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[int]string
	Input_Setup_Weights        map[string]float64
	Output_platetype           *wtype.LHPlate
	Output_major_group_layouts map[int][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[int]string
	Plate_lookup               map[string]string
	Stockconcs                 map[string]float64
	Policies                   *liquidhandling.LHPolicyRuleSet
	Input_order                []string
}
```

structure for defining a request to the liquid handler

#### func  AdvancedExecutionPlanner

```go
func AdvancedExecutionPlanner(request *LHRequest, parameters *liquidhandling.LHProperties) *LHRequest
```

#### func  BasicExecutionPlanner

```go
func BasicExecutionPlanner(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest
```
a default execution planner which relies on a call to code external to the Antha
project.

#### func  BasicLayoutAgent

```go
func BasicLayoutAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest
```
default layout: requests fill plates in column order

#### func  BasicSetupAgent

```go
func BasicSetupAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest
```
default setup agent

#### func  NewLHRequest

```go
func NewLHRequest() *LHRequest
```

#### func  Rationalise_Inputs

```go
func Rationalise_Inputs(lhr *LHRequest, lhp *liquidhandling.LHProperties, inputs ...*wtype.LHComponent) *LHRequest
```

#### func  Rationalise_Outputs

```go
func Rationalise_Outputs(lhr *LHRequest, lhp *liquidhandling.LHProperties, outputs ...*wtype.LHComponent) *LHRequest
```

#### func (*LHRequest) MarshalJSON

```go
func (req *LHRequest) MarshalJSON() ([]byte, error)
```

#### func (*LHRequest) UnmarshalJSON

```go
func (req *LHRequest) UnmarshalJSON(ar []byte) error
```

#### type Liquidhandler

```go
type Liquidhandler struct {
	Properties       *liquidhandling.LHProperties
	SetupAgent       func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	LayoutAgent      func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	ExecutionPlanner func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	PolicyManager    *LHPolicyManager
}
```

the liquid handler structure defines the interface to a particular liquid
handling platform. The structure holds the following items: - an LHRequest
structure defining the characteristics of the platform - a channel for
communicating with the liquid handler additionally three functions are defined
to implement platform-specific implementation requirements in each case the
LHRequest structure passed in has some additional information added and is then
passed out. Features which are already defined (e.g. by the scheduler or the
user) are respected as constraints and will be left unchanged. The three
functions define - setup (SetupAgent): How sources are assigned to plates and
plates to positions - layout (LayoutAgent): how experiments are assigned to
outputs - execution (ExecutionPlanner): generates instructions to implement the
required plan

The general mechanism by which requests which refer to specific items as opposed
to those which only state that an item of a particular kind is required is by
the definition of an 'inst' tag in the request structure with a guid. If this is
defined and valid it indicates that this item in the request (e.g. a plate,
stock etc.) is a specific instance. If this is absent then the GUID will either
be created or requested

#### func  Init

```go
func Init(properties *liquidhandling.LHProperties) *Liquidhandler
```
initialize the liquid handling structure

#### func (*Liquidhandler) Execute

```go
func (this *Liquidhandler) Execute(request *LHRequest) error
```
run the request via the driver

#### func (*Liquidhandler) ExecutionPlan

```go
func (this *Liquidhandler) ExecutionPlan(request *LHRequest) *LHRequest
```
make the instructions for executing this request

#### func (*Liquidhandler) GetInputs

```go
func (this *Liquidhandler) GetInputs(request *LHRequest) *LHRequest
```
request the inputs which are needed to run the plan, unless they have already
been requested

#### func (*Liquidhandler) GetPlates

```go
func (this *Liquidhandler) GetPlates(plates map[string]*wtype.LHPlate, major_layouts map[int][]string, ptype *wtype.LHPlate) map[string]*wtype.LHPlate
```
define which labware to use and request specific instances

#### func (*Liquidhandler) InitializeDriver

```go
func (this *Liquidhandler) InitializeDriver(request *LHRequest) error
```
TODO TODO TODO this call should not be here

#### func (*Liquidhandler) Layout

```go
func (this *Liquidhandler) Layout(request *LHRequest) *LHRequest
```
generate the output layout

#### func (*Liquidhandler) MakeSolutions

```go
func (this *Liquidhandler) MakeSolutions(request *LHRequest) *LHRequest
```
high-level function which requests planning and execution for an incoming set of
solutions

#### func (*Liquidhandler) Plan

```go
func (this *Liquidhandler) Plan(request *LHRequest)
```

#### func (*Liquidhandler) Setup

```go
func (this *Liquidhandler) Setup(request *LHRequest) *LHRequest
```
generate setup for the robot

#### func (*Liquidhandler) Tip_box_setup

```go
func (lh *Liquidhandler) Tip_box_setup(request *LHRequest) *LHRequest
```

    TASK: 	Determine number of tip boxes of each type
INPUT: instructions OUTPUT: arrays of tip boxes

#### type SLHRequest

```go
type SLHRequest struct {
	ID                         string
	Output_solutions           map[string]*wtype.LHSolution
	Input_solutions            map[string][]*wtype.LHComponent
	Plates                     map[string]*wtype.LHPlate
	Tips                       []*wtype.LHTipbox
	Locats                     []string
	Setup                      wtype.LHSetup
	InstructionSet             *liquidhandling.RobotInstructionSet
	Instructions               []liquidhandling.TerminalRobotInstruction
	Robotfn                    string
	Input_assignments          map[string][]string
	Output_plates              map[string]*wtype.LHPlate
	Input_platetypes           []*wtype.LHPlate
	Input_major_group_layouts  map[string][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[string]string
	Output_platetype           *wtype.LHPlate
	Output_major_group_layouts map[string][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[string]string
	Plate_lookup               map[string]string
	Stockconcs                 map[string]float64
	Policies                   *liquidhandling.LHPolicyRuleSet
}
```
