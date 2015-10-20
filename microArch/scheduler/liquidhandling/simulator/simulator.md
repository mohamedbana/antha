---
layout: default
type: api
navgroup: docs
shortname: liquidhandling/simulator
title: liquidhandling/simulator
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: liquidhandling/simulator
---
# simpleliquidhandler
--
    import "."


## Usage

```go
var (
	DEFAULT_SLEEP_TIME = 1 * time.Second
)
```

#### type LHCommandStatus

```go
type LHCommandStatus struct {
}
```


#### type LHPlate

```go
type LHPlate interface{}
```


#### type LHPositionState

```go
type LHPositionState interface{}
```


#### type LHProperties

```go
type LHProperties interface{}
```


#### type LHStatus

```go
type LHStatus interface{}
```


#### type SimpleLiquidHandler

```go
type SimpleLiquidHandler struct {
}
```


#### func (*SimpleLiquidHandler) Aspirate

```go
func (lh *SimpleLiquidHandler) Aspirate(volume float64, overstroke bool, head int, multi int) LHCommandStatus
```

#### func (*SimpleLiquidHandler) Dispense

```go
func (lh *SimpleLiquidHandler) Dispense(volume float64, blowout bool, head int, multi int) LHCommandStatus
```

#### func (*SimpleLiquidHandler) Finalize

```go
func (lh *SimpleLiquidHandler) Finalize() LHCommandStatus
```

#### func (*SimpleLiquidHandler) GetCapabilities

```go
func (lh *SimpleLiquidHandler) GetCapabilities() LHProperties
```

#### func (*SimpleLiquidHandler) GetCurrentPosition

```go
func (lh *SimpleLiquidHandler) GetCurrentPosition(head int) (string, LHCommandStatus)
```

#### func (*SimpleLiquidHandler) GetHeadState

```go
func (lh *SimpleLiquidHandler) GetHeadState(head int) (string, LHCommandStatus)
```

#### func (*SimpleLiquidHandler) GetPositionState

```go
func (lh *SimpleLiquidHandler) GetPositionState(position string) (string, LHCommandStatus)
```

#### func (*SimpleLiquidHandler) GetStatus

```go
func (lh *SimpleLiquidHandler) GetStatus() (LHStatus, LHCommandStatus)
```

#### func (*SimpleLiquidHandler) Go

```go
func (lh *SimpleLiquidHandler) Go() LHCommandStatus
```

#### func (*SimpleLiquidHandler) Initialize

```go
func (lh *SimpleLiquidHandler) Initialize() LHCommandStatus
```

#### func (*SimpleLiquidHandler) LoadTips

```go
func (lh *SimpleLiquidHandler) LoadTips(head, multi int) LHCommandStatus
```

#### func (*SimpleLiquidHandler) Move

```go
func (lh *SimpleLiquidHandler) Move(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type, head int) LHCommandStatus
```

#### func (*SimpleLiquidHandler) MoveExplicit

```go
func (lh *SimpleLiquidHandler) MoveExplicit(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type *LHPlate, head int) LHCommandStatus
```

#### func (*SimpleLiquidHandler) MoveRaw

```go
func (lh *SimpleLiquidHandler) MoveRaw(x, y, z float64) LHCommandStatus
```

#### func (*SimpleLiquidHandler) SetDriveSpeed

```go
func (lh *SimpleLiquidHandler) SetDriveSpeed(drive string, rate float64) LHCommandStatus
```

#### func (*SimpleLiquidHandler) SetPipetteSpeed

```go
func (lh *SimpleLiquidHandler) SetPipetteSpeed(rate float64)
```

#### func (*SimpleLiquidHandler) SetPositionState

```go
func (lh *SimpleLiquidHandler) SetPositionState(position string, state LHPositionState) LHCommandStatus
```

#### func (*SimpleLiquidHandler) Stop

```go
func (lh *SimpleLiquidHandler) Stop() LHCommandStatus
```

#### func (*SimpleLiquidHandler) UnloadTips

```go
func (lh *SimpleLiquidHandler) UnloadTips(head, multi int) LHCommandStatus
```
