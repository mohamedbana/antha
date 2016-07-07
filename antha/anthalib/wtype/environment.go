package wtype

import (
	"wunit"
)

type Environment struct {
	Temperature         wunit.Temperature
	Pressure            wunit.Pressure
	Humidity            float64
	MeanAirFlowVelocity wunit.Velocity
	/// etc
}
