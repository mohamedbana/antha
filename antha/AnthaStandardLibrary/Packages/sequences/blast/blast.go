package blast

import (
	"fmt"
	. "github.com/antha-lang/antha/internal/github.com/biogo/ncbi/blast"
	//"github.com/biogo/ncbi/blast_test"
	//"log"
	"github.com/antha-lang/antha/internal/github.com/mgutz/ansi"
	//"reflect"
	"strconv"
	"strings"
	"time"
)

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

func Print(desc string, value interface{}) (fmtd string) {

	fmtd = fmt.Sprint(ansi.Color(desc, "red"), value)
	return
}

/*
func Printfield( value interface{}) (fmtd string) {

	switch myValue := value.(type){
		case string:
		fmt.Println(myValue)
		case Hit:
		fmt.Printf("%+v\n", myValue)
		default:
		fmt.Println("Type not handled: ", reflect.TypeOf(value))
	}

	//a := &Hsp{Len: "afoo"}
	val := reflect.Indirect(reflect.ValueOf(value))
	desc := fmt.Sprint(val.Type().Field(0).Name)

	fmtd = fmt.Sprint(ansi.Color(desc, "red"), value)
	return
}
*/
func HitSummary(hits []Hit) (summary string, err error) {

	if len(hits) != 0 {

		/*for _, hit := range hits {
			for _, hsp := range hit.Hsps {
				fmt.Printf("%+v", hsp)

			}
		}*/
		seqlength := hits[0].Len

		identity := strconv.Itoa((*hits[0].Hsps[0].HspIdentity / len(hits[0].Hsps[0].QuerySeq)) * 100)  //+ "%"
		coverage := strconv.Itoa(len(hits[0].Hsps[0].QuerySeq) / len(hits[0].Hsps[0].SubjectSeq) * 100) // + "%"

		summary = fmt.Sprintln(ansi.Color("Hits: ", "red"), len(hits),
			//	Printfield(hits[0].Id),
			/*	Print("HspIdentity: ", strconv.Itoa(*hits[0].Hsps[0].HspIdentity)),
				Print("queryLen: ", len(hits[0].Hsps[0].QuerySeq)),
				Print("subjectLen: ", len(hits[0].Hsps[0].SubjectSeq)),
				Print("alignLen: ", *hits[0].Hsps[0].AlignLen),
				Print("Identity: ", identity),
				Print("coverage: ", coverage),*/
			//Print("HspIdentity", *hits[0].Hsps[0].HspIdentity),
			ansi.Color("Sequence length:", "red"), seqlength,
			ansi.Color("high scoring pairs for top match:", "red"), len(hits[0].Hsps),
			ansi.Color("Id:", "red"), hits[0].Id,
			ansi.Color("Definition:", "red"), hits[0].Def,
			ansi.Color("Accession:", "red"), hits[0].Accession,
			ansi.Color("Identity: ", "red"), identity, "%",
			ansi.Color("Coverage: ", "red"), coverage, "%",
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
