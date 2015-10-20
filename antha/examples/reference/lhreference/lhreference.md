---
layout: default
type: api
navgroup: docs
shortname: reference/lhreference
title: reference/lhreference
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: reference/lhreference
---
# lhreference
--
    import "."


## Usage

#### func  AddFeature

```go
func AddFeature(name string, param execute.ThreadParam, mapper execute.AsyncMapper, completer execute.AsyncCompleter, blocks *map[execute.ThreadID]*execute.AsyncBag, nkeys int, lock sync.Mutex)
```
candidate for refactoring out into execute

#### type InputBlock

```go
type InputBlock struct {
	A       *wtype.LHComponent
	B       *wtype.LHComponent
	Dest    *wtype.LHWell
	BlockID execute.ThreadID
	ID      execute.ThreadID
}
```


#### func  InputsFromJSON

```go
func InputsFromJSON(r io.Reader) (i *InputBlock)
```

#### func (*InputBlock) Map

```go
func (i *InputBlock) Map(m map[string]interface{}) interface{}
```

#### func (*InputBlock) ToJSON

```go
func (i *InputBlock) ToJSON(b bytes.Buffer)
```

#### type JSONBlock

```go
type JSONBlock struct {
	A_vol   *wunit.Volume
	B_vol   *wunit.Volume
	A       *wtype.LHComponent
	B       *wtype.LHComponent
	Dest    *wtype.LHWell
	BlockID *execute.ThreadID
	ID      *execute.ThreadID
}
```

JSON blocks are also required... not quite sure why though I'm sure we can
serialize the paramblock OK anyway

#### type LHReference

```go
type LHReference struct {
	flow.Component

	// these are data items
	A_vol <-chan execute.ThreadParam
	B_vol <-chan execute.ThreadParam

	A    <-chan execute.ThreadParam
	B    <-chan execute.ThreadParam
	Dest <-chan execute.ThreadParam

	Mixture chan<- execute.ThreadParam

	ParamBlocks map[execute.ThreadID]*execute.AsyncBag
	InputBlocks map[execute.ThreadID]*execute.AsyncBag
	PIBlocks    map[execute.ThreadID]*execute.AsyncBag
}
```


#### func  NewLHReference

```go
func NewLHReference() *LHReference
```

#### func (*LHReference) Complete

```go
func (lh *LHReference) Complete(val interface{})
```

#### func (*LHReference) OnA

```go
func (lh *LHReference) OnA(param execute.ThreadParam)
```

#### func (*LHReference) OnA_vol

```go
func (lh *LHReference) OnA_vol(param execute.ThreadParam)
```
ports for wiring into the network

#### func (*LHReference) OnB

```go
func (lh *LHReference) OnB(param execute.ThreadParam)
```

#### func (*LHReference) OnB_vol

```go
func (lh *LHReference) OnB_vol(param execute.ThreadParam)
```

#### func (*LHReference) OnDest

```go
func (lh *LHReference) OnDest(param execute.ThreadParam)
```

#### func (*LHReference) Setup

```go
func (lh *LHReference) Setup(v interface{})
```

#### func (*LHReference) Steps

```go
func (lh *LHReference) Steps(v interface{})
```

#### type OutputBlock

```go
type OutputBlock struct {
	// interestingly, Dest here comes out as part of SolOut
	SolOut  *wtype.LHSolution
	BlockID execute.ThreadID
	ID      execute.ThreadID
}
```


#### type PIBlock

```go
type PIBlock struct {
	Params  *ParamBlock
	Inputs  *InputBlock
	BlockID execute.ThreadID
	ID      execute.ThreadID
}
```


#### func (*PIBlock) Map

```go
func (pi *PIBlock) Map(m map[string]interface{}) interface{}
```

#### type ParamBlock

```go
type ParamBlock struct {
	A_vol   wunit.Volume
	B_vol   wunit.Volume
	BlockID execute.ThreadID
	ID      execute.ThreadID
}
```


#### func  ParamsFromJSON

```go
func ParamsFromJSON(r io.Reader) (p *ParamBlock)
```

#### func (*ParamBlock) Map

```go
func (p *ParamBlock) Map(m map[string]interface{}) interface{}
```

#### func (*ParamBlock) ToJSON

```go
func (p *ParamBlock) ToJSON() (b bytes.Buffer)
```
