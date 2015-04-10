package liquidhandling

// driver interface

type LiquidhandlingDriver interface {
	Move(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type, head int) LHCommandStatus
	MoveExplicit(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type *LHPlate, head int) LHCommandStatus
	MoveRaw(x, y, z float64) LHCommandStatus
	Aspirate(volume float64, overstroke bool, head int, multi int) LHCommandStatus
	Dispense(volume float64, blowout bool, head int, multi int) LHCommandStatus
	LoadTips(head, multi int) LHCommandStatus
	UnloadTips(head, multi int) LHCommandStatus
	SetPipetteSpeed(rate float64)
	SetDriveSpeed(drive string, rate float64) LHCommandStatus
	Stop() LHCommandStatus
	Go() LHCommandStatus
	Initialize() LHCommandStatus
	Finalize() LHCommandStatus
	SetPositionState(position string, state LHPositionState) LHCommandStatus
	GetCapabilities() LHProperties
	GetCurrentPosition(head int) (string, LHCommandStatus)
	GetPositionState(position string) (string, LHCommandStatus)
	GetHeadState(head int) (string, LHCommandStatus)
	GetStatus() (LHStatus, LHCommandStatus)
}

type ExtendedLiquidhandlingDriver interface {
	LiquidhandlingDriver
	UnloadHead(param int) LHCommandStatus
	LoadHead(param int) LHCommandStatus
	Wait(time float64) LHCommandStatus
	LightsOn() LHCommandStatus
	LightsOff() LHCommandStatus
	LoadAdaptor(param int) LHCommandStatus
	UnloadAdaptor(param int) LHCommandStatus
	// refactored into other interfaces?
	Open() LHCommandStatus
	Close() LHCommandStatus
}
