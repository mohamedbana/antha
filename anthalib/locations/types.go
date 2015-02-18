package locations

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
