// microArch/equipment/manual/cli/cli_test.go: Part of the Antha language
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
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRunWithReaderWriter(t *testing.T) {
	type SimpleTest struct {
		req   CLICommandRequest
		res   CLICommandResult
		input string
	}
	testSuite := make([]SimpleTest, 0)
	op1 := make(map[int]string, 2)
	op1[0] = "YES"
	op1[1] = "NO"
	testSuite = append(testSuite, SimpleTest{
		*NewCLICommandRequest("OK", op1, 0),
		*NewCLICommandResult(true, ""),
		"YES",
	})
	testSuite = append(testSuite, SimpleTest{
		*NewCLICommandRequest("OK??", op1, 0),
		*NewCLICommandResult(false, "Expecting YES, got NO."),
		"NO",
	})

	for _, st := range testSuite {
		//initialize our commandAsker
		ca := NewCLICommandAsker()

		reader := strings.NewReader(fmt.Sprintf("%s\n", st.input))
		writer := new(bytes.Buffer)

		ca.RunWithReaderWriter(reader, writer)
		ca.CommandQueue <- st.req
		res := <-ca.ResultQueue
		//		fmt.Println(res)
		if res.Result != st.res.Result {
			t.Errorf("With Command Test %v. Expecting result %v. got %v.", st.req, st.res.Result, res.Result)
		}
		if res.Answer != st.res.Answer {
			t.Errorf("With Command Test %v. Expecting answer %v. got %v.", st.req, st.res.Answer, res.Answer)
		}
	}
}

//TODO write more testcases, edgecases for empty options and so on...
