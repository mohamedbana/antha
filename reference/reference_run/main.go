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
/*
Sample JSON input:
{
    "Color": "white",
    "SleepTime": 100,
	"ID": 1
}
*/
package main

import (
	"encoding/json"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/reference"
	"github.com/antha-lang/goflow"
	"log"
	"os"
)

var (
	exitCode = 0
)

type ExampleApp struct {
	flow.Graph
}

func NewExampleApp() *ExampleApp {
	n := new(ExampleApp)
	n.InitGraphState()

	n.Add(reference.NewExample(), "example")

	n.MapInPort("Color", "example", "Color")
	n.MapInPort("SleepTime", "example", "SleepTime")
	n.MapOutPort("WellColor", "example", "WellColor")
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
	net := NewExampleApp()

	Colors := make(chan execute.ThreadParam)
	SleepTimes := make(chan execute.ThreadParam)
	WellColors := make(chan execute.ThreadParam)
	done := make(chan bool)

	net.SetInPort("Color", Colors)
	net.SetInPort("SleepTime", SleepTimes)
	net.SetOutPort("WellColor", WellColors)

	flow.RunNet(net)

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	log.SetOutput(os.Stderr)

	go func() {
		defer close(Colors)
		defer close(SleepTimes)
		for {
			var p reference.JSONBlock

			if err := dec.Decode(&p); err != nil {
				log.Println("Error decoding", err)
				return
			}
			log.Println(p)
			/* if no ID, print an error and skip this json object */
			if p.ID == nil {
				log.Println("Param without ID:", p)
				continue
			}
			if p.Color != nil {
				param := execute.ThreadParam{*(p.Color), *(p.ID)}
				Colors <- param
			}
			if p.SleepTime != nil {
				param := execute.ThreadParam{*(p.SleepTime), *(p.ID)}
				SleepTimes <- param
			}
		}
	}()

	go func() {
		defer close(done)
		for result := range WellColors {
			if err := enc.Encode(&result); err != nil {
				log.Println(err)
			}
		}
	}()

	<-net.Wait()
}
