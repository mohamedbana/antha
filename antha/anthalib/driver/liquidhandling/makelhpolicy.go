// /anthalib/driver/liquidhandling/makelhpolicy.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

func MakePolicies() map[string]LHPolicy {
	pols := make(map[string]LHPolicy)

	// what policies do we need?
	pols["water"] = MakeWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["foamy"] = MakeFoamyPolicy()
	pols["dna"] = MakeDNAPolicy()
	return pols
}

func MakeWaterPolicy() LHPolicy {
	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 1
	waterpolicy["CAN_MULTI"] = true
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["DSPZOFFSET"] = -0.5
	return waterpolicy
}

func MakeFoamyPolicy() LHPolicy {
	foamypolicy := make(LHPolicy, 5)
	foamypolicy["DSPREFERENCE"] = 1
	foamypolicy["TOUCHOFF"] = true
	foamypolicy["CAN_MULTI"] = true
	foamypolicy["CAN_MSA"] = false
	foamypolicy["CAN_SDD"] = true
	return foamypolicy
}

func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
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

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 5)
	glycerolpolicy["ASP_SPEED"] = 1.0
	glycerolpolicy["DSP_SPEED"] = 1.0
	glycerolpolicy["ASP_WAIT"] = 5
	glycerolpolicy["DSP_WAIT"] = 5
	glycerolpolicy["TOUCHOFF"] = true
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}

func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 4)
	solventpolicy["PRE_MIX"] = 3
	solventpolicy["DSPREFERENCE"] = 0
	solventpolicy["DSPZOFFSET"] = 1.0
	solventpolicy["NO_AIR_DISPENSE"] = true
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_VOLUME"] = 50
	dnapolicy["ASPSPEED"] = 0.06
	dnapolicy["DSPSPEED"] = 0.06
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 1.0
	dnapolicy["TIP_REUSE_LIMIT"] = 5
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy
}

func MakeDefaultPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 10)
	defaultpolicy["ASP_SPEED"] = 3.0
	defaultpolicy["DSP_SPEED"] = 3.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 1
	defaultpolicy["DSPZOFFSET"] = -0.5
	defaultpolicy["CAN_MULTI"] = true
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTOFFSET"] = -0.5
	defaultpolicy["BLOWOUTVOLUME"] = 200.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = false
	defaultpolicy["DEFAULTPIPETTESPEED"] = 1.0
	return defaultpolicy
}

func GetLHPolicyForTest() *LHPolicyRuleSet {
	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := NewLHPolicyRuleSet()

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		rule.AddCategoryConditionOn("LIQUIDCLASS", name)
		lhpr.AddRule(rule, policy)
	}

	return lhpr
}
