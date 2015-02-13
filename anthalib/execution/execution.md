---
layout: default
type: api
navgroup: docs
shortname: anthalib/execution
title: anthalib/execution
apidocs:
  published: 2014-11-14
  antha_version: 0.0.1
  package: anthalib/execution
---
# execution
--
    import "."


## Usage

#### func  GetUUID

```go
func GetUUID() string
```
this package wraps the uuid library appropriately by generating a V4 UUID

#### func  StartRuntime

```go
func StartRuntime()
```
not a true function: has the side-effect of starting runtime services

#### type AnthaConfig

```go
type AnthaConfig map[string]interface{}
```


#### func  NewAnthaConfig

```go
func NewAnthaConfig() *AnthaConfig
```

#### type ExecutionService

```go
type ExecutionService struct {
	StockMgr         *StockService
	SampleTracker    *SampleTrackerService
	Logger           *LogService
	Scheduler        *ScheduleService
	GarbageCollector *GarbageCollectionService
	// and also to the config
	Config *AnthaConfig
}
```

A structure which defines the execution context - holds pointers to the services

#### func  GetContext

```go
func GetContext() *ExecutionService
```

#### type GarbageCollectionRequest

```go
type GarbageCollectionRequest map[string]interface{}
```

map data structure defining a request for an object to enter the waste stream

#### type GarbageCollectionService

```go
type GarbageCollectionService struct {
	RequestsIn  chan GarbageCollectionRequest
	RequestsOut chan GarbageCollectionRequest
}
```

the garbage collector holds channels for communicating with the garbage
collection service provider

#### func  NewGarbageCollectionService

```go
func NewGarbageCollectionService() *GarbageCollectionService
```

#### func (*GarbageCollectionService) Init

```go
func (gcs *GarbageCollectionService) Init()
```
Initialize the garbage collection service the channels are given a default
capacity for 5 items before blocking

#### func (*GarbageCollectionService) RequestGarbageCollection

```go
func (gcs *GarbageCollectionService) RequestGarbageCollection(rin GarbageCollectionRequest) GarbageCollectionRequest
```
send an item to the waste stream

#### type LogRequest

```go
type LogRequest map[string]interface{}
```

data structure for defining a request to the logger

#### type LogService

```go
type LogService struct {
	RequestsIn  chan LogRequest
	RequestsOut chan LogRequest
}
```

structure defining a log service - holds channels for communicating with the
service

#### func  NewLogService

```go
func NewLogService() *LogService
```

#### func (*LogService) Init

```go
func (ls *LogService) Init()
```
initialize the log service the channels have capacity to buffer 5 requests by
default

#### func (*LogService) RequestLog

```go
func (ls *LogService) RequestLog(rin LogRequest) LogRequest
```
function for communicating with the log daemon

#### type SampleRequest

```go
type SampleRequest map[string]interface{}
```

data structure defining sample requests

#### type SampleTrackerService

```go
type SampleTrackerService struct {
	RequestsIn  chan SampleRequest
	RequestsOut chan SampleRequest
}
```

struct to hold channels for communicating with the sample tracker service

#### func  NewSampleTrackerService

```go
func NewSampleTrackerService() *SampleTrackerService
```

#### func (*SampleTrackerService) Init

```go
func (ss *SampleTrackerService) Init()
```
initialize the sample tracker agent the channels have capacity to buffer 5
requests before blocking

#### func (*SampleTrackerService) TrackSample

```go
func (ss *SampleTrackerService) TrackSample(rin SampleRequest) SampleRequest
```
send a request to the sample tracking service - e.g. update status or register a
new sample

#### type ScheduleRequest

```go
type ScheduleRequest map[string]interface{}
```

data structure for defining a request to communicate with the sceduler

#### type ScheduleService

```go
type ScheduleService struct {
	RequestsIn  chan ScheduleRequest
	RequestsOut chan ScheduleRequest
}
```

data structure for holding channels to communicate with the scheduler

#### func  NewScheduleService

```go
func NewScheduleService() *ScheduleService
```

#### func (*ScheduleService) Init

```go
func (ss *ScheduleService) Init()
```
start up the scheduler and make the communication channels the channels have
room to buffer 5 requests before blocking

#### func (*ScheduleService) RequestSchedule

```go
func (ss *ScheduleService) RequestSchedule(rin ScheduleRequest) ScheduleRequest
```
communicate a request to the scheduler and return the reply

#### type StockRequest

```go
type StockRequest map[string]interface{}
```

a map structure for defining requests to the stock manager

#### type StockService

```go
type StockService struct {
	RequestsIn  chan StockRequest
	RequestsOut chan StockRequest
}
```

structure to hold communication channels with the stock service

#### func  NewStockService

```go
func NewStockService() *StockService
```

#### func (*StockService) Init

```go
func (ss *StockService) Init()
```
Initialize the communication service The channels have the capacity to hold 5
requests before blocking

#### func (*StockService) RequestStock

```go
func (ss *StockService) RequestStock(rin StockRequest) StockRequest
```
Make a request to the stock agent and return the result
