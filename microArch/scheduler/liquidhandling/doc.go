// liquidhandling/doc.go: Part of the Antha language
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
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

// this is the back-end service which runs a liquid handler
// it needs to receive a packet of liquid-handling requests
// and make them happen

// this is defined using raw map[string] interface{} types to support
// flexibility in the precise composition of structures passed around.
// it may be of use later on to have proper structure definitions for these
// particular types and get compile-time checking on data types but
// it's possible we may actually need to map a given tag to a set of
// types rather than have it resolve down to a single data type

// here's a quick bit of schema for a liquid handler's device characteristics
// no item is mandatory but if included it must have the meaning described

// liquid handling device structure
// {
//     vartype: "liquidhandler"
// 	 model: string				// type of liquid handler
//	 manfr: string				// manufacturer
//	lhtype: "discrete" / "continuous"	// does this work on the basis of flow or pipetting?  I would put acoustic in the latter category
//	tptype: "fixed" / "disposable"		// only applies if lhtype=="discrete"
//	cmnvol: float				// minimum volume as presently configured: convenience term depending on how hard reconfiguration is
//	cmxvol: float				// current maximum volume
//	vlunit: string				// volume unit for interpreting all of the above
//	format: [] string			// what labware types are accepted?
//	cnfvol: [] {prms}			// sets of volume parameters defining sets of simultaneously active limits, e.g. some machines can do
//						// multiple mutually exclusive sets of volumes like 0.5-50, 10-200
//						// configurations should be named
//	curcnf: string				// the name of the current configuration of the machine
//	postns: [position]			// defines the current configuration of each position the unit has
//	nposns: int				// number of positions
// }

// parameter structure
// {
//     vartype: "parameters"
//	    id: guid
//	  name: string				// name of this config
//	minvol: float				// minimum volume
// 	maxvol: float				// maximum volume
//     volunit: string				// volume unit
// }
//

// position structure
//
// {
//   vartype:   "position"
//	 num:	int				// position number
//	name:   string				// name of this position
//	uid :   guid				// unique id of this location
//	maxh:   float				// maximum height allowed here
//	extra:  [device]			// has this list of extra devices mounted
// }
//

// liquid handling request structure
// {
//     	       vartype: "lhrequest"
// 		    id: string				// id of this request
//    output_solutions: [solution] 			// solutions to make, in order
//     input_solutions: [string] []solution		// inputs coming in
//		plates: {string:plate}			// plates coming in
//		locats: [solutionid:plate:well] 	// defines where the solutions will end up
//		 setup: [setup]				// defines where to put the various plates
//		instrx: [instruction]			// actual instructions
//
// }
//

// we need some way to string these requests together
// need to define how we handle different types of solution

// basically about sending a liquid handling policy over

// solution:
// {
//   vartype:      "solution"	//
//        id:      guid 	// temporary id
//	inst:      guid		// id of this solution
//	name:	   string	// name of this solution
//     order:      int		// which order this should be added if used as a component
//  components:    [component]	// what goes into this solution
//    container:   guid		// what is this in? -- only set if inst is set
// containertype:  string	// what sort of container is it in
//    welladdress: string	// actually a wellcoords structure
//   platetype:	   string	// what type of plate is this in
//
// }
//

// component:
// {
//   vartype:   "component"
//        id:   guid		// temporary id
//	inst:   guid		// only set if this is a specific liquid
//     order: 	int		// optional ordering in which component is to be added
//	name:   string		// name of this component
//	type: 	string		// type of liquid
//	 vol: 	float		// volume
//	conc:	float		// concentration
//     vunit:   string		// volume unit
//     cunit:   string		// concentration unit
//      tvol:   float		// total volume
// 	 loc:   string		// where this component is located
//	smax:	float		// maximum solubility
// }
//

// plate:	- a microplate
// {
//   vartype:	"plate"
//	inst:   guid		// only set if this is a specific plate
// 	  id:   guid		// id for this plate; gets re-set to same as inst once that is set
//       loc:   guid		// id of the location of this plate, if inst is set
//	name:   string		// user readable name for plate
//	type:   string		// plate type
//	mnfr:   string		// plate manufacturer
//    nwells:   int		// number of wells
//    wellsx:   int		// wells on horizontal axis
//    wellsy:   int		// wells on vertical axis
//     wells:   [guid]well	// structures representing wells
//    height:	float		// how tall the plate is
//     hunit:   string		// unit in which height is measured
// 	rows:	[][]string	// row-wise, by id
//	cols:	[][]string	// col-wise, by id
//wellcoords:	[string]string	// map from coords to ids
//  welltype:	well		// what type of well it has, for convenience
// }
//

// well: - a well in a microplate
// {
//   vartype: "well"
//	  id: guid
//	inst: guid		// only set if this is a specific well
// plateinst: guid		// the id of the plate this well is part of - as above this should only be set if inst is set
//   plateid: guid		// this can be set to identify a plate before actually assigning a specific instance
// platetype: string
//    coords: string		// which well is this
//	 vol: float		// how much this can hold
//     vunit: string		// unit of volume
//  contents: solution		// what this contains
//	rvol: float		// residual volume
//   currvol: float		// how much current volume
//     shape: int		// enum : 1=round,2=square
//    bottom: int		// enum : 1=flat,2=ushaped,3=conical
//      xdim: float64		// aligned with x axis for plate
//      ydim: float64		// aligned with y axis for plate
//      zdim: float64		// aligned with z axis for plate
//   bottomh: float64
//     dunit: string		// measurement unit
// }

// setup:
// {
//   vartype: setup
//    plates: map[guid]guid	// maps plate ids to location ids
//    inputs: map[guid]guid     // maps solution ids to wells
// }
//
