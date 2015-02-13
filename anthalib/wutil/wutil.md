---
layout: default
type: api
navgroup: docs
shortname: anthalib/wutil
title: anthalib/wutil
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: anthalib/wutil
---
# wutil
--
    import "."


## Usage

#### func  AlphaToNum

```go
func AlphaToNum(s string) int
```

#### func  Comp

```go
func Comp(s string) string
```

#### func  Copy

```go
func Copy(s1 map[string]interface{}) map[string]interface{}
```
this changes the UUID appropriately; more likely to be useful in practice also
ensures "inst" is not copied since we cannot physically duplicate anything

#### func  DecodeCoords

```go
func DecodeCoords(s string) (int, int)
```

#### func  DisplayMap

```go
func DisplayMap(s string, m map[string]interface{})
```

#### func  Dup

```go
func Dup(s1 map[string]interface{}) map[string]interface{}
```
an EXACT duplicate, including UUIDs which ordinarily should change

#### func  Error

```go
func Error(msg error)
```

#### func  FMax

```go
func FMax(floats []float64) float64
```

#### func  FMin

```go
func FMin(floats []float64) float64
```

#### func  GetFloat64FromMap

```go
func GetFloat64FromMap(m map[string]interface{}, k string) float64
```

#### func  GetIntFromMap

```go
func GetIntFromMap(m map[string]interface{}, k string) int
```

#### func  GetMapFromMap

```go
func GetMapFromMap(m map[string]interface{}, k string) map[string]interface{}
```

#### func  GetStringFromMap

```go
func GetStringFromMap(m map[string]interface{}, k string) string
```

#### func  MakeRankedList

```go
func MakeRankedList(arr []int) []int
```

#### func  MakeUnique

```go
func MakeUnique(a []string) []string
```
returns the same list of strings in the same order minus any duplicates i.e.
ordered by first appearance

#### func  Max

```go
func Max(ints []int) int
```

#### func  Min

```go
func Min(ints []int) int
```

#### func  NumToAlpha

```go
func NumToAlpha(n int) string
```

#### func  ParseFloat

```go
func ParseFloat(str string) float64
```

#### func  ParseInt

```go
func ParseInt(str string) int
```

#### func  ReadFastaSeqs

```go
func ReadFastaSeqs(fn string) []seq.Sequence
```

#### func  Rev

```go
func Rev(s string) string
```

#### func  RevComp

```go
func RevComp(s string) string
```

#### func  RoundInt

```go
func RoundInt(v float64) int
```

#### func  SeqToBioseq

```go
func SeqToBioseq(s seq.Sequence) string
```

#### func  Series

```go
func Series(start, end int) []int32
```

#### func  SortMapKeysBy

```go
func SortMapKeysBy(hashes []map[string]interface{}, key string)
```

#### func  Warn

```go
func Warn(s string)
```

#### func  WndsWith

```go
func WndsWith(s, sfx string) bool
```
