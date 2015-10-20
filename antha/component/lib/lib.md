---
layout: default
type: api
navgroup: docs
shortname: component/lib
title: component/lib
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: component/lib
---
# lib
--
    import "."


## Usage

#### func  GetComponents

```go
func GetComponents() []ComponentDesc
```

#### type ComponentDesc

```go
type ComponentDesc struct {
	Name        string
	Constructor func() interface{}
}
```
