package lib

import (
	"github.com/antha-lang/antha/antha/component/lib/Datacrunch"
	"github.com/antha-lang/antha/antha/component/lib/Evaporationrate"
	"github.com/antha-lang/antha/antha/component/lib/Kla"
	"github.com/antha-lang/antha/antha/component/lib/LookUpMolecule"
	"github.com/antha-lang/antha/antha/component/lib/MakeBuffer"
	"github.com/antha-lang/antha/antha/component/lib/PCR"
	"github.com/antha-lang/antha/antha/component/lib/Sum"
	"github.com/antha-lang/antha/antha/component/lib/SumVolume"
	"github.com/antha-lang/antha/antha/component/lib/Thawtime"
	"github.com/antha-lang/antha/antha/component/lib/Transformation"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISAssembly_design"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_alt"
	"github.com/antha-lang/antha/antha/component/lib/TypeIISConstructAssembly_sim"
)

type ComponentDesc struct {
	Name        string
	Constructor func() interface{}
}

func GetComponents() []ComponentDesc {
	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
	portMap["Datacrunch"] = make(map[string]bool)
	portMap["Datacrunch"]["Gene_name"] = true
	portMap["Datacrunch"]["Vmax"] = true
	portMap["Datacrunch"]["Vmaxunit"] = true
	portMap["Datacrunch"]["DNA_seq"] = true
	portMap["Datacrunch"]["Substrate_name"] = true
	portMap["Datacrunch"]["DNAConc"] = true
	portMap["Datacrunch"]["Vunit"] = true
	portMap["Datacrunch"]["SubstrateVol"] = true
	portMap["Datacrunch"]["Km"] = true
	portMap["Datacrunch"]["Kmunit"] = true
	portMap["Datacrunch"]["SubstrateConc"] = true
	portMap["Datacrunch"]["ProteinConc"] = true
	portMap["Datacrunch"]["S"] = true
	portMap["Datacrunch"]["Sunit"] = true

	portMap["Datacrunch"]["Orftrue"] = false
	portMap["Datacrunch"]["Status"] = false
	portMap["Datacrunch"]["V"] = false

	portMap["Evaporationrate"] = make(map[string]bool)
	portMap["Evaporationrate"]["Relativehumidity"] = true
	portMap["Evaporationrate"]["Airvelocity"] = true
	portMap["Evaporationrate"]["Executiontime"] = true
	portMap["Evaporationrate"]["Liquid"] = true
	portMap["Evaporationrate"]["Platetype"] = true
	portMap["Evaporationrate"]["Volumeperwell"] = true
	portMap["Evaporationrate"]["Pa"] = true
	portMap["Evaporationrate"]["Temp"] = true

	portMap["Evaporationrate"]["Status"] = false
	portMap["Evaporationrate"]["Evaporationrateestimate"] = false
	portMap["Evaporationrate"]["Evaporatedliquid"] = false
	portMap["Evaporationrate"]["Estimatedevaporationtime"] = false

	portMap["Kla"] = make(map[string]bool)
	portMap["Kla"]["D"] = true
	portMap["Kla"]["Platetype"] = true
	portMap["Kla"]["Liquid"] = true
	portMap["Kla"]["Shakertype"] = true
	portMap["Kla"]["Rpm"] = true
	portMap["Kla"]["Fillvolume"] = true
	portMap["Kla"]["TargetRE"] = true

	portMap["Kla"]["CalculatedKla"] = false
	portMap["Kla"]["Ncrit"] = false
	portMap["Kla"]["Status"] = false
	portMap["Kla"]["Flowstate"] = false
	portMap["Kla"]["Necessaryshakerspeed"] = false

	portMap["LookUpMolecule"] = make(map[string]bool)
	portMap["LookUpMolecule"]["Compoundlist"] = true
	portMap["LookUpMolecule"]["Compound"] = true

	portMap["LookUpMolecule"]["Compoundprops"] = false
	portMap["LookUpMolecule"]["List"] = false
	portMap["LookUpMolecule"]["Jsonstring"] = false
	portMap["LookUpMolecule"]["Status"] = false

	portMap["MakeBuffer"] = make(map[string]bool)
	portMap["MakeBuffer"]["Diluentvolume"] = true
	portMap["MakeBuffer"]["Diluentname"] = true
	portMap["MakeBuffer"]["Bufferstock"] = true
	portMap["MakeBuffer"]["InPlate"] = true
	portMap["MakeBuffer"]["FinalConcentration"] = true
	portMap["MakeBuffer"]["Bufferstockvolume"] = true
	portMap["MakeBuffer"]["Bufferstockconc"] = true
	portMap["MakeBuffer"]["FinalVolume"] = true
	portMap["MakeBuffer"]["Diluent"] = true
	portMap["MakeBuffer"]["OutPlate"] = true
	portMap["MakeBuffer"]["Buffername"] = true

	portMap["MakeBuffer"]["Status"] = false
	portMap["MakeBuffer"]["Buffer"] = false

	portMap["PCR"] = make(map[string]bool)
	portMap["PCR"]["Additiveconc"] = true
	portMap["PCR"]["TargetpolymeraseConcentration"] = true
	portMap["PCR"]["Denaturationtime"] = true
	portMap["PCR"]["DNTPS"] = true
	portMap["PCR"]["OutPlate"] = true
	portMap["PCR"]["Additives"] = true
	portMap["PCR"]["DNTPconc"] = true
	portMap["PCR"]["Annealingtime"] = true
	portMap["PCR"]["Extensiontemp"] = true
	portMap["PCR"]["FwdPrimer"] = true
	portMap["PCR"]["Buffer"] = true
	portMap["PCR"]["Template"] = true
	portMap["PCR"]["RevPrimerConc"] = true
	portMap["PCR"]["Templatevolume"] = true
	portMap["PCR"]["InitDenaturationtime"] = true
	portMap["PCR"]["Extensiontime"] = true
	portMap["PCR"]["RevPrimer"] = true
	portMap["PCR"]["PCRPolymerase"] = true
	portMap["PCR"]["ReactionVolume"] = true
	portMap["PCR"]["FwdPrimerConc"] = true
	portMap["PCR"]["Numberofcycles"] = true
	portMap["PCR"]["AnnealingTemp"] = true
	portMap["PCR"]["Finalextensiontime"] = true

	portMap["PCR"]["Reaction"] = false

	portMap["Thawtime"] = make(map[string]bool)
	portMap["Thawtime"]["Liquid"] = true
	portMap["Thawtime"]["Fillvolume"] = true
	portMap["Thawtime"]["Airvelocity"] = true
	portMap["Thawtime"]["SurfaceTemp"] = true
	portMap["Thawtime"]["BulkTemp"] = true
	portMap["Thawtime"]["Fudgefactor"] = true
	portMap["Thawtime"]["Platetype"] = true

	portMap["Thawtime"]["Thawtimeused"] = false
	portMap["Thawtime"]["Status"] = false
	portMap["Thawtime"]["Estimatedthawtime"] = false

	portMap["Transformation"] = make(map[string]bool)
	portMap["Transformation"]["Preplasmidtemp"] = true
	portMap["Transformation"]["Plateoutvolume"] = true
	portMap["Transformation"]["Reaction"] = true
	portMap["Transformation"]["CompetentCells"] = true
	portMap["Transformation"]["Recoverymedium"] = true
	portMap["Transformation"]["CompetentCellvolumeperassembly"] = true
	portMap["Transformation"]["Postplasmidtime"] = true
	portMap["Transformation"]["Postplasmidtemp"] = true
	portMap["Transformation"]["Recoveryvolume"] = true
	portMap["Transformation"]["Preplasmidtime"] = true
	portMap["Transformation"]["OutPlate"] = true
	portMap["Transformation"]["AgarPlate"] = true
	portMap["Transformation"]["Reactionvolume"] = true
	portMap["Transformation"]["Recoverytime"] = true
	portMap["Transformation"]["Recoverytemp"] = true

	portMap["Transformation"]["Platedculture"] = false

	portMap["TypeIISAssembly_design"] = make(map[string]bool)
	portMap["TypeIISAssembly_design"]["AssemblyStandard"] = true
	portMap["TypeIISAssembly_design"]["Level"] = true
	portMap["TypeIISAssembly_design"]["Vector"] = true
	portMap["TypeIISAssembly_design"]["PartMoClotypesinorder"] = true
	portMap["TypeIISAssembly_design"]["RestrictionsitetoAvoid"] = true
	portMap["TypeIISAssembly_design"]["Constructname"] = true
	portMap["TypeIISAssembly_design"]["Partsinorder"] = true

	portMap["TypeIISAssembly_design"]["PartswithOverhangs"] = false
	portMap["TypeIISAssembly_design"]["NewDNASequence"] = false
	portMap["TypeIISAssembly_design"]["Sitesfound"] = false
	portMap["TypeIISAssembly_design"]["BackupParts"] = false
	portMap["TypeIISAssembly_design"]["Warnings"] = false
	portMap["TypeIISAssembly_design"]["Status"] = false
	portMap["TypeIISAssembly_design"]["Simulationpass"] = false

	portMap["TypeIISConstructAssembly"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly"]["LigVol"] = true
	portMap["TypeIISConstructAssembly"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly"]["Water"] = true
	portMap["TypeIISConstructAssembly"]["InPlate"] = true
	portMap["TypeIISConstructAssembly"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly"]["ReVol"] = true
	portMap["TypeIISConstructAssembly"]["Parts"] = true
	portMap["TypeIISConstructAssembly"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly"]["PartNames"] = true
	portMap["TypeIISConstructAssembly"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly"]["Vector"] = true
	portMap["TypeIISConstructAssembly"]["RestrictionEnzyme"] = true
	portMap["TypeIISConstructAssembly"]["Buffer"] = true
	portMap["TypeIISConstructAssembly"]["PartVols"] = true
	portMap["TypeIISConstructAssembly"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly"]["Ligase"] = true
	portMap["TypeIISConstructAssembly"]["Atp"] = true
	portMap["TypeIISConstructAssembly"]["ReactionTemp"] = true

	portMap["TypeIISConstructAssembly"]["Reaction"] = false

	portMap["TypeIISConstructAssembly_alt"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly_alt"]["Buffer"] = true
	portMap["TypeIISConstructAssembly_alt"]["Water"] = true
	portMap["TypeIISConstructAssembly_alt"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly_alt"]["Vector"] = true
	portMap["TypeIISConstructAssembly_alt"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartConcs"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartNames"] = true
	portMap["TypeIISConstructAssembly_alt"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["LigVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly_alt"]["Ligase"] = true
	portMap["TypeIISConstructAssembly_alt"]["Atp"] = true
	portMap["TypeIISConstructAssembly_alt"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly_alt"]["InPlate"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartMinVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["Parts"] = true
	portMap["TypeIISConstructAssembly_alt"]["RestrictionEnzyme"] = true

	portMap["TypeIISConstructAssembly_alt"]["Reaction"] = false
	portMap["TypeIISConstructAssembly_alt"]["S"] = false

	portMap["TypeIISConstructAssembly_sim"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly_sim"]["ReVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["Parts"] = true
	portMap["TypeIISConstructAssembly_sim"]["Atp"] = true
	portMap["TypeIISConstructAssembly_sim"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["VectorConcentration"] = true
	portMap["TypeIISConstructAssembly_sim"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly_sim"]["RestrictionEnzyme"] = true
	portMap["TypeIISConstructAssembly_sim"]["Buffer"] = true
	portMap["TypeIISConstructAssembly_sim"]["Water"] = true
	portMap["TypeIISConstructAssembly_sim"]["Constructname"] = true
	portMap["TypeIISConstructAssembly_sim"]["PartConcs"] = true
	portMap["TypeIISConstructAssembly_sim"]["Partsinorder"] = true
	portMap["TypeIISConstructAssembly_sim"]["Vectordata"] = true
	portMap["TypeIISConstructAssembly_sim"]["Vector"] = true
	portMap["TypeIISConstructAssembly_sim"]["Ligase"] = true
	portMap["TypeIISConstructAssembly_sim"]["InPlate"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly_sim"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["LigVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly_sim"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly_sim"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly_sim"]["PartVols"] = true

	portMap["TypeIISConstructAssembly_sim"]["Reaction"] = false
	portMap["TypeIISConstructAssembly_sim"]["Status"] = false
	portMap["TypeIISConstructAssembly_sim"]["Simulationpass"] = false
	portMap["TypeIISConstructAssembly_sim"]["Molesperpart"] = false
	portMap["TypeIISConstructAssembly_sim"]["MolarratiotoVector"] = false
	portMap["TypeIISConstructAssembly_sim"]["NewDNASequence"] = false
	portMap["TypeIISConstructAssembly_sim"]["Sitesfound"] = false

	portMap["Sum"] = make(map[string]bool)
	portMap["Sum"]["A"] = true
	portMap["Sum"]["B"] = true

	portMap["Sum"]["Sum"] = false

	portMap["SumVolume"] = make(map[string]bool)
	portMap["SumVolume"]["A"] = true
	portMap["SumVolume"]["B"] = true
	portMap["SumVolume"]["C"] = true

	portMap["SumVolume"]["Sum"] = false
	portMap["SumVolume"]["Status"] = false

	c := make([]ComponentDesc, 0)
	c = append(c, ComponentDesc{Name: "Datacrunch", Constructor: Datacrunch.NewDatacrunch})
	c = append(c, ComponentDesc{Name: "Evaporationrate", Constructor: Evaporationrate.NewEvaporationrate})
	c = append(c, ComponentDesc{Name: "Kla", Constructor: Kla.NewKla})
	c = append(c, ComponentDesc{Name: "LookUpMolecule", Constructor: LookUpMolecule.NewLookUpMolecule})
	c = append(c, ComponentDesc{Name: "MakeBuffer", Constructor: MakeBuffer.NewMakeBuffer})
	c = append(c, ComponentDesc{Name: "PCR", Constructor: PCR.NewPCR})
	c = append(c, ComponentDesc{Name: "Thawtime", Constructor: Thawtime.NewThawtime})
	c = append(c, ComponentDesc{Name: "Transformation", Constructor: Transformation.NewTransformation})
	c = append(c, ComponentDesc{Name: "TypeIISAssembly_design", Constructor: TypeIISAssembly_design.NewTypeIISAssembly_design})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly", Constructor: TypeIISConstructAssembly.NewTypeIISConstructAssembly})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_alt", Constructor: TypeIISConstructAssembly_alt.NewTypeIISConstructAssembly_alt})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_sim", Constructor: TypeIISConstructAssembly_sim.NewTypeIISConstructAssembly_sim})
	c = append(c, ComponentDesc{Name: "Sum", Constructor: Sum.NewSum})
	c = append(c, ComponentDesc{Name: "SumVolume", Constructor: SumVolume.NewSumVolume})

	return c
}