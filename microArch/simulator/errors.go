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
    InfoSeverity SimulationErrorSeverity = iota
    WarningSeverity
    ErrorSeverity
)

severityNames := map[SimulationErrorSeverity]string {
    InfoSeverity: "info",
    WarningSeverity: "warn",
    ErrorSeverity: "err"
}

type SimulationError interface {
    Error() string
    ErrorSeverity() ErrorSeverity
    ErrorInstruction() interface {}
}

type SimulationError struct {
    Severity        SimulationErrorSeverity
    Instruction     interface{}
    What            string
}

func NewSimulationError(severity ErrorSeverity, what string, instruction interface{}) *SimulationError {
    var err SimulationError
    err.Severity = severity
    err.What = what
    err.Instruction = instruction
    return &err
}

func (self *SimulationError) Error() string {
    return fmt.Sprintf("(%s) %s: %s",
                simulationSeverityNames[self.Severity],
                self.ErrorTypeName(),
                self.What)
}

func (self *SimulationError) ErrorSeverity {
    return self.Severity
}

func (self *SimulationError) ErrorInstruction() string {
    return self.Instruction
}


