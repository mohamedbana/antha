package blast

import (
	"fmt"
	. "github.com/biogo/ncbi/blast"
	//"github.com/biogo/ncbi/blast_test"
	//"log"
	//"strings"
)

var (
	email = "ucbecrg@ucl.ac.uk"
	//email = "chrisrgrant@gmail.com"
	tool = "blast-example"
	//tool      = "blast-example1"
	params    Parameters
	putparams = PutParameters{Program: "blastn", Database: "nr"}
	getparams GetParameters
	page      = ""
	//query     = "X14032.1"
	//query     = "MSFSNYKVIAMPVLVANFVLGAATAWANENYPAKSAGYNQGDWVASFNFSKVYVGEELGDLNVGGGALPNADVSIGNDTTLTFDIAYFVSSNIAVDFFVGVPARAKFQGEKSISSLGRVSEVDYGPAILSLQYHYDSFERLYPYVGVGVGRVLFFDKTDGALSSFDIKDKWAPAFQVGLRYDLGNSWMLNSDVRYIPFKTDVTGTLGPVPVSTKIEVDPFILSLGASYVF"
	//query   = "atgagtttttctaattataaagtaatcgcgatgccggtgttggttgctaattttgttttgggggcggccactgcatgggcgaatgaaaattatccggcgaaatctgctggctataatcagggtgactgggtcgctagcttcaatttttctaaggtctatgtgggtgaggagcttggcgatctaaatgttggagggggggctttgccaaatgctgatgtaagtattggtaatgatacaacacttacgtttgatatcgcctattttgttagctcaaatatagcggtggatttttttgttggggtgccagctagggctaaatttcaaggtgagaaatcaatctcctcgctgggaagagtcagtgaagttgattacggccctgcaattctttcgcttcaatatcattacgatagctttgagcgactttatccatatgttggggttggtgttggtcgggtgctattttttgataaaaccgacggtgctttgagttcgtttgatattaaggataaatgggcgcctgcttttcaggttggccttagatatgaccttggtaactcatggatgctaaattcagatgtgcgttatattcctttcaaaacggacgtcacaggtactcttggcccggttcctgtttctactaaaattgaggttgatcctttcattctcagtcttggtgcgtcatatgttttctaa"
	retries = 1
	retry   = retries
)

func Blast(query string) (hits []Hit) {

	/*o, err := Example(query, retry, &putparams, &getparams)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(o)*/
	/*
		readcloser, err := RequestWebReadCloser(page, params, tool, email)
		if err != nil {
			log.Panic(err)
		}
	*/

	//
	r, err := Put(query, &putparams, tool, email)
	fmt.Println("RID=", RidString(*r))
	var o *Output
	for k := 0; k < retries; k++ {
		var s *SearchInfo
		s, err = r.SearchInfo(tool, email)
		fmt.Println(s.Status)
		/*for {
			if strings.Contains(s.Status, "WAITING") {
				s, err = r.SearchInfo(tool, email)
			} else {
				continue
			}
		}*/
		fmt.Println("hits?", s.HaveHits)
		if !s.HaveHits {
			continue
		}
		o, err = r.GetOutput(&getparams, tool, email)
		/*fmt.Println(len(o.Iterations))
		fmt.Println(len(o.Iterations[0].Hits))
		fmt.Println(len(o.Iterations[0].Hits[0].Hsps))
		fmt.Println(o.Iterations[0].Hits[0].Id)
		fmt.Println(o.Iterations[0].Hits[0].Def)
		fmt.Println(o.Iterations[0].Hits[0].Accession)
		fmt.Println(o.Iterations[0].Hits[0].Hsps[0].BitScore)
		fmt.Println(o.Iterations[0].Hits[0].Hsps[0].Score)
		for _, hit := range o.Iterations[0].Hits {
			fmt.Println(hit.Hsps[0].EValue)
		}*/
		if err == nil {
			break
		}

	}
	if o == nil {
		return hits
	}
	if len(o.Iterations) == 0 {
		return hits
	}
	if len(o.Iterations[0].Hits) == 0 {
		return hits
	}
	return o.Iterations[0].Hits
}
