---
layout: default
type: api
navgroup: docs
shortname: driver/movement
title: driver/movement
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: driver/movement
---
# movement
--
    import "."


## Usage

#### type MovementDriver

```go
type MovementDriver interface {
	//Initialize is a generic operation to initialize all the necessary data a device might need
	Initialize() driver.CommandStatus
	//Finalize is a generic operation to finalize the device driver
	Finalize() driver.CommandStatus
	//Move changes the location from one place to another
	Move(entity wtype.Entity, final wtype.Location) driver.CommandStatus
	//Stop is a generic operation to express the abortion of the movement
	Stop() driver.CommandStatus
	//Wait is a generic operation that will let a certain time pass before continuing with the execution
	Wait(time float64) driver.CommandStatus
}
```
