// conversion.go
// designed to convert units: e.g. kg/m^3 to g/L etc
package wunit

import ()

/*
type

func Splitunit(unit string)(numerators[]string, denominators[]string)

var conversiontable = map[string]map[string]float64{
	"density":map[string]float64{
		"g/L":
	}
}
*/

func MasstoVolume(m Mass, d Density) (v Volume) {
	mass := m.SIValue()
	density := d.SIValue()
	volume := mass / density // in m^3
	volume = volume * 1000   // in l
	v = NewVolume(mass, "l")
	return v
}
