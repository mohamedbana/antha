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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/internal/github.com/ghodss/yaml"
)

func MakePolicies() map[string]LHPolicy {
	pols := make(map[string]LHPolicy)

	// what policies do we need?
	pols["water"] = MakeWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["culturereuse"] = MakeCultureReusePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["foamy"] = MakeFoamyPolicy()
	pols["dna"] = MakeDNAPolicy()
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	pols["Paint"] = MakePaintPolicy()

	//      pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	pols["detergent"] = MakeDetergentPolicy()
	pols["load"] = MakeLoadPolicy()
	pols["loadlow"] = MakeLoadPolicy()
	pols["loadwater"] = MakeLoadWaterPolicy()
	pols["DispenseAboveLiquid"] = MakeDispenseAboveLiquidPolicy()
	//      pols["lysate"] = MakeLysatePolicy()

	/*policies, names := PolicyMaker(Allpairs, "DOE_run", false)
	for i, policy := range policies {
		pols[names[i]] = policy
	}
	*/
	if antha.Anthafileexists("ScreenLHPolicyDOE2.xlsx") {
		fmt.Println("found lhpolicy doe file")
		policies, names, err := PolicyMakerfromDesign("ScreenLHPolicyDOE2.xlsx", "DOE_run")

		for i, policy := range policies {
			pols[names[i]] = policy
		}
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("no lhpolicy doe file found")
	}
	return pols

}

func PolicyMakerfromDesign(dxdesignfilename string, prepend string) (policies []LHPolicy, names []string, err error) {
	runs, err := RunsFromDXDesign(filepath.Join(antha.Dirpath(), dxdesignfilename), []string{"Pre_MIX", "POST_MIX"})
	if err != nil {
		return policies, names, err
	}
	policies, names = PolicyMakerfromRuns(runs, prepend, false)
	return
}

func PolicyMaker(factors []DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	runs := AllCombinations(factors)

	policies, names = PolicyMakerfromRuns(runs, nameprepend, concatfactorlevelsinname)

	return
}

func PolicyMakerfromRuns(runs []Run, nameprepend string, concatfactorlevelsinname bool) (policies []LHPolicy, names []string) {

	names = make([]string, 0)
	policies = make([]LHPolicy, 0)

	//policy := make(LHPolicy, 0)
	policy := MakeDefaultPolicy()
	for _, run := range runs {
		for j, desc := range run.Factordescriptors {
			policy[desc] = run.Setpoints[j]
		}

		// raising runtime error when using concat == true
		if concatfactorlevelsinname {
			name := nameprepend
			for key, value := range policy {
				name = fmt.Sprint(name, "_", key, ":", value)

			}
			//	fmt.Println(name)
		} else {
			names = append(names, nameprepend+strconv.Itoa(run.RunNumber))
		}
		policies = append(policies, policy)
		//fmt.Println("len policy = ", len(policy))
		policy = MakeDefaultPolicy()
	}

	return
}

//func MakeLysatePolicy() LHPolicy {
//        lysatepolicy := make(LHPolicy, 6)
//        lysatepolicy["ASPSPEED"] = 1.0
//        lysatepolicy["DSPSPEED"] = 1.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["DSP_WAIT"] = 2.0
//        lysatepolicy["PRE_MIX"] = 5
//        lysatepolicy["CAN_MSA"]= false
//        return lysatepolicy
//}
//func MakeProteinPolicy() LHPolicy {
//        proteinpolicy := make(LHPolicy, 4)
//        proteinpolicy["DSPREFERENCE"] = 2
//        proteinpolicy["CAN_MULTI"] = true
//        proteinpolicy["PRE_MIX"] = 3
//        proteinpolicy["CAN_MSA"] = false
//        return proteinpolicy
//}

func MakePaintPolicy() LHPolicy {

	policy := make(LHPolicy, 12)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["ASP_SPEED"] = 1.5
	policy["DSP_SPEED"] = 1.5
	policy["ASP_WAIT"] = 1.0
	policy["DSP_WAIT"] = 1.0
	policy["PRE_MIX"] = 3
	policy["POST_MIX"] = 3
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true

	return policy
}

func MakeDispenseAboveLiquidPolicy() LHPolicy {

	policy := make(LHPolicy, 7)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASP_SPEED"] = 1.5
	policy["DSP_SPEED"] = 1.5
	policy["ASP_WAIT"] = 1.0
	policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false

	return policy
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
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	culturepolicy["CAN_MULTI"] = false
	culturepolicy["CAN_MSA"] = false
	culturepolicy["CAN_SDD"] = false
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["TIP_REUSE_LIMIT"] = 0
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

	return culturepolicy
}

func MakeCultureReusePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy["PRE_MIX"] = 2
	culturepolicy["ASPSPEED"] = 2.0
	culturepolicy["DSPSPEED"] = 2.0
	culturepolicy["CAN_MULTI"] = true
	culturepolicy["CAN_MSA"] = true
	culturepolicy["CAN_SDD"] = true
	culturepolicy["DSPREFERENCE"] = 0
	culturepolicy["DSPZOFFSET"] = 0.5
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["BLOWOUTVOLUME"] = 0.0
	culturepolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	culturepolicy["TOUCHOFF"] = false

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
	dnapolicy["POST_MIX_VOLUME"] = 50
	dnapolicy["POST_MIX"] = 3
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

func MakeDetergentPolicy() LHPolicy {
	detergentpolicy := make(LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["ASPSPEED"] = 1.0
	detergentpolicy["DSPSPEED"] = 1.0
	detergentpolicy["CAN_MULTI"] = false
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	return detergentpolicy
}
func MakeProteinPolicy() LHPolicy {
	proteinpolicy := make(LHPolicy, 10)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50
	proteinpolicy["ASPSPEED"] = 2.0
	proteinpolicy["DSPSPEED"] = 2.0
	proteinpolicy["CAN_MULTI"] = false
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	return proteinpolicy
}
func MakeLoadPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	return loadpolicy
}

func MakeLoadWaterPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 100
	return loadpolicy
}

func MakeLoadlowPolicy() LHPolicy {

	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 1.0
	loadpolicy["CAN_MULTI"] = false
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["DSPZOFFSET"] = 0.5
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	return loadpolicy
}

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 10)
	dnapolicy["POST_MIX"] = 4
	dnapolicy["POST_MIX_VOLUME"] = 75
	dnapolicy["ASPSPEED"] = 4.0
	dnapolicy["DSPSPEED"] = 4.0
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
	//jbp["TOUCHOFF"] = true
	return jbp
}

func MakeTOPolicy() LHPolicy {
	top := make(LHPolicy, 1)
	top["TOUCHOFF"] = true
	return top
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
