// microArch/logger/logger_test.go: Part of the Antha language
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

import (
	"bytes"
	"log"
	"sync"
	"testing"
)

var (
	logCounter     int
	measureCounter int
	sensorCounter  int
)

type testMiddleware struct{}

func (t testMiddleware) Log(l LogLevel, ts int64, s string, m string, e ...interface{}) { logCounter++ }
func (t testMiddleware) Measure(ts int64, s, m string, e ...interface{})                { measureCounter++ }
func (t testMiddleware) Sensor(ts int64, s, m string, e ...interface{})                 { sensorCounter++ }

func TestRegisterMiddleware(t *testing.T) {
	midCount := len(middlewares)
	RegisterMiddleware(testMiddleware{})
	if l := len(middlewares); l != 1+midCount {
		t.Error("middlewares size err. Expecting ", 1+midCount, " got ", l)
	}
}

func TestRegisterMiddlewareLock(t *testing.T) {
	midCount := len(middlewares)
	count := 20
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			RegisterMiddleware(testMiddleware{})
			wg.Done()
		}()
	}
	wg.Wait()
	if l := len(middlewares); l != count+midCount {
		t.Error("middlewares size err. Expecting ", count+midCount, " got ", l)
	}
}

func TestMiddlewareCalls(t *testing.T) {
	cleanMiddleware()
	//reset counter values
	logCounter = 0
	measureCounter = 0
	sensorCounter = 0
	RegisterMiddleware(testMiddleware{})
	Info("")
	Debug("")
	Warning("")
	Error("")
	if l := logCounter; l != 4 {
		t.Errorf("expected %d log calls, got %d", 4, l)
	}

	Measure("")
	if l := measureCounter; l != 1 {
		t.Errorf("expected %d measure calls, got %d", 1, l)
	}
	Sensor("")
	if l := sensorCounter; l != 1 {
		t.Error("expecting %d sensor calls, got %d", 1, l)
	}
}

func TestFatal(t *testing.T) {
	cleanMiddleware()
	panicked := false
	aTest := func() {
		defer func() {
			if res := recover(); res != nil {
				panicked = true
				return
			}
		}()
		Fatal("")
	}
	aTest()
	if !panicked {
		t.Error("expecting panic")
	}
}

func TestMiddlewareCallsSync(t *testing.T) {
	cleanMiddleware()
	//reset counter values
	logCounter = 0
	measureCounter = 0
	sensorCounter = 0

	syncCount := 5
	RegisterMiddleware(testMiddleware{})
	wg := sync.WaitGroup{}
	wg.Add(syncCount)

	for i := 0; i < syncCount; i++ {
		go func() {
			Info("")
			Debug("")
			Warning("")
			Error("")
			Measure("")
			Sensor("")
			wg.Done()
		}()
	}

	wg.Wait()
	if l := logCounter; l != 4*syncCount {
		t.Error("On Log. Expecting call count ", 4*syncCount, " got ", l)
	}
	if l := measureCounter; l != 1*syncCount {
		t.Error("On Measure. Expecting call count ", 1*syncCount, " got ", l)
	}
	if l := sensorCounter; l != 1*syncCount {
		t.Error("On Sensor. Expecting call count ", 1*syncCount, " got ", l)
	}
}

func TestDefaultLogger(t *testing.T) {
	cleanMiddleware() //MUST!
	x := make([]byte, 0)
	buffer := bytes.NewBuffer(x)
	log.SetOutput(buffer)
	Info("test")
	if len(buffer.Bytes()) == 0 {
		t.Errorf("default logger output empty. Expecting ", len("test"))
	}
	cleanMiddleware() //MUST!
	RegisterMiddleware(testMiddleware{})
	x = make([]byte, 0)
	buffer = bytes.NewBuffer(x)
	log.SetOutput(buffer) //even now, since we have something registered, we should not be going in
	Info("test")
	if len(buffer.Bytes()) != 0 {
		t.Errorf("unexpected output in log")
	}
}

func _TestStackExtra(t *testing.T) {
	cleanMiddleware()
	Fatal("FATAL TEST")
}
