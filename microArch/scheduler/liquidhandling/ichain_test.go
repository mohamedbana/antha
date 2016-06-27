package liquidhandling

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func TestIChain(t *testing.T) {
	fmt.Println("TEST 1")
	chain := NewIChain(nil)

	s := []string{"A", "B", "C", "D", "E", "F"}

	for _, k := range s {
		ins := wtype.NewLHInstruction()

		cmp := wtype.NewLHComponent()

		cmp.ID = k

		ins.AddComponent(cmp)
		chain.Add(ins)
	}

	chain.Print()

}
func TestIChain2(t *testing.T) {
	fmt.Println("TEST 2")
	chain := NewIChain(nil)

	s := []string{"A", "B", "C", "D", "E", "F"}

	cmps := make([]*wtype.LHComponent, 0, 1)

	for i, k := range s {

		cmp := wtype.NewLHComponent()

		cmp.ID = k
		if i != 0 {
			cmp.AddParent(s[i-1])
		}
		if i != len(s)-1 {
			cmp.AddDaughter(s[i+1])
		}

		cmps = append(cmps, cmp)
	}

	for i, k := range cmps {
		ins := wtype.NewLHInstruction()
		ins.AddComponent(k)
		if i != len(s)-1 {
			ins.AddProduct(cmps[i+1])
		}
		fmt.Println("DOING NODE ", k.ID, " WITH PARENT: ", k.ParentID, " AND PRODUCT ", ins.ProductID)
		chain.Add(ins)
	}

	chain.Print()

}

func TestIChain3(t *testing.T) {
	fmt.Println("TEST 2")
	chain := NewIChain(nil)

	s := []string{"A", "B", "C", "D", "E", "F"}

	cmps := make([]*wtype.LHComponent, 0, 1)

	for i, k := range s {

		cmp := wtype.NewLHComponent()

		cmp.ID = k
		if i != 0 {
			cmp.AddParent(s[i-1])
		}
		if i != len(s)-1 {
			cmp.AddDaughter(s[i+1])
		}

		cmps = append(cmps, cmp)
	}

	cmp := wtype.NewLHComponent()
	cmp.ID = "Z"
	cmp.AddParent("C")
	cmps = append(cmps, cmp)

	cmp = wtype.NewLHComponent()
	cmp.ID = "Y"
	cmps = append(cmps, cmp)

	for i, k := range cmps {
		ins := wtype.NewLHInstruction()
		ins.AddComponent(k)
		if i != len(s)-1 && cmp.ID != "Z" && cmp.ID != "Y" {
			ins.AddProduct(cmps[i+1])
		}
		chain.Add(ins)
	}

	chain.Print()

}
