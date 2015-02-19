package wtype

import (
	"fmt"
	"strconv"
)

// base types required for movement

// a Location is somewhere something can be
// it can recursively contain other Locations
type Location struct {
	AnthaObject
	Positions []*Location
	Container *Location
	Shape     int
}

func NewLocation(name string, nPositions int, shape int) *Location {
	nao := NewAnthaObject(name)
	positions := make([]*Location, nPositions)
	l := Location{nao, positions, nil, shape}
	l.Container = &l
	if nPositions > 0 {
		for i := 0; i < nPositions; i++ {
			l.Positions[i] = NewLocation(name+"_position_"+strconv.Itoa(i+1), 0, shape)
			l.Positions[i].Container = &l
		}
	}
	return &l
}

// defines when two locations are the same
// level = 0 means only the exact same place
// level = 1 means they share a parent etc
func SameLocation(l, m *Location, level int) bool {
	if level < 0 {
		panic(fmt.Sprintf("SameLocation: level parameter %d makes no sense", level))
	} else if level == 0 {
		if l.ID == m.ID {
			return true
		} else {
			return false
		}
	} else {
		return SameLocation(l.Container, m.Container, level-1)
	}
}
