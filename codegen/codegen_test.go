package codegen

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"testing"
)

func makeOneSample(masterMix ast.Node) (r []ast.Node) {
	cells := &ast.ApplyExpr{
		Func: "incubate",
		From: &ast.GatherExpr{
			Key:  "CellStart",
			From: &ast.NewExpr{From: &ast.UseExpr{Desc: "cells"}},
		},
	}
	r = append(r, cells)

	sample := &ast.ApplyExpr{
		Func: "mix",
		From: &ast.GatherExpr{
			From: &ast.ListExpr{
				From: []ast.Node{
					&ast.NewExpr{From: &ast.UseExpr{Desc: "water"}},
					masterMix,
					&ast.NewExpr{From: &ast.UseExpr{Desc: "part1"}},
					&ast.NewExpr{From: &ast.UseExpr{Desc: "part2"}},
				},
			},
		},
	}
	r = append(r, sample)

	sample1 := &ast.ApplyExpr{
		Func: "incubate",
		From: &ast.GatherExpr{
			Key:  "SampleStart",
			From: sample,
		},
	}
	r = append(r, sample1)

	sample2 := &ast.ApplyExpr{
		Func: "incubate",
		From: &ast.GatherExpr{
			Key:  "SampleStop",
			From: sample1,
		},
	}
	r = append(r, sample2)

	scells := &ast.ApplyExpr{
		Func: "mix",
		Near: cells,
		From: &ast.GatherExpr{
			From: &ast.ListExpr{
				From: []ast.Node{
					cells,
					sample2,
				},
			},
		},
	}
	r = append(r, scells)

	scells1 := &ast.ApplyExpr{
		Func: "incubate",
		From: &ast.GatherExpr{
			Key:  "PostPlasmid",
			From: scells,
		},
	}
	r = append(r, scells1)

	scells2 := &ast.ApplyExpr{
		Func: "mix",
		From: &ast.GatherExpr{
			From: &ast.ListExpr{
				From: []ast.Node{
					scells1,
					&ast.NewExpr{From: &ast.UseExpr{Desc: "recovery"}},
				},
			},
		},
	}
	r = append(r, scells2)

	rcells := &ast.ApplyExpr{
		Func: "incubate",
		From: &ast.GatherExpr{
			Key:  "Recovery",
			From: scells2,
		},
	}
	r = append(r, rcells)

	pcells := &ast.ApplyExpr{
		Opt:  "agarplate",
		Func: "mix",
		From: &ast.GatherExpr{
			From: &ast.ListExpr{
				From: []ast.Node{
					rcells,
				},
			},
		},
	}
	r = append(r, pcells)
	return
}

func TestWellFormed(t *testing.T) {
	t.Skip("tbd")
	// Example transformation protocol in AST form

	var nodes []ast.Node
	for i := 0; i < 4; i += 1 {
		// Premake master mix
		mix := &ast.ApplyExpr{
			Func: "mix",
			Gen:  0,
			From: &ast.GatherExpr{
				From: &ast.ListExpr{
					From: []ast.Node{
						&ast.NewExpr{From: &ast.UseExpr{Desc: "premix1"}},
						&ast.NewExpr{From: &ast.UseExpr{Desc: "premix2"}},
					},
				},
			},
		}
		nodes = append(nodes, makeOneSample(mix)...)
		nodes = append(nodes, mix)
	}

	if _, err := Compile(target.New(), &ast.BundleExpr{From: nodes}); err != nil {
		t.Fatal(err)
	}
}

func TestGenConstraint(t *testing.T) {
}

func TestNearConstraint(t *testing.T) {
}

func TestPrevConstraint(t *testing.T) {
}
