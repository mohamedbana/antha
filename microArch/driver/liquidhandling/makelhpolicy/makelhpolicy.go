// /anthalib/driver/liquidhandling/makelhpolicy/makelhpolicy.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package main

import (
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// transcribe this elsewhere; here are the policy items

// ASPENTRYSPEED	float64		mm/s	how fast to move vertically into the liquid
// DEFAULTZSPEED	float64		mm/s	Default up and down speed
// ASPXOFFSET		float64		mm	X offset when aspirating
// DSPXOFFSET		float64		mm	X offset when dispensing
// ASPYOFFSET		float64		mm	Y offset when aspirating
// DSPYOFFSET		float64		mm	Y offset when dispensing
// ASPZOFFSET		float64		mm	Z offset when aspirating
// DSPZOFFSET		float64		mm	Z offset when dispensing
// DSPENTRYSPEED	float64		mm/s	how fast to move vertically into the liquid
// PREMIX		int			n cycles to mix before aspirate	(needs additional params)
// POSTMIX		int			n cycles to mix after dispense	(needs additional params)
// TOUCHOFF		bool			touchoff?
// TOUCHOFFSET		float64		mm	touchoff offset in mm
// ASPSPEED		float64		ml/min	pipette speed for aspiration
// ASP_WAIT		float64		s	time to wait after aspirating
// DSP_WAIT		float64		s	time to wait after dispensing
// DSPREFERENCE		int		code	well top, well bottom etc. as ints 0: well bottom 1: well top 2: liquid level
// TIP_REUSE_LIMIT	int			how many times can we re-use a tip?
// CAN_MULTI		bool			can we use multichannel operations
// CAN_MSA		bool			can we do multi-source aspiration?
// CAN SDD		bool			can we do single-destination dispensing?
// NO_AIR_DISPENSE	bool			mandate no air dispensing can happen
// WASH_OUTER_TIP	bool			wash outside of tip
// WASH_INNER_TIP	bool			wash inside of tip

func MakePolicies() map[string]liquidhandling.LHPolicy {
	pols := make(map[string]liquidhandling.LHPolicy)

	// what policies do we need?
	pols["water"] = MakeWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["foamy"] = MakeFoamyPolicy()
        pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	return pols
}
func MakeLysatePolicy() liquidhandling.LHPolicy {
	lysatepolicy := make(liquidhandling.LHPolicy, 6)
	lysatepolicy["ASPSPEED"] = 1.0
	lysatepolicy["DSPSPEED"] = 1.0
	lysatepolicy["ASP_WAIT"] = 2
	lysatepolicy["DSP_WAIT"] = 2
	lysatepolicy["PRE_MIX"] = 5
	lysatepolicy["CAN_MSA"]= false 
	return lysatepolicy
}	
func MakeProteinPolicy() liquidhandling.LHPolicy {
	proteinpolicy := make(liquidhandling.LHPolicy, 4)
	proteinpolicy["DSPREFERENCE"] = 2
	proteinpolicy["CAN_MULTI"] = true
	proteinpolicy["PRE_MIX"] = 3
	proteinpolicy["CAN_MSA"] = false
	return proteinpolicy
}
func MakeWaterPolicy() liquidhandling.LHPolicy {
	waterpolicy := make(liquidhandling.LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 1
	waterpolicy["TOUCHOFF"] = false
	waterpolicy["CAN_MULTI"] = true
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["DSPZOFFSET"] = -0.5
	return waterpolicy
}

func MakeFoamyPolicy() liquidhandling.LHPolicy {
	foamypolicy := make(liquidhandling.LHPolicy, 5)
	foamypolicy["DSPREFERENCE"] = 1
	foamypolicy["TOUCHOFF"] = true
	foamypolicy["CAN_MULTI"] = true
	foamypolicy["CAN_MSA"] = false
	foamypolicy["CAN_SDD"] = true
	return foamypolicy
}

func MakeCulturePolicy() liquidhandling.LHPolicy {
	culturepolicy := make(liquidhandling.LHPolicy, 10)
	culturepolicy["ASPSPEED"] = 0.06
	culturepolicy["DSPSPEED"] = 0.06
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 1.0
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	return culturepolicy
}

func MakeGlycerolPolicy() liquidhandling.LHPolicy {
	glycerolpolicy := make(liquidhandling.LHPolicy, 5)
	glycerolpolicy["ASP_SPEED"] = 1.0
	glycerolpolicy["DSP_SPEED"] = 1.0
	glycerolpolicy["ASP_WAIT"] = 5
	glycerolpolicy["DSP_WAIT"] = 5
	glycerolpolicy["TOUCHOFF"] = true
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}

func MakeSolventPolicy() liquidhandling.LHPolicy {
	solventpolicy := make(liquidhandling.LHPolicy, 4)
	solventpolicy["PRE_MIX"] = 3
	solventpolicy["DSPREFERENCE"] = 0
	solventpolicy["DSPZOFFSET"] = 1.0
	solventpolicy["NO_AIR_DISPENSE"] = true
	return solventpolicy
}

func MakeDefaultPolicy() liquidhandling.LHPolicy {
	defaultpolicy := make(liquidhandling.LHPolicy, 10)
	defaultpolicy["ASP_SPEED"] = 3.0
	defaultpolicy["DSP_SPEED"] = 3.0
	defaultpolicy["TOUCHOFF"] = true
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 1
	defaultpolicy["DSPZOFFSET"] = -0.5
	defaultpolicy["CAN_MULTI"] = true
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["NO_AIR_DISPENSE"] = false
	return defaultpolicy
}

func main() {
	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := liquidhandling.NewLHPolicyRuleSet()

	for name, policy := range policies {
		rule := liquidhandling.NewLHPolicyRule(name)
		rule.AddCategoryConditionOn("LIQUIDCLASS", name)
		lhpr.AddRule(rule, policy)
	}

	//logger.Debug(fmt.Sprintln(lhpr))
	//str, _ := json.Marshal(lhpr)
	//logger.Debug(fmt.Sprintln(string(str)))
}
