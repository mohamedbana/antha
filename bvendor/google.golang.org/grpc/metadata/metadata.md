---
layout: default
type: api
navgroup: docs
shortname: grpc/metadata
title: grpc/metadata
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: grpc/metadata
---
# metadata
--
    import "github.com/antha-lang/antha/internal/google.golang.org/grpc/metadata"

Package metadata define the structure of the metadata supported by gRPC library.

## Usage

#### func  DecodeKeyValue

```go
func DecodeKeyValue(k, v string) (string, string, error)
```
DecodeKeyValue returns the original key and value corresponding to the encoded
data in k, v.

#### func  NewContext

```go
func NewContext(ctx context.Context, md MD) context.Context
```
NewContext creates a new context with md attached.

#### type MD

```go
type MD map[string][]string
```

MD is a mapping from metadata keys to values. Users should use the following two
convenience functions New and Pairs to generate MD.

#### func  FromContext

```go
func FromContext(ctx context.Context) (md MD, ok bool)
```
FromContext returns the MD in ctx if it exists.

#### func  New

```go
func New(m map[string]string) MD
```
New creates a MD from given key-value map.

#### func  Pairs

```go
func Pairs(kv ...string) MD
```
Pairs returns an MD formed by the mapping of key, value ... Pairs panics if
len(kv) is odd.

#### func (MD) Copy

```go
func (md MD) Copy() MD
```
Copy returns a copy of md.

#### func (MD) Len

```go
func (md MD) Len() int
```
Len returns the number of items in md.
