---
layout: default
type: api
navgroup: docs
shortname: incubating/manual
title: incubating/manual
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: incubating/manual
---
# manual
--
    import "."


## Usage

#### type ManualDriver

```go
type ManualDriver struct {
}
```


#### func  NewManualDriver

```go
func NewManualDriver() *ManualDriver
```
NewManualDriver returns a new instance of a manual driver pointing to the right
piece of equipment

#### func (*ManualDriver) Finalize

```go
func (m *ManualDriver) Finalize() driver.CommandStatus
```

#### func (*ManualDriver) Incubate

```go
func (m *ManualDriver) Incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) driver.CommandStatus
```

#### func (*ManualDriver) Initialize

```go
func (m *ManualDriver) Initialize() driver.CommandStatus
```

#### func (*ManualDriver) Move

```go
func (m *ManualDriver) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus
```
