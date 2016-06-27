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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	parser "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	biogo "github.com/biogo/ncbi/entrez"
)

var (
	email   = "no-reply@antha-lang.com"
	tool    = "entrez-biogo-antha"
	retries = 5
)

// This queries the selected database saving the record to file
// Database options are nucleotide, Protein, Gene. For full list see http://www.ncbi.nlm.nih.gov/books/NBK25497/table/chapter2.T._entrez_unique_identifiers_ui/?report=objectonly
// Return type includes but must match the database type. See http://www.ncbi.nlm.nih.gov/books/NBK25499/table/chapter4.T._valid_values_of__retmode_and/?report=objectonly
// Query can be any string but it is recommended to use GI number if one specific record is requred.
func RetrieveRecords(query string, database string, Max int, ReturnType string, out string) error {
	// query database

	h := biogo.History{}
	s, err := biogo.DoSearch(database, query, nil, &h, tool, email)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%d records found for your query.\n", s.Count)

	var of *os.File
	if out == "" {
		of = os.Stdout
	} else {
		/*if err := os.Mkdir(filepath.Dir(out), 0777); err != nil {
			return err
		}*/
		of, err = os.Create(filepath.Join(anthapath.Path(), out))
		if err != nil {
			return err
		}
		defer of.Close()
	}

	var (
		buf   = &bytes.Buffer{}
		p     = &biogo.Parameters{RetMax: Max, RetType: ReturnType, RetMode: "text"}
		bn, n int64
	)

	for p.RetStart = 0; p.RetStart < s.Count; p.RetStart += p.RetMax {
		fmt.Fprintf(os.Stdout, "Attempting to retrieve %d record(s).\n", p.RetMax)
		var t int
		for t = 0; t < retries; t++ {
			buf.Reset()
			s := time.Duration(1) * time.Second // limit queries to < 3 per second
			time.Sleep(s)

			var (
				r   io.ReadCloser
				_bn int64
			)
			r, err = biogo.Fetch(database, p, tool, email, &h)
			if err != nil {
				if r != nil {
					r.Close()
				}
				fmt.Fprintf(os.Stdout, "Failed to retrieve on attempt %d... error: %v retrying.\n", t, err)
				continue
			}
			_bn, err = io.Copy(buf, r)
			bn += _bn
			r.Close()
			if err == nil {
				break
			}
			fmt.Fprintf(os.Stdout, "Failed to buffer on attempt %d... error: %v retrying.\n", t, err)
		}
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Retrieved records with %d retries... writing out.\n", t)
		_n, err := io.Copy(of, buf)
		n += _n
		if err != nil {
			return err
		}

	}
	if bn != n {
		fmt.Fprintf(os.Stdout, "Writethrough mismatch: %d != %d\n", bn, n)
	}

	return nil
}

// This retrieves sequence of any type from any NCBI sequence database
func RetrieveSequence(id string, database string, filename string) (wtype.DNASequence, error) {
	if err := RetrieveRecords(id, database, 1, "gb", filename); err != nil {
		fmt.Println("RetrieveRecordsfail for", id)
		fmt.Println(err.Error())
		return wtype.DNASequence{}, err
	}

	contents, err := ioutil.ReadFile(filename)

	fmt.Println("ID:", id, "Contents:", string(contents))
	file := filepath.Join(anthapath.Path(), filename)
	seq, err := parser.GenbanktoFeaturelessDNASequence(file)
	if err != nil {
		fmt.Println("File:", file, "Error:", err.Error())
		return wtype.DNASequence{}, err
	}
	seq.Seq = strings.ToUpper(seq.Seq)

	return seq, nil
}

// This will retrieve vector using fasta or db
func RetrieveVector(id string) (wtype.DNASequence, error) {
	/*//first check if vector sequence is in fasta file
	if seq, err := parser.RetrieveSeqFromFASTA(id, filepath.Join(anthapath.Path(), "vectors.txt")); err != nil {
		// if not in refactor, check db*/
	return RetrieveSequence(id, "nucleotide", id+".gb")
	/*} else {
		return seq, nil
	}*/
}
