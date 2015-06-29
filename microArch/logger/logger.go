// /logger/logger.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

//package logger contains the data structures necessary for the logging
// of the messages in the antha platform.
package logger

import "fmt"

//LogLevel to be able to filter the log messages by importance
type LogLevel int

const (
	//INFO information level messages
	INFO = iota
	//DEBUG level messages should represent information useful for debugging purposes,
	// may contain technical details related to the code
	DEBUG
	//WARNING level messages represent a message used to identify a non fatal error. An
	// undesired situation that should be looked into
	WARNING
	//ERROR level messages should always be informed
	ERROR
)

//String gives a textual representation for a log level
func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

//LogEntry represents the information saved when a specific message wants to be saved to db
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

//String gives a textual string representation of a LogEntry including all the relevant information
func (l LogEntry) String() string {
	return fmt.Sprintf("LOG[%s] ID %s, SRC [%s]: %s {SIG: <%s>}", l.Level, l.ID, l.Source, l.Message, l.Signature)
}

//NewLogEntry creates a new LogEntry struct with the data given as parameters, source, level, message and signature
func NewLogEntry(source string, level LogLevel, message string, signature string) *LogEntry {
	ret := new(LogEntry)
	ret.Source = source
	ret.Level = level
	ret.Message = message
	ret.Signature = signature
	return ret
}

//SensorReadout describes a sensor readout
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

//String provides a text represntation of a SensorReadout
func (r SensorReadout) String() string {
	return fmt.Sprintf("SNS ID %s, SRC [%s]: %s {SIG: <%s>}", r.ID, r.Source, r.Data, r.Signature)
}

//NewSensorReadout creates a new sensor readout with the supplied parameters, source, data and signature
func NewSensorReadout(source, data, signature string) *SensorReadout {
	ret := new(SensorReadout)
	ret.Source = source
	ret.Data = data
	ret.Signature = signature
	return ret
}

//Telemetry wraps the telemetry information to be logged
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

//String provides a textual representation of a Telemetry struct
func (t Telemetry) String() string {
	return fmt.Sprintf("TLM ID %s, SRC [%s]: %s {SIG: <%s>}", t.ID, t.Source, t.Data, t.Signature)
}

//NewTelemetry instantiates a Telemetry struct with the given parameters
func NewTelemetry(source, data, signature string) *Telemetry {
	t := new(Telemetry)
	t.Source = source
	t.Data = data
	t.Signature = signature
	return t
}

//Logger represents the operations that a logger service must implement //TODO timestamping
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

	//GetLogById retrieves the LogEntry identified by the given id
	GetLogById(id string) (*LogEntry, error)
	//GetLogList retrieves a list of LogEntries
	GetLogList() ([]LogEntry, error)

	//TeleMeasure Saves a Telemetry read from a certain piece of equipment
	TeleMeasure(signal Telemetry) error

	//GetTelemetryById gets a reference to a Telemetry struct identified by the given id.
	GetTelemetryById(id string) (*Telemetry, error)
	//GetTelemetryList will provide a list of the telemetry data available
	GetTelemetryList() ([]Telemetry, error)

	//Sense Saves a Sensor Readout
	Sense(readout SensorReadout) error

	//GetSensorById gets a reference to a sensorreadout identified by the given id
	GetSensorById(id string) (*SensorReadout, error)
	//GetSensorList will provide a list of the sensor data available
	GetSensorList() ([]SensorReadout, error)

	//RegisterMiddleware attaches a log middleware to this log. Every LogEntry, Telemetry and Sensor Readout will be reported to it
	RegisterMiddleware(m LoggerMiddleware)
}

//LoggerMiddleware a means to react to specific log events
type LoggerMiddleware interface {
	//Log react to specific Log messages
	Log(entry LogEntry) error
	//Tele react to specific telemetry messages
	Tele(signal Telemetry) error
	//Sensor react to specific sensor readouts
	Sensor(readout SensorReadout) error
}
