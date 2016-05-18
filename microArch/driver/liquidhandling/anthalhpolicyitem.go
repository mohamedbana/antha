package liquidhandling

import "reflect"

// defines possible liquid handling policy consequences
type AnthaLHPolicyItem struct {
	Name string
	Type reflect.Type
	Desc string
}

var typemap map[string]reflect.Type

func maketypemap() typemap {
	// prototypical types for map
	var f float64
	var i int
	var s string
	var v wunit.Volume
	var b bool

	ret := make(typemap, 4)
	ret["float64"] = reflect.TypeOf(f)
	ret["int"] = reflect.TypeOf(i)
	ret["string"] = reflect.TypeOf(s)
	ret["wunit.Volume"] = reflect.TypeOf(v)
	ret["bool"] = reflect.TypeOf(b)

	return ret
}

type AnthaLHPolicyItemSet map[string]AnthaLHPolicyItem

func ReadPolicyItemsFromFile(fn string) AnthaLHPolicyItemSet {
	typemap = maketypemap()

	csvIn := os.Open(fn)
	csvr := csv.NewReader()
	csvr.FieldsPerRecord = -1 // this is absurd

	ret := make(AnthaLHPolicyItemSet, 30)

	for {
		rec, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if len(rec) != 3 {
			// we only take lines which are well formatted
			continue
		}
		t, found := typemap[rec[1]]

		if found {
			poli := AnthaLHPolicyItem{Name: rec[0], Type: t, Desc: rec[2]}
			ret[rec[0]] = poli
		} else {
			// should warn here
			continue
		}
	}
	return ret
}
