// microArch/logger/middleware/handler/handler_test.go: Part of the Antha language
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

package handler

import (
	"testing"

	"time"

	"github.com/antha-lang/antha/microArch/logger"
)

func TestLogHandler(t *testing.T) {
	counter := 0
	logHandlerFunc := func(level logger.LogLevel, ts int64, origin, message string, extra ...interface{}) {
		if message == "test" {
			counter++
		}
	}

	handler := NewHandler(logHandlerFunc, nil, nil, nil)

	handler.Log(logger.DEBUG, time.Now().Unix(), "test", "test")
	if counter != 1 {
		t.Error(
			"For", "counter",
			"expected", 1,
			"got", counter,
		)
	}
}

func TestLogHandlerMiddleware(t *testing.T) {
	counter := 0
	logHandlerFunc := func(level logger.LogLevel, ts int64, source, message string, extra ...interface{}) {
		if message == "test" {
			counter++
		}
	}

	handler := NewHandler(logHandlerFunc, nil, nil, nil)
	logger.RegisterMiddleware(handler)

	logger.Debug("test")
	if counter != 1 {
		t.Error(
			"For", "counter",
			"expected", 1,
			"got", counter,
		)
	}
}
