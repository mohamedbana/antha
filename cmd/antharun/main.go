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
// 2 Royal College St, London NW1 0NH UK

package main

import (
	"flag"
	"fmt"
	"github.com/antha-lang/antha/antha/component/lib"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
)

var (
	logFile        string
	parametersFile string
	workflowFile   string
	driverURI      string
	list           bool
)

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range lib.GetComponents() {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func run() error {
	if driverURI != "" {
		fe, err := NewRemoteFrontend(driverURI)
		if err != nil {
			return err
		}
		defer fe.Shutdown()
	} else {
		fmt.Println("Press [Enter] to load antha workflow with manual driver...")
		fmt.Println("Reminder: press [Control-X] to exit the workflow interface")
		if _, err := fmt.Scanln(); err != nil {
			return err
		}

		fe, err := NewCUIFrontend()
		if err != nil {
			return err
		}
		defer fe.Shutdown()
	}

	wdata, err := ioutil.ReadFile(workflowFile)
	if err != nil {
		return err
	}

	pdata, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		return err
	}

	ctx, err := makeContext()
	if err != nil {
		return err
	}

	w, err := execute.Run(ctx, execute.Options{
		WorkflowData: wdata,
		ParamData:    pdata,
	})
	if err != nil {
		return err
	}

	for k, v := range w.Outputs {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
}

func main() {
	flag.BoolVar(&list, "list", false, "list the available components")
	flag.StringVar(&parametersFile, "parameters", "parameters.yml", "parameters to workflow")
	flag.StringVar(&workflowFile, "workflow", "workflow.json", "workflow definition file")
	flag.StringVar(&logFile, "log", "", "log file")
	flag.StringVar(&driverURI, "driver", "", "uri where a grpc driver implementation listens")
	flag.Parse()

	if list {
		fmt.Println("Available Components")
		fmt.Println("====================")
		for _, v := range lib.GetComponents() {
			fmt.Printf("\t %q.\n", v.Name)
		}
		fmt.Println("====================")
	}

	if len(parametersFile) == 0 || len(workflowFile) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if driverURI == "" {
		if envDriver := os.Getenv("DRIVERURI"); len(envDriver) > 0 {
			driverURI = envDriver
		}
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
