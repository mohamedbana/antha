---
layout: default
type: api
navgroup: docs
shortname: examples/reference
title: examples/reference
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: examples/reference
---
# reference
--
    import "."


## Usage

#### type Example

```go
type Example struct {
	flow.Component                            // component "superclass" embedded
	Color          <-chan execute.ThreadParam // color to make this well
	SleepTime      <-chan execute.ThreadParam // amount of time to randomly wait
	WellColor      chan<- execute.ThreadParam // output color
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
execute.AsyncBag functions

#### func (*Example) Map

```go
func (e *Example) Map(m map[string]interface{}) interface{}
```
could handle mapping in the threadID better...

#### func (*Example) OnColor

```go
func (e *Example) OnColor(param execute.ThreadParam)
```

#### func (*Example) OnSleepTime

```go
func (e *Example) OnSleepTime(param execute.ThreadParam)
```

#### type JSONBlock

```go
type JSONBlock struct {
	Color     *string
	SleepTime *time.Duration
	WellColor *string
	ID        *execute.ThreadID
}
```


#### type ParamBlock

```go
type ParamBlock struct {
	Color     string
	SleepTime time.Duration
	WellColor string
	ID        execute.ThreadID
}
```

single execution thread variables with concrete types

#### func  ParamsFromJSON

```go
func ParamsFromJSON(r io.Reader) (p *ParamBlock)
```
helper generator function

#### func (*ParamBlock) ToJSON

```go
func (p *ParamBlock) ToJSON() (b bytes.Buffer)
```
support function for wire format

#### type WellColorParam

```go
type WellColorParam struct {
	WellColor string
}
```
