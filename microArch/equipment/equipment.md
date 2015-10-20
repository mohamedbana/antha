---
layout: default
type: api
navgroup: api
shortname: microArch/equipment
title: microArch/equipment
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: microArch/equipment
---
# equipment
--
    import "."


## Usage

#### type ActionDescription

```go
type ActionDescription struct {
	Action     action.Action
	ActionData string //TODO probably make a struct with proper fields
	Params     map[string]string
}
```

ActionDescription is the representation of an action that is going to be
executed and extends the particular Action with data on how to carry it out

#### func  NewActionDescription

```go
func NewActionDescription(action action.Action, data string, params map[string]string) *ActionDescription
```
NewActionDescription instantiates a new action with the given data

#### type Behaviour

```go
type Behaviour struct {
	Action      action.Action
	Constraints string //TODO probably make a struct with proper fields
}
```

Behaviour represents the capabilities of a piece of Equipment to perform an
action. This is related to how you can ask an equipment to carry out an action.
The actionData in ActionDescription should meet the constraints that a piece of
equipment describes in its behaviour

#### func  NewBehaviour

```go
func NewBehaviour(action action.Action, constraints string) *Behaviour
```
NewBehaviour will instantiate a new Behaviour matching the given action and
constraints

#### func (*Behaviour) Matches

```go
func (b *Behaviour) Matches(ac ActionDescription) bool
```
Matches checks whether the action description can be carried out by this
behaviour

#### type Equipment

```go
type Equipment interface {
	GetID() string
	//GetEquipmentDefinition returns a description of the equipment device in terms of
	// operations it can handle, restrictions, configuration options ...
	GetEquipmentDefinition()
	//Perform an action in the equipment. Actions might be transmitted in blocks to the equipment
	// The grouping of the actions (as a set, plate or whatever) is not performed at the equipment driver level
	// or is it?
	Do(actionDescription ActionDescription) error
	//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
	// of the piece of equipment.
	Can(ac ActionDescription) bool
	//Status should give a description of the current execution status and any future actions queued to the device
	Status() string
	//Init driver will be initialized when registered
	Init() error
	//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
	Shutdown() error
}
```

Equipment is something capable of performing different actions under different
restrictions and explaining what its capabilities and the restrictions on them
are
