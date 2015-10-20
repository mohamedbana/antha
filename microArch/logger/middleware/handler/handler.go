// microArch/logger/middleware/handler/handler.go: Part of the Antha language
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
	"github.com/antha-lang/antha/microArch/logger"
)

//logHandlerFunc level, unixtimestamp, origin, message, extra
type LogHandlerFunc func(logger.LogLevel, int64, string, string, ...interface{})
type MeasureHandlerFunc func(int64, string, string, ...interface{})
type SensorHandlerFunc func(int64, string, string, ...interface{})

type Handler struct {
	l LogHandlerFunc
	t MeasureHandlerFunc
	s SensorHandlerFunc
}

func NewHandler(l LogHandlerFunc, t MeasureHandlerFunc, s SensorHandlerFunc) *Handler {
	ret := new(Handler)
	ret.l = l
	ret.t = t
	ret.s = s
	return ret
}

//Log react to specific Log messages
func (h Handler) Log(level logger.LogLevel, ts int64, origin, message string, extra ...interface{}) {
	if h.l != nil {
		h.l(level, ts, origin, message, extra)
	}
}

//Tele react to specific telemetry messages
func (h Handler) Measure(ts int64, origin string, message string, extra ...interface{}) {
	if h.t != nil {
		h.t(ts, origin, message, extra)
	}
}

//Sensor react to specific sensor readouts
func (h Handler) Sensor(ts int64, origin string, message string, extra ...interface{}) {
	if h.s != nil {
		h.s(ts, origin, message, extra)
	}
}
