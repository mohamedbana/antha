// microArch/equipment/manual/cli/cli.go: Part of the Antha language
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

package cli
import (
	"fmt"
	"bufio"
	"os"
	"io"
	"bytes"
	"strings"
)

//CommandRequest represents a request to command line in which a question is made, a set of options are presented and
// one of those options must be selected. Expected is the option which should be selected for the command to execute
// successfully
type CLICommandRequest struct {
	Question	string
	Options		map[int]string
	Expected	int
}
func NewCLICommandRequest(question string, options map[int]string, expected int) *CLICommandRequest {
	cr := new(CLICommandRequest)
	cr.Question = question
	cr.Options = options
	cr.Expected = expected
	return cr
}
//CommandResult is the result of a CommandRequest. It has a bool representing whether the action was performed successfully
// and an Answer holding a string with a textual representation of the result should it be needed
type CLICommandResult struct {
	Result	bool
	Answer	string
}
func NewCLICommandResult(res bool, answer string) *CLICommandResult{
	cr := new(CLICommandResult)
	cr.Answer = answer
	cr.Result = res
	return cr
}
//CommandAsker will take CommandRequests, queue them and execute them, The output of each will be appended to an output
/// CommandResult queue
type CLICommandAsker struct {
	CommandQueue chan CLICommandRequest
	ResultQueue chan CLICommandResult
}

func NewCLICommandAsker() *CLICommandAsker {
	ca := new(CLICommandAsker)
	ca.CommandQueue = make(chan CLICommandRequest)
	ca.ResultQueue = make(chan CLICommandResult)
	return ca
}

func (ca *CLICommandAsker) RunCLI() error {
	reader := os.Stdin
	writer := os.Stdout
	return ca.RunWithReaderWriter(reader, writer)
}

func (ca *CLICommandAsker) Close() error {
	close(ca.CommandQueue)
	close(ca.ResultQueue)
	return nil
}

func (ca *CLICommandAsker) RunWithReaderWriter(reader io.Reader, writer io.Writer) error {
	bReader := bufio.NewReader(reader)
	bWriter := bufio.NewWriter(writer)
	go func(){
		for {  //TODO loop until an answer that matches the question options is given
			command := <-ca.CommandQueue
			bWriter.WriteString(fmt.Sprintf("%s : ", strings.TrimSpace(command.Question)))
			var optString bytes.Buffer
			for p, op := range command.Options {
				optString.WriteString(fmt.Sprintf(op))
				if p == len(command.Options) - 1 {
					optString.WriteString("?\n")
				} else {
					optString.WriteString("/")
				}
			}
			bWriter.WriteString(optString.String())
			bWriter.Flush()
			text, err := bReader.ReadString('\n')
			if err != nil {
				panic(err) //TODO
			}
			text = text[:len(text)-1] // strip the newline char
			res, why := evalAnswer(command, text)
			ca.ResultQueue <- *NewCLICommandResult(res, why)
		}
	}()
	return nil
}

func evalAnswer(req CLICommandRequest, res string) (bool, string) {
//	fmt.Println(res, " != ", req.Options[req.Expected], " -> ", res != req.Options[req.Expected])
	if strings.ToLower(res) != strings.ToLower(req.Options[req.Expected]) {
		return false, fmt.Sprintf("Expecting %s, got %s.", req.Options[req.Expected], res)
	}
	return true, ""
}
