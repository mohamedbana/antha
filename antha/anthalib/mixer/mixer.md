---
layout: default
type: api
navgroup: docs
shortname: anthalib/mixer
title: anthalib/mixer
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: anthalib/mixer
---
# mixer
--
    import "."


## Usage

#### func  Mix

```go
func Mix(components ...*wtype.LHComponent) *wtype.LHSolution
```
mix the specified wtype.LHComponents together and leave the destination TBD

#### func  MixInto

```go
func MixInto(destination *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution
```
mix the specified wtype.LHComponents together into the destination specified as
the first argument

#### func  Sample

```go
func Sample(l wtype.Liquid, v wunit.Volume) *wtype.LHComponent
```
take a sample of volume v from this liquid

#### func  SampleAll

```go
func SampleAll(l wtype.Liquid) *wtype.LHComponent
```
take all of this liquid

#### func  SampleForConcentration

```go
func SampleForConcentration(l wtype.Liquid, c wunit.Concentration) *wtype.LHComponent
```
take a sample of this liquid and aim for a particular concentration

#### func  SampleForTotalVolume

```go
func SampleForTotalVolume(l wtype.Liquid, v wunit.Volume) *wtype.LHComponent
```
take a sample of this liquid to be used to make the solution up to a particular
total volume
