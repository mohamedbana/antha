// conversion.go
// designed to convert units: e.g. kg/m^3 to g/L etc
package wunit

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

func VolumetoMass(v Volume, d Density) (m Mass) {
	//mass := m.SIValue()
	density := d.SIValue()

	volume := v.SIValue() //* 1000 // convert m^3 to l

	mass := volume * density // in m^3

	m = NewMass(mass, "kg")
	return m
}
