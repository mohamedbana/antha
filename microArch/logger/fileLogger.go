// /logger/fileLogger.go: Part of the Antha language
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
	"fmt"
	"os"

	"github.com/antha-lang/antha/microArch/config/file"
)

//AnthaFileLogger is a Logger that will store every logged information to a particular file
type AnthaFileLogger struct {
	ID         string
	FileName   string
	middleware []LoggerMiddleware
}

//AnthaFileLogger is an implementation of a Logger Interface that uses a file as storage. The file used
// as storage is defined as a configuration value
func NewAnthaFileLogger(id string) *AnthaFileLogger {
	ret := NewAnthaFileLoggerWithFilename(id, file.LOGFILENAME)
	ret.middleware = make([]LoggerMiddleware, 0)
	return ret
}

//NewAnthaFileLoggerWithFilename instantiates a new AnthaFileLogger with the given id and filename
// as output
func NewAnthaFileLoggerWithFilename(id string, filename string) *AnthaFileLogger {
	ret := new(AnthaFileLogger)
	ret.ID = id
	ret.FileName = filename
	return ret
}

//GetID returns the uniqued id that identifies this AnthaFileLogger. Ideally a uuid
func (l *AnthaFileLogger) GetID() string {
	return l.ID
}

//LogMessage creates a LogEntry and logs it with the given message and log level
func (l *AnthaFileLogger) LogMessage(message string, level LogLevel) error {
	return l.Log(LogEntry{Message: message, Level: level})
}

//Info saves a message with a Information level
func (l *AnthaFileLogger) Info(message string) error {
	return l.LogMessage(message, INFO)
}

//Debug reports the message to log with debug level
func (l *AnthaFileLogger) Debug(message string) error {
	return l.LogMessage(message, DEBUG)
}

//Warning reports the message as a warning
func (l *AnthaFileLogger) Warning(message string) error {
	return l.LogMessage(message, WARNING)
}

//Error reports the given message as an error
func (l *AnthaFileLogger) Error(message string) error {
	return l.LogMessage(message, ERROR)
}

//Log saves a log entry to the file, opens the file and flushes always. Appends to it
func (l *AnthaFileLogger) Log(entry LogEntry) error {
	for _, m := range l.middleware {
		m.Log(entry)
	}
	f, err := os.OpenFile(l.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n", entry)) //TODO, should we check for number of bytes written??
	return err
}

//TeleMeasure Saves a Telemetry read from a certain piece of equipment
func (l *AnthaFileLogger) TeleMeasure(signal Telemetry) error {
	for _, m := range l.middleware {
		m.Tele(signal)
	}
	f, err := os.OpenFile(l.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n", signal)) //TODO, should we check for number of bytes written??
	return err
}

//Sense Saves a Sensor Readout
func (l *AnthaFileLogger) Sense(readout SensorReadout) error {
	for _, m := range l.middleware {
		m.Sensor(readout)
	}
	f, err := os.OpenFile(l.FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n", readout)) //TODO, should we check for number of bytes written??
	return err
}

//GetLogList retrieves a list of LogEntries
func (l *AnthaFileLogger) GetLogList() ([]LogEntry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//GetTelemetryList will provide a list of the telemetry data available
func (l *AnthaFileLogger) GetTelemetryList() ([]Telemetry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//GetSensorList will provide a list of the sensor data available
func (l *AnthaFileLogger) GetSensorList() ([]SensorReadout, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//GetLogById retrieves the LogEntry identified by the given id
func (l *AnthaFileLogger) GetSensorById(id string) (*SensorReadout, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//GetTelemetryById gets a reference to a Telemetry struct identified by the given id.
func (l *AnthaFileLogger) GetTelemetryById(id string) (*Telemetry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//GetLogById retrieves the LogEntry identified by the given id
func (l *AnthaFileLogger) GetLogById(id string) (*LogEntry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

//RegisterMiddleware attaches a log middleware to this log. Every LogEntry, Telemetry and Sensor Readout will be reported to it
func (l *AnthaFileLogger) RegisterMiddleware(m LoggerMiddleware) {
	l.middleware = append(l.middleware, m)
}
