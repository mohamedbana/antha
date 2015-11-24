---
layout: default
type: api
navgroup: docs
shortname: grpc/grpclog
title: grpc/grpclog
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: grpc/grpclog
---
# grpclog
--
    import "github.com/antha-lang/antha/internal/google.golang.org/grpc/grpclog"

Package grpclog defines logging for grpc.

## Usage

#### func  Fatal

```go
func Fatal(args ...interface{})
```
Fatal is equivalent to Print() followed by a call to os.Exit() with a non-zero
exit code.

#### func  Fatalf

```go
func Fatalf(format string, args ...interface{})
```
Fatalf is equivalent to Printf() followed by a call to os.Exit() with a non-zero
exit code.

#### func  Fatalln

```go
func Fatalln(args ...interface{})
```
Fatalln is equivalent to Println() followed by a call to os.Exit()) with a
non-zero exit code.

#### func  Print

```go
func Print(args ...interface{})
```
Print prints to the logger. Arguments are handled in the manner of fmt.Print.

#### func  Printf

```go
func Printf(format string, args ...interface{})
```
Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.

#### func  Println

```go
func Println(args ...interface{})
```
Println prints to the logger. Arguments are handled in the manner of
fmt.Println.

#### func  SetLogger

```go
func SetLogger(l Logger)
```
SetLogger sets the logger that is used in grpc.

#### type Logger

```go
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}
```

Logger mimics golang's standard Logger as an interface.
