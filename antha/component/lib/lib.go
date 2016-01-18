package lib

import (
	"errors"
	"fmt"
	/*
		"github.com/antha-lang/antha/antha/component/lib/Aliquot"
		"github.com/antha-lang/antha/antha/component/lib/AliquotTo"
		"github.com/antha-lang/antha/antha/component/lib/Assaysetup"
		"github.com/antha-lang/antha/antha/component/lib/BlastSearch"
		"github.com/antha-lang/antha/antha/component/lib/BlastSearch_wtype"
		"github.com/antha-lang/antha/antha/component/lib/Colony_PCR"
		"github.com/antha-lang/antha/antha/component/lib/DNA_gel"
		"github.com/antha-lang/antha/antha/component/lib/Datacrunch"
		"github.com/antha-lang/antha/antha/component/lib/Evaporationrate"
		"github.com/antha-lang/antha/antha/component/lib/FindPartsthat"
		"github.com/antha-lang/antha/antha/component/lib/Iterative_assembly_design"
		"github.com/antha-lang/antha/antha/component/lib/Kla"
		"github.com/antha-lang/antha/antha/component/lib/LoadGel"
		"github.com/antha-lang/antha/antha/component/lib/LookUpMolecule"
		"github.com/antha-lang/antha/antha/component/lib/MakeBuffer"
		"github.com/antha-lang/antha/antha/component/lib/MakeMedia"
		"github.com/antha-lang/antha/antha/component/lib/Mastermix"
		"github.com/antha-lang/antha/antha/component/lib/Mastermix_reactions"
		"github.com/antha-lang/antha/antha/component/lib/MoClo_design"
		"github.com/antha-lang/antha/antha/component/lib/NewDNASequence"
		"github.com/antha-lang/antha/antha/component/lib/OD"
		"github.com/antha-lang/antha/antha/component/lib/PCR"
		"github.com/antha-lang/antha/antha/component/lib/Paintmix"
		"github.com/antha-lang/antha/antha/component/lib/Phytip_miniprep"
		"github.com/antha-lang/antha/antha/component/lib/PipetteImage"
		"github.com/antha-lang/antha/antha/component/lib/PipetteImage_CMYK"
		"github.com/antha-lang/antha/antha/component/lib/PipetteImage_living"
		"github.com/antha-lang/antha/antha/component/lib/PlateOut"
		"github.com/antha-lang/antha/antha/component/lib/Plotdata"
		"github.com/antha-lang/antha/antha/component/lib/Plotdata_spreadsheet"
		"github.com/antha-lang/antha/antha/component/lib/PreIncubation"
		"github.com/antha-lang/antha/antha/component/lib/Printname"
		"github.com/antha-lang/antha/antha/component/lib/ProtocolName_from_an_file"
		"github.com/antha-lang/antha/antha/component/lib/Recovery"
		"github.com/antha-lang/antha/antha/component/lib/RemoveRestrictionSites"
		"github.com/antha-lang/antha/antha/component/lib/RestrictionDigestion"
		"github.com/antha-lang/antha/antha/component/lib/RestrictionDigestion_conc"
		"github.com/antha-lang/antha/antha/component/lib/SDSprep"
		"github.com/antha-lang/antha/antha/component/lib/Scarfree_design"
		"github.com/antha-lang/antha/antha/component/lib/Scarfree_siteremove_orfcheck"
		"github.com/antha-lang/antha/antha/component/lib/SumVolume"
		"github.com/antha-lang/antha/antha/component/lib/Test"
		"github.com/antha-lang/antha/antha/component/lib/Thawtime"
		"github.com/antha-lang/antha/antha/component/lib/Transfer"
		"github.com/antha-lang/antha/antha/component/lib/Transformation"
		"github.com/antha-lang/antha/antha/component/lib/Transformation_complete"
		"github.com/antha-lang/antha/antha/component/lib/TypeIISAssembly_design"
		"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly"
		"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssemblyMMX"
		"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_alt"
		"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_sim"
	*/
	"github.com/antha-lang/antha/inject"
	"reflect"
)

var (
	components       []Component
	invalidComponent = errors.New("invalid component")
)

type ParamDesc struct {
	Name, Desc, Kind, Type string
}

type ComponentDesc struct {
	Desc   string
	Path   string
	Params []ParamDesc
}

type Component struct {
	Name        string
	Constructor func() interface{}
	Desc        ComponentDesc
}

// Helper function to add appropriate component information to the component
// library
func addComponent(desc Component) error {
	// Add type information if missing
	ts := make(map[string]string)

	type tdesc struct {
		Name string
		Type reflect.Type
	}

	add := func(name, t string) error {
		if _, seen := ts[name]; seen {
			return fmt.Errorf("parameter %q already seen", name)
		}
		ts[name] = t
		return nil
	}

	typeOf := func(i interface{}) ([]tdesc, error) {
		var tdescs []tdesc
		// Generated elements always have type *XXXOutput or *XXXInput
		t := reflect.TypeOf(i).Elem()
		if t.Kind() != reflect.Struct {
			return nil, invalidComponent
		}
		for i, l := 0, t.NumField(); i < l; i += 1 {
			tdescs = append(tdescs, tdesc{Name: t.Field(i).Name, Type: t.Field(i).Type})
		}
		return tdescs, nil
	}

	if r, ok := desc.Constructor().(inject.TypedRunner); !ok {
		return invalidComponent
	} else if inTypes, err := typeOf(r.Input()); err != nil {
		return err
	} else if outTypes, err := typeOf(r.Output()); err != nil {
		return err
	} else {
		for _, v := range append(inTypes, outTypes...) {
			if err := add(v.Name, v.Type.String()); err != nil {
				return err
			}
		}
	}

	for i, p := range desc.Desc.Params {
		t := &desc.Desc.Params[i].Type
		if len(*t) == 0 {
			*t = ts[p.Name]
		}
	}

	components = append(components, desc)
	return nil
}
func GetComponents() []Component {
	return components
}
