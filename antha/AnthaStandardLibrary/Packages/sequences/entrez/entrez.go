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

// package for querying all of NCBI databases
package entrez

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	biogo "github.com/antha-lang/antha/internal/github.com/biogo/ncbi/entrez"
)

type Fasta struct {
	Id   string
	Desc string
	Seq  string
}

var (
	email   = "no-reply@antha-lang.com"
	tool    = "biogo.example"
	h       = biogo.History{}
	retries = 5
)

// This queries the selected database saving the record to file
// Database options are nucleotide, Protein, Gene. For full list see http://www.ncbi.nlm.nih.gov/books/NBK25497/table/chapter2.T._entrez_unique_identifiers_ui/?report=objectonly
// Return type includes but must match the database type. See http://www.ncbi.nlm.nih.gov/books/NBK25499/table/chapter4.T._valid_values_of__retmode_and/?report=objectonly
// Query can be any string but it is recommended to use GI number if one specific record is requred.
func RetrieveRecords(query string, database string, Max int, ReturnType string, out string) {
	// query database
	s, err := biogo.DoSearch(database, query, nil, &h, tool, email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%d records found for your query.\n", s.Count)

	var of *os.File
	if out == "" {
		of = os.Stdout
	} else {
		out = fmt.Sprintf("%s%c%s", anthapath.Dirpath(), os.PathSeparator, out)
		of, err = os.Create(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer of.Close()
	}

	var (
		buf   = &bytes.Buffer{}
		p     = &biogo.Parameters{RetMax: Max, RetType: ReturnType, RetMode: "text"}
		bn, n int64
		//r     io.ReadCloser
		//_bn   int64
	)

	for i := 0; i < p.RetMax; i++ {
		fmt.Fprintf(os.Stderr, "Attempting to retrieve %d record(s).\n", p.RetMax)
		var t int
		for t = 0; t < retries; t++ {
			buf.Reset()
			var (
				r   io.ReadCloser
				_bn int64
			)
			r, err = biogo.Fetch(database, p, tool, email, &h)
			if err != nil {
				if r != nil {
					r.Close()
				}
				fmt.Fprintf(os.Stderr, "Failed to retrieve on attempt %d... error: %v retrying.\n", t, err)
				continue
			}
			_bn, err = io.Copy(buf, r)
			bn += _bn
			r.Close()
			if err == nil {
				break
			}
			fmt.Fprintf(os.Stderr, "Failed to buffer on attempt %d... error: %v retrying.\n", t, err)
		}
		if err != nil {
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Retrieved records with %d retries... writing out.\n", t)
		//x := buf.String()
		//y = wtype.MakeLinearDNASequence(query, x)
		_n, err := io.Copy(of, buf)
		n += _n
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	}
	if bn != n {
		fmt.Fprintf(os.Stderr, "Writethrough mismatch: %d != %d\n", bn, n)
	}

}
