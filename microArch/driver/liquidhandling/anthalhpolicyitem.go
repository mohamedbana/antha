package liquidhandling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/logger"
)

func GetPolicyConsequents() AnthaLHPolicyItemSet {
	return ReadPolicyItemsFromFile("lhpolicyconsequents.csv")
}

// defines possible liquid handling policy consequences
type AnthaLHPolicyItem struct {
	Name string
	Type reflect.Type
	Desc string
}

func (alhpi AnthaLHPolicyItem) TypeName() string {
	return alhpi.Type.Name()
}

var typemap map[string]reflect.Type

func maketypemap() map[string]reflect.Type {
	// prototypical types for map
	var f float64
	var i int
	var s string
	var v wunit.Volume
	var b bool

	ret := make(map[string]reflect.Type, 4)
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

	csvIn, err := os.Open(fn)

	if err != nil {
		logger.Fatal(fmt.Sprintf("Cannot find policy consequents file %s", fn))
	}

	csvr := csv.NewReader(csvIn)
	csvr.FieldsPerRecord = -1 // this is absurd

	ret := make(AnthaLHPolicyItemSet, 30)

	for {
		rec, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if len(rec) < 3 {
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

func (alhpis AnthaLHPolicyItemSet) TypeList() string {
	ks := make([]string, 0, len(alhpis))

	for k, _ := range alhpis {
		ks = append(ks, k)
	}

	sort.Strings(ks)

	s := ""

	for _, k := range ks {
		alhpi := alhpis[k]
		s += fmt.Sprintf("%s,%s,%s\n", k, alhpi.TypeName(), alhpi.Desc)
	}

	return s
}
