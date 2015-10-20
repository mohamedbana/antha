---
layout: default
type: api
navgroup: docs
shortname: anthalib/driver
title: anthalib/driver
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: anthalib/driver
---
# driver
--
    import "."


## Usage

```go
const (
	ERR int = iota
	OK
	NIM
)
```

#### type CommandStatus

```go
type CommandStatus struct {
	OK        bool
	Errorcode int
	Msg       string
}
```


#### type PositionState

```go
type PositionState map[string]interface{}
```


#### type Status

```go
type Status map[string]interface{}
```
