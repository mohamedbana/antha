// antha/AnthaStandardLibrary/Packages/enzymes/TypeIIs.go: Part of the Antha language
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

// Package for working with enzymes; in particular restriction enzymes
package enzymes

import (
	//"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/buffers"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	//"time"
)

var (
//radius = 0.005 // radius of surface in m

// map of liquid classes
/*
	DNApolymerases = map[string]map[string]float64{
		"Q5": map[string]float64{
			"activity_U/ml_assayconds":     50.0,
			"SperKb_upper":      30,
			"SperKb_lower":      20,          //0.62198 * pws / (pa - pws), // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"KBperSuncertainty": 0.01,        //0.62198 * pw / (pa - pw),   // equations not working
			"Fidelity":          0.000000001, //density, kg 􏰀/ m􏰁3
			"stockconc":         0.01,        //
			"workingconc":       0.0005,
			"extensiontemp":     72.0,
			"meltingtemp":       98.0,
		},
		"Taq": map[string]float64{
			"activity_U":        1.0,
			"SperKb_upper":      90,
			"SperKb_lower":      60,        //0.62198 * pws / (pa - pws), // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"KBperSuncertainty": 0.01,      //0.62198 * pw / (pa - pw),   // equations not working
			"Fidelity":          0.0000001, //density, kg 􏰀/ m􏰁3
			"stockconc":         0.01,      //
			"workingconc":       0.0005,
		},
	}
*/
/*Literaturetemps = map[string]map[string]wunit.Temperature{
	"SapI": map[string]wunit.Temperature{
		"extensiontemp": wunit.NewTemperature(72, "C"),
		"meltingtemp":   wunit.NewTemperature(98, "C"),
	},
	"Taq": map[string]wunit.Temperature{
		"extensiontemp": wunit.NewTemperature(68, "C"),
		"meltingtemp":   wunit.NewTemperature(95, "C"),
	},
}
*/
// what about cutting 3prime to recognition site?
)

/*
type DNAStock struct {
	ID int
	Location
	Storagetemp      wunit.Temperature
	Barcode          int
	Dateofproduction string
	ParentID         int
	DNAProperties
	Volume        wunit.Volume
	Concentration wunit.Concentration
	gpermol       float64
	Solution      bool
	buffer        Buffer
}

type Location struct {
	Room      string
	Equipment string
	Shelf     string
	Plate     string
	Well      string
}

type DNAProperties struct {
	wtype.DNASequence
	dam_methylation bool
	dcm_methylation bool
	CpG_methylation bool
	Doublestranded  bool
}

var Cutsmartbuffer = Buffer{

	{PotassiumAcetate50mM,
		TrisAcetate20mM,
		MagnesiumAcetate10mM,
		BSA100μgperml},
	{7.9, {25, "C"}},
}


*/
/*


type Formula struct {
	moieties []Moiety
}

type Moiety struct {
	Atoms         []Atom
	numberofatoms []int
}

type Atom struct {
	Name   string
	Symbol string
	Mass   float64
}


var KCH =Formula{
	Potassium,
	Acetate
}

var Potassium = Moiety{
	{K},
	{1}
}

var Methyl = Moiety{
	{C,H},
	{1,1}
}

var C = Atom{
	"Carbon",
	"C",
	12.01
}

var H = Atom{
	"Hydrogen",
	"C",
	1.0
}

var K = Atom{
	"Potassium",
	"K",
	39.1
}
*/

type CutPosition struct {
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
	CutDist                           int
	EndLength                         int
}

var SapI = wtype.RestrictionEnzyme{"GCTCTTC", 3, "SapI", "", 1, 4, "", []string{"N"}, []int{91, 1109, 1919, 1920}, "TypeIIs"}
var isoschizomers = []string{"BspQI", "LguI", "PciSI", "VpaK32I"}
var SapIenz = wtype.TypeIIs{SapI, "SapI", isoschizomers, 1, 4}
var BsaI = wtype.RestrictionEnzyme{"GGTCTC", 4, "BsaI", "Eco31I", 1, 5, "?(5)", []string{"N"}, []int{814, 1109, 1912, 1995, 1996}, "TypeIIs"}
var BsaIenz = wtype.TypeIIs{BsaI, "BsaI", []string{"none"}, 1, 5}

var BpiI = wtype.RestrictionEnzyme{"GAAGAC", 4, "BpiI", "BbvII", 2, 6, "", []string{"B"}, []int{718}, "TypeIIs"}
var BpiIenz = wtype.TypeIIs{BpiI, "BpiI", []string{"BbvII", "BbsI", "BpuAI", "BSTV2I"}, 2, 6}

var TypeIIsEnzymeproperties = map[string]wtype.TypeIIs{
	"SAPI": SapIenz,
	"BSAI": BsaIenz,
	"BPII": BpiIenz,
}

type Enzymeproperties struct {
	Class         string
	typeIIsenzyme wtype.TypeIIs
	storageBuffer buffers.SimpleBuffer
	assaybuffer   buffers.SimpleBuffer
	unitsperml    float64
	//Barcode      int
	Manufacturer string
	//Unitdefinition
	//reactionbuffer      buffermixture
	//storagebuffer       buffermixture
	heatinactivation bool
	damsensitivity   bool
	dcmsensitivity   bool
	CpGsensitivity   bool
	staractivity     bool
	//Cutsequence         wtype.DNASequence
	Recognitionsequence string //wtype.DNASequence
	/*Bestmodelatpresent  Stochasticmodel
	brendacode
	supplier
	barcode
	Stochasticmodel*/
}

/*
var Enzymes = map[string]Enzymeproperties{
	"R0569L": {SapIenz, SapIstoragebuffer, Cutsmartbuffer, 10000.0, "NEB", true, false, false, false, false, "GCTCTTC"}, // {Nm: "asasd", Seq: "asdasd"}, {Nm: "asasd", Seq: "asdasd"}},
}*/

/*var NewsapIstock = Restrictionstock{
	"R0569L",
	"R0569L 042141216121",
	"0421412",
	"NEB",
	wunit.Volume{125, "ul"},
	10000.0,
	Enzymes["R0569L"],
	time.Date(2016, time.December, 31, 23, 0, 0, 0, time.GMT),
	time.Date(2014, time.December, 31, 23, 0, 0, 0, time.GMT),
	wunit.Temperature{-20, "C"},
	SapIstoragebuffer,
	"PS-R0569S/L v1.0",
}*/
/*
type Restrictionstock struct {
	Productcode  string
	Barcode      string
	Lot          string
	Supplier     string
	Volumeleft   wunit.Volume
	ConcinUperml float64
	//Unitdefinition
	Enzymeproperties
	Shelflife     wunit.Time
	Dateproduced  wunit.Time // is this the right unit?
	Storagetemp   wunit.Temperature
	Storagebuffer string
	SpecVersion   string
}
*/
