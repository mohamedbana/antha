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
	"encoding/json"
	"flag"
	"fmt"
	"github.com/antha-lang/antha/antha/component/lib"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target"
	"io/ioutil"
	"log"
	"os"
)

const (
	jsonOutput   = "json"
	stringOutput = "pretty"
	defaultPort  = 50051
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

type opts struct {
	Frontend       string
	DriverURI      string
	ParametersFile string
	WorkflowFile   string
}

func runOne(opts opts) error {
	t := target.New()
	fe, err := NewFrontend(Options{
		Kind:   opts.Frontend,
		Target: t,
		URI:    opts.DriverURI,
	})
	if err != nil {
		return err
	}
	defer fe.Shutdown()

	wdata, err := ioutil.ReadFile(opts.WorkflowFile)
	if err != nil {
		return err
	}

	pdata, err := ioutil.ReadFile(opts.ParametersFile)
	if err != nil {
		return err
	}

	ctx, err := makeContext()
	if err != nil {
		return err
	}

	w, err := execute.Run(ctx, execute.Options{
		Target:       t,
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

type port struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
}

type component struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	InPorts     []port `json:"in_ports"`
	OutPorts    []port `json:"out_ports"`
}

func getComponents() ([]component, error) {
	var cs []component
	for _, v := range lib.GetComponents() {
		c := component{
			Id:          v.Name,
			Name:        v.Name,
			Description: v.Desc.Desc,
			Path:        v.Desc.Path,
		}
		for _, p := range v.Desc.Params {
			port := port{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Desc,
				Kind:        p.Kind,
			}
			switch p.Kind {
			case "Outputs", "Data":
				c.OutPorts = append(c.OutPorts, port)
			case "Inputs", "Parameters":
				c.InPorts = append(c.InPorts, port)
			default:
				return nil, fmt.Errorf("unknown parameter kind %q", p.Kind)
			}
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func printComponents(output string, cs []component) error {
	switch output {
	case jsonOutput:
		bs, err := json.Marshal(cs)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))
	default:
		for _, c := range cs {
			fmt.Printf("%s:\n", c.Name)
			fmt.Printf("\tInputs:\n")
			for _, p := range c.InPorts {
				fmt.Printf("\t\t%s %s\n", p.Name, p.Type)
			}
			fmt.Printf("\tOutputs:\n")
			for _, p := range c.OutPorts {
				fmt.Printf("\t\t%s %s\n", p.Name, p.Type)
			}
		}
	}
	return nil
}

func run() error {
	list := flag.Bool("list", false, "list the available components")
	output := flag.String("output", "pretty", "output format one of {json, pretty}")
	parametersFile := flag.String("parameters", "parameters.yml", "parameters to workflow")
	workflowFile := flag.String("workflow", "workflow.json", "workflow definition file")
	driver := flag.String("driver", "", "uri where remote grpc driver implementation listens")
	package_ := flag.String("package", "", "go package to spawn grpc driver from")
	frontend := flag.String("frontend", "debug", "kind of frontend one of {debug, cui, remote}")
	flag.Parse()

	if *list {
		if cs, err := getComponents(); err != nil {
			log.Fatal(err)
		} else if err := printComponents(*output, cs); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	if len(*driver) == 0 {
		if e := os.Getenv("DRIVERURI"); len(e) > 0 {
			*driver = e
		}
	}

	if len(*parametersFile) == 0 || len(*workflowFile) == 0 {
		return fmt.Errorf("missing parameters and/or workflow command line argument")
	}

	if len(*driver) > 0 && len(*package_) > 0 {
		return fmt.Errorf("command line argument driver is not compatible with package")
	}

	if len(*package_) > 0 {
		s, err := spawn(*package_, defaultPort)
		if s != nil {
			defer s.Close()
		}
		if err != nil {
			return err
		}
		if err := s.Command.Start(); err != nil {
			return err
		}
		defer func() {
			s.Command.Process.Kill()
		}()
		*driver = s.URI
	}

	if len(*driver) > 0 && *frontend == DEBUG {
		*frontend = REMOTE
	}

	if err := runOne(opts{
		Frontend:       *frontend,
		DriverURI:      *driver,
		ParametersFile: *parametersFile,
		WorkflowFile:   *workflowFile,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
