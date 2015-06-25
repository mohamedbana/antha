// /component/func.go: Part of the Antha language
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

package component

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
)

// Runs the component graph with json input and producing json output
func Run(js []byte, dec *json.Decoder, enc *json.Encoder, errChan chan<- error) <-chan int {
	done := make(chan int)
	graph, err := flow.ParseJSON(js)
	if err != nil {
		errChan <- err
		close(errChan)
		done <- 1
		return done
	}

	flow.RunNet(graph)

	for _, port := range graph.GetUnboundOutPorts() {
		ch := make(chan execute.ThreadParam)
		graph.SetOutPort(port.Port, ch)
		go func() {
			for a := range ch {
				if err := enc.Encode(&a); err != nil {
					errChan <- err
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
				if err != io.EOF {
					errChan <- err
				}
				return
			}
			if p.ID == nil { //TODO add error control in JSONBlock unmarshaling??
				errChan <- errors.New("No ID")
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
					param := execute.ThreadParam{
						Value: execute.JSONValue{Name: k, JSONString: string(sthg)},
						ID:    *p.ID,
						Error: *p.Error,
					}
					inPortMap[k] <- param
				}
			}
		}
	}()

	go func() {
		<-graph.Wait()
		close(errChan)
		done <- 1
	}()

	return done
}
