---
layout: default
type: api
navgroup: docs
shortname: driver/platereading
title: driver/platereading
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: driver/platereading
---
# platereading
--
    import "."


## Usage

#### type PlateReadingDriver

```go
type PlateReadingDriver interface {
	Initialize() driver.CommandStatus
	Finalize() driver.CommandStatus
	ReadPlate(matter, reading *interface{}) driver.CommandStatus
}
```
