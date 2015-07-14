package lib

import (
	"github.com/antha-lang/antha/antha/component/lib/Sum"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly"
)

type ComponentDesc struct {
	Name        string
	Constructor func() interface{}
}

func GetComponents() []ComponentDesc {
	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
	portMap["Sum"] = make(map[string]bool)

	portMap["TypeIISConstructAssembly"] = make(map[string]bool)

	c := make([]ComponentDesc, 0)
	c = append(c, ComponentDesc{Name: "Sum", Constructor: Sum.NewSum})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly", Constructor: TypeIISConstructAssembly.NewTypeIISConstructAssembly})

	return c
}
