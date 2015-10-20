---
layout: default
type: api
navgroup: api
shortname: microArch/logger
title: microArch/logger
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: microArch/logger
---
# logger
--
    import "."


## Usage

```go
const (
	INFO = iota
	DEBUG
	WARNING
	ERROR
)
```

#### func  SetLogger

```go
func SetLogger(l *Logger) error
```

#### type AnthaFileLogger

```go
type AnthaFileLogger struct {
	ID       string
	FileName string
}
```

AnthaFileLogger is an implementation of a Logger Interface that uses a file as
storage

#### func  NewAnthaFileLogger

```go
func NewAnthaFileLogger(id string) *AnthaFileLogger
```

#### func  NewAnthaFileLoggerWithFilename

```go
func NewAnthaFileLoggerWithFilename(id string, filename string) *AnthaFileLogger
```

#### func (*AnthaFileLogger) Debug

```go
func (l *AnthaFileLogger) Debug(message string) error
```

#### func (*AnthaFileLogger) Error

```go
func (l *AnthaFileLogger) Error(message string) error
```

#### func (*AnthaFileLogger) GetID

```go
func (l *AnthaFileLogger) GetID() string
```

#### func (*AnthaFileLogger) GetLogById

```go
func (l *AnthaFileLogger) GetLogById(id string) (*LogEntry, error)
```

#### func (*AnthaFileLogger) GetLogList

```go
func (l *AnthaFileLogger) GetLogList() ([]LogEntry, error)
```

#### func (*AnthaFileLogger) GetSensorById

```go
func (l *AnthaFileLogger) GetSensorById(id string) (*SensorReadout, error)
```

#### func (*AnthaFileLogger) GetSensorList

```go
func (l *AnthaFileLogger) GetSensorList() ([]SensorReadout, error)
```

#### func (*AnthaFileLogger) GetTelemetryById

```go
func (l *AnthaFileLogger) GetTelemetryById(id string) (*Telemetry, error)
```

#### func (*AnthaFileLogger) GetTelemetryList

```go
func (l *AnthaFileLogger) GetTelemetryList() ([]Telemetry, error)
```

#### func (*AnthaFileLogger) Info

```go
func (l *AnthaFileLogger) Info(message string) error
```

#### func (*AnthaFileLogger) Log

```go
func (l *AnthaFileLogger) Log(entry LogEntry) error
```
Log saves a log entry to the file, opens the file and flushes always. Appends to
it

#### func (*AnthaFileLogger) LogMessage

```go
func (l *AnthaFileLogger) LogMessage(message string, level LogLevel) error
```

#### func (*AnthaFileLogger) RegisterMiddleware

```go
func (l *AnthaFileLogger) RegisterMiddleware(m LoggerMiddleware)
```

#### func (*AnthaFileLogger) Sense

```go
func (l *AnthaFileLogger) Sense(readout SensorReadout) error
```

#### func (*AnthaFileLogger) TeleMeasure

```go
func (l *AnthaFileLogger) TeleMeasure(signal Telemetry) error
```

#### func (*AnthaFileLogger) Warning

```go
func (l *AnthaFileLogger) Warning(message string) error
```

#### type LogEntry

```go
type LogEntry struct {
	//uuid identifying the log entry
	ID string `json:"_id"`
	//Source uuid that identifies the origin of this message
	Source string
	//Level describes the type of message
	Level LogLevel
	//Message the actual message that is being logged
	Message string
	//Signature means to acknowledge the Source of this message
	Signature string
}
```

LogEntry represents the information saved when a specific message wants to be
saved to db

#### func  NewLogEntry

```go
func NewLogEntry(source string, level LogLevel, message string, signature string) *LogEntry
```

#### func (LogEntry) String

```go
func (l LogEntry) String() string
```

#### type LogLevel

```go
type LogLevel int
```

LogLevel to be able to filter the log messages by importance

#### func (LogLevel) String

```go
func (l LogLevel) String() string
```

#### type Logger

```go
type Logger interface {
	//Log is the basic operation, it saves a LogEntry, wraps a message that is the final information we want to save
	Log(entry LogEntry) error
	//Info Wraps a Log call with an Info level for the message
	Info(message string) error
	//Debug Wraps a Log call with an Info level for the message
	Debug(message string) error
	//Warning Wraps a Log call with an Info level for the message
	Warning(message string) error
	//Error Wraps a Log call with an Info level for the message
	Error(message string) error

	GetLogById(id string) (*LogEntry, error)
	GetLogList() ([]LogEntry, error)

	//TeleMeasure Saves a Telemetry read from a certain piece of equipment
	TeleMeasure(signal Telemetry) error

	GetTelemetryById(id string) (*Telemetry, error)
	GetTelemetryList() ([]Telemetry, error)

	//Sense Saves a Sensor Readout
	Sense(readout SensorReadout) error

	GetSensorById(id string) (*SensorReadout, error)
	GetSensorList() ([]SensorReadout, error)

	//Support Middleware
	RegisterMiddleware(m LoggerMiddleware)
}
```

Logger represents the operations that a logger service must implement //TODO
timestamping

#### func  GetLogger

```go
func GetLogger() *Logger
```

#### type LoggerMiddleware

```go
type LoggerMiddleware interface {
	//Log react to specific Log messages
	Log(entry LogEntry) error
	//Tele react to specific telemetry messages
	Tele(signal Telemetry) error
	//Sensor react to specific sensor readouts
	Sensor(readout SensorReadout) error
}
```

LoggerMiddleware a means to react to specific log events

#### type NilLogger

```go
type NilLogger struct {
}
```


#### func  NewNilLogger

```go
func NewNilLogger() *NilLogger
```

#### func (NilLogger) Debug

```go
func (l NilLogger) Debug(message string) error
```

#### func (NilLogger) Error

```go
func (l NilLogger) Error(message string) error
```

#### func (NilLogger) GetID

```go
func (l NilLogger) GetID() string
```

#### func (NilLogger) GetLogById

```go
func (l NilLogger) GetLogById(id string) (*LogEntry, error)
```

#### func (NilLogger) GetLogList

```go
func (l NilLogger) GetLogList() ([]LogEntry, error)
```

#### func (NilLogger) GetSensorById

```go
func (l NilLogger) GetSensorById(id string) (*SensorReadout, error)
```

#### func (NilLogger) GetSensorList

```go
func (l NilLogger) GetSensorList() ([]SensorReadout, error)
```

#### func (NilLogger) GetTelemetryById

```go
func (l NilLogger) GetTelemetryById(id string) (*Telemetry, error)
```

#### func (NilLogger) GetTelemetryList

```go
func (l NilLogger) GetTelemetryList() ([]Telemetry, error)
```

#### func (NilLogger) Info

```go
func (l NilLogger) Info(message string) error
```

#### func (NilLogger) Log

```go
func (l NilLogger) Log(entry LogEntry) error
```
Log saves a log entry to the file, opens the file and flushes always. Appends to
it

#### func (NilLogger) LogMessage

```go
func (l NilLogger) LogMessage(message string, level LogLevel) error
```

#### func (NilLogger) RegisterMiddleware

```go
func (l NilLogger) RegisterMiddleware(m LoggerMiddleware)
```

#### func (NilLogger) Sense

```go
func (l NilLogger) Sense(readout SensorReadout) error
```

#### func (NilLogger) TeleMeasure

```go
func (l NilLogger) TeleMeasure(signal Telemetry) error
```

#### func (NilLogger) Warning

```go
func (l NilLogger) Warning(message string) error
```

#### type SensorReadout

```go
type SensorReadout struct {
	//uuid identifying the sensor readout
	ID string `json:"_id"`
	//Source uuid that identifies the origin of this readout
	Source string
	//Data information piece that wants to be saved, in string format, can represent aggregated values if desired
	Data string
	//Signature means to acknowledge the origin of the readout is as claimed
	Signature string
}
```

SensorReadout describes a sensor readout

#### func  NewSensorReadout

```go
func NewSensorReadout(source, data, signature string) *SensorReadout
```

#### func (SensorReadout) String

```go
func (r SensorReadout) String() string
```

#### type Telemetry

```go
type Telemetry struct {
	//uuid identifying the telemetry
	ID string `json:"_id"`
	//Source identifies the procedence of the information
	Source string
	//Data telemetry measurement
	Data string
	//Signature acknowledges the origin and err check of the telemetry information
	Signature string
}
```

Telemetry wraps the telemetry information to be logged

#### func  NewTelemetry

```go
func NewTelemetry(source, data, signature string) *Telemetry
```

#### func (Telemetry) String

```go
func (t Telemetry) String() string
```
