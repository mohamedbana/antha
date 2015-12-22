package main

import "github.com/antha-lang/antha/antha/execute"
import "github.com/antha-lang/antha/flow"
import "os"
import "io"
import "encoding/json"
import "log"
import typeIISConstructAssemblyMMX "github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX"

var (
	exitCode = 0
)


type App struct {
    flow.Graph
}

func NewApp() *App {
    n := new(App)
    n.InitGraphState()

    n.Add(typeIISConstructAssemblyMMX.NewTypeIISConstructAssemblyMMX(), "TypeIISConstructAssemblyMMX")

	n.MapInPort("InactivationTemp", "TypeIISConstructAssemblyMMX", "InactivationTemp")
	n.MapInPort("InactivationTime", "TypeIISConstructAssemblyMMX", "InactivationTime")
	n.MapInPort("MMXVol", "TypeIISConstructAssemblyMMX", "MMXVol")
	n.MapInPort("MasterMix", "TypeIISConstructAssemblyMMX", "MasterMix")
	n.MapInPort("OutPlate", "TypeIISConstructAssemblyMMX", "OutPlate")
	n.MapInPort("OutputLocation", "TypeIISConstructAssemblyMMX", "OutputLocation")
	n.MapInPort("OutputPlateNum", "TypeIISConstructAssemblyMMX", "OutputPlateNum")
	n.MapInPort("OutputReactionName", "TypeIISConstructAssemblyMMX", "OutputReactionName")
	n.MapInPort("PartNames", "TypeIISConstructAssemblyMMX", "PartNames")
	n.MapInPort("PartVols", "TypeIISConstructAssemblyMMX", "PartVols")
	n.MapInPort("Parts", "TypeIISConstructAssemblyMMX", "Parts")
	n.MapInPort("ReactionTemp", "TypeIISConstructAssemblyMMX", "ReactionTemp")
	n.MapInPort("ReactionTime", "TypeIISConstructAssemblyMMX", "ReactionTime")
	n.MapInPort("ReactionVolume", "TypeIISConstructAssemblyMMX", "ReactionVolume")
	n.MapInPort("Water", "TypeIISConstructAssemblyMMX", "Water")

	n.MapOutPort("Reaction", "TypeIISConstructAssemblyMMX", "Reaction")


   return n
}

func referenceMain() {
    net := NewApp()

	InactivationTempChan := make(chan execute.ThreadParam)
    net.SetInPort("InactivationTemp", InactivationTempChan)
	InactivationTimeChan := make(chan execute.ThreadParam)
    net.SetInPort("InactivationTime", InactivationTimeChan)
	MMXVolChan := make(chan execute.ThreadParam)
    net.SetInPort("MMXVol", MMXVolChan)
	MasterMixChan := make(chan execute.ThreadParam)
    net.SetInPort("MasterMix", MasterMixChan)
	OutPlateChan := make(chan execute.ThreadParam)
    net.SetInPort("OutPlate", OutPlateChan)
	OutputLocationChan := make(chan execute.ThreadParam)
    net.SetInPort("OutputLocation", OutputLocationChan)
	OutputPlateNumChan := make(chan execute.ThreadParam)
    net.SetInPort("OutputPlateNum", OutputPlateNumChan)
	OutputReactionNameChan := make(chan execute.ThreadParam)
    net.SetInPort("OutputReactionName", OutputReactionNameChan)
	PartNamesChan := make(chan execute.ThreadParam)
    net.SetInPort("PartNames", PartNamesChan)
	PartVolsChan := make(chan execute.ThreadParam)
    net.SetInPort("PartVols", PartVolsChan)
	PartsChan := make(chan execute.ThreadParam)
    net.SetInPort("Parts", PartsChan)
	ReactionTempChan := make(chan execute.ThreadParam)
    net.SetInPort("ReactionTemp", ReactionTempChan)
	ReactionTimeChan := make(chan execute.ThreadParam)
    net.SetInPort("ReactionTime", ReactionTimeChan)
	ReactionVolumeChan := make(chan execute.ThreadParam)
    net.SetInPort("ReactionVolume", ReactionVolumeChan)
	WaterChan := make(chan execute.ThreadParam)
    net.SetInPort("Water", WaterChan)


	ReactionChan := make(chan execute.ThreadParam)
    net.SetOutPort("Reaction", ReactionChan)


    flow.RunNet(net)

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	log.SetOutput(os.Stderr)

	go func() {
		defer close(InactivationTempChan)
		defer close(InactivationTimeChan)
		defer close(MMXVolChan)
		defer close(MasterMixChan)
		defer close(OutPlateChan)
		defer close(OutputLocationChan)
		defer close(OutputPlateNumChan)
		defer close(OutputReactionNameChan)
		defer close(PartNamesChan)
		defer close(PartVolsChan)
		defer close(PartsChan)
		defer close(ReactionTempChan)
		defer close(ReactionTimeChan)
		defer close(ReactionVolumeChan)
		defer close(WaterChan)


		for {
			var p typeIISConstructAssemblyMMX.TypeIISConstructAssemblyMMXJSONBlock
			if err := dec.Decode(&p); err != nil {
				if err != io.EOF {
					log.Println("Error decoding", err)
				}
				return
			}
			//log.Print(p)
			if p.ID == nil {
				log.Println("Error, no ID")
				continue
			}
			if p.Error == nil {
				log.Println("Error, no error")
				continue
			}
			if p.InactivationTemp != nil {
				param := execute.ThreadParam{Value: *(p.InactivationTemp), ID: *(p.ID), Error: *(p.Error)}
				InactivationTempChan <- param
			}
			if p.InactivationTime != nil {
				param := execute.ThreadParam{Value: *(p.InactivationTime), ID: *(p.ID), Error: *(p.Error)}
				InactivationTimeChan <- param
			}
			if p.MMXVol != nil {
				param := execute.ThreadParam{Value: *(p.MMXVol), ID: *(p.ID), Error: *(p.Error)}
				MMXVolChan <- param
			}
			if p.MasterMix != nil {
				param := execute.ThreadParam{Value: *(p.MasterMix), ID: *(p.ID), Error: *(p.Error)}
				MasterMixChan <- param
			}
			if p.OutPlate != nil {
				param := execute.ThreadParam{Value: *(p.OutPlate), ID: *(p.ID), Error: *(p.Error)}
				OutPlateChan <- param
			}
			if p.OutputLocation != nil {
				param := execute.ThreadParam{Value: *(p.OutputLocation), ID: *(p.ID), Error: *(p.Error)}
				OutputLocationChan <- param
			}
			if p.OutputPlateNum != nil {
				param := execute.ThreadParam{Value: *(p.OutputPlateNum), ID: *(p.ID), Error: *(p.Error)}
				OutputPlateNumChan <- param
			}
			if p.OutputReactionName != nil {
				param := execute.ThreadParam{Value: *(p.OutputReactionName), ID: *(p.ID), Error: *(p.Error)}
				OutputReactionNameChan <- param
			}
			if p.PartNames != nil {
				param := execute.ThreadParam{Value: *(p.PartNames), ID: *(p.ID), Error: *(p.Error)}
				PartNamesChan <- param
			}
			if p.PartVols != nil {
				param := execute.ThreadParam{Value: *(p.PartVols), ID: *(p.ID), Error: *(p.Error)}
				PartVolsChan <- param
			}
			if p.Parts != nil {
				param := execute.ThreadParam{Value: *(p.Parts), ID: *(p.ID), Error: *(p.Error)}
				PartsChan <- param
			}
			if p.ReactionTemp != nil {
				param := execute.ThreadParam{Value: *(p.ReactionTemp), ID: *(p.ID), Error: *(p.Error)}
				ReactionTempChan <- param
			}
			if p.ReactionTime != nil {
				param := execute.ThreadParam{Value: *(p.ReactionTime), ID: *(p.ID), Error: *(p.Error)}
				ReactionTimeChan <- param
			}
			if p.ReactionVolume != nil {
				param := execute.ThreadParam{Value: *(p.ReactionVolume), ID: *(p.ID), Error: *(p.Error)}
				ReactionVolumeChan <- param
			}
			if p.Water != nil {
				param := execute.ThreadParam{Value: *(p.Water), ID: *(p.ID), Error: *(p.Error)}
				WaterChan <- param
			}

		}
	}()

	go func() {
		for sequence := range ReactionChan {
			if err := enc.Encode(&sequence); err != nil {
				log.Println(err)
			}
		}
	}()


	<-net.Wait()
}

func main() {
	referenceMain()
	os.Exit(exitCode)
}
