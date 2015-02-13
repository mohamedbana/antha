// anthalib//wopt/goldengate.go: Part of the Antha language
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

package main

import (
	"fmt"
	"math/rand"
	"math"
	"strings"
	"github.com/antha-lang/antha/anthalib/wtype"
	"time"
)

type BTSolution interface{
	BTScore()float64
	Reject()bool
	Accept()bool
	Next()BTSolution
	First()BTSolution
	ToString()string
	PrintNodes()
}

type GGSolution struct{
	Score float64
	CurNode int
	Nodes [][]string
	S map[string]float64
	P func(string, string)float64
}

func (ggs *GGSolution)PrintNodes(){
	fmt.Println("++++++++")

	for _,dbg:=range(ggs.Nodes){
		fmt.Println(dbg)
	}

	fmt.Println("++++++++")
}

func (ggs *GGSolution)BTScore()float64{
	return ggs.Score
}

func (ggs *GGSolution)Reject()bool{
	for i:=ggs.CurNode;i<len(ggs.Nodes);i++{
		if len(ggs.Nodes[i])==0{
			return true
		}
	}
	return false
}

func (ggs *GGSolution)Accept()bool{
	if(ggs.CurNode==len(ggs.Nodes)){
		return true
	}

	return false
}

func (ggs *GGSolution)Next()BTSolution{
	// return nil if there are no more nodes or no more candidates within this node

	if ggs.CurNode==len(ggs.Nodes) || len(ggs.Nodes[ggs.CurNode])==0{
		return nil
	}

	// create a new solution with
	// the curnode pointer advanced
	// and the score functions modified to account for 
	// the addition of the new interface
	// plus the scores revised
	// and return it
	// meanwhile in *this* node we delete the topmost candidate
	// from the list of candidates for this node

	var ifc string

	ifc,ggs.Nodes[ggs.CurNode]=ggs.Nodes[ggs.CurNode][0],ggs.Nodes[ggs.CurNode][1:len(ggs.Nodes[ggs.CurNode])]

	newnodes:=make([][]string, len(ggs.Nodes))

	for i,n:=range(ggs.Nodes){
		if(i<ggs.CurNode){
			newnodes[i]=n
		}else if(i==ggs.CurNode){
			n2:=make([]string,1)
			n2[0]=ifc
			newnodes[i]=n2
		}else{
			// remove any identical things
			n2:=make([]string,0,len(n))
			for _,s:=range n{
				if(s!=ifc){
					n2=append(n2,s)
				}
			}
			newnodes[i]=n2
		}
	}

	newS:=make(map[string]float64, len(ggs.S))

	for k,_:=range ggs.S{
		newS[k]=ggs.S[k]+ggs.P(ifc,k)
	}

	sc:=ggs.Score+ggs.S[ifc]

	ret:=GGSolution{sc, ggs.CurNode+1, newnodes, newS, ggs.P}

	return &ret
}

func (ggs *GGSolution)First()BTSolution{
	// synonymous with Next, to start with
	return ggs.Next()
}

func (ggs *GGSolution)ToString()string{
	s:=""
	for _,n:=range ggs.Nodes{
		s+=fmt.Sprintf("[")

		for _,i:=range n{
			s+=fmt.Sprintf("%s,",i)
		}

		s+=fmt.Sprintf("] ")
	}

	return s
}


// slight modification of the usual backtrack
// to allow score optimization
func bt(c BTSolution) BTSolution{
	if(c.Reject()){
		return nil
	}else if(c.Accept()){
		return c
	}

	best:=0.0
	var bestAt BTSolution

	s:=c.First()

	if s==nil{
		return s
	}

	for{
		t:=bt(s)
		if(t!=nil && t.BTScore()>best){
			best=t.BTScore()
			bestAt=t
		}
		s=s.Next()
		if s==nil{
			if best!=0.0{
				return bestAt
			}
			break
		}
	}
	return nil
}

// simple flat score for now

func make_unary_score(l int)map[string]float64{
	m:=make(map[string]float64, int(math.Pow(4,float64(l))))

	seqs:=all_dna_seqs_with_length(l)

	for _,s:=range seqs{
		m[s]=10.0
	}

	return m
}

func pairwise_score(s,t string)float64{
	// we go for a Hamming distance

	sc:=0.0

	for i:=0;i<len(s);i++{
		if s[i]==t[i]{
			sc-=1.0
		}
	}

	return sc
}


// interface to the optimized backtracking algorithm

func choose_interfaces(interfaces []string, unaryscore map[string]float64)(float64,[]string){
	// make the initial solution

	// convert the interface strings into nodes; in future we might sort

	nodes:=make([][]string, len(interfaces))

	for i,s:=range interfaces{
		n:=make([]string, 0, 1+len(s)/2)

		for k:=0;k<(len(s)/2)+1;k++{
			f:=s[k:k+len(s)/2]

			_,ok:=unaryscore[f]

			if strings.Contains(f, "N") || !ok{
				continue
			}
			n=append(n,f)
		}

		// we can break early here

		if len(n)==0{
			return -1.0,nil
		}

		nodes[i]=MakeUnique(n)
	}

	solution:=GGSolution{0.0, 0, nodes, unaryscore, pairwise_score}

	psaaaah:=bt(&solution)

	if psaaaah==nil{
		panic("No solutions")
	}

	psolution:=psaaaah.(*GGSolution)

	solution=*psolution

	stripnodes:=make([]string, len(solution.Nodes))

	for i,s:=range solution.Nodes{
		stripnodes[i]=s[0]
	}

	return  solution.Score, stripnodes 
}
