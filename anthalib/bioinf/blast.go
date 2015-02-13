// anthalib//bioinf/blast.go: Part of the Antha language
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

package bioinf

import (
	"os/exec"
	"fmt"
	"log"
	"bufio"
	"strings"
	"regexp"
	"github.com/antha-lang/antha/anthalib/wutil"
	"github.com/antha-lang/antha/anthalib/wtype"
	"code.google.com/p/biogo/io/seqio/fasta"
	"code.google.com/p/biogo/seq/linear"
	"code.google.com/p/biogo/seq"
)

func parse_blast_results(results string)wtype.BlastResults{
	ret:=wtype.NewBlastResults()
	tx:=strings.Split(results, "\n")
	parseScores:=false
	align_start:=-1

	for index,line:=range tx{
		lettline,_:=regexp.MatchString("total letters", line)
		if index==0{
			ret.Program=line
		}else if strings.HasPrefix(line, " Query="){
			tx2:=strings.Split(line, "= ")
			ret.Query=tx2[1]
		}else if strings.HasPrefix(line, "Database:"){
			tx2:=strings.Split(line, ": ")
			ret.DBname=tx2[1]
		}else if lettline{
			rx,_:=regexp.Compile("(\\d+,?)+\\s")
			loc:=rx.FindStringIndex(line)

			if loc!=nil{
				ret.DBSizeSeqs=wutil.ParseInt(line[loc[0]:loc[1]-1])

				loc=rx.FindStringIndex(line[loc[1]:len(line)-1])

				if loc!=nil{
					ret.DBSizeLetters=wutil.ParseInt(line[loc[0]:loc[1]-1])
				}

			}
		}else if strings.HasPrefix(line, "Sequences producing"){
			parseScores=true
		}else if(parseScores && strings.HasPrefix(line, ">")){
			// start of alignment section
			align_start=index
			break
		}else if(parseScores){
			mt,_:=regexp.MatchString("^\\s*$", line)

			if(mt){
				continue
			}
			hit:=NewBlastHit()

			hit.Name=line[0:69]
			hit.Score=wutil.ParseFloat(line[69:73])
			hit.Eval=wutil.ParseFloat(line[73:len(line)])
			ret.Hits=append(ret.Hits, hit)
		}
	}

	// now we move on to parsing the alignments
	index:=align_start
	alignmentcount:=0
	for(true){
		as:=NewAlignedSequence()
		for i:=index;i<len(tx);i++{
			// this means we've encountered another hit
			if((strings.HasPrefix(tx[i], ">") && as.Qstart!=-1) || strings.HasPrefix(tx[i], "  Database")){
				//fmt.Println(alignmentcount, " ", len(ret.Hits), " ", tx[i])
				ret.Hits[alignmentcount].Alignments=append(ret.Hits[alignmentcount].Alignments, as)
				alignmentcount+=1
				index=i
				break
			}

			if strings.HasPrefix(tx[i], " Score"){
				// if this alignment is initialized, time to get a new one 

				if as.Qstart!=-1{
					ret.Hits[alignmentcount].Alignments=append(ret.Hits[alignmentcount].Alignments, as)
					as=NewAlignedSequence()
				}
			}else if strings.HasPrefix(tx[i], "Query"){
				tx2:=strings.Split(tx[i], " ")
				if as.Qstart==-1{
					as.Qstart=wutil.ParseInt(tx2[1])
				}
				as.Qend=wutil.ParseInt(tx2[3])
				as.Qseq+=tx2[2]
			}else if strings.HasPrefix(tx[i], "Sbjct"){
				tx2:=strings.Split(tx[i], " ")

				if as.Sstart==-1{
					as.Sstart=wutil.ParseInt(tx2[1])
				}
				as.Send=wutil.ParseInt(tx2[3])
				as.Sseq+=tx2[2]
			}else if strings.HasPrefix(tx[i], " Identities"){
				rx,_:=regexp.Compile("\\d+%")
				match:=rx.FindStringIndex(tx[i])
				as.ID=wutil.ParseFloat(tx[i][match[0]:match[1]-1])
			}else if strings.HasPrefix(tx[i], " Strand"){
				tx2:=strings.Split(tx[i], " = ")
				tx3:=strings.Split(tx2[1], " / ")
				as.Qstrand=tx3[0]
				as.Sstrand=tx3[1]
			}
		}

		if(strings.HasPrefix(tx[index], "  Database")){
			break
		}
	}

	return ret
}

func RunBlastAndReturnResults(seq wutil.Bioseq)wutilBlastResults{
	res:=RunBlast(seq)
	blast_results:=parse_blast_results(res)
	return blast_results
}

func RunBlast(seq wutil.Bioseq)string{
	// TODO wrap these as environment requests for blast services
	// what we need is a service which wraps local vs. nonlocal 
	blast:="/Users/msadowski/software/blast-2.2.26_mac/bin/blastall"
	blastdb:="/Users/msadowski/data/patent_seqs/nrnl1"

	seqname:=wutil.makeseq("/tmp", "test", seq)

	// object lesson in using os.exec
	cmd:=exec.Command(blast, fmt.Sprintf("-d%s", blastdb), "-pblastn", fmt.Sprintf("-i%s", seqname))

	out,err:=cmd.CombinedOutput()

	// TODO improve error handling
	if err!=nil{
		log.Fatal(err)
	}
	return string(out)
}