package mixer

import "github.com/antha-lang/antha/antha/anthalib/wtype"

var (
	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0
	DefaultOpt                  = Opt{
		MaxPlates:            &defaultMaxPlates,
		MaxWells:             &defaultMaxWells,
		ResidualVolumeWeight: &defaultResidualVolumeWeight,
		InputPlateType:       []string{"pcrplate_skirted_riser"},
		OutputPlateType:      []string{"pcrplate_skirted_riser"},
		InputPlates:          []*wtype.LHPlate{},
	}
)

type Opt struct {
	MaxPlates            *float64
	MaxWells             *float64
	ResidualVolumeWeight *float64
	InputPlateType       []string
	OutputPlateType      []string
	TipType              []string
	// Three methods of populating Opt.InputPlates
	InputPlateFiles []string         // From filenames
	InputPlateData  [][]byte         // From contents of files
	InputPlates     []*wtype.LHPlate // Directly
	PlanningVersion *int
}

// Merge two configs together and return the result. Values in the argument
// override those in the receiver.
func (a Opt) Merge(x *Opt) Opt {
	if x == nil {
		return a
	}
	if x.MaxPlates != nil {
		a.MaxPlates = x.MaxPlates
	}
	if x.MaxWells != nil {
		a.MaxWells = x.MaxWells
	}
	if x.ResidualVolumeWeight != nil {
		a.ResidualVolumeWeight = x.ResidualVolumeWeight
	}
	if len(x.InputPlateType) != 0 {
		a.InputPlateType = x.InputPlateType
	}
	if len(x.OutputPlateType) != 0 {
		a.OutputPlateType = x.OutputPlateType
	}
	if len(x.TipType) != 0 {
		a.TipType = x.TipType
	}
	if len(x.InputPlateFiles) != 0 {
		a.InputPlateFiles = x.InputPlateFiles
	}
	if len(x.InputPlateData) != 0 {
		a.InputPlateData = x.InputPlateData
	}
	if len(x.InputPlates) != 0 {
		a.InputPlates = x.InputPlates
	}
	if x.PlanningVersion != nil {
		a.PlanningVersion = x.PlanningVersion
	}

	return a
}
