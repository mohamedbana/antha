package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type IChain struct {
	Parent *IChain
	Child  *IChain
	Values []*wtype.LHInstruction
}

func NewIChain(parent *IChain) *IChain {
	var it IChain
	it.Parent = parent
	it.Values = make([]*wtype.LHInstruction, 0, 1)
	return &it
}

func (it *IChain) ValueIDs() []string {
	r := make([]string, 0, 1)

	for _, v := range it.Values {
		r = append(r, v.ID)
	}
	return r
}

func (it *IChain) Add(ins *wtype.LHInstruction) {
	p := it.FindNodeFor(ins)
	p.Values = append(p.Values, ins)
}

func (it *IChain) GetChild() *IChain {
	if it.Child == nil {
		it.Child = NewIChain(it)
	}
	return it.Child
}

func (it *IChain) FindNodeFor(ins *wtype.LHInstruction) *IChain {
	// having now done this above it's pretty trivial

	if ins.ParentString() == "" {
		if it.Parent == nil {
			// this is the root, and we are in the right place
			return it
		}
	}

	if it.HasParentOf(ins) {
		return it.GetChild()
	}

	// if we're at the end, return it

	if len(it.Values) == 0 {
		return it
	}

	return it.GetChild().FindNodeFor(ins)
}

func (it *IChain) HasParentOf(ins *wtype.LHInstruction) bool {
	for _, v := range it.Values {
		for _, cmp := range v.Components {
			if ins.HasParent(cmp.ID) {
				return true
			}
		}
	}

	return false
}

func (it *IChain) HasChildOf(ins *wtype.LHInstruction) bool {
	for _, v := range it.Values {
		if v.HasParent(v.ProductID) {
			return true
		}
	}

	return false
}

func (it *IChain) Print() {
	fmt.Println("****")
	fmt.Println("\tPARENT NIL: ", it.Parent == nil)
	fmt.Println("\tINPUTS: ", len(it.InputIDs()))
	fmt.Println("\tPRODUCTS: ", len(it.ProductIDs()))
	if it.Child != nil {
		it.Child.Print()
	}
}

func (it *IChain) InputIDs() string {
	s := ""

	for _, ins := range it.Values {
		for _, c := range ins.Components {
			s += c.ID + "   "
		}
		s += ","
	}

	return s
}

func (it *IChain) ProductIDs() string {
	s := ""

	for _, ins := range it.Values {
		s += ins.ProductID + "   "
	}
	return s
}

func (it *IChain) Flatten() []string {
	var ret []string

	if it == nil {
		return ret
	}

	for _, v := range it.Values {
		ret = append(ret, v.ID)
	}

	ret = append(ret, it.Child.Flatten()...)

	return ret
}
