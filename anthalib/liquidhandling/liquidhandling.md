---
layout: default
type: api
navgroup: docs
shortname: anthalib/liquidhandling
title: anthalib/liquidhandling
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: anthalib/liquidhandling
---
# liquidhandling
--
    import "."

defines types for dealing with liquid handling requests

## Usage

```go
const (
	ASP int = iota
	DSP
	MOV
	LOD
	ULD
	TFR
)
```

#### func  Init

```go
func Init(properties *LHProperties) *liquidhandler
```
initialize the liquid handling structure

#### func  MakeConfigFile

```go
func MakeConfigFile(fn string, request LHRequest)
```

#### func  MakePlanFile

```go
func MakePlanFile(fn string, request LHRequest)
```

#### func  Make_layout

```go
func Make_layout()
```

#### func  PlateLookup

```go
func PlateLookup(rq LHRequest, id string) int
```
looks up where a plate is mounted on a liquid handler as expressed in a request

#### func  RunLiquidHandler

```go
func RunLiquidHandler(*chan *LHRequest)
```
tell the liquid handler to run

#### func  SimpleOutput

```go
func SimpleOutput(ti TransferInstruction, rq LHRequest) []RobotInstruction
```
placeholder function: the intention is to have flexible rewriting of transfers

#### func  SimpleTransfer

```go
func SimpleTransfer(posfrom, posto []int, wellfrom, wellto []string, amount []float64, unit []string, what []string, prms *LHParameter) []RobotInstruction
```

#### type AspirateInstruction

```go
type AspirateInstruction struct {
	Vol           float64
	Volunit       string
	Speed         float64
	Speedunit     string
	ComponentType string
}
```


#### func  Aspirate

```go
func Aspirate(vol float64, volunit string, speed float64, speedunit string, what string) AspirateInstruction
```

#### func (AspirateInstruction) GetParameter

```go
func (ins AspirateInstruction) GetParameter(s string) interface{}
```

#### func (AspirateInstruction) InstructionType

```go
func (ai AspirateInstruction) InstructionType() int
```

#### type DispenseInstruction

```go
type DispenseInstruction struct {
	Vol           float64
	Volunit       string
	Speed         float64
	Speedunit     string
	ComponentType string
}
```


#### func  Dispense

```go
func Dispense(vol float64, volunit string, speed float64, speedunit string, what string) DispenseInstruction
```

#### func (DispenseInstruction) GetParameter

```go
func (ins DispenseInstruction) GetParameter(s string) interface{}
```

#### func (DispenseInstruction) InstructionType

```go
func (di DispenseInstruction) InstructionType() int
```

#### type LHComponent

```go
type LHComponent struct {
	ID          string
	Inst        string
	Order       int
	Name        string
	Type        string
	Vol         float64
	Conc        float64
	Vunit       string
	Cunit       string
	Tvol        float64
	Loc         string
	Smax        float64
	Destination string
}
```

structure describing a liquid component and its desired properties

#### func  NewLHComponent

```go
func NewLHComponent() *LHComponent
```

#### type LHDevice

```go
type LHDevice struct {
	ID   string
	Name string
	Mnfr string
}
```


#### func  NewLHDevice

```go
func NewLHDevice(name, mfr string) *LHDevice
```

#### type LHParameter

```go
type LHParameter struct {
	ID      string
	Name    string
	Minvol  float64
	Maxvol  float64
	Volunit string
	Policy  LHPolicy
}
```

describes sets of parameters which can be used to create a configuration

#### func  NewLHParameter

```go
func NewLHParameter(name string, minvol, maxvol float64, volunit string) *LHParameter
```

#### type LHPlate

```go
type LHPlate struct {
	ID         string
	Inst       string
	Loc        string
	Name       string
	Type       string
	Mnfr       string
	WellsX     int
	WellsY     int
	Nwells     int
	Wells      map[string]*LHWell
	Height     float64
	Hunit      string
	Rows       [][]*LHWell
	Cols       [][]*LHWell
	Welltype   *LHWell
	Wellcoords map[string]*LHWell
}
```

structure describing a microplate this needs to be harmonised with the wtype
version

#### func  NewLHPlate

```go
func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell) *LHPlate
```

#### func (*LHPlate) MarshalJSON

```go
func (plate *LHPlate) MarshalJSON() ([]byte, error)
```

#### func (*LHPlate) UnmarshalJSON

```go
func (plate *LHPlate) UnmarshalJSON(b []byte) error
```

#### func (*LHPlate) Welldimensions

```go
func (plate *LHPlate) Welldimensions() *LHWellType
```

#### type LHPolicy

```go
type LHPolicy map[string]interface{}
```


#### type LHPosition

```go
type LHPosition struct {
	ID    string
	Name  string
	Num   int
	Extra []LHDevice
	Maxh  float64
}
```

describes a position on the liquid handling deck and its current state

#### func  NewLHPosition

```go
func NewLHPosition(position_number int, name string, maxh float64) *LHPosition
```

#### type LHProperties

```go
type LHProperties struct {
	ID                 string
	Nposns             int
	Positions          []*LHPosition
	Model              string
	Manfr              string
	LHType             string
	TPType             string
	Formats            []string
	Cnfvol             []*LHParameter
	CurrConf           *LHParameter
	Tip_preferences    []int
	Input_preferences  []int
	Output_preferences []int
}
```

describes a liquid handler, its capabilities and current state

#### func  NewLHProperties

```go
func NewLHProperties(num_positions int, model, manufacturer, lhtype, tptype string, formats []string) *LHProperties
```
constructor for the above

#### type LHRequest

```go
type LHRequest struct {
	ID                         string
	Output_solutions           map[string]*LHSolution
	Input_solutions            map[string][]*LHComponent
	Plates                     map[string]*LHPlate
	Tips                       []*LHTipbox
	Locats                     []string
	Setup                      LHSetup
	Instructions               []RobotInstruction
	Robotfn                    string
	Input_assignments          map[string][]string
	Output_assignments         []string
	Input_plates               map[string]*LHPlate
	Output_plates              map[string]*LHPlate
	Input_platetype            *LHPlate
	Input_major_group_layouts  map[int][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[int]string
	Output_platetype           *LHPlate
	Output_major_group_layouts map[int][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[int]string
	Plate_lookup               map[string]int
	Stockconcs                 map[string]float64
}
```

structure for defining a request to the liquid handler

#### func  AdvancedExecutionPlanner

```go
func AdvancedExecutionPlanner(request *LHRequest, parameters *LHProperties) *LHRequest
```

#### func  BasicExecutionPlanner

```go
func BasicExecutionPlanner(request *LHRequest, params *LHProperties) *LHRequest
```
a default execution planner which relies on a call to code external to the Antha
project.

#### func  BasicLayoutAgent

```go
func BasicLayoutAgent(request *LHRequest, params *LHProperties) *LHRequest
```
default layout: requests fill plates in column order

#### func  BasicSetupAgent

```go
func BasicSetupAgent(request *LHRequest, params *LHProperties) *LHRequest
```
default setup agent

#### func  NewLHRequest

```go
func NewLHRequest() *LHRequest
```

#### func (*LHRequest) MarshalJSON

```go
func (req *LHRequest) MarshalJSON() ([]byte, error)
```

#### func (*LHRequest) UnmarshalJSON

```go
func (req *LHRequest) UnmarshalJSON(ar []byte) error
```

#### type LHSetup

```go
type LHSetup map[string]interface{}
```


#### func  NewLHSetup

```go
func NewLHSetup() LHSetup
```

#### type LHSolution

```go
type LHSolution struct {
	ID               string
	Inst             string
	Name             string
	Order            int
	Components       []*LHComponent
	ContainerType    string
	Welladdress      string
	Platetype        string
	Vol              float64
	Type             string
	Conc             float64
	Tvol             float64
	Majorlayoutgroup int
	Minorlayoutgroup int
}
```

structure describing a solution: a combination of liquid components

#### func  NewLHSolution

```go
func NewLHSolution() *LHSolution
```

#### func (LHSolution) GetComponentVolume

```go
func (sol LHSolution) GetComponentVolume(key string) float64
```

#### type LHTip

```go
type LHTip struct {
	ID       string
	Mnfr     string
	Type     string
	Minvol   float64
	Maxvol   float64
	Curvol   float64
	Contents string
	Dirty    bool
}
```


#### func  NewLHTip

```go
func NewLHTip(manufacturer, tiptype string, minvol, maxvol float64) *LHTip
```

#### type LHTipbox

```go
type LHTipbox struct {
	ID    string
	Type  string
	Mnfr  string
	Nrows int
	Ncols int
	Tips  map[string]*LHTipholder
}
```


#### func  NewLHTipbox

```go
func NewLHTipbox(nrows, ncols int, manufacturer string, tiptype *LHTip) *LHTipbox
```

#### type LHTipholder

```go
type LHTipholder struct {
	ID       string
	ParentID string
	Contents []*LHTip
}
```


#### func  NewLHTipholder

```go
func NewLHTipholder(parentid string) *LHTipholder
```

#### type LHWell

```go
type LHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Platetype string
	Coords    string
	Vol       float64
	Vunit     string
	Contents  []*LHComponent
	Rvol      float64
	Currvol   float64
	Shape     int
	Bottom    int
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
}
```

structure representing a well on a microplate - description of a destination

#### func  NewLHWell

```go
func NewLHWell(platetype, plateid, crds string, vol, rvol float64, shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell
```
make a new well structure

#### func  NewLHWellCopy

```go
func NewLHWellCopy(template *LHWell) *LHWell
```

#### func (*LHWell) AddDimensions

```go
func (w *LHWell) AddDimensions(lhwt *LHWellType)
```

#### func (*LHWell) MarshalJSON

```go
func (well *LHWell) MarshalJSON() ([]byte, error)
```

#### func (*LHWell) UnmarshalJSON

```go
func (well *LHWell) UnmarshalJSON(ar []byte) error
```

#### type LHWellType

```go
type LHWellType struct {
	Vol     float64
	Vunit   string
	Rvol    float64
	Shape   int
	Bottom  int
	Xdim    float64
	Ydim    float64
	Zdim    float64
	Bottomh float64
	Dunit   string
}
```


#### type LoadInstruction

```go
type LoadInstruction struct {
}
```


#### func  Load

```go
func Load() LoadInstruction
```

#### func (LoadInstruction) GetParameter

```go
func (ins LoadInstruction) GetParameter(s string) interface{}
```

#### func (LoadInstruction) InstructionType

```go
func (li LoadInstruction) InstructionType() int
```

#### type MoveInstruction

```go
type MoveInstruction struct {
	Pos           int
	Well          string
	Height        int
	OffsetX       float64
	OffsetY       float64
	OffsetZ       float64
	ComponentType string
}
```


#### func  Move

```go
func Move(pos int, well string, height int, offsetX, offsetY, offsetZ float64, what string) MoveInstruction
```

#### func (MoveInstruction) GetParameter

```go
func (ins MoveInstruction) GetParameter(s string) interface{}
```

#### func (MoveInstruction) InstructionType

```go
func (mi MoveInstruction) InstructionType() int
```

#### type RobotInstruction

```go
type RobotInstruction interface {
	InstructionType() int
	GetParameter(name string) interface{}
}
```


#### type RobotOutputInterface

```go
type RobotOutputInterface struct {
	InstructionOutputs []string
}
```


#### func  NewOutputInterface

```go
func NewOutputInterface(filename string) RobotOutputInterface
```

#### func (RobotOutputInterface) Output

```go
func (self RobotOutputInterface) Output(ins RobotInstruction) string
```

#### func (RobotOutputInterface) ReplacePlaceholders

```go
func (self RobotOutputInterface) ReplacePlaceholders(s string, ins RobotInstruction) string
```

#### type SLHPlate

```go
type SLHPlate struct {
	ID             string
	Inst           string
	Loc            string
	Name           string
	Type           string
	Mnfr           string
	WellsX         int
	WellsY         int
	Nwells         int
	Height         float64
	Hunit          string
	Welltype       *LHWell
	Wellcoords     map[string]*LHWell
	Welldimensions *LHWellType
}
```

serializable, stripped-down version of the LHPlate

#### func (SLHPlate) FillPlate

```go
func (slhp SLHPlate) FillPlate(plate *LHPlate)
```

#### type SLHRequest

```go
type SLHRequest struct {
	ID                         string
	Output_solutions           map[string]*LHSolution
	Input_solutions            map[string][]*LHComponent
	Plates                     map[string]*LHPlate
	Tips                       []*LHTipbox
	Locats                     []string
	Setup                      LHSetup
	Instructions               []RobotInstruction
	Robotfn                    string
	Input_assignments          map[string][]string
	Output_plates              map[string]*LHPlate
	Input_platetype            *LHPlate
	Input_major_group_layouts  map[string][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         map[string]string
	Output_platetype           *LHPlate
	Output_major_group_layouts map[string][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        map[string]string
	Plate_lookup               map[string]int
	Stockconcs                 map[string]float64
}
```


#### type SLHWell

```go
type SLHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Coords    string
	Contents  []*LHComponent
	Currvol   float64
}
```


#### func (SLHWell) FillWell

```go
func (slw SLHWell) FillWell(lw *LHWell)
```

#### type TransferInstruction

```go
type TransferInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []float64 // this could be a Measurement
	VolumeUnit []string
	Prms       *LHParameter
}
```


#### func  Transfer

```go
func Transfer(what []string, pfrom, pto []string, wfrom, wto []string, v []float64, vu []string, prms *LHParameter) TransferInstruction
```

#### func (TransferInstruction) GetParameter

```go
func (ins TransferInstruction) GetParameter(s string) interface{}
```

#### func (TransferInstruction) InstructionType

```go
func (ti TransferInstruction) InstructionType() int
```

#### type TransferOutputFunc

```go
type TransferOutputFunc func(TransferInstruction) []RobotInstruction
```


#### type UnloadInstruction

```go
type UnloadInstruction struct {
}
```


#### func  Unload

```go
func Unload() UnloadInstruction
```

#### func (UnloadInstruction) GetParameter

```go
func (ins UnloadInstruction) GetParameter(s string) interface{}
```

#### func (UnloadInstruction) InstructionType

```go
func (ui UnloadInstruction) InstructionType() int
```
