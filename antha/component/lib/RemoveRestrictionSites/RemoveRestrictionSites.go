// This protocol is intended to check sequences for restriction sites and remove according to
// specified conditions

package RemoveRestrictionSites

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

// Input parameters for this protocol (data)

//wtype.DNASequence

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// i.e. parts to order

// Input Requirement specification
func (e *RemoveRestrictionSites) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *RemoveRestrictionSites) setup(p RemoveRestrictionSitesParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *RemoveRestrictionSites) steps(p RemoveRestrictionSitesParamBlock, r *RemoveRestrictionSitesResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	Sequence := wtype.MakeLinearDNASequence("Test", p.Sequencekey)

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// first lookup enzyme properties for all enzymes and make a new array
	enzlist := make([]wtype.RestrictionEnzyme, 0)
	for _, site := range p.RestrictionsitetoAvoid {
		enzsite := lookup.EnzymeLookup(site)
		enzlist = append(enzlist, enzsite)
	}

	// check for sites in the sequence
	sitesfound := enzymes.Restrictionsitefinder(Sequence, enzlist)

	// if no sites found skip to restriction map stage
	if len(sitesfound) == 0 {
		r.Warnings = "none"
		r.Status = "No sites found in sequence to remove so same sequence returned"
		r.SiteFreeSequence = Sequence
		r.Sitesfoundinoriginal = sitesfound

	} else {

		// make a list of sequences to avoid before modifying the sequence
		allsitestoavoid := make([]string, 0)

		// add all restriction sites (we need this step since the functions coming up require strings)
		for _, enzy := range enzlist {
			allsitestoavoid = append(allsitestoavoid, enzy.RecognitionSequence)
		}

		for _, site := range sitesfound {
			if site.Sitefound {

				var tempseq wtype.DNASequence
				var err error

				orfs := sequences.FindallORFs(Sequence.Seq)
				warnings = append(warnings, text.Print("orfs: ", orfs))
				features := sequences.ORFs2Features(orfs)

				//set up a boolean to change to true if a sequence is found in an ORF
				foundinorf := false
				//set up an index for each orf found with site within it (need enzyme name too but will recheck all anyway!)
				orfswithsites := make([]int, 0)

				if len(orfs) > 0 {
					for i, orf := range orfs {

						// change func to handle this step of making dnaseq first

						dnaseq := wtype.MakeLinearDNASequence("orf"+strconv.Itoa(i), orf.DNASeq)

						foundinorfs := enzymes.Restrictionsitefinder(dnaseq, enzlist) // won't work yet orf is actually type features

						for _, siteinorf := range foundinorfs {
							if siteinorf.Sitefound == true {
								foundinorf = true
							}
						}

						if foundinorf == true {

							warning := text.Print("sites found in orf"+dnaseq.Nm, orf)
							warnings = append(warnings, warning)
						}
					}
				}
				if p.RemoveifnotinORF {
					if foundinorf == false {
						tempseq, err = sequences.RemoveSite(Sequence, site.Enzyme, allsitestoavoid)
						if err != nil {
							warning := text.Print("removal of site failed! improve your algorithm!", err.Error())
							warnings = append(warnings, warning)

						}
						r.SiteFreeSequence = tempseq

						// all done if all sites are not in orfs!
						// make proper remove allsites func
					}
					if foundinorf == true {

						r.SiteFreeSequence, err = sequences.RemoveSitesOutsideofFeatures(Sequence, site.Enzyme.RecognitionSequence, sequences.ReplaceBycomplement, features)
						if err != nil {
							warnings = append(warnings, err.Error())
						}
					}
				} //		}else {
				if p.PreserveTranslatedseq {
					// make func to check codon and swap site to preserve aa sequence product
					for _, orfnumber := range orfswithsites {

						for _, position := range site.Positions("ALL") {
							orfcoordinates := sequences.MakeStartendPair(orfs[orfnumber].StartPosition, orfs[orfnumber].EndPosition)
							tempseq, _, _, err = sequences.ReplaceCodoninORF(tempseq, orfcoordinates, position, allsitestoavoid)
							if err != nil {
								warning := text.Print("removal of site from orf "+strconv.Itoa(orfnumber), " failed! improve your algorithm! "+err.Error())
								warnings = append(warnings, warning)
							}
						}

					}
				}

				r.SiteFreeSequence = tempseq
			}
		}
	}

	// Now let's find out the size of fragments we would get if digested with a common site cutter
	mapenz := lookup.EnzymeLookup(p.EnzymeforRestrictionmapping)

	r.FragmentSizesfromRestrictionmapping = enzymes.RestrictionMapper(Sequence, mapenz)

	// allow the data to be exported by capitalising the first letter of the variable
	r.Sitesfoundinoriginal = sitesfound

	r.Warnings = strings.Join(warnings, ";")

	// Print status
	if r.Status == "" {
		r.Status = fmt.Sprintln("Something went wrong!")
	} else {
		r.Status = fmt.Sprintln(
			text.Print("Warnings:", r.Warnings),
			text.Print("Sequence", Sequence),
			text.Print("Sitesfound", r.Sitesfoundinoriginal),
			text.Print("Test digestion sizes with"+p.EnzymeforRestrictionmapping, r.FragmentSizesfromRestrictionmapping),
		)
	}
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *RemoveRestrictionSites) analysis(p RemoveRestrictionSitesParamBlock, r *RemoveRestrictionSitesResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *RemoveRestrictionSites) validation(p RemoveRestrictionSitesParamBlock, r *RemoveRestrictionSitesResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *RemoveRestrictionSites) Complete(params interface{}) {
	p := params.(RemoveRestrictionSitesParamBlock)
	if p.Error {
		e.FragmentSizesfromRestrictionmapping <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.SiteFreeSequence <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Sitesfoundinoriginal <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Warnings <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(RemoveRestrictionSitesResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.FragmentSizesfromRestrictionmapping <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.SiteFreeSequence <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Sitesfoundinoriginal <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Warnings <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.FragmentSizesfromRestrictionmapping <- execute.ThreadParam{Value: r.FragmentSizesfromRestrictionmapping, ID: p.ID, Error: false}

	e.SiteFreeSequence <- execute.ThreadParam{Value: r.SiteFreeSequence, ID: p.ID, Error: false}

	e.Sitesfoundinoriginal <- execute.ThreadParam{Value: r.Sitesfoundinoriginal, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Warnings <- execute.ThreadParam{Value: r.Warnings, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *RemoveRestrictionSites) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *RemoveRestrictionSites) NewConfig() interface{} {
	return &RemoveRestrictionSitesConfig{}
}

func (e *RemoveRestrictionSites) NewParamBlock() interface{} {
	return &RemoveRestrictionSitesParamBlock{}
}

func NewRemoveRestrictionSites() interface{} { //*RemoveRestrictionSites {
	e := new(RemoveRestrictionSites)
	e.init()
	return e
}

// Mapper function
func (e *RemoveRestrictionSites) Map(m map[string]interface{}) interface{} {
	var res RemoveRestrictionSitesParamBlock
	res.Error = false || m["EnzymeforRestrictionmapping"].(execute.ThreadParam).Error || m["PreserveTranslatedseq"].(execute.ThreadParam).Error || m["RemoveifnotinORF"].(execute.ThreadParam).Error || m["RestrictionsitetoAvoid"].(execute.ThreadParam).Error || m["Sequencekey"].(execute.ThreadParam).Error

	vEnzymeforRestrictionmapping, is := m["EnzymeforRestrictionmapping"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RemoveRestrictionSitesJSONBlock
		json.Unmarshal([]byte(vEnzymeforRestrictionmapping.JSONString), &temp)
		res.EnzymeforRestrictionmapping = *temp.EnzymeforRestrictionmapping
	} else {
		res.EnzymeforRestrictionmapping = m["EnzymeforRestrictionmapping"].(execute.ThreadParam).Value.(string)
	}

	vPreserveTranslatedseq, is := m["PreserveTranslatedseq"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RemoveRestrictionSitesJSONBlock
		json.Unmarshal([]byte(vPreserveTranslatedseq.JSONString), &temp)
		res.PreserveTranslatedseq = *temp.PreserveTranslatedseq
	} else {
		res.PreserveTranslatedseq = m["PreserveTranslatedseq"].(execute.ThreadParam).Value.(bool)
	}

	vRemoveifnotinORF, is := m["RemoveifnotinORF"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RemoveRestrictionSitesJSONBlock
		json.Unmarshal([]byte(vRemoveifnotinORF.JSONString), &temp)
		res.RemoveifnotinORF = *temp.RemoveifnotinORF
	} else {
		res.RemoveifnotinORF = m["RemoveifnotinORF"].(execute.ThreadParam).Value.(bool)
	}

	vRestrictionsitetoAvoid, is := m["RestrictionsitetoAvoid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RemoveRestrictionSitesJSONBlock
		json.Unmarshal([]byte(vRestrictionsitetoAvoid.JSONString), &temp)
		res.RestrictionsitetoAvoid = *temp.RestrictionsitetoAvoid
	} else {
		res.RestrictionsitetoAvoid = m["RestrictionsitetoAvoid"].(execute.ThreadParam).Value.([]string)
	}

	vSequencekey, is := m["Sequencekey"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RemoveRestrictionSitesJSONBlock
		json.Unmarshal([]byte(vSequencekey.JSONString), &temp)
		res.Sequencekey = *temp.Sequencekey
	} else {
		res.Sequencekey = m["Sequencekey"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["EnzymeforRestrictionmapping"].(execute.ThreadParam).ID
	res.BlockID = m["EnzymeforRestrictionmapping"].(execute.ThreadParam).BlockID

	return res
}

func (e *RemoveRestrictionSites) OnEnzymeforRestrictionmapping(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("EnzymeforRestrictionmapping", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RemoveRestrictionSites) OnPreserveTranslatedseq(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PreserveTranslatedseq", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RemoveRestrictionSites) OnRemoveifnotinORF(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RemoveifnotinORF", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RemoveRestrictionSites) OnRestrictionsitetoAvoid(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RestrictionsitetoAvoid", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RemoveRestrictionSites) OnSequencekey(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Sequencekey", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type RemoveRestrictionSites struct {
	flow.Component                      // component "superclass" embedded
	lock                                sync.Mutex
	startup                             sync.Once
	params                              map[execute.ThreadID]*execute.AsyncBag
	EnzymeforRestrictionmapping         <-chan execute.ThreadParam
	PreserveTranslatedseq               <-chan execute.ThreadParam
	RemoveifnotinORF                    <-chan execute.ThreadParam
	RestrictionsitetoAvoid              <-chan execute.ThreadParam
	Sequencekey                         <-chan execute.ThreadParam
	FragmentSizesfromRestrictionmapping chan<- execute.ThreadParam
	SiteFreeSequence                    chan<- execute.ThreadParam
	Sitesfoundinoriginal                chan<- execute.ThreadParam
	Status                              chan<- execute.ThreadParam
	Warnings                            chan<- execute.ThreadParam
}

type RemoveRestrictionSitesParamBlock struct {
	ID                          execute.ThreadID
	BlockID                     execute.BlockID
	Error                       bool
	EnzymeforRestrictionmapping string
	PreserveTranslatedseq       bool
	RemoveifnotinORF            bool
	RestrictionsitetoAvoid      []string
	Sequencekey                 string
}

type RemoveRestrictionSitesConfig struct {
	ID                          execute.ThreadID
	BlockID                     execute.BlockID
	Error                       bool
	EnzymeforRestrictionmapping string
	PreserveTranslatedseq       bool
	RemoveifnotinORF            bool
	RestrictionsitetoAvoid      []string
	Sequencekey                 string
}

type RemoveRestrictionSitesResultBlock struct {
	ID                                  execute.ThreadID
	BlockID                             execute.BlockID
	Error                               bool
	FragmentSizesfromRestrictionmapping []int
	SiteFreeSequence                    wtype.DNASequence
	Sitesfoundinoriginal                []enzymes.Restrictionsites
	Status                              string
	Warnings                            string
}

type RemoveRestrictionSitesJSONBlock struct {
	ID                                  *execute.ThreadID
	BlockID                             *execute.BlockID
	Error                               *bool
	EnzymeforRestrictionmapping         *string
	PreserveTranslatedseq               *bool
	RemoveifnotinORF                    *bool
	RestrictionsitetoAvoid              *[]string
	Sequencekey                         *string
	FragmentSizesfromRestrictionmapping *[]int
	SiteFreeSequence                    *wtype.DNASequence
	Sitesfoundinoriginal                *[]enzymes.Restrictionsites
	Status                              *string
	Warnings                            *string
}

func (c *RemoveRestrictionSites) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("EnzymeforRestrictionmapping", "string", "EnzymeforRestrictionmapping", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PreserveTranslatedseq", "bool", "PreserveTranslatedseq", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RemoveifnotinORF", "bool", "RemoveifnotinORF", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RestrictionsitetoAvoid", "[]string", "RestrictionsitetoAvoid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sequencekey", "string", "Sequencekey", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("FragmentSizesfromRestrictionmapping", "[]int", "FragmentSizesfromRestrictionmapping", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("SiteFreeSequence", "wtype.DNASequence", "SiteFreeSequence", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sitesfoundinoriginal", "[]enzymes.Restrictionsites", "Sitesfoundinoriginal", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Warnings", "string", "Warnings", true, true, nil, nil))

	ci := execute.NewComponentInfo("RemoveRestrictionSites", "RemoveRestrictionSites", "", false, inp, outp)

	return ci
}
