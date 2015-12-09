package execute

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/workflow"
	"golang.org/x/net/context"
	"reflect"
)

type constructor func(string) (interface{}, error)

var (
	ptipbox         wtype.LHTipbox
	pplate          wtype.LHPlate
	pcomponent      wtype.LHComponent
	unknownParam    = errors.New("unknown parameter")
	cannotConstruct = errors.New("cannot construct parameter")
	nilValue        reflect.Value
	constructors    = map[reflect.Type]constructor{
		reflect.TypeOf(ptipbox): func(x string) (interface{}, error) {
			return constructOrError(factory.GetTipByType(x))
		},
		reflect.TypeOf(pplate): func(x string) (interface{}, error) {
			return constructOrError(factory.GetPlateByType(x))
		},
		reflect.TypeOf(pcomponent): func(x string) (interface{}, error) {
			return constructOrError(factory.GetComponentByType(x))
		},
	}
)

func constructOrError(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, cannotConstruct
	}
	return v, nil
}

// Structure of parameter data
type ParamsData struct {
	Parameters map[string]map[string][]byte
}

type Options struct {
	WorkflowData []byte
	ParamData    []byte
}

func findConstructor(typ reflect.Type) constructor {
	// XXX: consider supporting convertible types too
	return constructors[typ]
}

func unmarshalOne(value reflect.Value, data []byte) (reflect.Value, error) {
	typ := value.Type()
	newValue := reflect.New(typ).Interface()
	origErr := json.Unmarshal(data, newValue)
	if origErr != nil {
		// Try to run constructor instead
		var carg string
		if construct := findConstructor(typ); construct == nil {
			return nilValue, origErr
		} else if err := json.Unmarshal(data, &carg); err != nil {
			return nilValue, fmt.Errorf("%s: %s", err, origErr)
		} else if v, err := construct(carg); err != nil {
			return nilValue, fmt.Errorf("%s: %s", err, origErr)
		} else {
			newValue = v
		}
	}
	return reflect.ValueOf(newValue).Elem(), nil
}

func unmarshal(value reflect.Value, data []byte) (reflect.Value, error) {
	typ := value.Type()
	switch typ.Kind() {
	case reflect.Slice:
		raw := make([]json.RawMessage, 0)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		s := reflect.MakeSlice(typ, 0, 0)
		for idx, bs := range raw {
			svalue := reflect.Zero(typ.Elem())
			if idx < value.Len() {
				svalue = reflect.ValueOf(value.Index(idx).Interface())
			}
			v, err := unmarshal(svalue, bs)
			if err != nil {
				return nilValue, err
			}
			s = reflect.Append(s, v)
		}
		return s, nil
	case reflect.Map:
		raw := make(map[string]json.RawMessage)
		if err := json.Unmarshal(data, &raw); err != nil {
			return nilValue, err
		}
		m := reflect.MakeMap(typ)
		for k, bs := range raw {
			kvalue := reflect.ValueOf(k)
			mvalue := value.MapIndex(kvalue)
			if mvalue == nilValue {
				mvalue = reflect.Zero(typ.Elem())
			} else {
				mvalue = reflect.ValueOf(mvalue.Interface())
			}
			v, err := unmarshal(mvalue, bs)
			if err != nil {
				return nilValue, err
			}
			m.SetMapIndex(kvalue, v)
		}
		return m, nil
	case reflect.Struct:
		return unmarshalOne(value, data)
	case reflect.Ptr:
		etyp := typ.Elem()
		if c := findConstructor(etyp); c == nil {
			return unmarshal(reflect.Zero(etyp), data)
		} else if v, err := unmarshalOne(reflect.Zero(etyp), data); err != nil {
			return nilValue, err
		} else {
			return v.Addr(), nil
		}
	default:
		return unmarshalOne(value, data)
	}
}

func setParam(w *workflow.Workflow, process, name string, data []byte, in map[string]interface{}) error {
	prototype, ok := in[name]
	if !ok {
		return unknownParam
	}

	value, err := unmarshal(reflect.ValueOf(prototype), data)
	if err != nil {
		return err
	}
	return w.SetParam(workflow.Port{Process: process, Port: name}, value)
}

func setParams(ctx context.Context, data []byte, w *workflow.Workflow) error {
	var params ParamsData
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}
	for process, params := range params.Parameters {
		c, err := w.FuncName(process)
		if err != nil {
			return fmt.Errorf("cannot get component for process %q: %s", process, err)
		}
		runner, err := inject.Find(ctx, inject.NameQuery{Repo: c})
		if err != nil {
			return fmt.Errorf("unknown component %q: %s", c, err)
		}
		cr, ok := runner.(*inject.CheckedRunner)
		if !ok {
			return fmt.Errorf("cannot get type information for component %q", c)
		}
		in := inject.MakeValue(cr.In)
		for name, value := range params {
			if err := setParam(w, process, name, value, in); err != nil {
				return fmt.Errorf("cannot assign parameter %q of process %q to %s: %s", name, process, string(value), err)
			}
		}
	}
	return nil
}

func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData})
	if err != nil {
		return nil, err
	}

	if err := setParams(parent, opt.ParamData, w); err != nil {
		return nil, err
	}

	// XXX: Setup trace
	if err := w.Run(parent); err != nil {
		return nil, err
	}

	return w, nil

	//cfData, err := ioutil.ReadFile(parametersFile)
	//if err != nil {
	//	return err
	//}

	//cf, err := util.NewConfig(cfData, wf)
	//if err != nil {
	//	return err
	//}
	//if _, ok := cf.Config["JOBID"]; !ok {
	//	cf.Config["JOBID"] = "default"
	//}
}
