// ph.go
package buffers

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"time"
)

/*const (
	reftemp wunit.Temperature{25,"C"}
)
*/
type PHperdegC float64

type PHMeasurement struct {
	Component       *wtype.LHComponent
	Location        *wtype.LHPlate
	PH              float64
	Temp            wunit.Temperature
	TempCorrected   *float64
	RefTemp         *wunit.Temperature
	TempCoefficient *PHperdegC
	Adjusted        *float64
	Adjustedwith    *wtype.LHComponent
}

func (ph *PHMeasurement) TempCompensation(reftemp wunit.Temperature, tempcoefficientforsolution PHperdegC) (compensatedph float64) {

	ph.RefTemp = &reftemp //.SIValue()
	ph.TempCoefficient = &tempcoefficientforsolution

	tempdiff := ph.Temp.SIValue() - ph.RefTemp.SIValue()

	compensatedph = ph.PH + (float64(tempcoefficientforsolution) * tempdiff)
	ph.TempCorrected = &compensatedph
	return
}

// placeholder

/*func MeasurePH(*wtype.LHSolution) (measurement float64) {
	return 7.0
}*/

func MeasurePH(*wtype.LHSolution) (measured PHMeasurement) {
	measured = PHMeasurement{nil, nil, 0.0, wunit.NewTemperature(0.0, "C"), nil, nil, nil, nil, nil}
	return
}

// this should be performed on an LHComponent
// currently (wrongly) assumes only acid or base will be needed
func (ph *PHMeasurement) AdjustpH(ph_setpoint float64, ph_tolerance float64, ph_setPointTemp wunit.Temperature, Acid *wtype.LHComponent, Base *wtype.LHComponent) (adjustedsol wtype.LHComponent, newph PHMeasurement, componentadded wtype.LHComponent, err error) {

	pHmax := ph_setpoint + ph_tolerance
	pHmin := ph_setpoint - ph_tolerance

	//sammake([]wtype.LHComponent,0)

	if ph.PH > pHmax {
		// calculate concentration of solution needed first, for now we'll add 10ul at a time until adjusted
		for {
			//newphmeasurement = ph
			acidsamp := mixer.Sample(Acid, wunit.NewVolume(10, "ul"))
			temporary := mixer.MixInto(ph.Location, ph.Component, acidsamp)
			time.Sleep(10 * time.Second)
			newphmeasurement := MeasurePH(temporary)
			if newphmeasurement.PH > pHmax {
				*ph = newphmeasurement
				//}
				//if {
				//ph.PH < pHmin
				//	continue
			} else {
				adjustedsol = *ph.Component
				newph = *ph
				componentadded = *Acid
				err = nil
				return
			}

		}
	}
	// basically just a series of sample, stir, wait and recheck pH

	if ph.PH < pHmin {
		for {
			//newphmeasurement = ph
			basesamp := mixer.Sample(Base, wunit.NewVolume(10, "ul"))
			temporary := mixer.MixInto(ph.Location, ph.Component, basesamp)
			time.Sleep(10 * time.Second)
			newphmeasurement := MeasurePH(temporary)
			if newphmeasurement.PH > pHmax {
				*ph = newphmeasurement
			} else {
				adjustedsol = *ph.Component
				newph = *ph
				componentadded = *Base
				err = nil
				return
			}

		}
	}
	//adjustedsol = ph.Component, newph = ph, componentadded = Acid,
	err = fmt.Errorf("Something went wrong here!")
	return
}
