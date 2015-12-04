package inject

import (
	"testing"
)

func TestAssign(t *testing.T) {
	type X struct {
		A int
		B string
	}
	type SuperX struct {
		A int
		B string
		C int
	}
	type SubX struct {
		B int
	}

	var x1, x2 X
	var ipx2 interface{} = &x2
	var ix2 interface{} = x2
	var superx SuperX
	var subx SubX

	if err := AssignableTo(&x1, &x2); err != nil {
		t.Error(err)
	}
	if err := AssignableTo(x1, &x2); err != nil {
		t.Error(err)
	}
	if err := AssignableTo(x1, x2); err == nil {
		t.Errorf("expected error assigning to value")
	}

	if err := AssignableTo(x1, ipx2); err != nil {
		t.Error(err)
	}

	if err := AssignableTo(x1, ix2); err == nil {
		t.Errorf("expected error assigning to value")
	}

	if err := AssignableTo(x1, &superx); err != nil {
		t.Error(err)
	}

	if err := AssignableTo(x1, &subx); err == nil {
		t.Errorf("expected error assigning to value")
	}

	x1.A = 3
	x1.B = "hello"
	if err := Assign(x1, &x2); err != nil {
		t.Error(err)
	} else if x2.A != 3 || x2.B != "hello" {
		t.Errorf("expecting %v got %v", x1, x2)
	}
}
