---
layout: default
type: api
navgroup: docs
shortname: antha/execute
title: antha/execute
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: antha/execute
---
# execute
--
    import "."

support package with wrapper classes for marshalling parameters into elements

## Usage

#### type AnthaElement

```go
type AnthaElement interface {
	// contains filtered or unexported methods
}
```

type to allow generic access to Antha Elements

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
adds value and returns true if the bag was fired TODO: Should the completion be
wrapped in a sync.Once in case there are duplicate params flowing through the
network with the same threadID?

#### func (*AsyncBag) Init

```go
func (a *AsyncBag) Init(keys int, completer AsyncCompleter, mapper AsyncMapper)
```
makes a new AsyncBag which requires keys to fire f

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

#### type BlockConfig

```go
type BlockConfig struct {
	BlockID ThreadID
	Threads map[string]string
}
```


#### type ComponentInfo

```go
type ComponentInfo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Subgraph    bool       `json:"subgraph"`
	InPorts     []PortInfo `json:"inPorts"`
	OutPorts    []PortInfo `json:"outPorts"`
}
```

ComponentInfo describes a protocol as a fbp component

#### func  NewComponentInfo

```go
func NewComponentInfo(pkgName string, description string, icon string, subgraph bool, inPorts []PortInfo, outPorts []PortInfo) *ComponentInfo
```
NewComponentInfo returns a new ComponentInfo initalized with the given
information

#### type GraphConnection

```go
type GraphConnection struct {
	Data interface{} `json:",omitempty"`
	Src  struct {
		Process string
		Port    string
	} `json:",omitempty"`
	Tgt struct {
		Process string
		Port    string
	}
	Metadata struct {
		Buffer int `json:",omitempty"`
	} `json:",omitempty"`
}
```

GraphConnection describes a connection between two processes inside a
GraphDescription

#### type GraphDescription

```go
type GraphDescription struct {
	Properties struct {
		Name string
	}
	Processes map[string]struct {
		Component string
		Metadata  struct {
			Sync                 bool   `json:",omitempty"`
			PoolSize             int64  `json:",omitempty"`
			NameSpaceClass       string `json:",omitempty"`
			ComponentLibraryName string `json:",omitempty"`
		} `json:",omitempty"`
	}
	Connections []GraphConnection
	Exports     []struct {
		Private string
		Public  string
	}
	InPorts map[string]struct {
		Process string
		Port    string
	}
	OutPorts map[string]struct {
		Process string
		Port    string
	}
}
```

GraphDescription describes a fbp graph.

#### type JSONBlock

```go
type JSONBlock struct {
	ID     *ThreadID
	Error  *bool
	Values map[string]interface{}
}
```

JSONBlock holds information from a JSON string in a key value fashion except for
ID and Error keys can be used to unmarshal unknown structs and decide types on
the fly

#### func (*JSONBlock) MarshalJSON

```go
func (j *JSONBlock) MarshalJSON() ([]byte, error)
```
MarshalJSON builds a json string containing a JSONBlock structure in which ID
and Error are added explicitly and the rest of fields are added from a key/value
pair map

#### func (*JSONBlock) UnmarshalJSON

```go
func (j *JSONBlock) UnmarshalJSON(in []byte) error
```
UnmarshalJson JSONBlock from a json string saving ID and Error and the rest of
information as a key value pair

#### type JSONValue

```go
type JSONValue struct {
	Name       string
	JSONString string
}
```

JSONValue holds information for a pair key value inside a JSONBlock

#### type PortInfo

```go
type PortInfo struct {
	Id          string        `json:"id"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Addressable bool          `json:"addressable"` // ignored
	Required    bool          `json:"required"`
	Values      []interface{} `json:"values"`  // ignored
	Default     interface{}   `json:"default"` // ignored
}
```

PortInfo describes a port from a ComponentInfo

#### func  NewPortInfo

```go
func NewPortInfo(id string, portInfoType string, description string, addressable bool, required bool, values []interface{}, defaultValue interface{}) *PortInfo
```
NewPortInfo returns a new PortInfo struct initialized with the given values

#### type ThreadID

```go
type ThreadID string
```


#### type ThreadParam

```go
type ThreadParam struct {
	Value   interface{}
	ID      ThreadID
	BlockID ThreadID
	Error   bool
}
```
