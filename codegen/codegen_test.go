package codegen

import (
	"github.com/antha-lang/antha/target"
	"testing"
)

// TODO(ddn): add gathers after mix-sample to make mix-plate

func makeOneSample(masterMix AstNode) (r []AstNode) {
	cells := &ApplyExpr{
		Func: "incubate0",
		From: &GatherExpr{
			Key:  "CellStart",
			From: &NewExpr{From: &UseExpr{Desc: "cells"}},
		},
	}
	r = append(r, cells)

	sample := &ApplyExpr{
		Func: "mix",
		From: &ListExpr{
			From: []AstNode{
				&NewExpr{From: &UseExpr{Desc: "water"}},
				masterMix,
				&NewExpr{From: &UseExpr{Desc: "part1"}},
				&NewExpr{From: &UseExpr{Desc: "part2"}},
			},
		},
	}
	r = append(r, sample)

	sample1 := &ApplyExpr{
		Func: "incubate1",
		From: &GatherExpr{
			Key:  "SampleStart",
			From: sample,
		},
	}
	r = append(r, sample1)

	sample2 := &ApplyExpr{
		Func: "incubate2",
		From: &GatherExpr{
			Key:  "SampleStop",
			From: sample1,
		},
	}
	r = append(r, sample2)

	scells := &ApplyExpr{
		Func: "mix",
		Near: cells,
		From: &ListExpr{
			From: []AstNode{
				cells,
				sample2,
			},
		},
	}
	r = append(r, scells)

	scells1 := &ApplyExpr{
		Func: "incubate3",
		From: &GatherExpr{
			Key:  "PostPlasmid",
			From: scells,
		},
	}
	r = append(r, scells1)

	scells2 := &ApplyExpr{
		Func: "mix",
		From: &ListExpr{
			From: []AstNode{
				scells1,
				&NewExpr{From: &UseExpr{Desc: "recovery"}},
			},
		},
	}
	r = append(r, scells2)

	rcells := &ApplyExpr{
		Func: "incubate4",
		From: &GatherExpr{
			Key:  "Recovery",
			From: scells2,
		},
	}
	r = append(r, rcells)

	pcells := &ApplyExpr{
		Opt:  "agarplate",
		Func: "mix",
		From: &ListExpr{
			From: []AstNode{
				rcells,
			},
		},
	}
	r = append(r, pcells)
	return
}

func TestWellFormed(t *testing.T) {
	t.Skip("tbd")
	// Example transformation protocol in AST form

	var nodes []AstNode
	for i := 0; i < 4; i += 1 {
		// Premake master mix
		mix := &ApplyExpr{
			Func: "mix",
			Gen:  0,
			From: &ListExpr{
				From: []AstNode{
					&NewExpr{From: &UseExpr{Desc: "premix1"}},
					&NewExpr{From: &UseExpr{Desc: "premix2"}},
				},
			},
		}
		nodes = append(nodes, makeOneSample(mix)...)
		nodes = append(nodes, mix)
	}

	if _, err := Compile(target.New(), &BundleExpr{From: nodes}); err != nil {
		t.Fatal(err)
	}
}

func TestGenConstraint(t *testing.T) {
}

func TestNearConstraint(t *testing.T) {
}

func TestPrevConstraint(t *testing.T) {
}
