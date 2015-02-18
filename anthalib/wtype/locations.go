package wtype

import "fmt"

// base types required for movement

// a Location is somewhere something can be
// it can recursively contain other Locations
type Location struct {
	ID        string
	Name      string
	Positions []*Location
	Container *Location
	Shape     int
}

// defines when two locations are the same
// level = 0 means only the exact same place
// level = 1 means they share a parent etc
func SameLocation(l, m *Location, level int) bool {
	if level < 0 {
		panic(fmt.Sprintf("SameLocation: level parameter $d makes no sense".level))
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
