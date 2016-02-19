package wtype

import "strings"

type platelocation struct {
	ID     string
	Coords WellCoords
}

func (pc platelocation) ToString() string {
	return pc.ID + ":" + pc.Coords.FormatA1()
}

func PlateLocationFromString(s string) platelocation {
	tx := strings.Split(s, ":")
	return platelocation{tx[0], MakeWellCoords(tx[1])}
}
