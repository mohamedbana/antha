---
layout: default
type: api
navgroup: docs
shortname: anthalib/wunit
title: anthalib/wunit
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: anthalib/wunit
---
# wunit
--
    import "."


## Usage

```go
const (
	M = -30 + (iota * 3)
	G = -30 + (iota * 3)
	T = -30 + (iota * 3)
	P = -30 + (iota * 3)
	E = -30 + (iota * 3)
	Z = -30 + (iota * 3)
	Y = -30 + (iota * 3)
)
```

```go
const (
	TopNode = iota
	UnitPlusPrefixNode
	LeafNode
)
```

#### func  GetPrefixLib

```go
func GetPrefixLib(fn string) (*(map[string]SIPrefix), error)
```
deserialize JSON prefix library

#### func  GetUnitLib

```go
func GetUnitLib(fn string) (*(map[string]GenericUnit), error)
```
deserialize JSON unit library

#### func  MakePrefices

```go
func MakePrefices() map[string]SIPrefix
```
make the prefix structure

#### func  Make_units

```go
func Make_units() map[string]GenericUnit
```
generate an initial unit library

#### func  NewWFloat

```go
func NewWFloat(v float64) wfloat
```
wrap a float in the wvalue structure

#### func  NewWInt

```go
func NewWInt(v int) wint
```
wrap an int in the wvalue structure

#### func  NewWString

```go
func NewWString(v string) wstring
```

#### func  PrefixDiv

```go
func PrefixDiv(x string, y string) string
```
divide one prefix by another take care: there are no checks for going out of
bounds e.g. Z/z will give an error!

#### func  PrefixMul

```go
func PrefixMul(x string, y string) string
```
multiply two prefix values take care: there are no checks for going out of
bounds e.g. Z*Z will generate an error!

#### func  ReverseLookupPrefix

```go
func ReverseLookupPrefix(i int) string
```
helper function for reverse lookup of prefix

#### func  RoundInt

```go
func RoundInt(v float64) int
```

#### type Amount

```go
type Amount struct {
	ConcreteMeasurement
}
```

mole

#### func  NewAmount

```go
func NewAmount(v float64, unit string) Amount
```
generate a new Amount in moles

#### func (*Amount) Quantity

```go
func (a *Amount) Quantity() Measurement
```
defines Amount to be a SubstanceQuantity

#### type Angle

```go
type Angle struct {
	ConcreteMeasurement
}
```

angle

#### func  NewAngle

```go
func NewAngle(v float64, unit string) Angle
```
generate a new angle unit

#### type Area

```go
type Area struct {
	ConcreteMeasurement
}
```

area

#### func  NewArea

```go
func NewArea(v float64, unit string) Area
```
make an area unit

#### type BaseUnit

```go
type BaseUnit interface {
	// unit name
	Name() string
	// unit symbol
	Symbol() string
	// multiply by this to get SI value
	// nb this should be a function since we actually need
	// an affine transformation
	BaseSIConversionFactor() float64 // this can be calculated in many cases
	// if we convert to the SI units what is the appropriate unit symbol
	BaseSIUnit() string // if we use the above, what unit do we get?
	// print this
	ToString() string
}
```

structure defining a base unit

#### type Concentration

```go
type Concentration struct {
	ConcreteMeasurement
}
```

defines a concentration unit

#### func  NewConcentration

```go
func NewConcentration(v float64, unit string) Concentration
```
make a new concentration in SI units... either M/l or kg/l

#### func (*Concentration) MarshalJSON

```go
func (m *Concentration) MarshalJSON() ([]byte, error)
```

#### func (*Concentration) UnmarshalJSON

```go
func (m *Concentration) UnmarshalJSON(b []byte) error
```

#### type ConcreteMeasurement

```go
type ConcreteMeasurement struct {
	// the raw value
	Mvalue float64
	// the relevant units
	Munit *GenericPrefixedUnit
}
```

structure implementing the Measurement interface

#### func  NewMeasurement

```go
func NewMeasurement(v float64, prefix string, unit string) ConcreteMeasurement
```
helper function for creating a new measurement

#### func  NewPMeasurement

```go
func NewPMeasurement(v float64, pu string) ConcreteMeasurement
```

#### func (*ConcreteMeasurement) Add

```go
func (cm *ConcreteMeasurement) Add(m Measurement)
```

#### func (*ConcreteMeasurement) ConvertTo

```go
func (cm *ConcreteMeasurement) ConvertTo(p PrefixedUnit) float64
```
convert to a different unit nb this is NOT destructive

#### func (*ConcreteMeasurement) ConvertToString

```go
func (cm *ConcreteMeasurement) ConvertToString(s string) float64
```

#### func (*ConcreteMeasurement) EqualTo

```go
func (cm *ConcreteMeasurement) EqualTo(m Measurement) bool
```

#### func (*ConcreteMeasurement) EqualToFloat

```go
func (cm *ConcreteMeasurement) EqualToFloat(f float64) bool
```

#### func (*ConcreteMeasurement) GreaterThan

```go
func (cm *ConcreteMeasurement) GreaterThan(m Measurement) bool
```

#### func (*ConcreteMeasurement) GreaterThanFloat

```go
func (cm *ConcreteMeasurement) GreaterThanFloat(f float64) bool
```

#### func (*ConcreteMeasurement) LessThan

```go
func (cm *ConcreteMeasurement) LessThan(m Measurement) bool
```

#### func (*ConcreteMeasurement) LessThanFloat

```go
func (cm *ConcreteMeasurement) LessThanFloat(f float64) bool
```

#### func (*ConcreteMeasurement) RawValue

```go
func (cm *ConcreteMeasurement) RawValue() float64
```
value without conversion

#### func (*ConcreteMeasurement) SIValue

```go
func (cm *ConcreteMeasurement) SIValue() float64
```
value when converted to SI units

#### func (*ConcreteMeasurement) SetValue

```go
func (cm *ConcreteMeasurement) SetValue(v float64) float64
```
set the value of this measurement

#### func (*ConcreteMeasurement) Subtract

```go
func (cm *ConcreteMeasurement) Subtract(m Measurement)
```

#### func (*ConcreteMeasurement) ToString

```go
func (cm *ConcreteMeasurement) ToString() string
```

#### func (*ConcreteMeasurement) Unit

```go
func (cm *ConcreteMeasurement) Unit() PrefixedUnit
```
get unit with prefix

#### type Density

```go
type Density struct {
	ConcreteMeasurement
}
```

a structure which defines a density

#### func  NewDensity

```go
func NewDensity(v float64, unit string) Density
```
make a new density structure in SI units

#### type Energy

```go
type Energy struct {
	ConcreteMeasurement
}
```

this is really Mass(Length/Time)^2

#### func  NewEnergy

```go
func NewEnergy(v float64, unit string) Energy
```
make a new energy unit

#### type FlowRate

```go
type FlowRate struct {
	ConcreteMeasurement
}
```


#### func  NewFlowRate

```go
func NewFlowRate(v float64, unit string) FlowRate
```

#### type Force

```go
type Force struct {
	ConcreteMeasurement
}
```

a Force

#### func  NewForce

```go
func NewForce(v float64, unit string) Force
```
a new force in Newtons

#### type GenericPrefixedUnit

```go
type GenericPrefixedUnit struct {
	GenericUnit
	SPrefix SIPrefix
}
```

the generic prefixed unit structure

#### func  NewPrefixedUnit

```go
func NewPrefixedUnit(prefix string, unit string) *GenericPrefixedUnit
```
helper function to make it easier to make a new unit with prefix directly

#### func  ParsePrefixedUnit

```go
func ParsePrefixedUnit(unit string) *GenericPrefixedUnit
```

#### func (*GenericPrefixedUnit) BaseSIConversionFactor

```go
func (gpu *GenericPrefixedUnit) BaseSIConversionFactor() float64
```
multiplier to convert to SI base unit... for composites this is the ratio of the
base units for the dimensions in question e.g. kg/l for concentration

#### func (*GenericPrefixedUnit) BaseSISymbol

```go
func (gpu *GenericPrefixedUnit) BaseSISymbol() string
```
symbol for unit after conversion to base si unit

#### func (*GenericPrefixedUnit) ConvertTo

```go
func (gpu *GenericPrefixedUnit) ConvertTo(p2 PrefixedUnit) float64
```
gives the conversion factor from one prefixed unit to another

#### func (*GenericPrefixedUnit) Prefix

```go
func (gpu *GenericPrefixedUnit) Prefix() SIPrefix
```

#### func (*GenericPrefixedUnit) PrefixedSymbol

```go
func (gpu *GenericPrefixedUnit) PrefixedSymbol() string
```
symbol with prefix

#### func (*GenericPrefixedUnit) RawSymbol

```go
func (gpu *GenericPrefixedUnit) RawSymbol() string
```
symbol without prefix

#### func (*GenericPrefixedUnit) Symbol

```go
func (gpu *GenericPrefixedUnit) Symbol() string
```
symbol with prefix

#### type GenericUnit

```go
type GenericUnit struct {
	StrName             string
	StrSymbol           string
	FltConversionfactor float64
	StrBaseUnit         string
}
```

structure for defining a generic unit

#### func  UnitBySymbol

```go
func UnitBySymbol(sym string) GenericUnit
```
look up unit by symbol

#### func (*GenericUnit) BaseSIConversionFactor

```go
func (gu *GenericUnit) BaseSIConversionFactor() float64
```

#### func (*GenericUnit) BaseSIUnit

```go
func (gu *GenericUnit) BaseSIUnit() string
```

#### func (*GenericUnit) Name

```go
func (gu *GenericUnit) Name() string
```

#### func (*GenericUnit) Symbol

```go
func (gu *GenericUnit) Symbol() string
```

#### func (*GenericUnit) ToString

```go
func (gu *GenericUnit) ToString() string
```

#### type Length

```go
type Length struct {
	ConcreteMeasurement
}
```

length

#### func  EZLength

```go
func EZLength(v float64) Length
```

#### func  NewLength

```go
func NewLength(v float64, unit string) Length
```
make a length

#### type Mass

```go
type Mass struct {
	ConcreteMeasurement
}
```

mass

#### func  NewMass

```go
func NewMass(v float64, unit string) Mass
```

#### func (*Mass) Quantity

```go
func (m *Mass) Quantity() Measurement
```
func NewMass(v float64, unit string) Mass {

    if unit != "kg" && unit != "g" {
    	panic("Can't make masses which aren't in grammes or kilograms")
    }

    var t Mass

    if unit == "kg" {
    	t = Mass{NewMeasurement(v, "k", "g")}
    } else {
    	t = Mass{NewMeasurement(v, "", "g")}
    }
    return t

}

defines mass to be a SubstanceQuantity

#### type Measurement

```go
type Measurement interface {
	// the value in base SI units
	SIValue() float64
	// the value in the current units
	RawValue() float64
	// unit plus prefix
	Unit() PrefixedUnit
	// set the value, this must be thread-safe
	// returns old value
	SetValue(v float64) float64
	// convert units
	ConvertTo(p PrefixedUnit) float64
	// wrapper for above
	ConvertToString(s string) float64
	// add to this measurement
	Add(m Measurement)
	// subtract from this measurement
	Subtract(m Measurement)
	// comparison operators
	LessThan(m Measurement) bool
	GreaterThan(m Measurement) bool
	EqualTo(m Measurement) bool

	// A nice string representation
	ToString() string
}
```

fundamental representation of a value in the system

#### type MeasurementLimits

```go
type MeasurementLimits struct {
	Limits map[string]NDManifold
}
```

for holding minima, maxima etc. e.g. ml.Limits["max"] =>
{"length"->{4.0,{"Metres", "m", 1.0}}}...

#### type NDManifold

```go
type NDManifold struct {
	Dimensions map[string]Measurement
}
```

N-dimensional manifold structure

#### type NodeType

```go
type NodeType uint8
```


#### type PNode

```go
type PNode struct {
	Name     string
	Type     NodeType
	Up       *PNode
	Children []*PNode
	Value    interface{}
}
```


#### func  NewNode

```go
func NewNode(name string, typ NodeType, cap uint8) *PNode
```

#### func (*PNode) AddChild

```go
func (pn *PNode) AddChild(child *PNode)
```

#### type Prefix_Net

```go
type Prefix_Net int
```


#### type PrefixedUnit

```go
type PrefixedUnit interface {
	BaseUnit
	// the prefix of the unit
	Prefix() SIPrefix
	// the symbol including prefix
	PrefixedSymbol() string
	// the symbol excluding prefix
	RawSymbol() string
	// appropriate unit if we ask for SI values
	BaseSISymbol() string
	// returns conversion factor from *this* unit to the other
	ConvertTo(pu PrefixedUnit) float64
}
```

a unit with an SI prefix

#### type Pressure

```go
type Pressure struct {
	ConcreteMeasurement
}
```

a Pressure structure

#### func  NewPressure

```go
func NewPressure(v float64, unit string) Pressure
```
make a new pressure in Pascals

#### type SIPrefix

```go
type SIPrefix struct {
	// prefix name
	Name string
	// meaning in base 10
	Value float64
}
```

structure defining an SI prefix

#### func  SIPrefixBySymbol

```go
func SIPrefixBySymbol(symbol string) SIPrefix
```
helper function to allow lookup of prefix

#### type SIPrefixedUnit

```go
type SIPrefixedUnit struct {
	Stack   []*PNode
	Text    string
	TreeTop *PNode
	CurNode *PNode
}
```


#### func (*SIPrefixedUnit) AddNodeToCurrent

```go
func (p *SIPrefixedUnit) AddNodeToCurrent(node *PNode)
```

#### func (*SIPrefixedUnit) AddNodeToStack

```go
func (p *SIPrefixedUnit) AddNodeToStack(node *PNode)
```

#### func (*SIPrefixedUnit) AddStackToNode

```go
func (p *SIPrefixedUnit) AddStackToNode(node *PNode)
```

#### func (*SIPrefixedUnit) AddUnit

```go
func (p *SIPrefixedUnit) AddUnit(s string)
```

#### func (*SIPrefixedUnit) AddUnitPlusPrefixNode

```go
func (p *SIPrefixedUnit) AddUnitPlusPrefixNode()
```

#### func (*SIPrefixedUnit) AddUnitPrefix

```go
func (p *SIPrefixedUnit) AddUnitPrefix(s string)
```

#### func (*SIPrefixedUnit) Init

```go
func (p *SIPrefixedUnit) Init(text []byte)
```

#### func (*SIPrefixedUnit) NewStack

```go
func (p *SIPrefixedUnit) NewStack()
```

#### func (*SIPrefixedUnit) PopStack

```go
func (p *SIPrefixedUnit) PopStack() (node *PNode)
```

#### func (*SIPrefixedUnit) PopStackAndAddTo

```go
func (p *SIPrefixedUnit) PopStackAndAddTo(node *PNode)
```

#### type SIPrefixedUnitGrammar

```go
type SIPrefixedUnitGrammar struct {
	SIPrefixedUnit

	Buffer string

	Parse func(rule ...int) error
	Reset func()
}
```


#### func (*SIPrefixedUnitGrammar) Execute

```go
func (p *SIPrefixedUnitGrammar) Execute()
```

#### func (*SIPrefixedUnitGrammar) Highlighter

```go
func (p *SIPrefixedUnitGrammar) Highlighter()
```

#### func (*SIPrefixedUnitGrammar) Init

```go
func (p *SIPrefixedUnitGrammar) Init()
```

#### func (*SIPrefixedUnitGrammar) PrintSyntaxTree

```go
func (p *SIPrefixedUnitGrammar) PrintSyntaxTree()
```

#### type SpecificHeatCapacity

```go
type SpecificHeatCapacity struct {
	ConcreteMeasurement
}
```

a structure which defines a specific heat capacity

#### func  NewSpecificHeatCapacity

```go
func NewSpecificHeatCapacity(v float64, unit string) SpecificHeatCapacity
```
make a new specific heat capacity structure in SI units

#### type SubstanceQuantity

```go
type SubstanceQuantity interface {
	Quantity() Measurement
}
```

mass or mole

#### type Temperature

```go
type Temperature struct {
	ConcreteMeasurement
}
```

temperature

#### func  NewTemperature

```go
func NewTemperature(v float64, unit string) Temperature
```
make a temperature

#### func (*Temperature) MarshalJSON

```go
func (m *Temperature) MarshalJSON() ([]byte, error)
```

#### func (*Temperature) UnmarshalJSON

```go
func (m *Temperature) UnmarshalJSON(b []byte) error
```

#### type Time

```go
type Time struct {
	ConcreteMeasurement
}
```

time

#### func  NewTime

```go
func NewTime(v float64, unit string) Time
```
make a time unit

#### func (*Time) MarshalJSON

```go
func (m *Time) MarshalJSON() ([]byte, error)
```

#### func (*Time) UnmarshalJSON

```go
func (m *Time) UnmarshalJSON(b []byte) error
```

#### type Volume

```go
type Volume struct {
	ConcreteMeasurement
}
```

volume -- strictly speaking of course this is length^3

#### func  CopyVolume

```go
func CopyVolume(v *Volume) *Volume
```

#### func  NewVolume

```go
func NewVolume(v float64, unit string) Volume
```
make a volume

#### func (*Volume) MarshalJSON

```go
func (m *Volume) MarshalJSON() ([]byte, error)
```

#### func (*Volume) UnmarshalJSON

```go
func (m *Volume) UnmarshalJSON(b []byte) error
```
