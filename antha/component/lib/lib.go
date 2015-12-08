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
	"github.com/antha-lang/antha/antha/component/lib/RemoveRestrictionSites"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_design"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_siteremove_orfcheck"
	"github.com/antha-lang/antha/antha/component/lib/SumVolume"
	"github.com/antha-lang/antha/antha/component/lib/Thawtime"
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
	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
	portMap["Aliquot"] = make(map[string]bool)
	portMap["Aliquot"]["InPlate"] = true
	portMap["Aliquot"]["NumberofAliquots"] = true
	portMap["Aliquot"]["OutPlate"] = true
	portMap["Aliquot"]["Solution"] = true
	portMap["Aliquot"]["SolutionVolume"] = true
	portMap["Aliquot"]["VolumePerAliquot"] = true

	portMap["Aliquot"]["Aliquots"] = false

	portMap["Assaysetup"] = make(map[string]bool)
	portMap["Assaysetup"]["Buffer"] = true
	portMap["Assaysetup"]["Enzyme"] = true
	portMap["Assaysetup"]["EnzymeVolume"] = true
	portMap["Assaysetup"]["NumberofReactions"] = true
	portMap["Assaysetup"]["OutPlate"] = true
	portMap["Assaysetup"]["Substrate"] = true
	portMap["Assaysetup"]["SubstrateVolume"] = true
	portMap["Assaysetup"]["TotalVolume"] = true

	portMap["Assaysetup"]["Reactions"] = false
	portMap["Assaysetup"]["Status"] = false

	portMap["BlastSearch"] = make(map[string]bool)
	portMap["BlastSearch"]["AnthaSeq"] = true

	portMap["BlastSearch"]["Hits"] = false

	portMap["DNA_gel"] = make(map[string]bool)
	portMap["DNA_gel"]["DNAgel"] = true
	portMap["DNA_gel"]["DNAgelrunvolume"] = true
	portMap["DNA_gel"]["DNAladder"] = true
	portMap["DNA_gel"]["Loadingdye"] = true
	portMap["DNA_gel"]["Loadingdyeinsample"] = true
	portMap["DNA_gel"]["Loadingdyevolume"] = true
	portMap["DNA_gel"]["Sampletotest"] = true

	portMap["DNA_gel"]["Loadedgel"] = false

	portMap["Datacrunch"] = make(map[string]bool)
	portMap["Datacrunch"]["DNAConc"] = true
	portMap["Datacrunch"]["DNA_seq"] = true
	portMap["Datacrunch"]["Gene_name"] = true
	portMap["Datacrunch"]["Km"] = true
	portMap["Datacrunch"]["Kmunit"] = true
	portMap["Datacrunch"]["ProteinConc"] = true
	portMap["Datacrunch"]["S"] = true
	portMap["Datacrunch"]["SubstrateConc"] = true
	portMap["Datacrunch"]["SubstrateVol"] = true
	portMap["Datacrunch"]["Substrate_name"] = true
	portMap["Datacrunch"]["Sunit"] = true
	portMap["Datacrunch"]["Vmax"] = true
	portMap["Datacrunch"]["Vmaxunit"] = true
	portMap["Datacrunch"]["Vunit"] = true

	portMap["Datacrunch"]["Orftrue"] = false
	portMap["Datacrunch"]["Status"] = false
	portMap["Datacrunch"]["V"] = false

	portMap["Evaporationrate"] = make(map[string]bool)
	portMap["Evaporationrate"]["Airvelocity"] = true
	portMap["Evaporationrate"]["Executiontime"] = true
	portMap["Evaporationrate"]["Liquid"] = true
	portMap["Evaporationrate"]["Pa"] = true
	portMap["Evaporationrate"]["Platetype"] = true
	portMap["Evaporationrate"]["Relativehumidity"] = true
	portMap["Evaporationrate"]["Temp"] = true
	portMap["Evaporationrate"]["Volumeperwell"] = true

	portMap["Evaporationrate"]["Estimatedevaporationtime"] = false
	portMap["Evaporationrate"]["Evaporatedliquid"] = false
	portMap["Evaporationrate"]["Evaporationrateestimate"] = false
	portMap["Evaporationrate"]["Status"] = false

	portMap["FindPartsthat"] = make(map[string]bool)
	portMap["FindPartsthat"]["OnlyreturnAvailableParts"] = true
	portMap["FindPartsthat"]["OnlyreturnWorkingparts"] = true
	portMap["FindPartsthat"]["Partdescriptions"] = true
	portMap["FindPartsthat"]["Parts"] = true
	portMap["FindPartsthat"]["Parttypes"] = true

	portMap["FindPartsthat"]["FulllistBackupParts"] = false
	portMap["FindPartsthat"]["Status"] = false
	portMap["FindPartsthat"]["Warnings"] = false

	portMap["Iterative_assembly_design"] = make(map[string]bool)
	portMap["Iterative_assembly_design"]["ApprovedEnzymes"] = true
	portMap["Iterative_assembly_design"]["Constructname"] = true
	portMap["Iterative_assembly_design"]["Seqsinorder"] = true
	portMap["Iterative_assembly_design"]["Vector"] = true

	portMap["Iterative_assembly_design"]["BackupEnzymes"] = false
	portMap["Iterative_assembly_design"]["EnzymeUsed"] = false
	portMap["Iterative_assembly_design"]["NewDNASequence"] = false
	portMap["Iterative_assembly_design"]["PartswithOverhangs"] = false
	portMap["Iterative_assembly_design"]["Simulationpass"] = false
	portMap["Iterative_assembly_design"]["Status"] = false
	portMap["Iterative_assembly_design"]["Warnings"] = false

	portMap["Kla"] = make(map[string]bool)
	portMap["Kla"]["D"] = true
	portMap["Kla"]["Fillvolume"] = true
	portMap["Kla"]["Liquid"] = true
	portMap["Kla"]["Platetype"] = true
	portMap["Kla"]["Rpm"] = true
	portMap["Kla"]["Shakertype"] = true
	portMap["Kla"]["TargetRE"] = true

	portMap["Kla"]["CalculatedKla"] = false
	portMap["Kla"]["Flowstate"] = false
	portMap["Kla"]["Ncrit"] = false
	portMap["Kla"]["Necessaryshakerspeed"] = false
	portMap["Kla"]["Status"] = false

	portMap["LookUpMolecule"] = make(map[string]bool)
	portMap["LookUpMolecule"]["Compound"] = true
	portMap["LookUpMolecule"]["Compoundlist"] = true

	portMap["LookUpMolecule"]["Compoundprops"] = false
	portMap["LookUpMolecule"]["Jsonstring"] = false
	portMap["LookUpMolecule"]["List"] = false
	portMap["LookUpMolecule"]["Status"] = false

	portMap["MakeBuffer"] = make(map[string]bool)
	portMap["MakeBuffer"]["Buffername"] = true
	portMap["MakeBuffer"]["Bufferstock"] = true
	portMap["MakeBuffer"]["Bufferstockconc"] = true
	portMap["MakeBuffer"]["Bufferstockvolume"] = true
	portMap["MakeBuffer"]["Diluent"] = true
	portMap["MakeBuffer"]["Diluentname"] = true
	portMap["MakeBuffer"]["Diluentvolume"] = true
	portMap["MakeBuffer"]["FinalConcentration"] = true
	portMap["MakeBuffer"]["FinalVolume"] = true
	portMap["MakeBuffer"]["InPlate"] = true
	portMap["MakeBuffer"]["OutPlate"] = true

	portMap["MakeBuffer"]["Buffer"] = false
	portMap["MakeBuffer"]["Status"] = false

	portMap["MakeMedia"] = make(map[string]bool)
	portMap["MakeMedia"]["LiqComponentVolumes"] = true
	portMap["MakeMedia"]["LiqComponents"] = true
	portMap["MakeMedia"]["Name"] = true
	portMap["MakeMedia"]["PH_setPoint"] = true
	portMap["MakeMedia"]["PH_setPointTemp"] = true
	portMap["MakeMedia"]["PH_tolerance"] = true
	portMap["MakeMedia"]["SolidComponentDensities"] = true
	portMap["MakeMedia"]["SolidComponentMasses"] = true
	portMap["MakeMedia"]["SolidComponents"] = true
	portMap["MakeMedia"]["TotalVolume"] = true
	portMap["MakeMedia"]["Vessel"] = true
	portMap["MakeMedia"]["Water"] = true

	portMap["MakeMedia"]["Media"] = false
	portMap["MakeMedia"]["Status"] = false

	portMap["MoClo_design"] = make(map[string]bool)
	portMap["MoClo_design"]["AssemblyStandard"] = true
	portMap["MoClo_design"]["Constructname"] = true
	portMap["MoClo_design"]["Level"] = true
	portMap["MoClo_design"]["PartMoClotypesinorder"] = true
	portMap["MoClo_design"]["Partsinorder"] = true
	portMap["MoClo_design"]["Vector"] = true

	portMap["MoClo_design"]["NewDNASequence"] = false
	portMap["MoClo_design"]["PartswithOverhangs"] = false
	portMap["MoClo_design"]["Simulationpass"] = false
	portMap["MoClo_design"]["Status"] = false
	portMap["MoClo_design"]["Warnings"] = false

	portMap["NewDNASequence"] = make(map[string]bool)
	portMap["NewDNASequence"]["DNA_seq"] = true
	portMap["NewDNASequence"]["Gene_name"] = true
	portMap["NewDNASequence"]["Linear"] = true
	portMap["NewDNASequence"]["Plasmid"] = true
	portMap["NewDNASequence"]["SingleStranded"] = true

	portMap["NewDNASequence"]["DNA"] = false
	portMap["NewDNASequence"]["DNAwithORFs"] = false
	portMap["NewDNASequence"]["Status"] = false

	portMap["OD"] = make(map[string]bool)
	portMap["OD"]["Blank_absorbance"] = true
	portMap["OD"]["Diluent"] = true
	portMap["OD"]["Diluent_volume"] = true
	portMap["OD"]["Heightof100ulinm"] = true
	portMap["OD"]["ODplate"] = true
	portMap["OD"]["ODtoDCWconversionfactor"] = true
	portMap["OD"]["Sample_volume"] = true
	portMap["OD"]["Sampletotest"] = true
	portMap["OD"]["Wlength"] = true

	portMap["OD"]["Blankcorrected_absorbance"] = false
	portMap["OD"]["Estimateddrycellweight_conc"] = false
	portMap["OD"]["OD"] = false
	portMap["OD"]["Sample_absorbance"] = false

	portMap["PCR"] = make(map[string]bool)
	portMap["PCR"]["Additiveconc"] = true
	portMap["PCR"]["Additives"] = true
	portMap["PCR"]["AnnealingTemp"] = true
	portMap["PCR"]["Annealingtime"] = true
	portMap["PCR"]["Buffer"] = true
	portMap["PCR"]["DNTPS"] = true
	portMap["PCR"]["DNTPconc"] = true
	portMap["PCR"]["Denaturationtime"] = true
	portMap["PCR"]["Extensiontemp"] = true
	portMap["PCR"]["Extensiontime"] = true
	portMap["PCR"]["Finalextensiontime"] = true
	portMap["PCR"]["FwdPrimer"] = true
	portMap["PCR"]["FwdPrimerConc"] = true
	portMap["PCR"]["InitDenaturationtime"] = true
	portMap["PCR"]["Numberofcycles"] = true
	portMap["PCR"]["OutPlate"] = true
	portMap["PCR"]["PCRPolymerase"] = true
	portMap["PCR"]["ReactionVolume"] = true
	portMap["PCR"]["RevPrimer"] = true
	portMap["PCR"]["RevPrimerConc"] = true
	portMap["PCR"]["TargetpolymeraseConcentration"] = true
	portMap["PCR"]["Template"] = true
	portMap["PCR"]["Templatevolume"] = true

	portMap["PCR"]["Reaction"] = false

	portMap["Paintmix"] = make(map[string]bool)
	portMap["Paintmix"]["Colour1"] = true
	portMap["Paintmix"]["Colour1vol"] = true
	portMap["Paintmix"]["Colour2"] = true
	portMap["Paintmix"]["Colour2vol"] = true
	portMap["Paintmix"]["Numberofcopies"] = true
	portMap["Paintmix"]["OutPlate"] = true

	portMap["Paintmix"]["NewColours"] = false
	portMap["Paintmix"]["Status"] = false

	portMap["Phytip_miniprep"] = make(map[string]bool)
	portMap["Phytip_miniprep"]["Airstep"] = true
	portMap["Phytip_miniprep"]["Blotcycles"] = true
	portMap["Phytip_miniprep"]["Blottime"] = true
	portMap["Phytip_miniprep"]["Capturestep"] = true
	portMap["Phytip_miniprep"]["Cellpellet"] = true
	portMap["Phytip_miniprep"]["Drytime"] = true
	portMap["Phytip_miniprep"]["Elutionstep"] = true
	portMap["Phytip_miniprep"]["Equilibrationstep"] = true
	portMap["Phytip_miniprep"]["Lysisstep"] = true
	portMap["Phytip_miniprep"]["Phytips"] = true
	portMap["Phytip_miniprep"]["Precipitationstep"] = true
	portMap["Phytip_miniprep"]["Resuspensionstep"] = true
	portMap["Phytip_miniprep"]["Tips"] = true
	portMap["Phytip_miniprep"]["Vacuum"] = true
	portMap["Phytip_miniprep"]["Vacuumstrength"] = true
	portMap["Phytip_miniprep"]["Washsteps"] = true

	portMap["Phytip_miniprep"]["PlasmidDNAsolution"] = false

	portMap["Printname"] = make(map[string]bool)
	portMap["Printname"]["Name"] = true

	portMap["Printname"]["Fullname"] = false

	portMap["RemoveRestrictionSites"] = make(map[string]bool)
	portMap["RemoveRestrictionSites"]["EnzymeforRestrictionmapping"] = true
	portMap["RemoveRestrictionSites"]["PreserveTranslatedseq"] = true
	portMap["RemoveRestrictionSites"]["RemoveifnotinORF"] = true
	portMap["RemoveRestrictionSites"]["RestrictionsitetoAvoid"] = true
	portMap["RemoveRestrictionSites"]["Sequencekey"] = true

	portMap["RemoveRestrictionSites"]["FragmentSizesfromRestrictionmapping"] = false
	portMap["RemoveRestrictionSites"]["SiteFreeSequence"] = false
	portMap["RemoveRestrictionSites"]["Sitesfoundinoriginal"] = false
	portMap["RemoveRestrictionSites"]["Status"] = false
	portMap["RemoveRestrictionSites"]["Warnings"] = false

	portMap["Scarfree_design"] = make(map[string]bool)
	portMap["Scarfree_design"]["Constructname"] = true
	portMap["Scarfree_design"]["Enzymename"] = true
	portMap["Scarfree_design"]["ORFstoConfirm"] = true
	portMap["Scarfree_design"]["Seqsinorder"] = true
	portMap["Scarfree_design"]["Vector"] = true

	portMap["Scarfree_design"]["NewDNASequence"] = false
	portMap["Scarfree_design"]["ORFmissing"] = false
	portMap["Scarfree_design"]["PartswithOverhangs"] = false
	portMap["Scarfree_design"]["Simulationpass"] = false
	portMap["Scarfree_design"]["Status"] = false
	portMap["Scarfree_design"]["Warnings"] = false

	portMap["Scarfree_siteremove_orfcheck"] = make(map[string]bool)
	portMap["Scarfree_siteremove_orfcheck"]["Constructname"] = true
	portMap["Scarfree_siteremove_orfcheck"]["Enzymename"] = true
	portMap["Scarfree_siteremove_orfcheck"]["ORFstoConfirm"] = true
	portMap["Scarfree_siteremove_orfcheck"]["RemoveproblemRestrictionSites"] = true
	portMap["Scarfree_siteremove_orfcheck"]["Seqsinorder"] = true
	portMap["Scarfree_siteremove_orfcheck"]["Vector"] = true

	portMap["Scarfree_siteremove_orfcheck"]["NewDNASequence"] = false
	portMap["Scarfree_siteremove_orfcheck"]["ORFmissing"] = false
	portMap["Scarfree_siteremove_orfcheck"]["PartswithOverhangs"] = false
	portMap["Scarfree_siteremove_orfcheck"]["Simulationpass"] = false
	portMap["Scarfree_siteremove_orfcheck"]["Status"] = false
	portMap["Scarfree_siteremove_orfcheck"]["Warnings"] = false

	portMap["SumVolume"] = make(map[string]bool)
	portMap["SumVolume"]["A"] = true
	portMap["SumVolume"]["B"] = true
	portMap["SumVolume"]["C"] = true

	portMap["SumVolume"]["Status"] = false
	portMap["SumVolume"]["Sum"] = false

	portMap["Thawtime"] = make(map[string]bool)
	portMap["Thawtime"]["Airvelocity"] = true
	portMap["Thawtime"]["BulkTemp"] = true
	portMap["Thawtime"]["Fillvolume"] = true
	portMap["Thawtime"]["Fudgefactor"] = true
	portMap["Thawtime"]["Liquid"] = true
	portMap["Thawtime"]["Platetype"] = true
	portMap["Thawtime"]["SurfaceTemp"] = true

	portMap["Thawtime"]["Estimatedthawtime"] = false
	portMap["Thawtime"]["Status"] = false
	portMap["Thawtime"]["Thawtimeused"] = false

	portMap["Transformation"] = make(map[string]bool)
	portMap["Transformation"]["AgarPlate"] = true
	portMap["Transformation"]["CompetentCells"] = true
	portMap["Transformation"]["CompetentCellvolumeperassembly"] = true
	portMap["Transformation"]["OutPlate"] = true
	portMap["Transformation"]["Plateoutvolume"] = true
	portMap["Transformation"]["Postplasmidtemp"] = true
	portMap["Transformation"]["Postplasmidtime"] = true
	portMap["Transformation"]["Preplasmidtemp"] = true
	portMap["Transformation"]["Preplasmidtime"] = true
	portMap["Transformation"]["Reaction"] = true
	portMap["Transformation"]["Reactionvolume"] = true
	portMap["Transformation"]["Recoverymedium"] = true
	portMap["Transformation"]["Recoverytemp"] = true
	portMap["Transformation"]["Recoverytime"] = true
	portMap["Transformation"]["Recoveryvolume"] = true

	portMap["Transformation"]["Platedculture"] = false

	portMap["TypeIISAssembly_design"] = make(map[string]bool)
	portMap["TypeIISAssembly_design"]["AssemblyStandard"] = true
	portMap["TypeIISAssembly_design"]["Constructname"] = true
	portMap["TypeIISAssembly_design"]["Level"] = true
	portMap["TypeIISAssembly_design"]["PartMoClotypesinorder"] = true
	portMap["TypeIISAssembly_design"]["Partsinorder"] = true
	portMap["TypeIISAssembly_design"]["RestrictionsitetoAvoid"] = true
	portMap["TypeIISAssembly_design"]["Vector"] = true

	portMap["TypeIISAssembly_design"]["BackupParts"] = false
	portMap["TypeIISAssembly_design"]["NewDNASequence"] = false
	portMap["TypeIISAssembly_design"]["PartswithOverhangs"] = false
	portMap["TypeIISAssembly_design"]["Simulationpass"] = false
	portMap["TypeIISAssembly_design"]["Sitesfound"] = false
	portMap["TypeIISAssembly_design"]["Status"] = false
	portMap["TypeIISAssembly_design"]["Warnings"] = false

	portMap["TypeIISConstructAssembly"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly"]["Atp"] = true
	portMap["TypeIISConstructAssembly"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly"]["Buffer"] = true
	portMap["TypeIISConstructAssembly"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly"]["InPlate"] = true
	portMap["TypeIISConstructAssembly"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly"]["LigVol"] = true
	portMap["TypeIISConstructAssembly"]["Ligase"] = true
	portMap["TypeIISConstructAssembly"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly"]["OutputReactionName"] = true
	portMap["TypeIISConstructAssembly"]["PartNames"] = true
	portMap["TypeIISConstructAssembly"]["PartVols"] = true
	portMap["TypeIISConstructAssembly"]["Parts"] = true
	portMap["TypeIISConstructAssembly"]["ReVol"] = true
	portMap["TypeIISConstructAssembly"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssembly"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly"]["RestrictionEnzyme"] = true
	portMap["TypeIISConstructAssembly"]["Vector"] = true
	portMap["TypeIISConstructAssembly"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly"]["Water"] = true

	portMap["TypeIISConstructAssembly"]["Reaction"] = false

	portMap["TypeIISConstructAssemblyMMX"] = make(map[string]bool)
	portMap["TypeIISConstructAssemblyMMX"]["InPlate"] = true
	portMap["TypeIISConstructAssemblyMMX"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssemblyMMX"]["InactivationTime"] = true
	portMap["TypeIISConstructAssemblyMMX"]["MMXVol"] = true
	portMap["TypeIISConstructAssemblyMMX"]["MasterMix"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutPlate"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutputReactionName"] = true
	portMap["TypeIISConstructAssemblyMMX"]["PartNames"] = true
	portMap["TypeIISConstructAssemblyMMX"]["PartVols"] = true
	portMap["TypeIISConstructAssemblyMMX"]["Parts"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionTime"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssemblyMMX"]["Vector"] = true
	portMap["TypeIISConstructAssemblyMMX"]["VectorVol"] = true
	portMap["TypeIISConstructAssemblyMMX"]["Water"] = true

	portMap["TypeIISConstructAssemblyMMX"]["Reaction"] = false

	portMap["TypeIISConstructAssembly_alt"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly_alt"]["Atp"] = true
	portMap["TypeIISConstructAssembly_alt"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["Buffer"] = true
	portMap["TypeIISConstructAssembly_alt"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["InPlate"] = true
	portMap["TypeIISConstructAssembly_alt"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly_alt"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly_alt"]["LigVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["Ligase"] = true
	portMap["TypeIISConstructAssembly_alt"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartConcs"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartMinVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["PartNames"] = true
	portMap["TypeIISConstructAssembly_alt"]["Parts"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly_alt"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly_alt"]["RestrictionEnzyme"] = true
	portMap["TypeIISConstructAssembly_alt"]["Vector"] = true
	portMap["TypeIISConstructAssembly_alt"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly_alt"]["Water"] = true

	portMap["TypeIISConstructAssembly_alt"]["Reaction"] = false
	portMap["TypeIISConstructAssembly_alt"]["S"] = false

	portMap["TypeIISConstructAssembly_sim"] = make(map[string]bool)
	portMap["TypeIISConstructAssembly_sim"]["Atp"] = true
	portMap["TypeIISConstructAssembly_sim"]["AtpVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["Buffer"] = true
	portMap["TypeIISConstructAssembly_sim"]["BufferVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["Constructname"] = true
	portMap["TypeIISConstructAssembly_sim"]["InPlate"] = true
	portMap["TypeIISConstructAssembly_sim"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssembly_sim"]["InactivationTime"] = true
	portMap["TypeIISConstructAssembly_sim"]["LigVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["Ligase"] = true
	portMap["TypeIISConstructAssembly_sim"]["OutPlate"] = true
	portMap["TypeIISConstructAssembly_sim"]["PartConcs"] = true
	portMap["TypeIISConstructAssembly_sim"]["PartVols"] = true
	portMap["TypeIISConstructAssembly_sim"]["Parts"] = true
	portMap["TypeIISConstructAssembly_sim"]["Partsinorder"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionTime"] = true
	portMap["TypeIISConstructAssembly_sim"]["ReactionVolume"] = true
	portMap["TypeIISConstructAssembly_sim"]["RestrictionEnzyme"] = true
	portMap["TypeIISConstructAssembly_sim"]["Vector"] = true
	portMap["TypeIISConstructAssembly_sim"]["VectorConcentration"] = true
	portMap["TypeIISConstructAssembly_sim"]["VectorVol"] = true
	portMap["TypeIISConstructAssembly_sim"]["Vectordata"] = true
	portMap["TypeIISConstructAssembly_sim"]["Water"] = true

	portMap["TypeIISConstructAssembly_sim"]["MolarratiotoVector"] = false
	portMap["TypeIISConstructAssembly_sim"]["Molesperpart"] = false
	portMap["TypeIISConstructAssembly_sim"]["NewDNASequence"] = false
	portMap["TypeIISConstructAssembly_sim"]["Reaction"] = false
	portMap["TypeIISConstructAssembly_sim"]["Simulationpass"] = false
	portMap["TypeIISConstructAssembly_sim"]["Sitesfound"] = false
	portMap["TypeIISConstructAssembly_sim"]["Status"] = false

	c := make([]ComponentDesc, 0)
	c = append(c, ComponentDesc{Name: "Aliquot", Constructor: Aliquot.NewAliquot})
	c = append(c, ComponentDesc{Name: "Assaysetup", Constructor: Assaysetup.NewAssaysetup})
	c = append(c, ComponentDesc{Name: "BlastSearch", Constructor: BlastSearch.NewBlastSearch})
	c = append(c, ComponentDesc{Name: "DNA_gel", Constructor: DNA_gel.NewDNA_gel})
	c = append(c, ComponentDesc{Name: "Datacrunch", Constructor: Datacrunch.NewDatacrunch})
	c = append(c, ComponentDesc{Name: "Evaporationrate", Constructor: Evaporationrate.NewEvaporationrate})
	c = append(c, ComponentDesc{Name: "FindPartsthat", Constructor: FindPartsthat.NewFindPartsthat})
	c = append(c, ComponentDesc{Name: "Iterative_assembly_design", Constructor: Iterative_assembly_design.NewIterative_assembly_design})
	c = append(c, ComponentDesc{Name: "Kla", Constructor: Kla.NewKla})
	c = append(c, ComponentDesc{Name: "LookUpMolecule", Constructor: LookUpMolecule.NewLookUpMolecule})
	c = append(c, ComponentDesc{Name: "MakeBuffer", Constructor: MakeBuffer.NewMakeBuffer})
	c = append(c, ComponentDesc{Name: "MakeMedia", Constructor: MakeMedia.NewMakeMedia})
	c = append(c, ComponentDesc{Name: "MoClo_design", Constructor: MoClo_design.NewMoClo_design})
	c = append(c, ComponentDesc{Name: "NewDNASequence", Constructor: NewDNASequence.NewNewDNASequence})
	c = append(c, ComponentDesc{Name: "OD", Constructor: OD.NewOD})
	c = append(c, ComponentDesc{Name: "PCR", Constructor: PCR.NewPCR})
	c = append(c, ComponentDesc{Name: "Paintmix", Constructor: Paintmix.NewPaintmix})
	c = append(c, ComponentDesc{Name: "Phytip_miniprep", Constructor: Phytip_miniprep.NewPhytip_miniprep})
	c = append(c, ComponentDesc{Name: "Printname", Constructor: Printname.NewPrintname})
	c = append(c, ComponentDesc{Name: "RemoveRestrictionSites", Constructor: RemoveRestrictionSites.NewRemoveRestrictionSites})
	c = append(c, ComponentDesc{Name: "Scarfree_design", Constructor: Scarfree_design.NewScarfree_design})
	c = append(c, ComponentDesc{Name: "Scarfree_siteremove_orfcheck", Constructor: Scarfree_siteremove_orfcheck.NewScarfree_siteremove_orfcheck})
	c = append(c, ComponentDesc{Name: "SumVolume", Constructor: SumVolume.NewSumVolume})
	c = append(c, ComponentDesc{Name: "Thawtime", Constructor: Thawtime.NewThawtime})
	c = append(c, ComponentDesc{Name: "Transformation", Constructor: Transformation.NewTransformation})
	c = append(c, ComponentDesc{Name: "TypeIISAssembly_design", Constructor: TypeIISAssembly_design.NewTypeIISAssembly_design})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly", Constructor: TypeIISConstructAssembly.NewTypeIISConstructAssembly})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssemblyMMX", Constructor: TypeIISConstructAssemblyMMX.NewTypeIISConstructAssemblyMMX})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_alt", Constructor: TypeIISConstructAssembly_alt.NewTypeIISConstructAssembly_alt})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_sim", Constructor: TypeIISConstructAssembly_sim.NewTypeIISConstructAssembly_sim})

	return c
}
