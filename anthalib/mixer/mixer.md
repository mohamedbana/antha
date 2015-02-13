---
layout: default
type: api
navgroup: docs
shortname: anthalib/mixer
title: anthalib/mixer
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: anthalib/mixer
---
# mixer
--
    import "."


## Usage

#### type SampleComponent

```go
type SampleComponent map[string]interface{}
```

mix needs to define the interface with liquid handling in order to do this it
has to make the appropriate liquid handling request structure the functions in
this package mostly convert the wtype representations into the simpler map
representations which are the interface with the network layer

#### func  Mix

```go
func Mix(components ...SampleComponent) SampleComponent
```
mix the specified SampleComponents together the destination will be the location
of the first sample

#### func  MixInto

```go
func MixInto(destination wtype.LiquidContainer, components ...SampleComponent) SampleComponent
```
mix the specified SampleComponents together into the destination specified as
the first argument

#### func  Sample

```go
func Sample(l wtype.Liquid, v wunit.Volume) SampleComponent
```
take a sample of volume v from this liquid

#### func  SampleAll

```go
func SampleAll(l wtype.Liquid) SampleComponent
```
take all of this liquid

#### func  SampleForConcentration

```go
func SampleForConcentration(l wtype.Liquid, c wunit.Concentration) SampleComponent
```
take a sample of this liquid and aim for a particular concentration

#### func  SampleForTotalVolume

```go
func SampleForTotalVolume(l wtype.Liquid, v wunit.Volume) SampleComponent
```
take a sample of this liquid to be used to make the solution up to a particular
total volume
