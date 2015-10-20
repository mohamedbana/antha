---
layout: default
type: api
navgroup: docs
shortname: antha/component
title: antha/component
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: antha/component
---
# component
--
    import "."


## Usage

#### func  Run

```go
func Run(js []byte, dec *json.Decoder, enc *json.Encoder, errChan chan<- error) <-chan int
```
Runs the component graph with json input and producing json output

#### type Component

```go
type Component interface {
	ComponentInfo() *execute.ComponentInfo

	Complete(interface{})

	// Returns a new ComponentConfig
	NewConfig() interface{}

	// Returns a new ComponentParamBlock
	NewParamBlock() interface{}

	// Takes a map of ThreadParams which contain execute.JSONValues or concrete values
	// and returns a ComponentParamBlock
	Map(map[string]interface{}) interface{}
}
```
