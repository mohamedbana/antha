---
layout: default
type: api
navgroup: docs
shortname: anthalib/wconstraints
title: anthalib/wconstraints
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: anthalib/wconstraints
---
# wconstraints
--
    import "."


## Usage

#### type Constraint

```go
type Constraint interface {
	// constraints return true when violated
	Check(i interface{}) bool
}
```

Interface to define a constraint -- the main requirement is that the constraint
can be checked

#### type Constraints

```go
type Constraints interface {
	Add(c Constraint)
	Remove(c Constraint)
}
```

general constraint handling interface

#### type TempConstraint

```go
type TempConstraint struct {
	Min wunit.Temperature
	Max wunit.Temperature
}
```

Temperature constraint -- sample must be kept within the specified limits
defines a max and a min temperature

#### func (TempConstraint) Check

```go
func (tc TempConstraint) Check(i interface{}) bool
```
Test the constraint against the current temperature

#### type TimeConstraint

```go
type TimeConstraint struct {
	Start  time.Time
	Length time.Duration
}
```

time constraint - simply keeps track of a duration by recording when the
constraint was created and how long it is

#### func (*TimeConstraint) Check

```go
func (tc *TimeConstraint) Check(i interface{}) bool
```
test whether time is up
