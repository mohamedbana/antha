// anthalib//wutil/goldengate.go: Part of the Antha language
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

package wutil

import (
	"github.com/antha-lang/antha/anthalib/wtype"
	"math/rand"
	"fmt"
)

// returns the same list of strings in the same order minus any duplicates
// i.e. ordered by first appearance
func MakeUnique(a []string)[]string{
	r:=make([]string, 0, len(a))

	m:=make(map[string]int,len(a))

	for _,s:=range a{
		_,ok:=m[s]

		if ok{
			continue
		}

		m[s]=1
		r=append(r,s)
	}

	return r
}


func Rev(s string)string{
	r:=""

	for i:=len(s)-1;i>=0;i--{
		r+=string(s[i])
	}

	return r
}

func Comp(s string)string{
	r:=""

	m:=map[string]string{
		"A":"T",
		"T":"A",
		"U":"A",
		"C":"G",
		"G":"C",
	}

	for _,c:=range s{
		r+=m[string(c)]
	}

	return r
}

func RevComp(s string)string{
	return Comp(Rev(s))
}

func random_dna_seq(leng int)string{
	s:=""
	for i:=0;i<leng;i++{
		s+=random_char("ACTG")
	}
	return s
}

func all_dna_seqs_with_length(l int)[]string{
	if l==0{
		return []string {""}
	}

	s:=all_dna_seqs_with_length(l-1)

	r:=make([]string, 0, 4*len(s))

	for _,s2:=range s{
		for _,c:=range "ACTG"{
			r=append(r, s2+string(c))
		}
	}
	return r
}

func random_char(chars string)string{
	return string(chars[rand.Intn(len(chars))])
}

func makeABunchaRandomSeqs(n_seq_sets, seqs_per_set, min_len, len_var int)[][]wtype.DNASequence{
	var seqs [][]wtype.DNASequence

	seqs=make([][]wtype.DNASequence, n_seq_sets)

	for i:=0;i<n_seq_sets;i++{
		seqs[i]=make([]wtype.DNASequence, seqs_per_set)
		for j:=0;j<seqs_per_set;j++{
			seqs[i][j]=wtype.DNASequence{fmt.Sprintf("SEQ%04d", i*seqs_per_set+j+1),random_dna_seq(rand.Intn(len_var)+min_len)}
		}
	}
	return seqs
}
