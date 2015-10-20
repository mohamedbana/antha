---
layout: default
type: api
navgroup: docs
shortname: driver/manual
title: driver/manual
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: driver/manual
---
# manual
--
    import "."


## Usage

#### type ManualDriver

```go
type ManualDriver interface {
	Init()
	Message(message string) driver.CommandStatus
}
```
