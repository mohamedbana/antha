---
layout: default
type: api
navgroup: docs
shortname: driver/liquidhandling
title: driver/liquidhandling
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: driver/liquidhandling
---
# liquidhandling
--
    import "github.com/antha-lang/antha/microArch/driver/liquidhandling"


## Usage

```go
const (
	LHP_AND int = iota
	LHP_OR
)
```

```go
const (
	TFR int = iota // Transfer
	CTF            // Channel Transfer
	SCB            // Single channel transfer block
	MCB            // Multi channel transfer block
	SCT            // Single channel transfer
	MCT            // multi channel transfer
	CCC            // ChangeChannelCharacteristics
	LDT            // Load Tips + Move
	UDT            // Unload Tips + Move
	RST            // Reset
	CHA            // ChangeAdaptor
	ASP            // Aspirate
	DSP            // Dispense
	BLO            // Blowout
	PTZ            // Reset pistons
	MOV            // Move
	MRW            // Move Raw
	LOD            // Load Tips
	ULD            // Unload Tips
	SUK            // Suck
	BLW            // Blow
	SPS            // Set Pipette Speed
	SDS            // Set Drive Speed
	INI            // Initialize
	FIN            // Finalize
	WAI            // Wait
	LON            // Lights On
	LOF            // Lights Off
	OPN            // Open
	CLS            // Close
	LAD            // Load Adaptor
	UAD            // Unload Adaptor
	MMX            // Move and Mix
	MIX            // Mix
)
```

```go
var RobotParameters = []string{"HEAD", "CHANNEL", "LIQUIDCLASS", "POSTO", "WELLFROM", "WELLTO", "REFERENCE", "VOLUME", "VOLUNT", "FROMPLATETYPE", "WELLFROMVOLUME", "POSFROM", "WELLTOVOLUME", "TOPLATETYPE", "MULTI", "WHAT", "LLF", "PLT", "TOWELLVOLUME", "OFFSETX", "OFFSETY", "OFFSETZ", "TIME", "SPEED"}
```

```go
var Robotinstructionnames = []string{"TFR", "CTF", "SCB", "MCB", "SCT", "MCT", "CCC", "LDT", "UDT", "RST", "CHA", "ASP", "DSP", "BLO", "PTZ", "MOV", "MRW", "LOD", "ULD", "SUK", "BLW", "SPS", "SDS", "INI", "FIN", "WAI", "LON", "LOF", "OPN", "CLS", "LAD", "UAD", "MIX"}
```

#### func  ChooseChannel

```go
func ChooseChannel(vol *wunit.Volume, prms *LHProperties) (*wtype.LHChannelParameter, string)
```

#### func  GetMultiSet

```go
func GetMultiSet(a []string, channelmulti int, fromplatemulti int, toplatemulti int) [][]int
```

#### func  GetNextSet

```go
func GetNextSet(a []string, channelmulti int, fromplatemulti int, toplatemulti int) ([]int, []string)
```

#### func  InsToString

```go
func InsToString(ins RobotInstruction) string
```

#### func  MakePolicies

```go
func MakePolicies() map[string]LHPolicy
```

#### func  MinMinVol

```go
func MinMinVol(channels []*wtype.LHChannelParameter) wunit.Volume
```

#### func  TransferVolumes

```go
func TransferVolumes(Vol, Min, Max wunit.Volume) []wunit.Volume
```

#### func  ValidateLHProperties

```go
func ValidateLHProperties(props *LHProperties) (bool, string)
```

#### type AspirateInstruction

```go
type AspirateInstruction struct {
	Type       int
	Head       int
	Volume     []*wunit.Volume
	Overstroke bool
	Multi      int
	Plt        []string
	What       []string
	LLF        []bool
}
```


#### func  NewAspirateInstruction

```go
func NewAspirateInstruction() *AspirateInstruction
```

#### func (*AspirateInstruction) Generate

```go
func (ins *AspirateInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*AspirateInstruction) GetParameter

```go
func (ins *AspirateInstruction) GetParameter(name string) interface{}
```

#### func (*AspirateInstruction) InstructionType

```go
func (ins *AspirateInstruction) InstructionType() int
```

#### func (*AspirateInstruction) OutputTo

```go
func (ins *AspirateInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type BlowInstruction

```go
type BlowInstruction struct {
	Type       int
	Head       int
	What       []string
	PltTo      []string
	WellTo     []string
	Volume     []*wunit.Volume
	TPlateType []string
	TVolume    []*wunit.Volume
	Prms       *wtype.LHChannelParameter
	Multi      int
}
```


#### func  NewBlowInstruction

```go
func NewBlowInstruction() *BlowInstruction
```

#### func (*BlowInstruction) AddTransferParams

```go
func (ins *BlowInstruction) AddTransferParams(tp TransferParams)
```

#### func (*BlowInstruction) Generate

```go
func (ins *BlowInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*BlowInstruction) GetParameter

```go
func (ins *BlowInstruction) GetParameter(name string) interface{}
```

#### func (*BlowInstruction) InstructionType

```go
func (ins *BlowInstruction) InstructionType() int
```

#### func (*BlowInstruction) Params

```go
func (scti *BlowInstruction) Params() MultiTransferParams
```

#### type BlowoutInstruction

```go
type BlowoutInstruction struct {
	Type   int
	Head   int
	Volume []*wunit.Volume
	Multi  int
	Plt    []string
	What   []string
	LLF    []bool
}
```


#### func  NewBlowoutInstruction

```go
func NewBlowoutInstruction() *BlowoutInstruction
```

#### func (*BlowoutInstruction) Generate

```go
func (ins *BlowoutInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*BlowoutInstruction) GetParameter

```go
func (ins *BlowoutInstruction) GetParameter(name string) interface{}
```

#### func (*BlowoutInstruction) InstructionType

```go
func (ins *BlowoutInstruction) InstructionType() int
```

#### func (*BlowoutInstruction) OutputTo

```go
func (ins *BlowoutInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type ChangeAdaptorInstruction

```go
type ChangeAdaptorInstruction struct {
	Type           int
	Head           int
	DropPosition   string
	GetPosition    string
	OldAdaptorType string
	NewAdaptorType string
}
```


#### func  NewChangeAdaptorInstruction

```go
func NewChangeAdaptorInstruction(head int, droppos, getpos, oldad, newad string) *ChangeAdaptorInstruction
```

#### func (*ChangeAdaptorInstruction) Generate

```go
func (ins *ChangeAdaptorInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*ChangeAdaptorInstruction) GetParameter

```go
func (ins *ChangeAdaptorInstruction) GetParameter(name string) interface{}
```

#### func (*ChangeAdaptorInstruction) InstructionType

```go
func (ins *ChangeAdaptorInstruction) InstructionType() int
```

#### type CloseInstruction

```go
type CloseInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewCloseInstruction

```go
func NewCloseInstruction() *CloseInstruction
```

#### func (*CloseInstruction) Generate

```go
func (ins *CloseInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*CloseInstruction) GetParameter

```go
func (ins *CloseInstruction) GetParameter(name string) interface{}
```

#### func (*CloseInstruction) InstructionType

```go
func (ins *CloseInstruction) InstructionType() int
```

#### func (*CloseInstruction) OutputTo

```go
func (ins *CloseInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type DispenseInstruction

```go
type DispenseInstruction struct {
	Type   int
	Head   int
	Volume []*wunit.Volume
	Multi  int
	Plt    []string
	What   []string
	LLF    []bool
}
```


#### func  NewDispenseInstruction

```go
func NewDispenseInstruction() *DispenseInstruction
```

#### func (*DispenseInstruction) Generate

```go
func (ins *DispenseInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*DispenseInstruction) GetParameter

```go
func (ins *DispenseInstruction) GetParameter(name string) interface{}
```

#### func (*DispenseInstruction) InstructionType

```go
func (ins *DispenseInstruction) InstructionType() int
```

#### func (*DispenseInstruction) OutputTo

```go
func (ins *DispenseInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type ExtendedLiquidhandlingDriver

```go
type ExtendedLiquidhandlingDriver interface {
	LiquidhandlingDriver
	SetPositionState(position string, state driver.PositionState) driver.CommandStatus
	GetCapabilities() (LHProperties, driver.CommandStatus)
	GetCurrentPosition(head int) (string, driver.CommandStatus)
	GetPositionState(position string) (string, driver.CommandStatus)
	GetHeadState(head int) (string, driver.CommandStatus)
	GetStatus() (driver.Status, driver.CommandStatus)
	UpdateMetaData(props *LHProperties) driver.CommandStatus
	UnloadHead(param int) driver.CommandStatus
	LoadHead(param int) driver.CommandStatus
	LightsOn() driver.CommandStatus
	LightsOff() driver.CommandStatus
	LoadAdaptor(param int) driver.CommandStatus
	UnloadAdaptor(param int) driver.CommandStatus
	// refactored into other interfaces?
	Open() driver.CommandStatus
	Close() driver.CommandStatus
	Message(level int, title, text string, showcancel bool) driver.CommandStatus
}
```


#### type FinalizeInstruction

```go
type FinalizeInstruction struct {
	Type int
}
```


#### func  NewFinalizeInstruction

```go
func NewFinalizeInstruction() *FinalizeInstruction
```

#### func (*FinalizeInstruction) Generate

```go
func (ins *FinalizeInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*FinalizeInstruction) GetParameter

```go
func (ins *FinalizeInstruction) GetParameter(name string) interface{}
```

#### func (*FinalizeInstruction) InstructionType

```go
func (ins *FinalizeInstruction) InstructionType() int
```

#### func (*FinalizeInstruction) OutputTo

```go
func (ins *FinalizeInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type InitializeInstruction

```go
type InitializeInstruction struct {
	Type int
}
```


#### func  NewInitializeInstruction

```go
func NewInitializeInstruction() *InitializeInstruction
```

#### func (*InitializeInstruction) Generate

```go
func (ins *InitializeInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*InitializeInstruction) GetParameter

```go
func (ins *InitializeInstruction) GetParameter(name string) interface{}
```

#### func (*InitializeInstruction) InstructionType

```go
func (ins *InitializeInstruction) InstructionType() int
```

#### func (*InitializeInstruction) OutputTo

```go
func (ins *InitializeInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type LHCategoryCondition

```go
type LHCategoryCondition struct {
	Category string
}
```


#### func (LHCategoryCondition) IsEqualTo

```go
func (lhcc LHCategoryCondition) IsEqualTo(other LHCondition) bool
```

#### func (LHCategoryCondition) Match

```go
func (lhcc LHCategoryCondition) Match(v interface{}) bool
```

#### func (LHCategoryCondition) Type

```go
func (lhcc LHCategoryCondition) Type() string
```

#### type LHCondition

```go
type LHCondition interface {
	Match(interface{}) bool
	Type() string
	IsEqualTo(LHCondition) bool
}
```


#### type LHNumericCondition

```go
type LHNumericCondition struct {
	Upper float64
	Lower float64
}
```


#### func (LHNumericCondition) IsEqualTo

```go
func (lhnc LHNumericCondition) IsEqualTo(other LHCondition) bool
```

#### func (LHNumericCondition) Match

```go
func (lhnc LHNumericCondition) Match(v interface{}) bool
```

#### func (LHNumericCondition) Type

```go
func (lhnc LHNumericCondition) Type() string
```

#### type LHPolicy

```go
type LHPolicy map[string]interface{}
```

this structure defines parameters

#### func  DupLHPolicy

```go
func DupLHPolicy(in LHPolicy) LHPolicy
```

#### func  MakeCulturePolicy

```go
func MakeCulturePolicy() LHPolicy
```

#### func  MakeDNAPolicy

```go
func MakeDNAPolicy() LHPolicy
```

#### func  MakeDefaultPolicy

```go
func MakeDefaultPolicy() LHPolicy
```

#### func  MakeFoamyPolicy

```go
func MakeFoamyPolicy() LHPolicy
```

#### func  MakeGlycerolPolicy

```go
func MakeGlycerolPolicy() LHPolicy
```

#### func  MakeJBPolicy

```go
func MakeJBPolicy() LHPolicy
```

#### func  MakeLVExtraPolicy

```go
func MakeLVExtraPolicy() LHPolicy
```

#### func  MakeSolventPolicy

```go
func MakeSolventPolicy() LHPolicy
```

#### func  MakeWaterPolicy

```go
func MakeWaterPolicy() LHPolicy
```

#### func (LHPolicy) MergeWith

```go
func (lhp LHPolicy) MergeWith(other LHPolicy) LHPolicy
```
clobber everything in here with the other policy then return the merged copy

#### type LHPolicyRule

```go
type LHPolicyRule struct {
	Name       string
	Conditions []LHVariableCondition
	Priority   int
	Type       int // AND =0 OR = 1
}
```

conditions are ANDed together there is no chaining

#### func  NewLHPolicyRule

```go
func NewLHPolicyRule(name string) LHPolicyRule
```

#### func (*LHPolicyRule) AddCategoryConditionOn

```go
func (lhpr *LHPolicyRule) AddCategoryConditionOn(variable, category string)
```

#### func (*LHPolicyRule) AddNumericConditionOn

```go
func (lhpr *LHPolicyRule) AddNumericConditionOn(variable string, low, up float64)
```

#### func (LHPolicyRule) Check

```go
func (lhpr LHPolicyRule) Check(ins RobotInstruction) bool
```

#### func (LHPolicyRule) HasCondition

```go
func (lhpr LHPolicyRule) HasCondition(cond LHVariableCondition) bool
```

#### func (LHPolicyRule) IsEqualTo

```go
func (lhpr LHPolicyRule) IsEqualTo(other LHPolicyRule) bool
```
this just looks for the same conditions, doesn't matter if the rules lead to
different outcomes... not sure if this quite gives us the right behaviour but
let's plough on for now

#### type LHPolicyRuleSet

```go
type LHPolicyRuleSet struct {
	Policies map[string]LHPolicy
	Rules    map[string]LHPolicyRule
}
```


#### func  CloneLHPolicyRuleSet

```go
func CloneLHPolicyRuleSet(parent *LHPolicyRuleSet) *LHPolicyRuleSet
```

#### func  GetLHPolicyForTest

```go
func GetLHPolicyForTest() *LHPolicyRuleSet
```

#### func  LoadLHPoliciesFrom

```go
func LoadLHPoliciesFrom(filename string) *LHPolicyRuleSet
```

#### func  LoadLHPoliciesFromFile

```go
func LoadLHPoliciesFromFile() (*LHPolicyRuleSet, error)
```

#### func  NewLHPolicyRuleSet

```go
func NewLHPolicyRuleSet() *LHPolicyRuleSet
```

#### func (*LHPolicyRuleSet) AddRule

```go
func (lhpr *LHPolicyRuleSet) AddRule(rule LHPolicyRule, consequent LHPolicy)
```

#### func (LHPolicyRuleSet) GetEquivalentRuleTo

```go
func (lhpr LHPolicyRuleSet) GetEquivalentRuleTo(rule LHPolicyRule) string
```

#### func (LHPolicyRuleSet) GetPolicyFor

```go
func (lhpr LHPolicyRuleSet) GetPolicyFor(ins RobotInstruction) LHPolicy
```

#### func (*LHPolicyRuleSet) MergeWith

```go
func (lhpr *LHPolicyRuleSet) MergeWith(other *LHPolicyRuleSet)
```

#### type LHProperties

```go
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
	Driver               LiquidhandlingDriver        `gotopb:"-"`
	CurrConf             *wtype.LHChannelParameter   // TODO: initialise
	Cnfvol               []*wtype.LHChannelParameter // TODO: initialise
	Layout               map[string]wtype.Coordinates
	MaterialType         material.MaterialType
}
```

describes a liquid handler, its capabilities and current state probably needs
splitting up to separate out the state information from the properties
information

#### func  NewLHProperties

```go
func NewLHProperties(num_positions int, model, manufacturer, lhtype, tiptype string, layout map[string]wtype.Coordinates) *LHProperties
```
constructor for the above

#### func (*LHProperties) AddPlate

```go
func (lhp *LHProperties) AddPlate(pos string, plate *wtype.LHPlate)
```

#### func (*LHProperties) AddTipBox

```go
func (lhp *LHProperties) AddTipBox(tipbox *wtype.LHTipbox)
```

#### func (*LHProperties) AddTipBoxTo

```go
func (lhp *LHProperties) AddTipBoxTo(pos string, tipbox *wtype.LHTipbox)
```

#### func (*LHProperties) AddTipWaste

```go
func (lhp *LHProperties) AddTipWaste(tipwaste *wtype.LHTipwaste)
```

#### func (*LHProperties) AddTipWasteTo

```go
func (lhp *LHProperties) AddTipWasteTo(pos string, tipwaste *wtype.LHTipwaste)
```

#### func (*LHProperties) AddWash

```go
func (lhp *LHProperties) AddWash(wash *wtype.LHPlate)
```

#### func (*LHProperties) AddWashTo

```go
func (lhp *LHProperties) AddWashTo(pos string, wash *wtype.LHPlate)
```

#### func (*LHProperties) AddWasteTo

```go
func (lhp *LHProperties) AddWasteTo(pos string, waste *wtype.LHPlate)
```

#### func (*LHProperties) DropDirtyTips

```go
func (lhp *LHProperties) DropDirtyTips(channel *wtype.LHChannelParameter, multi int) (wells, positions, boxtypes []string)
```

#### func (*LHProperties) Dup

```go
func (lhp *LHProperties) Dup() *LHProperties
```

#### func (*LHProperties) GetCleanTips

```go
func (lhp *LHProperties) GetCleanTips(tiptype string, channel *wtype.LHChannelParameter, mirror bool, multi int) (wells, positions, boxtypes []string)
```

#### func (*LHProperties) GetMaterialType

```go
func (lhp *LHProperties) GetMaterialType() material.MaterialType
```
GetMaterialType implement stockableMaterial

#### func (*LHProperties) RemoveTipBoxes

```go
func (lhp *LHProperties) RemoveTipBoxes()
```

#### type LHVariableCondition

```go
type LHVariableCondition struct {
	TestVariable string
	Condition    LHCondition
}
```


#### func  NewLHVariableCondition

```go
func NewLHVariableCondition(testvariable string) LHVariableCondition
```

#### func (LHVariableCondition) Check

```go
func (lhvc LHVariableCondition) Check(ins RobotInstruction) bool
```

#### func (LHVariableCondition) IsEqualTo

```go
func (lhvc LHVariableCondition) IsEqualTo(other LHVariableCondition) bool
```

#### func (*LHVariableCondition) SetCategoric

```go
func (lhvc *LHVariableCondition) SetCategoric(category string)
```

#### func (*LHVariableCondition) SetNumeric

```go
func (lhvc *LHVariableCondition) SetNumeric(low, up float64)
```

#### func (*LHVariableCondition) UnmarshalJSON

```go
func (lh *LHVariableCondition) UnmarshalJSON(data []byte) error
```

#### type LightsOffInstruction

```go
type LightsOffInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewLightsOffInstruction

```go
func NewLightsOffInstruction() *LightsOffInstruction
```

#### func (*LightsOffInstruction) Generate

```go
func (ins *LightsOffInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*LightsOffInstruction) GetParameter

```go
func (ins *LightsOffInstruction) GetParameter(name string) interface{}
```

#### func (*LightsOffInstruction) InstructionType

```go
func (ins *LightsOffInstruction) InstructionType() int
```

#### func (*LightsOffInstruction) OutputTo

```go
func (ins *LightsOffInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type LightsOnInstruction

```go
type LightsOnInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewLightsOnInstruction

```go
func NewLightsOnInstruction() *LightsOnInstruction
```

#### func (*LightsOnInstruction) Generate

```go
func (ins *LightsOnInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*LightsOnInstruction) GetParameter

```go
func (ins *LightsOnInstruction) GetParameter(name string) interface{}
```

#### func (*LightsOnInstruction) InstructionType

```go
func (ins *LightsOnInstruction) InstructionType() int
```

#### func (*LightsOnInstruction) OutputTo

```go
func (ins *LightsOnInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type LiquidhandlingDriver

```go
type LiquidhandlingDriver interface {
	Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus
	MoveRaw(head int, x, y, z float64) driver.CommandStatus
	Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
	LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
	SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus
	SetDriveSpeed(drive string, rate float64) driver.CommandStatus
	Stop() driver.CommandStatus
	Go() driver.CommandStatus
	Initialize() driver.CommandStatus
	Finalize() driver.CommandStatus
	Wait(time float64) driver.CommandStatus
	Mix(head int, volume []float64, platetype []string, cycles []int, multi int, what []string, blowout []bool) driver.CommandStatus
	ResetPistons(head, channel int) driver.CommandStatus
	AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus
	RemoveAllPlates() driver.CommandStatus
	RemovePlateAt(position string) driver.CommandStatus
}
```


#### type LoadAdaptorInstruction

```go
type LoadAdaptorInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewLoadAdaptorInstruction

```go
func NewLoadAdaptorInstruction() *LoadAdaptorInstruction
```

#### func (*LoadAdaptorInstruction) Generate

```go
func (ins *LoadAdaptorInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*LoadAdaptorInstruction) GetParameter

```go
func (ins *LoadAdaptorInstruction) GetParameter(name string) interface{}
```

#### func (*LoadAdaptorInstruction) InstructionType

```go
func (ins *LoadAdaptorInstruction) InstructionType() int
```

#### func (*LoadAdaptorInstruction) OutputTo

```go
func (ins *LoadAdaptorInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type LoadTipsInstruction

```go
type LoadTipsInstruction struct {
	Type       int
	Head       int
	Pos        []string
	Well       []string
	Channels   []int
	TipType    []string
	HolderType []string
	Multi      int
}
```


#### func  NewLoadTipsInstruction

```go
func NewLoadTipsInstruction() *LoadTipsInstruction
```

#### func (*LoadTipsInstruction) Generate

```go
func (ins *LoadTipsInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*LoadTipsInstruction) GetParameter

```go
func (ins *LoadTipsInstruction) GetParameter(name string) interface{}
```

#### func (*LoadTipsInstruction) InstructionType

```go
func (ins *LoadTipsInstruction) InstructionType() int
```

#### func (*LoadTipsInstruction) OutputTo

```go
func (ins *LoadTipsInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type LoadTipsMoveInstruction

```go
type LoadTipsMoveInstruction struct {
	Type       int
	Head       int
	Well       []string
	FPosition  []string
	FPlateType []string
	Multi      int
}
```


#### func  NewLoadTipsMoveInstruction

```go
func NewLoadTipsMoveInstruction() *LoadTipsMoveInstruction
```

#### func (*LoadTipsMoveInstruction) Generate

```go
func (ins *LoadTipsMoveInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*LoadTipsMoveInstruction) GetParameter

```go
func (ins *LoadTipsMoveInstruction) GetParameter(name string) interface{}
```

#### func (*LoadTipsMoveInstruction) InstructionType

```go
func (ins *LoadTipsMoveInstruction) InstructionType() int
```

#### type MixInstruction

```go
type MixInstruction struct {
	Type      int
	Head      int
	Volume    []*wunit.Volume
	PlateType []string
	What      []string
	Blowout   []bool
	Multi     int
	Cycles    []int
}
```


#### func  NewMixInstruction

```go
func NewMixInstruction() *MixInstruction
```

#### func (*MixInstruction) Generate

```go
func (ins *MixInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MixInstruction) GetParameter

```go
func (ins *MixInstruction) GetParameter(name string) interface{}
```

#### func (*MixInstruction) InstructionType

```go
func (mi *MixInstruction) InstructionType() int
```

#### func (*MixInstruction) OutputTo

```go
func (mi *MixInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type MoveInstruction

```go
type MoveInstruction struct {
	Type      int
	Head      int
	Pos       []string
	Plt       []string
	Well      []string
	WVolume   []*wunit.Volume
	Reference []int
	OffsetX   []float64
	OffsetY   []float64
	OffsetZ   []float64
}
```


#### func  NewMoveInstruction

```go
func NewMoveInstruction() *MoveInstruction
```

#### func (*MoveInstruction) Generate

```go
func (ins *MoveInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MoveInstruction) GetParameter

```go
func (ins *MoveInstruction) GetParameter(name string) interface{}
```

#### func (*MoveInstruction) InstructionType

```go
func (ins *MoveInstruction) InstructionType() int
```

#### func (*MoveInstruction) OutputTo

```go
func (ins *MoveInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type MoveMixInstruction

```go
type MoveMixInstruction struct {
	Type      int
	Head      int
	Plt       []string
	Well      []string
	Volume    []*wunit.Volume
	PlateType []string
	FVolume   []*wunit.Volume
	Cycles    []int
	What      []string
	Blowout   []bool
	Multi     int
	Prms      map[string]interface{}
}
```


#### func  NewMoveMixInstruction

```go
func NewMoveMixInstruction() *MoveMixInstruction
```

#### func (*MoveMixInstruction) Generate

```go
func (ins *MoveMixInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MoveMixInstruction) GetParameter

```go
func (ins *MoveMixInstruction) GetParameter(name string) interface{}
```

#### func (*MoveMixInstruction) InstructionType

```go
func (ins *MoveMixInstruction) InstructionType() int
```

#### type MoveRawInstruction

```go
type MoveRawInstruction struct {
	Type       int
	Head       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewMoveRawInstruction

```go
func NewMoveRawInstruction() *MoveRawInstruction
```

#### func (*MoveRawInstruction) Generate

```go
func (ins *MoveRawInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MoveRawInstruction) GetParameter

```go
func (ins *MoveRawInstruction) GetParameter(name string) interface{}
```

#### func (*MoveRawInstruction) InstructionType

```go
func (ins *MoveRawInstruction) InstructionType() int
```

#### func (*MoveRawInstruction) OutputTo

```go
func (ins *MoveRawInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type MultiChannelBlockInstruction

```go
type MultiChannelBlockInstruction struct {
	Type       int
	What       [][]string
	PltFrom    [][]string
	PltTo      [][]string
	WellFrom   [][]string
	WellTo     [][]string
	Volume     [][]*wunit.Volume
	FPlateType [][]string
	TPlateType [][]string
	FVolume    [][]*wunit.Volume
	TVolume    [][]*wunit.Volume
	Multi      int
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewMultiChannelBlockInstruction

```go
func NewMultiChannelBlockInstruction() *MultiChannelBlockInstruction
```

#### func (*MultiChannelBlockInstruction) AddTransferParams

```go
func (ins *MultiChannelBlockInstruction) AddTransferParams(mct MultiTransferParams)
```

#### func (*MultiChannelBlockInstruction) Generate

```go
func (ins *MultiChannelBlockInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MultiChannelBlockInstruction) GetParameter

```go
func (ins *MultiChannelBlockInstruction) GetParameter(name string) interface{}
```

#### func (*MultiChannelBlockInstruction) InstructionType

```go
func (ins *MultiChannelBlockInstruction) InstructionType() int
```

#### type MultiChannelTransferInstruction

```go
type MultiChannelTransferInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
	Multi      int
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewMultiChannelTransferInstruction

```go
func NewMultiChannelTransferInstruction() *MultiChannelTransferInstruction
```

#### func (*MultiChannelTransferInstruction) Generate

```go
func (ins *MultiChannelTransferInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*MultiChannelTransferInstruction) GetParameter

```go
func (ins *MultiChannelTransferInstruction) GetParameter(name string) interface{}
```

#### func (*MultiChannelTransferInstruction) InstructionType

```go
func (ins *MultiChannelTransferInstruction) InstructionType() int
```

#### func (*MultiChannelTransferInstruction) Params

```go
func (scti *MultiChannelTransferInstruction) Params(k int) TransferParams
```

#### type MultiTransferParams

```go
type MultiTransferParams struct {
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
	Channel    *wtype.LHChannelParameter
}
```


#### func  NewMultiTransferParams

```go
func NewMultiTransferParams(multi int) MultiTransferParams
```

#### type OpenInstruction

```go
type OpenInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewOpenInstruction

```go
func NewOpenInstruction() *OpenInstruction
```

#### func (*OpenInstruction) Generate

```go
func (ins *OpenInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*OpenInstruction) GetParameter

```go
func (ins *OpenInstruction) GetParameter(name string) interface{}
```

#### func (*OpenInstruction) InstructionType

```go
func (ins *OpenInstruction) InstructionType() int
```

#### func (*OpenInstruction) OutputTo

```go
func (ins *OpenInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type PTZInstruction

```go
type PTZInstruction struct {
	Type    int
	Head    int
	Channel int
}
```


#### func  NewPTZInstruction

```go
func NewPTZInstruction() *PTZInstruction
```

#### func (*PTZInstruction) Generate

```go
func (ins *PTZInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*PTZInstruction) GetParameter

```go
func (ins *PTZInstruction) GetParameter(name string) interface{}
```

#### func (*PTZInstruction) InstructionType

```go
func (ins *PTZInstruction) InstructionType() int
```

#### func (*PTZInstruction) OutputTo

```go
func (ins *PTZInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type ResetInstruction

```go
type ResetInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewResetInstruction

```go
func NewResetInstruction() *ResetInstruction
```

#### func (*ResetInstruction) AddMultiTransferParams

```go
func (ins *ResetInstruction) AddMultiTransferParams(mtp MultiTransferParams)
```

#### func (*ResetInstruction) AddTransferParams

```go
func (ins *ResetInstruction) AddTransferParams(tp TransferParams)
```

#### func (*ResetInstruction) Generate

```go
func (ins *ResetInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*ResetInstruction) GetParameter

```go
func (ins *ResetInstruction) GetParameter(name string) interface{}
```

#### func (*ResetInstruction) InstructionType

```go
func (ins *ResetInstruction) InstructionType() int
```

#### type RobotInstruction

```go
type RobotInstruction interface {
	InstructionType() int
	GetParameter(name string) interface{}
	Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
}
```


#### func  DropTips

```go
func DropTips(tiptype string, params *LHProperties, channel *wtype.LHChannelParameter, multi int) RobotInstruction
```

#### func  GetTips

```go
func GetTips(tiptype string, params *LHProperties, channel *wtype.LHChannelParameter, multi int, mirror bool) RobotInstruction
```

#### type RobotInstructionSet

```go
type RobotInstructionSet struct {
}
```


#### func  NewRobotInstructionSet

```go
func NewRobotInstructionSet(p RobotInstruction) *RobotInstructionSet
```

#### func (*RobotInstructionSet) Add

```go
func (ri *RobotInstructionSet) Add(ins RobotInstruction)
```

#### func (*RobotInstructionSet) Generate

```go
func (ri *RobotInstructionSet) Generate(lhpr *LHPolicyRuleSet, lhpm *LHProperties) []RobotInstruction
```

#### func (*RobotInstructionSet) ToString

```go
func (ri *RobotInstructionSet) ToString(level int) string
```

#### type SetDriveSpeedInstruction

```go
type SetDriveSpeedInstruction struct {
	Type  int
	Drive string
	Speed float64
}
```


#### func  NewSetDriveSpeedInstruction

```go
func NewSetDriveSpeedInstruction() *SetDriveSpeedInstruction
```

#### func (*SetDriveSpeedInstruction) Generate

```go
func (ins *SetDriveSpeedInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*SetDriveSpeedInstruction) GetParameter

```go
func (ins *SetDriveSpeedInstruction) GetParameter(name string) interface{}
```

#### func (*SetDriveSpeedInstruction) InstructionType

```go
func (ins *SetDriveSpeedInstruction) InstructionType() int
```

#### func (*SetDriveSpeedInstruction) OutputTo

```go
func (ins *SetDriveSpeedInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type SetPipetteSpeedInstruction

```go
type SetPipetteSpeedInstruction struct {
	Type    int
	Head    int
	Channel int
	Speed   float64
}
```


#### func  NewSetPipetteSpeedInstruction

```go
func NewSetPipetteSpeedInstruction() *SetPipetteSpeedInstruction
```

#### func (*SetPipetteSpeedInstruction) Generate

```go
func (ins *SetPipetteSpeedInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*SetPipetteSpeedInstruction) GetParameter

```go
func (ins *SetPipetteSpeedInstruction) GetParameter(name string) interface{}
```

#### func (*SetPipetteSpeedInstruction) InstructionType

```go
func (ins *SetPipetteSpeedInstruction) InstructionType() int
```

#### func (*SetPipetteSpeedInstruction) OutputTo

```go
func (ins *SetPipetteSpeedInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type SingleChannelBlockInstruction

```go
type SingleChannelBlockInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewSingleChannelBlockInstruction

```go
func NewSingleChannelBlockInstruction() *SingleChannelBlockInstruction
```

#### func (*SingleChannelBlockInstruction) AddTransferParams

```go
func (ins *SingleChannelBlockInstruction) AddTransferParams(mct TransferParams)
```

#### func (*SingleChannelBlockInstruction) Generate

```go
func (ins *SingleChannelBlockInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*SingleChannelBlockInstruction) GetParameter

```go
func (ins *SingleChannelBlockInstruction) GetParameter(name string) interface{}
```

#### func (*SingleChannelBlockInstruction) InstructionType

```go
func (ins *SingleChannelBlockInstruction) InstructionType() int
```

#### type SingleChannelTransferInstruction

```go
type SingleChannelTransferInstruction struct {
	Type       int
	What       string
	PltFrom    string
	PltTo      string
	WellFrom   string
	WellTo     string
	Volume     *wunit.Volume
	FPlateType string
	TPlateType string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewSingleChannelTransferInstruction

```go
func NewSingleChannelTransferInstruction() *SingleChannelTransferInstruction
```

#### func (*SingleChannelTransferInstruction) Generate

```go
func (ins *SingleChannelTransferInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*SingleChannelTransferInstruction) GetParameter

```go
func (ins *SingleChannelTransferInstruction) GetParameter(name string) interface{}
```

#### func (*SingleChannelTransferInstruction) InstructionType

```go
func (ins *SingleChannelTransferInstruction) InstructionType() int
```

#### func (*SingleChannelTransferInstruction) Params

```go
func (scti *SingleChannelTransferInstruction) Params() TransferParams
```

#### type StateChangeInstruction

```go
type StateChangeInstruction struct {
	Type     int
	OldState *wtype.LHChannelParameter
	NewState *wtype.LHChannelParameter
}
```


#### func  NewStateChangeInstruction

```go
func NewStateChangeInstruction(oldstate, newstate *wtype.LHChannelParameter) *StateChangeInstruction
```

#### func (*StateChangeInstruction) Generate

```go
func (ins *StateChangeInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*StateChangeInstruction) GetParameter

```go
func (ins *StateChangeInstruction) GetParameter(name string) interface{}
```

#### func (*StateChangeInstruction) InstructionType

```go
func (ins *StateChangeInstruction) InstructionType() int
```

#### type SuckInstruction

```go
type SuckInstruction struct {
	Type       int
	Head       int
	What       []string
	PltFrom    []string
	WellFrom   []string
	Volume     []*wunit.Volume
	FPlateType []string
	FVolume    []*wunit.Volume
	Prms       *wtype.LHChannelParameter
	Multi      int
	Overstroke bool
}
```


#### func  NewSuckInstruction

```go
func NewSuckInstruction() *SuckInstruction
```

#### func (*SuckInstruction) AddTransferParams

```go
func (ins *SuckInstruction) AddTransferParams(tp TransferParams)
```

#### func (*SuckInstruction) Generate

```go
func (ins *SuckInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*SuckInstruction) GetParameter

```go
func (ins *SuckInstruction) GetParameter(name string) interface{}
```

#### func (*SuckInstruction) InstructionType

```go
func (ins *SuckInstruction) InstructionType() int
```

#### type TerminalRobotInstruction

```go
type TerminalRobotInstruction interface {
	RobotInstruction
	OutputTo(driver LiquidhandlingDriver)
}
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
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FPlateWX   []int
	FPlateWY   []int
	TPlateWX   []int
	TPlateWY   []int
	FVolume    []*wunit.Volume
	TVolume    []*wunit.Volume
}
```


#### func  NewTransferInstruction

```go
func NewTransferInstruction(what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []*wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int) *TransferInstruction
```

#### func (*TransferInstruction) Generate

```go
func (ins *TransferInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*TransferInstruction) GetParallelSetsFor

```go
func (ins *TransferInstruction) GetParallelSetsFor(channel *wtype.LHChannelParameter) [][]int
```

#### func (*TransferInstruction) GetParameter

```go
func (ins *TransferInstruction) GetParameter(name string) interface{}
```

#### func (*TransferInstruction) InstructionType

```go
func (ins *TransferInstruction) InstructionType() int
```

#### func (*TransferInstruction) ParamSet

```go
func (ti *TransferInstruction) ParamSet(n int) TransferParams
```

#### func (*TransferInstruction) ToString

```go
func (ti *TransferInstruction) ToString() string
```

#### type TransferParams

```go
type TransferParams struct {
	What       string
	PltFrom    string
	PltTo      string
	WellFrom   string
	WellTo     string
	Volume     *wunit.Volume
	FPlateType string
	TPlateType string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Channel    *wtype.LHChannelParameter
}
```


#### func (TransferParams) ToString

```go
func (tp TransferParams) ToString() string
```

#### type UnloadAdaptorInstruction

```go
type UnloadAdaptorInstruction struct {
	Type       int
	What       []string
	PltFrom    []string
	PltTo      []string
	WellFrom   []string
	WellTo     []string
	Volume     []*wunit.Volume
	FPlateType []string
	TPlateType []string
	FVolume    *wunit.Volume
	TVolume    *wunit.Volume
	Prms       *wtype.LHChannelParameter
}
```


#### func  NewUnloadAdaptorInstruction

```go
func NewUnloadAdaptorInstruction() *UnloadAdaptorInstruction
```

#### func (*UnloadAdaptorInstruction) Generate

```go
func (ins *UnloadAdaptorInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*UnloadAdaptorInstruction) GetParameter

```go
func (ins *UnloadAdaptorInstruction) GetParameter(name string) interface{}
```

#### func (*UnloadAdaptorInstruction) InstructionType

```go
func (ins *UnloadAdaptorInstruction) InstructionType() int
```

#### func (*UnloadAdaptorInstruction) OutputTo

```go
func (ins *UnloadAdaptorInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type UnloadTipsInstruction

```go
type UnloadTipsInstruction struct {
	Type       int
	Head       int
	Channels   []int
	TipType    []string
	HolderType []string
	Multi      int
	Pos        []string
	Well       []string
}
```


#### func  NewUnloadTipsInstruction

```go
func NewUnloadTipsInstruction() *UnloadTipsInstruction
```

#### func (*UnloadTipsInstruction) Generate

```go
func (ins *UnloadTipsInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*UnloadTipsInstruction) GetParameter

```go
func (ins *UnloadTipsInstruction) GetParameter(name string) interface{}
```

#### func (*UnloadTipsInstruction) InstructionType

```go
func (ins *UnloadTipsInstruction) InstructionType() int
```

#### func (*UnloadTipsInstruction) OutputTo

```go
func (ins *UnloadTipsInstruction) OutputTo(driver LiquidhandlingDriver)
```

#### type UnloadTipsMoveInstruction

```go
type UnloadTipsMoveInstruction struct {
	Type       int
	Head       int
	PltTo      []string
	WellTo     []string
	TPlateType []string
	Multi      int
}
```


#### func  NewUnloadTipsMoveInstruction

```go
func NewUnloadTipsMoveInstruction() *UnloadTipsMoveInstruction
```

#### func (*UnloadTipsMoveInstruction) Generate

```go
func (ins *UnloadTipsMoveInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*UnloadTipsMoveInstruction) GetParameter

```go
func (ins *UnloadTipsMoveInstruction) GetParameter(name string) interface{}
```

#### func (*UnloadTipsMoveInstruction) InstructionType

```go
func (ins *UnloadTipsMoveInstruction) InstructionType() int
```

#### type VolumeSet

```go
type VolumeSet struct {
	Vols []*wunit.Volume
}
```


#### func  NewVolumeSet

```go
func NewVolumeSet(n int) VolumeSet
```

#### func (VolumeSet) Add

```go
func (vs VolumeSet) Add(v *wunit.Volume)
```

#### func (VolumeSet) GetACopy

```go
func (vs VolumeSet) GetACopy() []*wunit.Volume
```

#### func (VolumeSet) MaxMultiTransferVolume

```go
func (vs VolumeSet) MaxMultiTransferVolume() *wunit.Volume
```

#### func (VolumeSet) SetEqualTo

```go
func (vs VolumeSet) SetEqualTo(v *wunit.Volume)
```

#### func (VolumeSet) Sub

```go
func (vs VolumeSet) Sub(v *wunit.Volume) []*wunit.Volume
```

#### type WaitInstruction

```go
type WaitInstruction struct {
	Type int
	Time float64
}
```


#### func  NewWaitInstruction

```go
func NewWaitInstruction() *WaitInstruction
```

#### func (*WaitInstruction) Generate

```go
func (ins *WaitInstruction) Generate(policy *LHPolicyRuleSet, prms *LHProperties) []RobotInstruction
```

#### func (*WaitInstruction) GetParameter

```go
func (ins *WaitInstruction) GetParameter(name string) interface{}
```

#### func (*WaitInstruction) InstructionType

```go
func (ins *WaitInstruction) InstructionType() int
```

#### func (*WaitInstruction) OutputTo

```go
func (ins *WaitInstruction) OutputTo(driver LiquidhandlingDriver)
```
