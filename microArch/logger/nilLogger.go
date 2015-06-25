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


type NilLogger struct {
}

func NewNilLogger() *NilLogger {
	ret := new(NilLogger)
	return ret
}

func (l NilLogger) GetID() string {
	return ""
}

func (l NilLogger) LogMessage(message string, level LogLevel) error {
	return l.doNothing()
}
func (l NilLogger) Info(message string) error {
	return l.doNothing()
}
func (l NilLogger) Debug(message string) error {
	return l.doNothing()
}
func (l NilLogger) Warning(message string) error {
	return l.doNothing()
}
func (l NilLogger) Error(message string) error {
	return l.doNothing()
}
//Log saves a log entry to the file, opens the file and flushes always. Appends to it
func (l NilLogger) Log(entry LogEntry) error {
	return l.doNothing()
}

func (l NilLogger) TeleMeasure(signal Telemetry) error {
	return l.doNothing()
}

func (l NilLogger) Sense(readout SensorReadout) error{
	return l.doNothing()
}

func (l NilLogger) doNothing() error{
	return nil
}

func (l NilLogger) GetLogList() ([]LogEntry, error){
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) GetTelemetryList() ([]Telemetry, error){
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) GetSensorList() ([]SensorReadout, error){
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) GetSensorById(id string) (*SensorReadout, error){
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) GetTelemetryById(id string) (*Telemetry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) GetLogById(id string) (*LogEntry, error) {
	return nil, errors.New("Nil Logger does not implement this operation")
}

func (l NilLogger) RegisterMiddleware(m LoggerMiddleware) {
	return
}
