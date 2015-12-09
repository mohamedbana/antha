package execute

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	type Value string
	var x Value
	golden := Value("hello")

	if v, err := unmarshal(reflect.ValueOf(x), []byte(`"hello"`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s, golden) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestInt(t *testing.T) {
	type Value int
	var x Value
	golden := Value(1)
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`1`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s, golden) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestStruct(t *testing.T) {
	type Value struct {
		A string
		B int
	}
	var x Value
	golden := Value{A: "hello", B: 1}

	if v, err := unmarshal(reflect.ValueOf(x), []byte(`{"A": "hello", "B": 1}`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s, golden) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestMap(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	type Value map[string]Elem
	var x Value
	golden := Value{
		"A": Elem{A: "hello", B: 1},
		"B": Elem{A: "hello", B: 2},
	}
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`{"A": {"A": "hello", "B": 1}, "B": {"A": "hello", "B": 2} }`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s, golden) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestSlice(t *testing.T) {
	type Elem struct {
		A string
		B int
	}
	type Value []Elem
	var x Value
	golden := Value{
		Elem{A: "hello", B: 1},
		Elem{A: "hello", B: 2},
	}
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`[ {"A": "hello", "B": 1}, {"A": "hello", "B": 2} ]`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s, golden) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestConstruct(t *testing.T) {
	var x *wtype.LHTipbox
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`"CyBio250Tipbox"`)); err != nil {
		t.Fatal(err)
	} else if _, ok := v.Interface().(*wtype.LHTipbox); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	}
}

func TestConstructMapFailure(t *testing.T) {
	type Elem struct {
		A string
		T *wtype.LHTipbox
	}
	type Value map[string]Elem
	var x Value
	if _, err := unmarshal(reflect.ValueOf(x), []byte(`{"A": {"A": "hello", "T": "CyBio250Tipbox"} }`)); err == nil {
		t.Fatal("expecting failure but got success")
	}
}

func TestConstructMap(t *testing.T) {
	type Value map[string]interface{}
	x := Value{
		"A": &wtype.LHTipbox{},
		"B": 0,
		"C": "",
	}
	golden := Value{
		"A": factory.GetTipboxByType("CyBio250Tipbox"),
		"B": 1,
		"C": "hello",
	}
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`{"A": "CyBio250Tipbox", "B": 1, "C": "hello" }`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if !reflect.DeepEqual(s["B"], golden["B"]) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if !reflect.DeepEqual(s["C"], golden["C"]) {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa, ok := golden["A"].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if bb, ok := s["A"].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}

func TestConstructSlice(t *testing.T) {
	type Value []interface{}
	x := Value{
		&wtype.LHTipbox{},
		&wtype.LHPlate{},
	}
	golden := Value{
		factory.GetTipboxByType("CyBio250Tipbox"),
		factory.GetPlateByType("pcrplate_with_cooler"),
	}
	if v, err := unmarshal(reflect.ValueOf(x), []byte(`[ "CyBio250Tipbox", "pcrplate_with_cooler" ]`)); err != nil {
		t.Fatal(err)
	} else if s, ok := v.Interface().(Value); !ok {
		t.Fatalf("expecting %T but got %T instead", x, v.Interface())
	} else if len(s) != 2 {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa, ok := golden[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if bb, ok := s[0].(*wtype.LHTipbox); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa, ok := golden[1].(*wtype.LHPlate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if bb, ok := s[1].(*wtype.LHPlate); !ok {
		t.Errorf("expecting %v but got %v instead", golden, s)
	} else if aa.Type != bb.Type {
		t.Errorf("expecting %v but got %v instead", golden, s)
	}
}
