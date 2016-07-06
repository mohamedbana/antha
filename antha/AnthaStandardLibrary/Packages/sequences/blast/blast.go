// Part of the Antha language
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

// Package for performing blast queries
package blast

import (
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	. "github.com/biogo/ncbi/blast"
	"github.com/mgutz/ansi"
)

// package for interacting with the ncbi BLAST service

var (
	email     = "no-reply@antha-lang.com"
	tool      = "blast-biogo-antha"
	params    Parameters
	putparams = PutParameters{Program: "blastn", Megablast: true, Database: "nr"}
	getparams GetParameters
	page      = ""
	//query     = "X14032.1"
	//query     = "MSFSNYKVIAMPVLVANFVLGAATAWANENYPAKSAGYNQGDWVASFNFSKVYVGEELGDLNVGGGALPNADVSIGNDTTLTFDIAYFVSSNIAVDFFVGVPARAKFQGEKSISSLGRVSEVDYGPAILSLQYHYDSFERLYPYVGVGVGRVLFFDKTDGALSSFDIKDKWAPAFQVGLRYDLGNSWMLNSDVRYIPFKTDVTGTLGPVPVSTKIEVDPFILSLGASYVF"
	//query   = "atgagtttttctaattataaagtaatcgcgatgccggtgttggttgctaattttgttttgggggcggccactgcatgggcgaatgaaaattatccggcgaaatctgctggctataatcagggtgactgggtcgctagcttcaatttttctaaggtctatgtgggtgaggagcttggcgatctaaatgttggagggggggctttgccaaatgctgatgtaagtattggtaatgatacaacacttacgtttgatatcgcctattttgttagctcaaatatagcggtggatttttttgttggggtgccagctagggctaaatttcaaggtgagaaatcaatctcctcgctgggaagagtcagtgaagttgattacggccctgcaattctttcgcttcaatatcattacgatagctttgagcgactttatccatatgttggggttggtgttggtcgggtgctattttttgataaaaccgacggtgctttgagttcgtttgatattaaggataaatgggcgcctgcttttcaggttggccttagatatgaccttggtaactcatggatgctaaattcagatgtgcgttatattcctttcaaaacggacgtcacaggtactcttggcccggttcctgtttctactaaaattgaggttgatcctttcattctcagtcttggtgcgtcatatgttttctaa"
	retries = 1
	retry   = retries
)

func RerunRIDstring(rid string) (o *Output, err error) {

	r := NewRid(rid)

	if r != nil {
		fmt.Println("RID=", r.String())

		//var o *Output
		for k := 0; k < retries; k++ {
			var s *SearchInfo
			s, err = r.SearchInfo(tool, email)
			fmt.Println(s.Status)

			fmt.Println("hits?", s.HaveHits)
			if !s.HaveHits {
				continue
			}
			o, err = r.GetOutput(&getparams, tool, email)

			if err == nil {
				break
			}

		}
	} else {
		err = fmt.Errorf("r == nil")
	}

	return
}

func RerunRID(r *Rid) (o *Output, err error) {

	if r != nil {
		fmt.Println("RID=", r.String())

		//var o *Output
		for k := 0; k < retries; k++ {
			var s *SearchInfo
			s, err = r.SearchInfo(tool, email)
			fmt.Println(s.Status)

			fmt.Println("hits?", s.HaveHits)
			if !s.HaveHits {
				continue
			}
			o, err = r.GetOutput(&getparams, tool, email)

			if err == nil {
				break
			}

		}
	} else {
		err = fmt.Errorf("r == nil")
	}

	return
}

func HitSummary(hits []Hit) (summary string, err error) {

	if len(hits) != 0 {

		/*for _, hit := range hits {
			for _, hsp := range hit.Hsps {
				fmt.Printf("%+v", hsp)

			}
		}*/
		seqlength := hits[0].Len

		identity := strconv.Itoa((*hits[0].Hsps[0].HspIdentity/len(hits[0].Hsps[0].QuerySeq))*100) + "%"
		coverage := strconv.Itoa(len(hits[0].Hsps[0].QuerySeq)/len(hits[0].Hsps[0].SubjectSeq)*100) + "%"

		summary = fmt.Sprintln(ansi.Color("Hits: ", "red"), len(hits),
			//	Printfield(hits[0].Id),
			text.Print("HspIdentity: ", strconv.Itoa(*hits[0].Hsps[0].HspIdentity)),
			text.Print("queryLen: ", len(hits[0].Hsps[0].QuerySeq)),
			text.Print("subjectLen: ", len(hits[0].Hsps[0].SubjectSeq)),
			text.Print("alignLen: ", *hits[0].Hsps[0].AlignLen),
			text.Print("Identity: ", identity),
			text.Print("coverage: ", coverage),
			//Print("HspIdentity", *hits[0].Hsps[0].HspIdentity),
			ansi.Color("Sequence length:", "red"), seqlength,
			ansi.Color("high scoring pairs for top match:", "red"), len(hits[0].Hsps),
			ansi.Color("Id:", "red"), hits[0].Id,
			ansi.Color("Definition:", "red"), hits[0].Def,
			ansi.Color("Accession:", "red"), hits[0].Accession,
			//ansi.Color("Identity: ", "red"), identity, "%",
			//ansi.Color("Coverage: ", "red"), coverage, "%",
			ansi.Color("Bitscore", "red"), hits[0].Hsps[0].BitScore,
			ansi.Color("Score", "red"), hits[0].Hsps[0].Score,
			ansi.Color("EValue", "red"), hits[0].Hsps[0].EValue)
	} else {
		summary = "No hits!"
		err = fmt.Errorf(summary)
	}
	return
}

func MegaBlastP(query string) (hits []Hit, err error) {

	putparams = PutParameters{Program: "blastp", Megablast: true, Database: "nr"}

	o, err := SimpleBlast(query)
	if err != nil {
		return
	}
	hits, err = Hits(o)
	if err != nil {
		return
	}

	return
}

func MegaBlastN(query string) (hits []Hit, err error) {

	putparams = PutParameters{Program: "blastn", Megablast: true, Database: "nr"}

	o, err := SimpleBlast(query)
	if err != nil {
		return
	}
	hits, err = Hits(o)
	//fmt.Println(hits)
	if err != nil {
		return
	}
	return
}

func SimpleBlast(query string) (o *Output, err error) {

	r, err := Put(query, &putparams, tool, email)
	fmt.Println("RID=", r.String())
	fmt.Println("Submitting request to BLAST server, please wait")
	//var o *Output
	for k := 0; k < retries; k++ {
		var s *SearchInfo
		s, err = r.SearchInfo(tool, email)
		fmt.Println(s.Status)

		fmt.Println("hits?", s.HaveHits)
		if s.HaveHits == true {
			o, err = r.GetOutput(&getparams, tool, email)
			continue
		} else if strings.Contains(s.Status, "WAITING") == true {
			for {
				if strings.Contains(s.Status, "WAITING") == true {
					fmt.Println("waiting 1 min to rerun RID:", r.String())
					time.Sleep(1 * time.Minute)
					s, err = r.SearchInfo(tool, email)
					o, err = RerunRID(r)

				} else {
					return
				}
			}

			if err == nil {
				break
			}

		}

	}

	return
}
func Hits(o *Output) (hits []Hit, err error) {

	if o == nil {
		err = fmt.Errorf("output == nil")
		return
	}
	if len(o.Iterations) == 0 {
		err = fmt.Errorf("len(output.Iterations) == 0")
		return
	}
	if len(o.Iterations[0].Hits) == 0 {
		err = fmt.Errorf("len(output.Iterations[0].Hits) == 0")
		return
	}
	hits = o.Iterations[0].Hits

	return
}
