package wtype

import (
	"fmt"
	"strconv"
)

// base types required for movement

// a Location is somewhere something can be
// it can recursively contain other Locations
type Location interface {
	Location_ID() string
	Location_Name() string
	Positions() []Location
	Container() Location
	Shape() Shape
}

type ConcreteLocation struct {
	AnthaObject
	Psns []*ConcreteLocation
	Cntr *ConcreteLocation
	Shap Shape
}

func (cl *ConcreteLocation) Location_ID() string {
	return cl.ID
}
func (cl *ConcreteLocation) Location_Name() string {
	return cl.Name
}
func (cl *ConcreteLocation) Positions() []Location {
	s := make([]Location, len(cl.Psns))
	for i, v := range cl.Psns {
		s[i] = Location(v)
	}
	return s
}
func (cl *ConcreteLocation) Container() Location {
	return cl.Cntr
}
func (cl *ConcreteLocation) Shape() Shape {
	return cl.Shap
}

func NewLocation(name string, nPositions int, shape Shape) Location {
	nao := NewAnthaObject(name)
	positions := make([]*ConcreteLocation, nPositions)
	l := ConcreteLocation{nao, positions, nil, shape}
	l.Cntr = &l
	if nPositions > 0 {
		for i := 0; i < nPositions; i++ {
			l.Psns[i] = NewLocation(name+"_position_"+strconv.Itoa(i+1), 0, shape).(*ConcreteLocation)
			l.Psns[i].Cntr = &l
		}
	}
	return &l
}

// defines when two locations are the same
// level = 0 means only the exact same place
// level = 1 means they share a parent etc
func SameLocation(l, m Location, level int) bool {
	if level < 0 {
		panic(fmt.Sprintf("SameLocation: level parameter %d makes no sense", level))
	} else if level == 0 {
		if l.Location_ID() == m.Location_ID() {
			return true
		} else {
			return false
		}
	} else {
		return SameLocation(l.Container(), m.Container(), level-1)
	}
}
