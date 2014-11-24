// antha/reference/reference_run/main.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
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

package main

import (
	"github.com/antha-lang/antha/reference"
	"fmt"
	"github.com/Synthace/goflow"
	"os"
	"time"
)

var (
	exitCode = 0
)

type Printer struct {
	flow.Component
	Line <-chan reference.ThreadParam
}

func (p *Printer) OnLine(line reference.ThreadParam) {
	fmt.Println(line.ID)
	fmt.Println(line.Value.(string))
}

type ExampleApp struct {
	flow.Graph
}

func NewExampleApp() *ExampleApp {
	n := new(ExampleApp)
	n.InitGraphState()

	n.Add(reference.NewExample(), "example")
	n.Add(new(Printer), "printer")

	n.Connect("example", "WellColor", "printer", "Line")

	n.MapInPort("Color", "example", "Color")
	n.MapInPort("SleepTime", "example", "SleepTime")
	return n
}

func main() {
	// call gofmtMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	referenceMain()
	os.Exit(exitCode)
}

func referenceMain() {
	var color string = "White"
	var sleep time.Duration = 1
	net := NewExampleApp()

	colors := make(chan reference.ThreadParam)
	sleeps := make(chan reference.ThreadParam)

	net.SetInPort("Color", colors)
	net.SetInPort("SleepTime", sleeps)

	flow.RunNet(net)

	go func() {
		colors <- reference.ThreadParam{color, "1"}
		colors <- reference.ThreadParam{color, "2"}
		colors <- reference.ThreadParam{color, "3"}
		close(colors)
	}()

	// intentionally unbalanced
	go func() {
		sleeps <- reference.ThreadParam{sleep, "1"}
		sleeps <- reference.ThreadParam{sleep, "2"}
		sleeps <- reference.ThreadParam{sleep, "3"}
		close(sleeps)
	}()

	<-net.Wait()

}