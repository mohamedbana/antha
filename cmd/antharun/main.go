// /antharun/main.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	logFile        string
	parametersFile string
	workflowFile   string
)

func initWorkflow(*Workflow) {}

func run() error {
	wfData, err := ioutil.ReadFile(workflowFile)
	if err != nil {
		return err
	}

	wf, err := NewWorkflow(wfData)
	if err != nil {
		return err
	}

	cfData, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		return err
	}

	cf, err := NewConfig(cfData, wf)
	if err != nil {
		return err
	}

	fmt.Println("Press [Enter] to load antha workflow with manual driver...")
	fmt.Println("Reminder: press [Control-X] to exit the workflow interface")
	if _, err := fmt.Scanln(); err != nil {
		return err
	}

	fe, err := NewFrontend()
	if err != nil {
		return err
	}
	defer fe.Shutdown()

	msgs, err := wf.Run(cf)
	if err != nil {
		return err
	}

	for _, m := range msgs {
		fe.SendAlert(fmt.Sprintf("Output: %v\n", m))
	}

	return nil
}

func main() {
	flag.StringVar(&parametersFile, "parameters", "", "parameters to workflow")
	flag.StringVar(&workflowFile, "workflow", "", "workflow definition file")
	flag.StringVar(&logFile, "log", "", "log file")
	flag.Parse()

	if len(parametersFile) == 0 || len(workflowFile) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
