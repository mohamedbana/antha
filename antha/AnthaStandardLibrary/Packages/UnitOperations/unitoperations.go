// Example syntax
package UnitOperations

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"time"
)

type Pellet struct {
	wtype.Physical
}

type Culture struct {
	wtype.Suspension
	wtype.Organism
}

func Separate(culture Culture) (supernatant *wtype.LHComponent, pellet Pellet) {

	return
}

type Chromstep struct {
	Pipetstep
	Column
	Buffer *wtype.LHComponent
	//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
}

type Pipetstep struct {
	Name          string
	Volume        wunit.Volume //transfer volume and process volume?
	Aspiraterate  wunit.FlowRate
	Aspiratepause time.Duration
	Dispenserate  wunit.FlowRate
	Dispensepause time.Duration
	Cycles        int
	//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
}

type Column struct {
	Beadsize           wunit.Length
	Stationaryphase    string
	Separationproperty string
	Diameter           wunit.Length
	Height             wunit.Length
	Packedvolume       wunit.Volume
}

type Phytips struct { // interface?
	Column
	Tip wtype.LHTip
}

// may already be functions for aspirate and dispense in anthalib
func Aspirate(column Column, mixture *wtype.LHSolution, volume wunit.Volume, aspiraterate wunit.FlowRate) (aspiratedcolumn Column, aspiratedsolution *wtype.LHSolution) {

	return
}

func Dispense(column Column, mixture *wtype.LHSolution, volume wunit.Volume, aspiraterate wunit.FlowRate) (dispensedcolumn Column, dispensedsolution *wtype.LHSolution) {

	return
}

func PhysicaltoComponent(pellet *wtype.Physical) (component *wtype.LHComponent) {
	// placeholder
	return
}

func Resuspend(pellet *wtype.Physical, step Chromstep, column Column) (output_c *wtype.LHComponent, processedcolumn Column) {

	var output *wtype.LHSolution
	input := PhysicaltoComponent(pellet)
	samples := make([]*wtype.LHComponent, 0)
	samples = append(samples, step.Buffer, input)
	mixture := mixer.Mix(samples...)
	for i := 0; i < step.Pipetstep.Cycles; i++ {

		aspiratedcolumn, aspiratedsolution := Aspirate(column, mixture, step.Volume, step.Aspiraterate)
		time.Sleep(step.Aspiratepause)
		_, output = Dispense(aspiratedcolumn, aspiratedsolution, step.Volume, step.Dispenserate)
		time.Sleep(step.Dispensepause)
	}

	processedcolumn = column
	output_c = wtype.SolutionToComponent(output)
	return output_c, processedcolumn
}
func Chromatography(input *wtype.LHComponent, step Chromstep, column Column) (output_c *wtype.LHComponent, processedcolumn Column) {

	var output *wtype.LHSolution
	samples := make([]*wtype.LHComponent, 0)
	samples = append(samples, input, step.Buffer)
	mixture := mixer.Mix(samples...)
	for i := 0; i < step.Pipetstep.Cycles; i++ {

		aspiratedcolumn, aspiratedsolution := Aspirate(column, mixture, step.Volume, step.Aspiraterate)
		time.Sleep(step.Aspiratepause)
		_, output = Dispense(aspiratedcolumn, aspiratedsolution, step.Volume, step.Dispenserate)
		time.Sleep(step.Dispensepause)
	}

	processedcolumn = column
	output_c = wtype.SolutionToComponent(output)
	return output_c, processedcolumn
}

func Blot(column Column, blotcycles int, blottime time.Duration) (blottedcolumn Column) {
	// placeholder for moving column to blot paper and dabbing and then waiting
	for i := 0; i < blotcycles; i++ {
		time.Sleep(blottime)
	}
	return blottedcolumn
}

func Dry(tips Column /*wtype.LHTip*/, Drytime time.Duration, Vacuumstrength float64) (drytips Column) {
	// set vacuum manifold to vacuum strength
	///move tips to vacuum position
	time.Sleep(Drytime)
	drytips = tips
	return drytips
}

/*
func Equilibration (Input Watersolution, step Chromstep, column Column) readycolumn Column {

		readycolumn = step Chromstep.cycles *
	            	(
						column ( aspirate mixture(
									step Chromstep.volume,
									step Chromstep.aspiraterate
									)
									wait step Chromstep.aspiratepause
									dispense mixture(
										step Chromstep.volume,
										step Chromstep.dispenserate
									)
									wait step Chromstep.dispensepause
								)
					)
					return readycolumn
}

func AirChromatography (Input Gas, step Chromstep, column Column) readycolumn Column {
	  mixture := mix (Input, step Chromstep.buffer)
		Output = step Chromstep.cycles *
	            	(
						column ( aspirate mixture(
								step Chromstep.volume,
								step Chromstep.aspiraterate
									)
								wait step Chromstep.aspiratepause
									dispense mixture(
										step Chromstep.volume,
										step Chromstep.dispenserate
									)
									wait step Chromstep.dispensepause
								)
					)
					return readycolumn
}


func Resuspension (Input Pellet, step pipetstep) output Suspension {
	  mixture := mix (step Pipetstep.buffer, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Lysis (Input Suspension, step pipetstep) output Lysate {
	  mixture := mix (step, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Precipitation (Input Suspension, step pipetstep) output Precipitate {
	  mixture := mix (step, Input)
		Output = step Pipetstep.cycles *
	            	(
						aspirate mixture(
								step Pipetstep.volume,
								step Pipetstep.aspiraterate
							)
								wait step Pipetstep.aspiratepause
									dispense mixture(
										step Pipetstep.volume,
										step Pipetstep.dispenserate
								)
								wait step Pipetstep.dispensepause
					)
					return output
}

func Growthcurve () {

}
*/
