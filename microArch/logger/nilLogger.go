// /logger/nilLogger.go: Part of the Antha language
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

package logger

import (
	"errors"
)

//NilLogger is a special type of Logger that does nothing with the output,
// the middleware and the id. It basically has no functionality as a logger
type NilLogger struct {
}

//NewNilLogger instantiates a new NilLogger struct
func NewNilLogger() *NilLogger {
	ret := new(NilLogger)
	return ret
}

//GetID returns nothing for NilLogger
func (l NilLogger) GetID() string {
	return ""
}

//Info Wraps a Log call with an Info level for the message
func (l NilLogger) Info(message string) error {
	return l.doNothing()
}

//Debug Wraps a Log call with an Info level for the message
func (l NilLogger) Debug(message string) error {
	return l.doNothing()
}

//Warning Wraps a Log call with an Info level for the message
func (l NilLogger) Warning(message string) error {
	return l.doNothing()
}

//Error Wraps a Log call with an Info level for the message
func (l NilLogger) Error(message string) error {
	return l.doNothing()
}

//Log is the basic operation, it saves a LogEntry, wraps a message that is the final information we want to save
func (l NilLogger) Log(entry LogEntry) error {
	return l.doNothing()
}

//TeleMeasure Saves a Telemetry read from a certain piece of equipment
func (l NilLogger) TeleMeasure(signal Telemetry) error {
	return l.doNothing()
}

//Sense Saves a Sensor Readout
func (l NilLogger) Sense(readout SensorReadout) error {
	return l.doNothing()
}

//doNothing is a method used to annotate the rest of the methods that do not have an actual implementation
func (l NilLogger) doNothing() error {
	return nil
}

//GetLogList retrieves a list of LogEntries
func (l NilLogger) GetLogList() ([]LogEntry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//GetTelemetryList will provide a list of the telemetry data available
func (l NilLogger) GetTelemetryList() ([]Telemetry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//GetSensorList will provide a list of the sensor data available
func (l NilLogger) GetSensorList() ([]SensorReadout, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//GetSensorById gets a reference to a sensorreadout identified by the given id
func (l NilLogger) GetSensorById(id string) (*SensorReadout, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//GetTelemetryById gets a reference to a Telemetry struct identified by the given id.
func (l NilLogger) GetTelemetryById(id string) (*Telemetry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//GetLogById retrieves the LogEntry identified by the given id
func (l NilLogger) GetLogById(id string) (*LogEntry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

//RegisterMiddleware attaches a log middleware to this log. Every LogEntry, Telemetry and Sensor Readout will be reported to it
func (l NilLogger) RegisterMiddleware(m LoggerMiddleware) {
	return
}
