---
layout: default
type: api
navgroup: docs
shortname: driver/incubating
title: driver/incubating
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: driver/incubating
---
# incubating
--
    import "."


## Usage

#### type IncubatingDriver

```go
type IncubatingDriver interface {
	Initialize() driver.CommandStatus
	Finalize() driver.CommandStatus
	Incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) driver.CommandStatus
}
```
