// Example syntax
package unitoperations

func Separate (culture Culture) supernatant Watersolution, pellet Pellet {
	
}

type Chromstep struct {
		Pipetstep 
		Column
		//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
	}
	
	type Pipetstep struct {
		Name "string"
		volume.Volume //transfer volume and process volume?
		aspiraterate Rate
		aspiratepause Time
		dispenserate Rate
		dispensepause Time
		cycles int
		//Mobilephase ... this is really the buffer; should separate out chromstep from pipette step
	}
	

	
	type Column struct{
		Beadsize
		stationaryphase
		separationproperty
		Diameter length
		Height length
		Packedvolume Volume
		
	}
	
	type Phytips struct { // interface?
		Column
		Tip
	}
	

)

func Chromatography (Input Watersolution, step Chromstep, column Column) output Watersolution {
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
					return output
}

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

