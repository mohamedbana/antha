package lib

import (
	"github.com/antha-lang/antha/antha/component/lib/Aliquot"
	"github.com/antha-lang/antha/antha/component/lib/AliquotTo"
	"github.com/antha-lang/antha/antha/component/lib/Assaysetup"
	"github.com/antha-lang/antha/antha/component/lib/BlastSearch"
	"github.com/antha-lang/antha/antha/component/lib/Colony_PCR"
	"github.com/antha-lang/antha/antha/component/lib/DNA_gel"
	"github.com/antha-lang/antha/antha/component/lib/Datacrunch"
	"github.com/antha-lang/antha/antha/component/lib/Evaporationrate"
	"github.com/antha-lang/antha/antha/component/lib/FindPartsthat"
	"github.com/antha-lang/antha/antha/component/lib/Iterative_assembly_design"
	"github.com/antha-lang/antha/antha/component/lib/Kla"
	"github.com/antha-lang/antha/antha/component/lib/LoadGel"
	"github.com/antha-lang/antha/antha/component/lib/LookUpMolecule"
	"github.com/antha-lang/antha/antha/component/lib/MakeBuffer"
	"github.com/antha-lang/antha/antha/component/lib/MakeMedia"
	"github.com/antha-lang/antha/antha/component/lib/Mastermix"
	"github.com/antha-lang/antha/antha/component/lib/Mastermix_reactions"
	"github.com/antha-lang/antha/antha/component/lib/MoClo_design"
	"github.com/antha-lang/antha/antha/component/lib/NewDNASequence"
	"github.com/antha-lang/antha/antha/component/lib/OD"
	"github.com/antha-lang/antha/antha/component/lib/PCR"
	"github.com/antha-lang/antha/antha/component/lib/Paintmix"
	"github.com/antha-lang/antha/antha/component/lib/Phytip_miniprep"
	"github.com/antha-lang/antha/antha/component/lib/PlateOut"
	"github.com/antha-lang/antha/antha/component/lib/PreIncubation"
	"github.com/antha-lang/antha/antha/component/lib/Printname"
	"github.com/antha-lang/antha/antha/component/lib/ProtocolName_from_an_file"
	"github.com/antha-lang/antha/antha/component/lib/Recovery"
	"github.com/antha-lang/antha/antha/component/lib/RemoveRestrictionSites"
	"github.com/antha-lang/antha/antha/component/lib/RestrictionDigestion"
	"github.com/antha-lang/antha/antha/component/lib/RestrictionDigestion_conc"
	"github.com/antha-lang/antha/antha/component/lib/SDSprep"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_design"
	"github.com/antha-lang/antha/antha/component/lib/Scarfree_siteremove_orfcheck"
	"github.com/antha-lang/antha/antha/component/lib/SumVolume"
	"github.com/antha-lang/antha/antha/component/lib/Thawtime"
	"github.com/antha-lang/antha/antha/component/lib/Transfer"
	"github.com/antha-lang/antha/antha/component/lib/Transformation"
	"github.com/antha-lang/antha/antha/component/lib/Transformation_complete"
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

	portMap["AliquotTo"] = make(map[string]bool)
	portMap["AliquotTo"]["InPlate"] = true
	portMap["AliquotTo"]["NumberofAliquots"] = true
	portMap["AliquotTo"]["OutPlate"] = true
	portMap["AliquotTo"]["Solution"] = true
	portMap["AliquotTo"]["SolutionVolume"] = true
	portMap["AliquotTo"]["VolumePerAliquot"] = true

	portMap["AliquotTo"]["Aliquots"] = false

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

	portMap["Colony_PCR"] = make(map[string]bool)
	portMap["Colony_PCR"]["Additiveconc"] = true
	portMap["Colony_PCR"]["Additives"] = true
	portMap["Colony_PCR"]["AnnealingTemp"] = true
	portMap["Colony_PCR"]["Annealingtime"] = true
	portMap["Colony_PCR"]["Buffer"] = true
	portMap["Colony_PCR"]["DNTPS"] = true
	portMap["Colony_PCR"]["DNTPconc"] = true
	portMap["Colony_PCR"]["Denaturationtime"] = true
	portMap["Colony_PCR"]["Extensiontemp"] = true
	portMap["Colony_PCR"]["Extensiontime"] = true
	portMap["Colony_PCR"]["Finalextensiontime"] = true
	portMap["Colony_PCR"]["FwdPrimer"] = true
	portMap["Colony_PCR"]["FwdPrimerConc"] = true
	portMap["Colony_PCR"]["InitDenaturationtime"] = true
	portMap["Colony_PCR"]["Numberofcycles"] = true
	portMap["Colony_PCR"]["OutPlate"] = true
	portMap["Colony_PCR"]["PCRPolymerase"] = true
	portMap["Colony_PCR"]["ReactionVolume"] = true
	portMap["Colony_PCR"]["RevPrimer"] = true
	portMap["Colony_PCR"]["RevPrimerConc"] = true
	portMap["Colony_PCR"]["TargetpolymeraseConcentration"] = true
	portMap["Colony_PCR"]["Template"] = true
	portMap["Colony_PCR"]["Templatevolume"] = true

	portMap["Colony_PCR"]["Reaction"] = false

	portMap["DNA_gel"] = make(map[string]bool)
	portMap["DNA_gel"]["DNAgel"] = true
	portMap["DNA_gel"]["DNAgelrunvolume"] = true
	portMap["DNA_gel"]["InPlate"] = true
	portMap["DNA_gel"]["Loadingdye"] = true
	portMap["DNA_gel"]["Loadingdyeinsample"] = true
	portMap["DNA_gel"]["Loadingdyevolume"] = true
	portMap["DNA_gel"]["Mixingpolicy"] = true
	portMap["DNA_gel"]["Samplenames"] = true
	portMap["DNA_gel"]["Samplenumber"] = true
	portMap["DNA_gel"]["Sampletotest"] = true
	portMap["DNA_gel"]["Water"] = true
	portMap["DNA_gel"]["Watervol"] = true

	portMap["DNA_gel"]["Loadedsamples"] = false

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

	portMap["LoadGel"] = make(map[string]bool)
	portMap["LoadGel"]["GelPlate"] = true
	portMap["LoadGel"]["InPlate"] = true
	portMap["LoadGel"]["LoadVolume"] = true
	portMap["LoadGel"]["Protein"] = true
	portMap["LoadGel"]["SampleName"] = true
	portMap["LoadGel"]["Water"] = true
	portMap["LoadGel"]["WaterName"] = true
	portMap["LoadGel"]["WaterVolume"] = true

	portMap["LoadGel"]["RunSolution"] = false
	portMap["LoadGel"]["Status"] = false

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

	portMap["Mastermix"] = make(map[string]bool)
	portMap["Mastermix"]["AliquotbyRow"] = true
	portMap["Mastermix"]["Buffer"] = true
	portMap["Mastermix"]["Inplate"] = true
	portMap["Mastermix"]["NumberofMastermixes"] = true
	portMap["Mastermix"]["OtherComponentVolumes"] = true
	portMap["Mastermix"]["OtherComponents"] = true
	portMap["Mastermix"]["OutPlate"] = true
	portMap["Mastermix"]["TotalVolumeperMastermix"] = true

	portMap["Mastermix"]["Mastermixes"] = false
	portMap["Mastermix"]["Status"] = false

	portMap["Mastermix_reactions"] = make(map[string]bool)
	portMap["Mastermix_reactions"]["AliquotbyRow"] = true
	portMap["Mastermix_reactions"]["ComponentVolumesperReaction"] = true
	portMap["Mastermix_reactions"]["Components"] = true
	portMap["Mastermix_reactions"]["Inplate"] = true
	portMap["Mastermix_reactions"]["NumberofMastermixes"] = true
	portMap["Mastermix_reactions"]["OutPlate"] = true
	portMap["Mastermix_reactions"]["Reactionspermastermix"] = true
	portMap["Mastermix_reactions"]["TopUpBuffer"] = true
	portMap["Mastermix_reactions"]["TotalVolumeperreaction"] = true
	portMap["Mastermix_reactions"]["VolumetoLeaveforDNAperreaction"] = true

	portMap["Mastermix_reactions"]["Mastermixes"] = false
	portMap["Mastermix_reactions"]["Status"] = false

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

	portMap["PlateOut"] = make(map[string]bool)
	portMap["PlateOut"]["AgarPlate"] = true
	portMap["PlateOut"]["Diluent"] = true
	portMap["PlateOut"]["DilutionX"] = true
	portMap["PlateOut"]["IncubationTemp"] = true
	portMap["PlateOut"]["IncubationTime"] = true
	portMap["PlateOut"]["Plateoutvolume"] = true
	portMap["PlateOut"]["RecoveredCells"] = true

	portMap["PlateOut"]["Platedculture"] = false

	portMap["PreIncubation"] = make(map[string]bool)
	portMap["PreIncubation"]["CompetentCells"] = true
	portMap["PreIncubation"]["CompetentCellvolumeperassembly"] = true
	portMap["PreIncubation"]["OutPlate"] = true
	portMap["PreIncubation"]["Preplasmidtemp"] = true
	portMap["PreIncubation"]["Preplasmidtime"] = true

	portMap["PreIncubation"]["ReadyCompCells"] = false

	portMap["Printname"] = make(map[string]bool)
	portMap["Printname"]["Name"] = true

	portMap["Printname"]["Fullname"] = false

	portMap["ProtocolName_from_an_file"] = make(map[string]bool)
	portMap["ProtocolName_from_an_file"]["InputVariable"] = true
	portMap["ProtocolName_from_an_file"]["OutPlate"] = true
	portMap["ProtocolName_from_an_file"]["ParameterVariableAsValuewithunit"] = true
	portMap["ProtocolName_from_an_file"]["ParameterVariableAsint"] = true
	portMap["ProtocolName_from_an_file"]["ParameterVariablestring"] = true

	portMap["ProtocolName_from_an_file"]["OutputData"] = false
	portMap["ProtocolName_from_an_file"]["PhysicalOutput"] = false

	portMap["Recovery"] = make(map[string]bool)
	portMap["Recovery"]["AgarPlate"] = true
	portMap["Recovery"]["OutPlate"] = true
	portMap["Recovery"]["Recoverymedium"] = true
	portMap["Recovery"]["Recoverytemp"] = true
	portMap["Recovery"]["Recoverytime"] = true
	portMap["Recovery"]["Recoveryvolume"] = true
	portMap["Recovery"]["Transformedcells"] = true

	portMap["Recovery"]["RecoveredCells"] = false

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

	portMap["RestrictionDigestion"] = make(map[string]bool)
	portMap["RestrictionDigestion"]["BSAoptional"] = true
	portMap["RestrictionDigestion"]["BSAvol"] = true
	portMap["RestrictionDigestion"]["Buffer"] = true
	portMap["RestrictionDigestion"]["BufferVol"] = true
	portMap["RestrictionDigestion"]["DNAName"] = true
	portMap["RestrictionDigestion"]["DNASolution"] = true
	portMap["RestrictionDigestion"]["DNAVol"] = true
	portMap["RestrictionDigestion"]["EnzSolutions"] = true
	portMap["RestrictionDigestion"]["EnzVolumestoadd"] = true
	portMap["RestrictionDigestion"]["EnzymeNames"] = true
	portMap["RestrictionDigestion"]["InPlate"] = true
	portMap["RestrictionDigestion"]["InactivationTemp"] = true
	portMap["RestrictionDigestion"]["InactivationTime"] = true
	portMap["RestrictionDigestion"]["OutPlate"] = true
	portMap["RestrictionDigestion"]["ReactionTemp"] = true
	portMap["RestrictionDigestion"]["ReactionTime"] = true
	portMap["RestrictionDigestion"]["ReactionVolume"] = true
	portMap["RestrictionDigestion"]["Water"] = true

	portMap["RestrictionDigestion"]["Reaction"] = false

	portMap["RestrictionDigestion_conc"] = make(map[string]bool)
	portMap["RestrictionDigestion_conc"]["BSAoptional"] = true
	portMap["RestrictionDigestion_conc"]["BSAvol"] = true
	portMap["RestrictionDigestion_conc"]["Buffer"] = true
	portMap["RestrictionDigestion_conc"]["BufferConcX"] = true
	portMap["RestrictionDigestion_conc"]["DNAConc"] = true
	portMap["RestrictionDigestion_conc"]["DNAMassperReaction"] = true
	portMap["RestrictionDigestion_conc"]["DNAName"] = true
	portMap["RestrictionDigestion_conc"]["DNASolution"] = true
	portMap["RestrictionDigestion_conc"]["DesiredConcinUperml"] = true
	portMap["RestrictionDigestion_conc"]["EnzSolutions"] = true
	portMap["RestrictionDigestion_conc"]["EnzymeNames"] = true
	portMap["RestrictionDigestion_conc"]["InPlate"] = true
	portMap["RestrictionDigestion_conc"]["InactivationTemp"] = true
	portMap["RestrictionDigestion_conc"]["InactivationTime"] = true
	portMap["RestrictionDigestion_conc"]["OutPlate"] = true
	portMap["RestrictionDigestion_conc"]["ReactionTemp"] = true
	portMap["RestrictionDigestion_conc"]["ReactionTime"] = true
	portMap["RestrictionDigestion_conc"]["ReactionVolume"] = true
	portMap["RestrictionDigestion_conc"]["StockReConcinUperml"] = true
	portMap["RestrictionDigestion_conc"]["Water"] = true

	portMap["RestrictionDigestion_conc"]["Reaction"] = false

	portMap["SDSprep"] = make(map[string]bool)
	portMap["SDSprep"]["Buffer"] = true
	portMap["SDSprep"]["BufferName"] = true
	portMap["SDSprep"]["BufferStockConc"] = true
	portMap["SDSprep"]["BufferVolume"] = true
	portMap["SDSprep"]["DenatureTemp"] = true
	portMap["SDSprep"]["DenatureTime"] = true
	portMap["SDSprep"]["FinalConcentration"] = true
	portMap["SDSprep"]["InPlate"] = true
	portMap["SDSprep"]["OutPlate"] = true
	portMap["SDSprep"]["Protein"] = true
	portMap["SDSprep"]["ReactionVolume"] = true
	portMap["SDSprep"]["SampleName"] = true
	portMap["SDSprep"]["SampleVolume"] = true

	portMap["SDSprep"]["LoadSample"] = false
	portMap["SDSprep"]["Status"] = false

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

	portMap["Transfer"] = make(map[string]bool)
	portMap["Transfer"]["LiquidVolume"] = true
	portMap["Transfer"]["Liquidname"] = true
	portMap["Transfer"]["OutPlate"] = true
	portMap["Transfer"]["Startingsolution"] = true

	portMap["Transfer"]["FinalSolution"] = false
	portMap["Transfer"]["Status"] = false

	portMap["Transformation"] = make(map[string]bool)
	portMap["Transformation"]["CompetentCellvolumeperassembly"] = true
	portMap["Transformation"]["OutPlate"] = true
	portMap["Transformation"]["Postplasmidtemp"] = true
	portMap["Transformation"]["Postplasmidtime"] = true
	portMap["Transformation"]["Reaction"] = true
	portMap["Transformation"]["Reactionvolume"] = true
	portMap["Transformation"]["ReadyCompCells"] = true

	portMap["Transformation"]["Transformedcells"] = false

	portMap["Transformation_complete"] = make(map[string]bool)
	portMap["Transformation_complete"]["AgarPlate"] = true
	portMap["Transformation_complete"]["CompetentCells"] = true
	portMap["Transformation_complete"]["CompetentCellvolumeperassembly"] = true
	portMap["Transformation_complete"]["OutPlate"] = true
	portMap["Transformation_complete"]["Plateoutvolume"] = true
	portMap["Transformation_complete"]["Postplasmidtemp"] = true
	portMap["Transformation_complete"]["Postplasmidtime"] = true
	portMap["Transformation_complete"]["Preplasmidtemp"] = true
	portMap["Transformation_complete"]["Preplasmidtime"] = true
	portMap["Transformation_complete"]["Reaction"] = true
	portMap["Transformation_complete"]["Reactionvolume"] = true
	portMap["Transformation_complete"]["Recoverymedium"] = true
	portMap["Transformation_complete"]["Recoverytemp"] = true
	portMap["Transformation_complete"]["Recoverytime"] = true
	portMap["Transformation_complete"]["Recoveryvolume"] = true

	portMap["Transformation_complete"]["Platedculture"] = false

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
	portMap["TypeIISConstructAssemblyMMX"]["InactivationTemp"] = true
	portMap["TypeIISConstructAssemblyMMX"]["InactivationTime"] = true
	portMap["TypeIISConstructAssemblyMMX"]["MMXVol"] = true
	portMap["TypeIISConstructAssemblyMMX"]["MasterMix"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutPlate"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutputLocation"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutputPlateNum"] = true
	portMap["TypeIISConstructAssemblyMMX"]["OutputReactionName"] = true
	portMap["TypeIISConstructAssemblyMMX"]["PartNames"] = true
	portMap["TypeIISConstructAssemblyMMX"]["PartVols"] = true
	portMap["TypeIISConstructAssemblyMMX"]["Parts"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionTemp"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionTime"] = true
	portMap["TypeIISConstructAssemblyMMX"]["ReactionVolume"] = true
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
	c = append(c, ComponentDesc{Name: "AliquotTo", Constructor: AliquotTo.NewAliquotTo})
	c = append(c, ComponentDesc{Name: "Assaysetup", Constructor: Assaysetup.NewAssaysetup})
	c = append(c, ComponentDesc{Name: "BlastSearch", Constructor: BlastSearch.NewBlastSearch})
	c = append(c, ComponentDesc{Name: "Colony_PCR", Constructor: Colony_PCR.NewColony_PCR})
	c = append(c, ComponentDesc{Name: "DNA_gel", Constructor: DNA_gel.NewDNA_gel})
	c = append(c, ComponentDesc{Name: "Datacrunch", Constructor: Datacrunch.NewDatacrunch})
	c = append(c, ComponentDesc{Name: "Evaporationrate", Constructor: Evaporationrate.NewEvaporationrate})
	c = append(c, ComponentDesc{Name: "FindPartsthat", Constructor: FindPartsthat.NewFindPartsthat})
	c = append(c, ComponentDesc{Name: "Iterative_assembly_design", Constructor: Iterative_assembly_design.NewIterative_assembly_design})
	c = append(c, ComponentDesc{Name: "Kla", Constructor: Kla.NewKla})
	c = append(c, ComponentDesc{Name: "LoadGel", Constructor: LoadGel.NewLoadGel})
	c = append(c, ComponentDesc{Name: "LookUpMolecule", Constructor: LookUpMolecule.NewLookUpMolecule})
	c = append(c, ComponentDesc{Name: "MakeBuffer", Constructor: MakeBuffer.NewMakeBuffer})
	c = append(c, ComponentDesc{Name: "MakeMedia", Constructor: MakeMedia.NewMakeMedia})
	c = append(c, ComponentDesc{Name: "MakeMedia", Constructor: MakeMedia.NewMakeMedia})
	c = append(c, ComponentDesc{Name: "Mastermix", Constructor: Mastermix.NewMastermix})
	c = append(c, ComponentDesc{Name: "Mastermix_reactions", Constructor: Mastermix_reactions.NewMastermix_reactions})
	c = append(c, ComponentDesc{Name: "MoClo_design", Constructor: MoClo_design.NewMoClo_design})
	c = append(c, ComponentDesc{Name: "NewDNASequence", Constructor: NewDNASequence.NewNewDNASequence})
	c = append(c, ComponentDesc{Name: "OD", Constructor: OD.NewOD})
	c = append(c, ComponentDesc{Name: "PCR", Constructor: PCR.NewPCR})
	c = append(c, ComponentDesc{Name: "Paintmix", Constructor: Paintmix.NewPaintmix})
	c = append(c, ComponentDesc{Name: "Phytip_miniprep", Constructor: Phytip_miniprep.NewPhytip_miniprep})
	c = append(c, ComponentDesc{Name: "PlateOut", Constructor: PlateOut.NewPlateOut})
	c = append(c, ComponentDesc{Name: "PreIncubation", Constructor: PreIncubation.NewPreIncubation})
	c = append(c, ComponentDesc{Name: "Printname", Constructor: Printname.NewPrintname})
	c = append(c, ComponentDesc{Name: "ProtocolName_from_an_file", Constructor: ProtocolName_from_an_file.NewProtocolName_from_an_file})
	c = append(c, ComponentDesc{Name: "Recovery", Constructor: Recovery.NewRecovery})
	c = append(c, ComponentDesc{Name: "RemoveRestrictionSites", Constructor: RemoveRestrictionSites.NewRemoveRestrictionSites})
	c = append(c, ComponentDesc{Name: "RestrictionDigestion", Constructor: RestrictionDigestion.NewRestrictionDigestion})
	c = append(c, ComponentDesc{Name: "RestrictionDigestion_conc", Constructor: RestrictionDigestion_conc.NewRestrictionDigestion_conc})
	c = append(c, ComponentDesc{Name: "SDSprep", Constructor: SDSprep.NewSDSprep})
	c = append(c, ComponentDesc{Name: "Scarfree_design", Constructor: Scarfree_design.NewScarfree_design})
	c = append(c, ComponentDesc{Name: "Scarfree_siteremove_orfcheck", Constructor: Scarfree_siteremove_orfcheck.NewScarfree_siteremove_orfcheck})
	c = append(c, ComponentDesc{Name: "SumVolume", Constructor: SumVolume.NewSumVolume})
	c = append(c, ComponentDesc{Name: "Thawtime", Constructor: Thawtime.NewThawtime})
	c = append(c, ComponentDesc{Name: "Transfer", Constructor: Transfer.NewTransfer})
	c = append(c, ComponentDesc{Name: "Transformation", Constructor: Transformation.NewTransformation})
	c = append(c, ComponentDesc{Name: "Transformation_complete", Constructor: Transformation_complete.NewTransformation_complete})
	c = append(c, ComponentDesc{Name: "TypeIISAssembly_design", Constructor: TypeIISAssembly_design.NewTypeIISAssembly_design})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly", Constructor: TypeIISConstructAssembly.NewTypeIISConstructAssembly})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssemblyMMX", Constructor: TypeIISConstructAssemblyMMX.NewTypeIISConstructAssemblyMMX})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_alt", Constructor: TypeIISConstructAssembly_alt.NewTypeIISConstructAssembly_alt})
	c = append(c, ComponentDesc{Name: "TypeIISConstructAssembly_sim", Constructor: TypeIISConstructAssembly_sim.NewTypeIISConstructAssembly_sim})

	return c
}
