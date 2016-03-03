package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/logger"
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
	pstr := ins.ParentString()
	if pstr == "" {
		// instruction is a root... belongs in the root node
		if it.Parent == nil {
			return it
		} else {
			// should not be here!
			logger.Fatal("Improper use of IChain")
		}
	}

	// is this node the root? We know pstr is not the right place for this
	// so we keep searching

	if it.Parent == nil {
		return it.GetChild().FindNodeFor(ins)
	}

	// so now we're here we know neither ins nor it is a root

	// if there are no values here we return this node - we're at the leaf

	if len(it.Values) == 0 {
		return it
	}

	Ihasparent := it.HasParentOf(ins)

	Chaschild := false

	if it.Child != nil && it.Child.HasChildOf(ins) {
		Chaschild = true
	}

	if Ihasparent || Chaschild {
		if it.Child != nil {
			// need a new link
			ch := it.Child
			it.Child = NewIChain(it)
			it.Child.Child = ch
			return ch
		} else {
			// add a new link
			return it.GetChild()
		}
	}

	// so we have neither the parent of ins in this node,
	// nor the child of ins in our child node, so we just
	// carry on the search

	return it.GetChild().FindNodeFor(ins)

}

func (it *IChain) HasParentOf(ins *wtype.LHInstruction) bool {
	for _, v := range it.Values {
		if ins.HasParent(v.ProductID) {
			return true
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
	fmt.Println("PARENT: ", it.Parent)
	fmt.Println("\tValues:", it.Values)
	if it.Child != nil {
		it.Child.Print()
	}
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
