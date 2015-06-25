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

//AnthaFileLogger is an implementation of a Logger Interface that uses a file as storage
type AnthaFileLogger struct {
	ID         string
	FileName   string
	middleware []LoggerMiddleware
}

func NewAnthaFileLogger(id string) *AnthaFileLogger {
	ret := NewAnthaFileLoggerWithFilename(id, file.LOGFILENAME)
	ret.middleware = make([]LoggerMiddleware, 0)
	return ret
}

func NewAnthaFileLoggerWithFilename(id string, filename string) *AnthaFileLogger {
	ret := new(AnthaFileLogger)
	ret.ID = id
	ret.FileName = filename
	return ret
}

func (l *AnthaFileLogger) GetID() string {
	return l.ID
}

func (l *AnthaFileLogger) LogMessage(message string, level LogLevel) error {
	return l.Log(LogEntry{Message: message, Level: level})
}
func (l *AnthaFileLogger) Info(message string) error {
	return l.LogMessage(message, INFO)
}
func (l *AnthaFileLogger) Debug(message string) error {
	return l.LogMessage(message, DEBUG)
}
func (l *AnthaFileLogger) Warning(message string) error {
	return l.LogMessage(message, WARNING)
}
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

func (l *AnthaFileLogger) GetLogList() ([]LogEntry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) GetTelemetryList() ([]Telemetry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) GetSensorList() ([]SensorReadout, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) GetSensorById(id string) (*SensorReadout, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) GetTelemetryById(id string) (*Telemetry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) GetLogById(id string) (*LogEntry, error) {
	return nil, errors.New("File Logger does not implement this operation")
}

func (l *AnthaFileLogger) RegisterMiddleware(m LoggerMiddleware) {
	l.middleware = append(l.middleware, m)
}
