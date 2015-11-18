---
layout: default
type: api
navgroup: docs
shortname: liquidhandling/manual
title: liquidhandling/manual
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: liquidhandling/manual
---
# manual
--
    import "github.com/antha-lang/antha/microArch/driver/liquidhandling/manual"


## Usage

#### type Aggregator

```go
type Aggregator struct {
	//Calls is the list of actionDescriptions
	Calls []equipment.ActionDescription
}
```

Aggregator has a list of actions in which it tries to stablish relations to make
them more human friendly

#### func  NewAggregator

```go
func NewAggregator() *Aggregator
```
NewAggregator instantiates a new aggregator

#### type ManualDriver

```go
type ManualDriver struct {
}
```

ManualDriver represents a piece of equipment to substitute a human performing
the different actions possible in the antha ecosystem. I can substitute a piece
of equipment or help with the usage of a non-connected device

#### func  NewManualDriver

```go
func NewManualDriver() *ManualDriver
```
NewManualDriver returns a new instance of a manual driver pointing to the right
piece of equipment

#### func (*ManualDriver) AddPlateTo

```go
func (m *ManualDriver) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus
```

#### func (*ManualDriver) Aspirate

```go
func (m *ManualDriver) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
```

#### func (*ManualDriver) Dispense

```go
func (m *ManualDriver) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus
```

#### func (*ManualDriver) Finalize

```go
func (m *ManualDriver) Finalize() driver.CommandStatus
```

#### func (*ManualDriver) GetCapabilities

```go
func (m *ManualDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus)
```
TODO implement changing capabilities for manual driver

#### func (*ManualDriver) GetCurrentPosition

```go
func (m *ManualDriver) GetCurrentPosition(head int) (string, driver.CommandStatus)
```

#### func (*ManualDriver) GetHeadState

```go
func (m *ManualDriver) GetHeadState(head int) (string, driver.CommandStatus)
```

#### func (*ManualDriver) GetPositionState

```go
func (m *ManualDriver) GetPositionState(position string) (string, driver.CommandStatus)
```

#### func (*ManualDriver) GetStatus

```go
func (m *ManualDriver) GetStatus() (driver.Status, driver.CommandStatus)
```

#### func (*ManualDriver) Go

```go
func (m *ManualDriver) Go() driver.CommandStatus
```

#### func (*ManualDriver) Initialize

```go
func (m *ManualDriver) Initialize() driver.CommandStatus
```

#### func (*ManualDriver) LoadTips

```go
func (m *ManualDriver) LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
```

#### func (*ManualDriver) Mix

```go
func (m *ManualDriver) Mix(head int, volume []float64, fvolume []float64, platetype []string, cycles []int, multi int, prms map[string]interface{}) driver.CommandStatus
```

#### func (*ManualDriver) Move

```go
func (m *ManualDriver) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus
```

#### func (*ManualDriver) MoveExplicit

```go
func (m *ManualDriver) MoveExplicit(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []*wtype.LHPlate, head int) driver.CommandStatus
```

#### func (*ManualDriver) MoveRaw

```go
func (m *ManualDriver) MoveRaw(head int, x, y, z float64) driver.CommandStatus
```

#### func (*ManualDriver) RemoveAllPlates

```go
func (m *ManualDriver) RemoveAllPlates() driver.CommandStatus
```

#### func (*ManualDriver) RemovePlateAt

```go
func (m *ManualDriver) RemovePlateAt(position string) driver.CommandStatus
```

#### func (*ManualDriver) ResetPistons

```go
func (m *ManualDriver) ResetPistons(head, channel int) driver.CommandStatus
```

#### func (*ManualDriver) SetDriveSpeed

```go
func (m *ManualDriver) SetDriveSpeed(drive string, rate float64) driver.CommandStatus
```

#### func (*ManualDriver) SetPipetteSpeed

```go
func (m *ManualDriver) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus
```

#### func (*ManualDriver) SetPositionState

```go
func (m *ManualDriver) SetPositionState(position string, state driver.PositionState) driver.CommandStatus
```

#### func (*ManualDriver) Stop

```go
func (m *ManualDriver) Stop() driver.CommandStatus
```

#### func (*ManualDriver) UnloadTips

```go
func (m *ManualDriver) UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus
```

#### func (*ManualDriver) UpdateMetaData

```go
func (m *ManualDriver) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus
```

#### func (*ManualDriver) Wait

```go
func (m *ManualDriver) Wait(time float64) driver.CommandStatus
```

#### type TranslateDictionary

```go
type TranslateDictionary struct {
}
```

TranslateDictionary keeps a dictionary of numerated items of the type prefix and
gives names and holds a lookup table the names are `prefix num count`

#### func  NewTranslateDictionary

```go
func NewTranslateDictionary(prefix string) *TranslateDictionary
```
NewTranslateDictionary creates a new dictionary with the given prefix
