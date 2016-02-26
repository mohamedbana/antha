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
		if it.Parent == nil {
			return it
		} else {
			// should not be here!
			logger.Fatal("Improper use of IChain")
		}
	} else {
		if it == nil {
			logger.Fatal("IT shouldn't be nil")
		}
		for _, v := range it.Values {
			// true if any component used by ins is *this*
			if ins.HasParent(v.ProductID) {
				return it.GetChild()
			}
		}

		return it.GetChild().FindNodeFor(ins)
	}
	// unreachable: pstr either is or isn't ""
	return nil
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
