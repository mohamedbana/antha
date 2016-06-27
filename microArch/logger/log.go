// microArch/logger/log.go: Part of the Antha language
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

package logger

import "log"

// Middleware using log package
type LogMiddleware struct {
	log *log.Logger
}

func NewLogMiddleware() *LogMiddleware {
	return new(LogMiddleware)
}

func (m *LogMiddleware) Log(level LogLevel, ts int64, source string, message string, extra ...interface{}) {
	m.log.Println(level, ts, source, message, extra)
}

func (m *LogMiddleware) Measure(ts int64, source string, message string, extra ...interface{}) {
	m.log.Println(ts, source, message, extra)
}

func (m *LogMiddleware) Sensor(ts int64, source string, message string, extra ...interface{}) {
	m.log.Println(ts, source, message, extra)
}

func (m *LogMiddleware) Data(ts int64, data interface{}, extra ...interface{}) {
	m.log.Printf("%d | %+v | %+v\n", ts, data, extra)
}
