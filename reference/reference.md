---
layout: default
type: api
navgroup: docs
shortname: antha/reference
title: antha/reference
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: antha/reference
---
# reference
--
    import "."


## Usage

#### type AsyncBag

```go
type AsyncBag struct {
}
```

Simple structure to coordinate the asynchronous aggregation of multiple values
that have to be fired together

#### func (*AsyncBag) AddValue

```go
func (a *AsyncBag) AddValue(key string, value interface{}) bool
```
adds value and returns true if the bag was fired TODO: Should the competion be
wrapped in a sync.Once in case there are duplicate params flowing through the
network with the same threadID?

#### type AsyncCompleter

```go
type AsyncCompleter interface {
	Complete(interface{})
}
```

support function to fire when a full bag of values has arrived

#### type AsyncMapper

```go
type AsyncMapper interface {
	Map(map[string]interface{}) interface{}
}
```

support function to map into a concrete struct

#### type Example

```go
type Example struct {
	flow.Component                    // component "superclass" embedded
	Color          <-chan ThreadParam // color to make this well
	SleepTime      <-chan ThreadParam // amount of time to randomly wait
	WellColor      chan<- ThreadParam // output color
}
```

channel interfaces with threadID grouped types

#### func  NewExample

```go
func NewExample() *Example
```

#### func (*Example) Complete

```go
func (e *Example) Complete(params interface{})
```
AsyncBag functions

#### func (*Example) Map

```go
func (e *Example) Map(m map[string]interface{}) interface{}
```
could handle mapping in the threadID better...

#### func (*Example) OnColor

```go
func (e *Example) OnColor(param ThreadParam)
```

#### func (*Example) OnSleepTime

```go
func (e *Example) OnSleepTime(param ThreadParam)
```

#### type ThreadID

```go
type ThreadID string
```


#### type ThreadParam

```go
type ThreadParam struct {
	Value interface{}
	ID    ThreadID
}
```
