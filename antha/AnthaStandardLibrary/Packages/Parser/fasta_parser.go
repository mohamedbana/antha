// antha/AnthaStandardLibrary/Packages/Parser/fasta_parser.go: Part of the Antha language
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

package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Fasta struct {
	Id   string
	Desc string
	Seq  string
}

func Build_fasta(header string, seq bytes.Buffer) (Record Fasta) {
	fields := strings.SplitN(header, " ", 2)

	var record Fasta

	if len(fields) > 1 {
		record.Id = fields[0]
		record.Desc = "`" + fields[1] + "`"
	} else {
		record.Id = fields[0]
		record.Desc = ""
	}

	record.Seq = seq.String()

	Record = record

	return Record
}

func FastaParse(fastaFh io.Reader) chan Fasta {

	outputChannel := make(chan Fasta)

	scanner := bufio.NewScanner(fastaFh)
	// scanner.Split(bufio.ScanLines)
	header := ""
	var seq bytes.Buffer

	go func() {
		// Loop over the letters in inputString
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			// line := scanner.Text()

			if line[0] == '>' {
				// If we stored a previous identifier, get the DNA string and map to the
				// identifier and clear the string
				if header != "" {
					// outputChannel <- build_fasta(header, seq.String())
					outputChannel <- Build_fasta(header, seq)
					header = ""
					seq.Reset()
				}

				// Standard FASTA identifiers look like: ">id desc"
				header = line[1:]
			} else {
				// Append here since multi-line DNA strings are possible
				seq.WriteString(line)
			}

		}

		outputChannel <- Build_fasta(header, seq)

		// Close the output channel, so anything that loops over it
		// will know that it is finished.
		close(outputChannel)
	}()

	return outputChannel
}

func Fastatocsv(inputfilename string, outputfileprefix string) (csvfile *os.File, status string) {

	status = "incomplete"

	fastaFh, err := os.Open(inputfilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fastaFh.Close()

	//csvfile, err := os.Create(outputfilename)
	csvfile, err = ioutil.TempFile(os.TempDir(), "csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close()

	records := make([][]string, 0)
	seq := make([]string, 0)
	seq = []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)
	for record := range FastaParse(fastaFh) {
		plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"
		if strings.Contains(record.Desc, "Plasmid") || strings.Contains(record.Id, "Circular") || strings.Contains(record.Id, "Vector") {
			plasmidstatus = "TRUE"
		}
		if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
			seqtype = "AA"
		}

		if strings.Contains(record.Desc, "Class:") {
			uptoclass := strings.Index(record.Desc, "Class:")
			prefix := uptoclass + len("class:")
			class = record.Desc[prefix:]
		}
		seq = []string{record.Id, record.Seq, plasmidstatus, seqtype, class}
		records = append(records, seq)
	}

	writer := csv.NewWriter(csvfile)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	status = "complete"
	writer.Flush()
	return csvfile, status
}
