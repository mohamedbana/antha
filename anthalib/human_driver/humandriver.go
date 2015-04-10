package human_driver

import (
	"fmt"
	"github.com/antha-lang/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/anthalib/wtype"
	"regexp"
	"strings"
	"sync"
)

// the driver has to implement the liquid handling driver interface which is defined in liquidhandling

type HumanDriver struct {
	InstructionOutputs []string
	movedefaults       map[string]interface{}
	aspdefaults        map[string]interface{}
	dispdefaults       map[string]interface{}
	loaddefaults       map[string]interface{}
	unloaddefaults     map[string]interface{}
	lock               *sync.Mutex
}

func NewHumanDriver(filename string) *HumanDriver {
	var hd HumanDriver

	hd.InstructionOutputs = make([]string, 6)

	hd.InstructionOutputs[0] = "Aspirate volume %-6.2f_VOLUME%"
	hd.InstructionOutputs[1] = "Dispense volume %-6.2f_VOLUME%"
	hd.InstructionOutputs[2] = "Move to position %d_POSITION% well %s_WELL% height %d_HEIGHT% Offsets: %-6.2f_OFFSETX% X %-6.2f_OFFSETY% Y %-6.2f_OFFSETZ% Z"
	hd.InstructionOutputs[3] = "Load tips"
	hd.InstructionOutputs[4] = "Unload tips"
	hd.InstructionOutputs[5] = "Should not be seen"

	return &hd
}

func (self *HumanDriver) Output(ins liquidhandling.RobotInstruction) string {
	// we get the appropriate output type
	s := self.InstructionOutputs[ins.InstructionType()]
	// then change the meta string into the appropriate one
	return self.ReplacePlaceholders(s, ins)
}

func (self *HumanDriver) ReplacePlaceholders(s string, ins liquidhandling.RobotInstruction) string {
	rx, _ := regexp.Compile("%-?(\\d+(\\.\\d)?)?[defgs]_[A-Za-z]+%")
	loc := rx.FindIndex([]byte(s))

	if loc == nil {
		return s
	}

	// we have a match

	match := s[loc[0] : loc[1]-1]
	pre := s[0:loc[0]]
	post := s[loc[1]:len(s)]

	// make the replacement
	var replacement string
	tx := strings.Split(match, "_")

	/*
		switch(tx[1]){
			case "VOLUME" : replacement=fmt.Sprintf(tx[0], ins.Vol)
			case "SPEED"  : replacement=fmt.Sprintf(tx[0], ins.Speed)
			case "POS"    : replacement=fmt.Sprintf(tx[0], ins.Pos)
			case "WELL"   : replacement=fmt.Sprintf(tx[0], ins.Well)
			case "HEIGHT" : replacement=fmt.Sprintf(tx[0], ins.Height)
			case "OFFSETX": replacement=fmt.Sprintf(tx[0], ins.OffsetX)
			case "OFFSETY": replacement=fmt.Sprintf(tx[0], ins.OffsetY)
			case "OFFSETZ": replacement=fmt.Sprintf(tx[0], ins.OffsetZ)
			default: raiseError(fmt.Sprintf("Illegal parameter: %s", tx[1]))
		}
	*/

	// simpler way

	replacement = fmt.Sprintf(tx[0], ins.GetParameter(tx[1]))

	// rebuild the string
	s = pre + replacement + post

	// now do the next one
	return self.ReplacePlaceholders(s, ins)
}

func getmovedefaults() map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func getaspdefaults() map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func getdispdefaults() map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func getloaddefaults() map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func getunloaddefaults() map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func getWellCollection(wellcoords string, n int) string {
	wc := wtype.MakeWellCoordsXY(wellcoords)
	s := 8*(wc.X-1) + wc.Y
	return fmt.Sprintf("%d:%d", s, s+n)
}

// liquid handling driver functions
func (gd *HumanDriver) Move(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type, head int) liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()

	var lhc liquidhandling.LHCommandStatus
	lhc.OK = true
	lhc.Errorcode = liquidhandling.OK

	return lhc
}

func (gd *HumanDriver) MoveExplicit(deckposition string, wellcoords string, reference int, offsetX, offsetY, offsetZ float64, plate_type *liquidhandling.LHPlate, head int) liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}

func (gd *HumanDriver) MoveRaw(x, y, z float64) liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}

func (gd *HumanDriver) Aspirate(volume float64, overstroke bool, head int, multi int) liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}
func (gd *HumanDriver) Dispense(volume float64, blowout bool, head int, multi int) liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}

func (gd *HumanDriver) LoadTips(head, multi int) liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}

func (gd *HumanDriver) UnloadTips(head, multi int) liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}

func (gd *HumanDriver) SetPipetteSpeed(rate float64) liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}
func (gd *HumanDriver) SetDriveSpeed(drive string, rate float64) liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}
func (gd *HumanDriver) Stop() liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}
func (gd *HumanDriver) Go() liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}
func (gd *HumanDriver) Initialize() liquidhandling.LHCommandStatus {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}
func (gd *HumanDriver) Finalize() liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{true, liquidhandling.OK, ""}
	return lhc
}

func (gd *HumanDriver) SetPositionState(position string, state liquidhandling.LHPositionState) liquidhandling.LHCommandStatus {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhc
}
func (gd *HumanDriver) GetCapabilities() (liquidhandling.LHProperties, liquidhandling.LHCommandStatus) {
	var lhp liquidhandling.LHProperties
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return lhp, lhc
}
func (gd *HumanDriver) GetCurrentPosition(head int) (string, liquidhandling.LHCommandStatus) {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return "", lhc
}
func (gd *HumanDriver) GetPositionState(position string) (string, liquidhandling.LHCommandStatus) {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return "", lhc
}
func (gd *HumanDriver) GetHeadState(head int) (string, liquidhandling.LHCommandStatus) {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	return "", lhc
}
func (gd *HumanDriver) GetStatus() (liquidhandling.LHStatus, liquidhandling.LHCommandStatus) {
	lhc := liquidhandling.LHCommandStatus{false, liquidhandling.NIM, "Not Implemented"}
	lhs := liquidhandling.LHStatus{}
	return lhs, lhc
}
