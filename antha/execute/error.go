// antha/execute/error.go: Part of the Antha language
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

package execute

import (
	"fmt"
	"strings"
	"sync"
)

var errors *ExecutionError
var lock sync.Mutex

type RuntimeError struct {
	BaseError interface{}
	Stack     []byte
}

func (a *RuntimeError) Error() string {
	return fmt.Sprintf("%s at:\n%s", a.BaseError, string(a.Stack))
}

// Simplify component control flow by storing errors during component execution
// in a global
type ExecutionError struct {
	errors []interface{}
}

func (a *ExecutionError) Add(err interface{}) {
	a.errors = append(a.errors, err)
}

func (a *ExecutionError) Error() string {
	var lines []string
	for _, v := range a.errors {
		switch v := v.(type) {
		case error:
			lines = append(lines, v.Error())
		case string:
			lines = append(lines, v)
		default:
			lines = append(lines, fmt.Sprintf("%s", v))
		}
	}
	return strings.Join(lines, "\n")
}

// Reset state of global error list
func TakeErrors() error {
	lock.Lock()
	defer lock.Unlock()
	err := errors
	errors = nil
	if err != nil {
		return err
	}
	return nil
}

// Add error to global error list
func AddError(err interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if errors == nil {
		errors = &ExecutionError{[]interface{}{err}}
	} else {
		errors.Add(err)
	}
}
