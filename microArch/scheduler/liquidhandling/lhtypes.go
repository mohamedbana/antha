// liquidhandling/lhtypes.Go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// structure for defining a request to the liquid handler
type LHRequest struct {
	ID                         string
	BlockID                    wtype.BlockID
	BlockName                  string
	Output_solutions           map[string]*wtype.LHSolution
	Input_solutions            map[string][]*wtype.LHComponent
	Plates                     map[string]*wtype.LHPlate
	Tips                       []*wtype.LHTipbox
	Tip_Type                   *wtype.LHTipbox
	Locats                     []string
	Setup                      wtype.LHSetup
	InstructionSet             *liquidhandling.RobotInstructionSet
	Instructions               []liquidhandling.TerminalRobotInstruction
	Robotfn                    string
	Outputfn                   string
	Input_assignments          map[string][]string
	Output_assignments         []string
	Input_plates               map[string]*wtype.LHPlate
	Output_plates              map[string]*wtype.LHPlate
	Input_platetypes           []*wtype.LHPlate
	Input_major_group_layouts  [][]string
	Input_minor_group_layouts  [][]string
	Input_plate_layout         []string
	Input_Setup_Weights        map[string]float64
	Output_platetype           *wtype.LHPlate
	Output_major_group_layouts [][]string
	Output_minor_group_layouts [][]string
	Output_plate_layout        []string
	Plate_lookup               map[string]string
	Stockconcs                 map[string]float64
	Policies                   *liquidhandling.LHPolicyRuleSet
	Input_order                []string
}

// this function checks requests so we can see early on whether or not they
// are going to cause problems
// TODO: much of this will need to change as the system evolves;
// this is something which must be carefully checked whenever changes
// are made downstream
func ValidateLHRequest(rq *LHRequest) (bool, string) {
	if rq.Output_platetype == nil {
		return false, "No output plate type specified"
	}

	if len(rq.Input_platetypes) == 0 {
		return false, "No input plate types specified"
	}

	if rq.Tip_Type == nil {
		return false, "No tip type specified"
	}

	if rq.Policies == nil {
		return false, "No policies specified"
	}

	return true, "OK"
}

func NewLHRequest() *LHRequest {
	var lhr LHRequest
	lhr.ID = wtype.GetUUID()
	lhr.Output_solutions = make(map[string]*wtype.LHSolution)
	lhr.Input_solutions = make(map[string][]*wtype.LHComponent)
	lhr.Plates = make(map[string]*wtype.LHPlate)
	lhr.Tips = make([]*wtype.LHTipbox, 0, 1)
	lhr.Locats = make([]string, 0, 1)
	lhr.Input_plates = make(map[string]*wtype.LHPlate)
	lhr.Input_platetypes = make([]*wtype.LHPlate, 0, 2)
	lhr.Input_Setup_Weights = make(map[string]float64)
	lhr.Output_plates = make(map[string]*wtype.LHPlate)
	lhr.Input_major_group_layouts = make([][]string, 0, 1)
	lhr.Input_minor_group_layouts = make([][]string, 0, 1)
	lhr.Output_major_group_layouts = make([][]string, 0, 1)
	lhr.Output_minor_group_layouts = make([][]string, 0, 1)
	lhr.Output_plate_layout = make([]string, 0, 1)
	lhr.Plate_lookup = make(map[string]string)
	lhr.Stockconcs = make(map[string]float64)
	lhr.Input_order = make([]string, 0)
	return &lhr
}

type LHPolicyManager struct {
	SystemPolicies *liquidhandling.LHPolicyRuleSet
	UserPolicies   *liquidhandling.LHPolicyRuleSet
}

func (mgr *LHPolicyManager) MergePolicies(protocolpolicies *liquidhandling.LHPolicyRuleSet) *liquidhandling.LHPolicyRuleSet {
	ret := liquidhandling.CloneLHPolicyRuleSet(mgr.SystemPolicies)

	// things coming in take precedence over things already there
	ret.MergeWith(mgr.UserPolicies)
	ret.MergeWith(protocolpolicies)

	return ret
}
