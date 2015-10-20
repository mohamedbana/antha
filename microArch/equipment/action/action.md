---
layout: default
type: api
navgroup: api
shortname: equipment/action
title: equipment/action
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: equipment/action
---
# action
--
    import "."


## Usage

```go
const (
	NONE = iota
	MESSAGE
	LH_SETUP
	LH_MOVE
	LH_MOVE_EXPLICIT
	LH_MOVE_RAW
	LH_ASPIRATE
	LH_DISPENSE
	LH_LOAD_TIPS
	LH_UNLOAD_TIPS
	LH_SET_PIPPETE_SPEED
	LH_SET_DRIVE_SPEED
	LH_STOP
	LH_SET_POSITION_STATE
	LH_RESET_PISTONS
	LH_WAIT
	LH_MIX
	LH_ADD_PLATE
	LH_REMOVE_PLATE
	LH_REMOVE_ALL_PLATES //?? maybe not necessary
	MLH_CHANGE_TIPS
	IN_INCUBATE
	IN_INCUBATE_SHAKE
)
```

#### type Action

```go
type Action int
```

Action describes a particular function that an equipment can perform. It is the
same concept as an interface, but since it is a representation of the real
world, we cannot model it with an actual one

#### func (Action) String

```go
func (a Action) String() string
```
