package lib

import (
	"github.com/antha-lang/antha/antha/component/lib/Aliquot"
	"github.com/antha-lang/antha/antha/component/lib/Assaysetup"
	"github.com/antha-lang/antha/antha/component/lib/BlastSearch"
	"github.com/antha-lang/antha/antha/component/lib/DNA_gel"
	"github.com/antha-lang/antha/antha/component/lib/Datacrunch"
	"github.com/antha-lang/antha/antha/component/lib/Evaporationrate"
	"github.com/antha-lang/antha/antha/component/lib/FindPartsthat"
	"github.com/antha-lang/antha/antha/component/lib/Iterative_assembly_design"
	"github.com/antha-lang/antha/antha/component/lib/Kla"
	"github.com/antha-lang/antha/antha/component/lib/LookUpMolecule"
	"github.com/antha-lang/antha/antha/component/lib/MakeBuffer"
	"github.com/antha-lang/antha/antha/component/lib/MakeMedia"
	"github.com/antha-lang/antha/antha/component/lib/MoClo_design"
	"github.com/antha-lang/antha/antha/component/lib/NewDNASequence"
	"github.com/antha-lang/antha/antha/component/lib/OD"
	"github.com/antha-lang/antha/antha/component/lib/PCR"
	"github.com/antha-lang/antha/antha/component/lib/Paintmix"
	"github.com/antha-lang/antha/antha/component/lib/Phytip_miniprep"
	"github.com/antha-lang/antha/antha/component/lib/Printname"
	"github.com/antha-lang/antha/antha/component/lib/ProtocolName_from_an_file"
	"github.com/antha-lang/antha/antha/component/lib/RemoveRestrictionSites"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_design"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_siteremove_orfcheck"
	"github.com/antha-lang/antha/antha/component/lib/SumVolume"
	"github.com/antha-lang/antha/antha/component/lib/Thawtime"
	"github.com/antha-lang/antha/antha/component/lib/Transfer"
	"github.com/antha-lang/antha/antha/component/lib/Transformation"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISAssembly_design"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssemblyMMX"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_alt"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_sim"
)

type ComponentDesc struct {
	Name        string
	Constructor func() interface{}
}

func GetComponents() []ComponentDesc {
	c := make([]ComponentDesc, 0)
	c = append(c, ComponentDesc{Name: "BlastSearch", Constructor: BlastSearch.New})
	c = append(c, ComponentDesc{Name: "FindPartsthat", Constructor: FindPartsthat.New})
	c = append(c, ComponentDesc{Name: "NewDNASequence", Constructor: NewDNASequence.New})
	c = append(c, ComponentDesc{Name: "RemoveRestrictionSites", Constructor: RemoveRestrictionSites.New})
	c = append(c, ComponentDesc{Name: "Iterative_assembly_design", Constructor: Iterative_assembly_design.New})
	c = append(c, ComponentDesc{Name: "MoClo_design", Constructor: MoClo_design.New})
	c = append(c, ComponentDesc{Name: "Scarfree_design", Constructor: Scarfree_design.New})
	c = append(c, ComponentDesc{Name: "Scarfree_siteremove_orfcheck", Constructor: Scarfree_siteremove_orfcheck.New})
	c = append(c, ComponentDesc{Name: "TypeIISAssembly_design", Constructor: TypeIISAssembly_design.New})
	c = append(c, ComponentDesc{Name: "Datacrunch", Constructor: Datacrunch.New})
	c = append(c, ComponentDesc{Name: "LookUpMolecule", Constructor: LookUpMolecule.New})
	c = append(c, ComponentDesc{Name: "Printname", Constructor: Printname.New})
	c = append(c, ComponentDesc{Name: "SumVolume", Constructor: SumVolume.New})
	c = append(c, ComponentDesc{Name: "Aliquot", Constructor: Aliquot.New})
	c = append(c, ComponentDesc{Name: "Assaysetup", Constructor: Assaysetup.New})
	c = append(c, ComponentDesc{Name: "Paintmix", Constructor: Paintmix.New})
	c = append(c, ComponentDesc{Name: "DNA_gel", Constructor: DNA_gel.New})
	c = append(c, ComponentDesc{Name: "MakeBuffer", Constructor: MakeBuffer.New})
	c = append(c, ComponentDesc{Name: "MakeMedia", Constructor: MakeMedia.New})
	c = append(c, ComponentDesc{Name: "OD", Constructor: OD.New})
	c = append(c, ComponentDesc{Name: "PCR", Constructor: PCR.New})
	c = append(c, ComponentDesc{Name: "Phytip_miniprep", Constructor: Phytip_miniprep.New})
	c = append(c, ComponentDesc{Name: "Transfer", Constructor: Transfer.New})
	c = append(c, ComponentDesc{Name: "Transformation", Constructor: Transformation.New})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly", Constructor: TypeIISConstructAssembly.New})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_alt", Constructor: TypeIISConstructAssembly_alt.New})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssemblyMMX", Constructor: TypeIISConstructAssemblyMMX.New})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_sim", Constructor: TypeIISConstructAssembly_sim.New})
	c = append(c, ComponentDesc{Name: "ProtocolName_from_an_file", Constructor: ProtocolName_from_an_file.New})
	c = append(c, ComponentDesc{Name: "Evaporationrate", Constructor: Evaporationrate.New})
	c = append(c, ComponentDesc{Name: "Kla", Constructor: Kla.New})
	c = append(c, ComponentDesc{Name: "Thawtime", Constructor: Thawtime.New})

	return c
}
