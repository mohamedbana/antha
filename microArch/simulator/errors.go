// /anthalib/simulator/errors.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

/* Simulation Errors and Warnings */

package simulator

import (
    "fmt"
)

//SimulationErrorSeverity how serious is the error
type ErrorSeverity int
const (
    SeverityNone ErrorSeverity = iota
    SeverityInfo
    SeverityWarning
    SeverityError
)

var severityNames = map[ErrorSeverity]string {
    SeverityInfo: "info",
    SeverityWarning: "warn",
    SeverityError: "err",
}

type SimulationError struct {
    severity        ErrorSeverity
    fname           string
    message         string
}

func NewSimulationError(severity ErrorSeverity, message string, function_name string) *SimulationError {
    err := SimulationError{severity, function_name, message}
    return &err
}

func (self *SimulationError) Error() string {
    return fmt.Sprintf("(%s) %s: %s",
                severityNames[self.Severity()],
                self.fname,
                self.message)
}

func (self *SimulationError) Severity() ErrorSeverity {
    return self.severity
}

func (self *SimulationError) FunctionName() string {
    return self.fname
}

//ErrorReporter 
type ErrorReporter struct {
    errors []*SimulationError
    worst_error *SimulationError
}

func NewErrorReporter() *ErrorReporter {
    er := ErrorReporter{}
    er.errors = make([]*SimulationError, 0)
    er.worst_error = nil

    return &er
}

func (self *ErrorReporter) GetErrors() []*SimulationError {
    return self.errors
}

func (self *ErrorReporter) GetWorstError() *SimulationError {
    return self.worst_error
}

func (self *ErrorReporter) GetErrorSeverity() ErrorSeverity {
    if self.worst_error != nil {
        return self.worst_error.Severity()
    }
    return SeverityNone
}

func (self *ErrorReporter) HasError() bool {
    return self.GetErrorSeverity() == SeverityError
}

func (self *ErrorReporter) HasWarning() bool {
    return self.GetErrorSeverity() >= SeverityWarning
}

func (self *ErrorReporter) AddSimulationError(se *SimulationError) {
    self.errors = append(self.errors, se)
    if se.Severity() > self.GetErrorSeverity() {
        self.worst_error = se
    }
}

func (self *ErrorReporter) AddError(fname, msg string) {
    se := NewSimulationError(SeverityError, msg, fname)
    self.AddSimulationError(se)
}

func (self *ErrorReporter) AddWarning(fname, msg string) {
    se := NewSimulationError(SeverityWarning, msg, fname)
    self.AddSimulationError(se)
}

func (self *ErrorReporter) AddInfo(fname, msg string) {
    se := NewSimulationError(SeverityInfo, msg, fname)
    self.AddSimulationError(se)
}

func (self *ErrorReporter) AddErrorf(fname, format string, args ...interface{}) {
    self.AddError(fname, fmt.Sprintf(format, args...))
}

func (self *ErrorReporter) AddWarningf(fname, format string, args ...interface{}) {
    self.AddWarning(fname, fmt.Sprintf(format, args...))
}

func (self *ErrorReporter) AddInfof(fname, format string, args ...interface{}) {
    self.AddInfo(fname, fmt.Sprintf(format, args...))
}


