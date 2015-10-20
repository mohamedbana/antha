---
layout: default
type: api
navgroup: api
shortname: equipment/manual
title: equipment/manual
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: equipment/manual
---
# manual
--
    import "."


## Usage

#### type AnthaManual

```go
type AnthaManual struct {
	ID         string
	Behaviours []equipment.Behaviour
	Cui        cli.CUI
}
```


#### func  NewAnthaManual

```go
func NewAnthaManual(id string) *AnthaManual
```

#### func (AnthaManual) Can

```go
func (e AnthaManual) Can(b equipment.ActionDescription) bool
```

#### func (AnthaManual) Do

```go
func (e AnthaManual) Do(actionDescription equipment.ActionDescription) error
```

#### func (AnthaManual) GetEquipmentDefinition

```go
func (e AnthaManual) GetEquipmentDefinition()
```

#### func (AnthaManual) GetID

```go
func (e AnthaManual) GetID() string
```

#### func (AnthaManual) Init

```go
func (e AnthaManual) Init() error
```

#### func (AnthaManual) Shutdown

```go
func (e AnthaManual) Shutdown() error
```

#### func (AnthaManual) Status

```go
func (e AnthaManual) Status() string
```
