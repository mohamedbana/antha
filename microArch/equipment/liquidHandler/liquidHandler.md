---
layout: default
type: api
navgroup: api
shortname: equipment/liquidHandler
title: equipment/liquidHandler
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: equipment/liquidHandler
---
# liquidHandler
--
    import "."


## Usage

#### type AnthaLiquidHandler

```go
type AnthaLiquidHandler struct {
	ID         string
	Behaviours []equipment.Behaviour
}
```


#### func  NewAnthaLiquidHandler

```go
func NewAnthaLiquidHandler(id string) *AnthaLiquidHandler
```

#### func (AnthaLiquidHandler) Can

```go
func (e AnthaLiquidHandler) Can(b equipment.ActionDescription) bool
```

#### func (AnthaLiquidHandler) Do

```go
func (e AnthaLiquidHandler) Do(actionDescription equipment.ActionDescription) error
```

#### func (AnthaLiquidHandler) GetEquipmentDefinition

```go
func (e AnthaLiquidHandler) GetEquipmentDefinition()
```

#### func (AnthaLiquidHandler) GetID

```go
func (e AnthaLiquidHandler) GetID() string
```

#### func (AnthaLiquidHandler) Status

```go
func (e AnthaLiquidHandler) Status() string
```
