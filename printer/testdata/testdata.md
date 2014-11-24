---
layout: default
type: api
navgroup: docs
shortname: printer/testdata
title: printer/testdata
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: printer/testdata
---
# parser
--
    import "."


## Usage

```go
const (
	PackageClauseOnly uint = 1 << iota // parsing stops after package clause
	ImportsOnly                        // parsing stops after import declarations
	ParseComments                      // parse comments and add them to AST
	Trace                              // print a trace of parsed productions
	DeclarationErrors                  // report declaration errors
)
```
The mode parameter to the Parse* functions is a set of flags (or 0). They
control the amount of source code parsed and other optional parser
functionality.
