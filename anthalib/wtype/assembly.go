// anthalib//wtype/assembly.go: Part of the Antha language
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

package wtype

// base assembly type - it's essentially a union, even though this is a dirty word in 
// golang

type Assembly struct{
	Components []Assembly
	Sequences  []DNASequence
}

func NewAssembly()Assembly{
	components:=make([]Assembly, 0, 1)
	seqs:=make([]DNASequence, 0, 1)
	return Assembly{components, seqs}
}

func NewAssemblyFromSeqs(seqs []DNASequence)Assembly{
	components:=make([]Assembly, 0, 1)
	return Assembly{components, seqs}
}

func NewAssemblyFromComponents(seqs [][]DNASequence)Assembly{
	components:=make([]Assembly, len(seqs))
	seqs2:=make([]DNASequence,0,1)

	for i,set:=range seqs{
		components[i]=NewAssemblyFromSeqs(set)
	}

	return Assembly{components,seqs2}
}

// add the prefix to everything in the beginning set of sequences
func (ass *Assembly)Prepend(s string){
	if len(ass.Sequences)!=0{
		for i,_:=range ass.Sequences{
			ass.Sequences[i].Prepend(s)
		}
	}else{
		ass.Components[0].Prepend(s)
	}
}

// add the suffix to the end of all sequences in the last set
func(ass *Assembly)Append(s string){
	if len(ass.Sequences)!=0{
		for i,_:=range ass.Sequences{
			ass.Sequences[i].Append(s)
		}
	}else{
		ass.Components[len(ass.Components)-1].Append(s)
	}
}

// true if this is just a collection of sequences
func (ass *Assembly)IsBlock()bool{
	if len(ass.Sequences)==0{
		return false
	}
	return true
}

func (ass *Assembly)Generalised_Substring(start, length int)string{
	if ass.IsBlock(){
		return Generalised_Substring(DNAToStrings(ass.Sequences), start, length)
	} else{
		if start==0{
			return ass.Components[0].Generalised_Substring(start, length)
		} else{
			return ass.Components[len(ass.Components)-1].Generalised_Substring(start, length)
		}
	}
}

func Generalised_Substring(seqs []string, start, end int)string{
	s:=make([]byte, end)

	b:=0

	for i:=0;i<len(seqs);i++{
		if start==1{
			b=len(seqs[i])-end
		}

		for j:=0;j<end;j++{
			if i==0{
				s[j]=seqs[i][j+b]
			}else if(seqs[i][j+b]!=s[j]){
				s[j]=byte('N') 
			}
		}

	}

	return string(s)
}

func DNAToStrings(seqs []DNASequence)[]string{
	ret:=make([]string, len(seqs))

	for i,s:=range seqs{
		ret[i]=s.Sequence()
	}

	return ret
}