---
layout: default
type: api
navgroup: api
shortname: microArch/equipmentManager
title: microArch/equipmentManager
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: microArch/equipmentManager
---
# equipmentManager
--
    import "."


## Usage

#### func  SetEquipmentManager

```go
func SetEquipmentManager(e *EquipmentManager)
```

#### type AnthaEquipmentManager

```go
type AnthaEquipmentManager struct {
	//ID the uuid that identifies this piece of equipment
	ID string
}
```

AnthaEquipmentManager implements the EquipmentManager interface using an
EquipmentREgistry as the means to keep track of the different equipment that is
registered in the system at a certain time.

#### func  NewAnthaEquipmentManager

```go
func NewAnthaEquipmentManager(id string) *AnthaEquipmentManager
```
NewAnthaScheduler returns a new AnthaEquipmentManager identified by id

#### func (*AnthaEquipmentManager) CheckAvailability

```go
func (e *AnthaEquipmentManager) CheckAvailability(equipmentQuery string)
```

#### func (*AnthaEquipmentManager) GetActionCandidate

```go
func (e *AnthaEquipmentManager) GetActionCandidate(actionQuery equipment.ActionDescription) *equipment.Equipment
```

#### func (*AnthaEquipmentManager) GetEquipment

```go
func (e *AnthaEquipmentManager) GetEquipment(equipmentQuery string) []equipment.Equipment
```

#### func (*AnthaEquipmentManager) GetID

```go
func (e *AnthaEquipmentManager) GetID() string
```
GetID returns the string id representing a manager, usually a uuid

#### func (*AnthaEquipmentManager) RegisterEquipment

```go
func (e *AnthaEquipmentManager) RegisterEquipment(equipment *equipment.Equipment) error
```
@implements EquipmentManager

#### func (*AnthaEquipmentManager) Shutdown

```go
func (e *AnthaEquipmentManager) Shutdown() error
```

#### type EquipmentManager

```go
type EquipmentManager interface {
	//get a list of available equipment based on specific equipment constraints
	GetEquipment(equipmentQuery string) []equipment.Equipment
	//get a list of available equipment based on an action we want to perform
	// the difference with the previous resides in where our logic for this translation resides
	GetActionCandidate(actionQuery equipment.ActionDescription) *equipment.Equipment
	//get an equipment availability, for scheduling purposes
	CheckAvailability(equipmentQuery string) //including timings??

	//RegisterEquipment is called by an equipment service to notify the equipmentManager of its existence
	RegisterEquipment(equipment *equipment.Equipment) error

	//Shutdown performs the shutdown operation on all registered devices that are still available
	Shutdown() error
}
```

EquipmentManager represents the set of operations that an equipment manager must
fulfil. It is the equipmentManager duty to ensure that all pieces of equipment
have implemented the necessary initializations

#### func  GetEquipmentManager

```go
func GetEquipmentManager() *EquipmentManager
```

#### type EquipmentRegistry

```go
type EquipmentRegistry interface {
	//RegisterEquipment inserts the given piece of equipment into this registry
	RegisterEquipment(eq *equipment.Equipment) error
	//GetEquipmentByID gets us a piece of equipment identified by a certain ID
	GetEquipmentByID(id string) *equipment.Equipment
	//ListEquipment wil return a list of all the pieces of equipment this registry knows of
	ListEquipment() map[string]equipment.Equipment
}
```

EquipmentRegistry Stores Equipment information related to a particular execution
environment

#### type MemoryEquipmentRegistry

```go
type MemoryEquipmentRegistry struct {
	EquipmentList map[string]equipment.Equipment
}
```

EquipmentRegistry in memory implementation of an equipmentRegistry

#### func  NewMemoryEquipmentRegistry

```go
func NewMemoryEquipmentRegistry() *MemoryEquipmentRegistry
```
NewMemoryEquipmentRegistry instantiates a new memory registry

#### func (*MemoryEquipmentRegistry) GetEquipmentByID

```go
func (reg *MemoryEquipmentRegistry) GetEquipmentByID(id string) *equipment.Equipment
```

#### func (*MemoryEquipmentRegistry) ListEquipment

```go
func (reg *MemoryEquipmentRegistry) ListEquipment() map[string]equipment.Equipment
```

#### func (*MemoryEquipmentRegistry) RegisterEquipment

```go
func (reg *MemoryEquipmentRegistry) RegisterEquipment(eq *equipment.Equipment) error
```
