package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"

	"github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Graph definition file missing")
		return
	}

	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
	portMap["TypeIISConstructAssemblyMMX"] = make(map[string]bool)
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true
	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = true

	portMap["TypeIISConstructAssemblyMMX"]["TypeIISConstructAssemblyMMX"] = false

	flow.Register("TypeIISConstructAssemblyMMX", TypeIISConstructAssemblyMMX.NewTypeIISConstructAssemblyMMX)

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	log.SetOutput(os.Stderr)

	graph := flow.LoadJSON(flag.Arg(0))
	if graph == nil {
		fmt.Println("empty graph")
		return
	}
	flow.RunNet(graph)

	//<-graph.Ready()

	for _, port := range graph.GetUnboundOutPorts() {
		ch := make(chan execute.ThreadParam)
		graph.SetOutPort(port.Port, ch)
		go func() {
			for a := range ch {
				if err := enc.Encode(&a); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	inPortMap := make(map[string]chan execute.ThreadParam)
	for _, port := range graph.GetUnboundInPorts() {
		ch := make(chan execute.ThreadParam)
		inPortMap[port.Port] = ch
		graph.SetInPort(port.Port, ch)
	}
	go func() {
		for _, ch := range inPortMap {
			defer close(ch)
		}

		for {
			var p execute.JSONBlock
			if err := dec.Decode(&p); err != nil {
				log.Println("Error decoding", err)
				return
			}
			if p.ID == nil { //TODO add error control in JSONBlock unmarshaling??
				log.Println("Error, no ID")
				continue
			}
			for k, v := range p.Values {
				tmp := make(map[string]interface{})
				tmp[k] = v
				sthg, err := json.Marshal(&tmp)
				if err != nil {
					continue
				}
				if _, exists := inPortMap[k]; exists {
					param := execute.ThreadParam{Value: execute.JSONValue{Name: k, JSONString: string(sthg)}, ID: *p.ID, Error: *p.Error}
					inPortMap[k] <- param
				}
			}
		}
	}()

	<-graph.Wait()
}
