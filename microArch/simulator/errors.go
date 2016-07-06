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
    instruction     interface{}
    message         string
}

func NewSimulationError(severity ErrorSeverity, message string, instruction interface{}) *SimulationError {
    var err SimulationError
    err.severity = severity
    err.message = message
    err.instruction = instruction
    return &err
}

func (self *SimulationError) Error() string {
    return fmt.Sprintf("(%s) %s",
                severityNames[self.Severity()],
                self.message)
}

func (self *SimulationError) Severity() ErrorSeverity {
    return self.severity
}

func (self *SimulationError) Instruction() interface{} {
    return self.instruction
}


