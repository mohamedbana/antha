package lib

var (
	components []Component
)

type ParamDesc struct {
	Name, Desc, Kind string
}

type ComponentDesc struct {
	Desc   string
	Params []ParamDesc
}

type Component struct {
	Name        string
	Constructor func() interface{}
	Desc        ComponentDesc
}

func addComponent(desc Component) {
	components = append(components, desc)
}

func GetComponents() []Component {
	return components
}
