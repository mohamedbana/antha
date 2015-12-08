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
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/internal/github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

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
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	return pols
}

func MakeWaterPolicy() LHPolicy {
	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	waterpolicy["CAN_MULTI"] = true
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["DSPZOFFSET"] = 0.5
	return waterpolicy
}

func MakeFoamyPolicy() LHPolicy {
	foamypolicy := make(LHPolicy, 5)
	foamypolicy["DSPREFERENCE"] = 0
	foamypolicy["TOUCHOFF"] = true
	foamypolicy["CAN_MULTI"] = false
	foamypolicy["CAN_MSA"] = false
	foamypolicy["CAN_SDD"] = true
	return foamypolicy
}

func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["ASPSPEED"] = 0.5
	culturepolicy["DSPSPEED"] = 0.5
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	return culturepolicy
}

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 5)
	glycerolpolicy["ASP_SPEED"] = 1.5
	glycerolpolicy["DSP_SPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}

func MakeViscousPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 4)
	glycerolpolicy["ASP_SPEED"] = 1.5
	glycerolpolicy["DSP_SPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	//glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	return glycerolpolicy
}
func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 4)
	solventpolicy["PRE_MIX"] = 3
	solventpolicy["DSPREFERENCE"] = 0
	solventpolicy["DSPZOFFSET"] = 0.5
	solventpolicy["NO_AIR_DISPENSE"] = true
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_VOLUME"] = 50
	dnapolicy["ASPSPEED"] = 2.0
	dnapolicy["DSPSPEED"] = 2.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	return dnapolicy
}

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_VOLUME"] = 50
	dnapolicy["ASPSPEED"] = 2.0
	dnapolicy["DSPSPEED"] = 2.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
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
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MULTI"] = false
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
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	return defaultpolicy
}

func MakeJBPolicy() LHPolicy {
	jbp := make(LHPolicy, 1)
	jbp["JUSTBLOWOUT"] = true
	return jbp
}

func MakeLVExtraPolicy() LHPolicy {
	lvep := make(LHPolicy, 2)
	lvep["EXTRA_ASP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	lvep["EXTRA_DISP_VOLUME"] = wunit.NewVolume(0.5, "ul")
	return lvep
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

	// add a specific case for transfers of water to dry wells
	// nb for this to really work I think we still need to make sure well volumes
	// are being properly kept in sync

	rule := NewLHPolicyRule("BlowOutToEmptyWells")
	rule.AddNumericConditionOn("WELLTOVOLUME", 0.0, 1.0)
	rule.AddCategoryConditionOn("LIQUIDCLASS", "water")
	pol := MakeJBPolicy()
	lhpr.AddRule(rule, pol)

	// a further refinement: for low volumes we need to add extra volume
	// for aspirate and dispense

	rule = NewLHPolicyRule("ExtraVolumeForLV")
	rule.AddNumericConditionOn("VOLUME", 0.0, 20.0)
	pol = MakeLVExtraPolicy()
	lhpr.AddRule(rule, pol)

	return lhpr
}

func LoadLHPoliciesFromFile() (*LHPolicyRuleSet, error) {
	lhPoliciesFileName := os.Getenv("ANTHA_LHPOLICIES_FILE")
	if lhPoliciesFileName == "" {
		return nil, fmt.Errorf("Env variable ANTHA_LHPOLICIES_FILE not set")
	}
	contents, err := ioutil.ReadFile(lhPoliciesFileName)
	if err != nil {
		return nil, err
	}
	lhprs := NewLHPolicyRuleSet()
	lhprs.Policies = make(map[string]LHPolicy)
	lhprs.Rules = make(map[string]LHPolicyRule)
	//	err = readYAML(contents, lhprs)
	err = readJSON(contents, lhprs)
	if err != nil {
		return nil, err
	}
	return lhprs, nil
}

func readYAML(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := yaml.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}

func readJSON(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := json.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}
