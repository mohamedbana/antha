package inject

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"reflect"
)

var notStructOrValue = errors.New("not pointer to struct or map[string]interface{}")
var stringType = reflect.TypeOf("")

// Input and output of injectable functions. Implementation of named and typed
// function parameters.
type Value map[string]interface{}

func MakeValue(x interface{}) map[string]interface{} {
	return structs.Map(x)
}

func isStringMap(value reflect.Value) bool {
	typ := value.Type()
	if typ.Kind() != reflect.Map {
		return false
	} else if typ.Key() != stringType {
		return false
	} else {
		return true
	}
}

func makeMapFields(value reflect.Value) map[string]reflect.Value {
	fields := make(map[string]reflect.Value)
	for _, key := range value.MapKeys() {
		k := key.Interface().(string)
		v := value.MapIndex(key)
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		fields[k] = v
	}
	return fields
}

func makeStructFields(value reflect.Value) map[string]reflect.Value {
	typ := value.Type()
	fields := make(map[string]reflect.Value)
	for i, n := 0, typ.NumField(); i < n; i += 1 {
		sf := typ.Field(i)
		f := value.Field(i)
		fields[sf.Name] = f
	}
	return fields
}

func makeFields(value reflect.Value) (map[string]reflect.Value, error) {
	typ := value.Type()
	switch k := typ.Kind(); {
	case k == reflect.Ptr:
		return makeFields(value.Elem())
	case k == reflect.Interface:
		return makeFields(value.Elem())
	case k == reflect.Struct:
		return makeStructFields(value), nil
	case isStringMap(value):
		return makeMapFields(value), nil
	}
	return nil, notStructOrValue
}

func assign(from, to interface{}, set bool) error {
	toValue := reflect.ValueOf(to)
	toFields, err := makeFields(toValue)
	if err != nil {
		return err
	}

	fromFields, err := makeFields(reflect.ValueOf(from))
	if err != nil {
		return err
	}

	isToMap := isStringMap(toValue)
	for name, v := range fromFields {
		toV, ok := toFields[name]

		if !ok {
			return fmt.Errorf("missing field %q", name)
		} else if fromT, toT := v.Type(), toV.Type(); !fromT.AssignableTo(toT) {
			return fmt.Errorf("field %q of type %s not assignable to %s", name, fromT, toT)
		} else if !isToMap && !toV.CanSet() {
			return fmt.Errorf("cannot set field %q", name)
		} else if !set {
			continue
		}

		if isToMap {
			toValue.SetMapIndex(reflect.ValueOf(name), v)
		} else {
			toV.Set(v)
		}
	}

	return nil
}

// Return if src is assignable to dst. Typing rule is as follows: (1) every
// field of src must have a field of the same name in dst, (2) the type of the
// src field must be golang assignable to the dst field, and (3) the dst fields
// must be golang settable (i.e., have an address).
func AssignableTo(src, dst interface{}) error {
	return assign(src, dst, false)
}

// Assign values from Value or struct to Value or struct. If src is not
// AssignableTo dst, return an error.
func Assign(src, dst interface{}) error {
	return assign(src, dst, true)
}
