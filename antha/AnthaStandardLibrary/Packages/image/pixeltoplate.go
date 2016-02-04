// pixeltoplate.go
// image.go
package image

import (
	"fmt"
	goimage "image"
	"image/color"
	"image/color/palette"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/internal/github.com/disintegration/imaging"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Colour palette to use // this would relate to a map of components of these available colours in factor

var AvailablePalettes = map[string]color.Palette{
	"Palette1":               palettefromMap(Colourcomponentmap), //Chosencolourpalette,
	"WebSafe":                palette.WebSafe,                    //websafe,
	"Plan9":                  palette.Plan9,
	"ProteinPaintboxVisible": palettefromMap(ProteinPaintboxmap),
	"ProteinPaintboxUV":      palettefromMap(UVProteinPaintboxmap),
}

var AvailableComponentmaps = map[string]map[color.Color]string{
	"Palette1":               Colourcomponentmap, //Chosencolourpalette,
	"ProteinPaintboxVisible": ProteinPaintboxmap,
	"ProteinPaintboxUV":      UVProteinPaintboxmap,
}

var Visibleequivalentmaps = map[string]map[color.Color]string{
	"ProteinPaintboxUV": ProteinPaintboxmap,
}

func ColourtoCMYK(colour color.Color) (cmyk color.CMYK) {
	r, g, b, _ := colour.RGBA()
	cmyk.C, cmyk.M, cmyk.Y, cmyk.K = color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	return
}

var Chosencolourpalette color.Palette = availablecolours //palette.WebSafe
//var websafe color.Palette = palette.WebSafe
var availablecolours = []color.Color{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}, //white
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)},  //blue
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)},  // yellow
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)},   // green
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // red
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)},   // black
}

// map of RGB colour to description for use as key in crossreferencing colour to component in other maps
var Colourcomponentmap = map[color.Color]string{
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "white",
	color.RGBA{R: uint8(13), G: uint8(105), B: uint8(171), A: uint8(255)}:  "blue",
	color.RGBA{R: uint8(245), G: uint8(205), B: uint8(47), A: uint8(255)}:  "yellow",
	color.RGBA{R: uint8(75), G: uint8(151), B: uint8(74), A: uint8(255)}:   "green",
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)}:   "red",
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(0)}:         "black",
}

func palettefromMap(colourmap map[color.Color]string) (palette color.Palette) {

	array := make([]color.Color, 0)

	for key, _ := range colourmap {

		array = append(array, key)
	}

	palette = array
	return

}

func reversepalettemap(colourmap map[color.Color]string) (stringmap map[string]color.Color, err error) {

	stringmap = make(map[string]color.Color, len(colourmap))

	for key, value := range colourmap {

		_, ok := stringmap[value]
		if ok == true {
			alreadyinthere := stringmap[value]

			err = fmt.Errorf("attempt to add value", key, "for key", value, "to stringmap", stringmap, "failed due to duplicate entry", alreadyinthere)
		} else {
			stringmap[value] = key
		}
		fmt.Println("key:", key, "value", value)
	}
	return
}

var ProteinPaintboxmap = map[color.Color]string{
	// under visible light

	// Chromogenic proteins
	color.RGBA{R: uint8(70), G: uint8(105), B: uint8(172), A: uint8(255)}:  "BlitzenBlue",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:   "DreidelTeal",
	color.RGBA{R: uint8(107), G: uint8(80), B: uint8(140), A: uint8(255)}:  "VirginiaViolet",
	color.RGBA{R: uint8(120), G: uint8(76), B: uint8(190), A: uint8(255)}:  "VixenPurple",
	color.RGBA{R: uint8(77), G: uint8(11), B: uint8(137), A: uint8(255)}:   "TinselPurple",
	color.RGBA{R: uint8(82), G: uint8(35), B: uint8(119), A: uint8(255)}:   "MaccabeePurple",
	color.RGBA{R: uint8(152), G: uint8(76), B: uint8(128), A: uint8(255)}:  "DonnerMagenta",
	color.RGBA{R: uint8(159), G: uint8(25), B: uint8(103), A: uint8(255)}:  "CupidPink",
	color.RGBA{R: uint8(206), G: uint8(89), B: uint8(142), A: uint8(255)}:  "SeraphinaPink",
	color.RGBA{R: uint8(215), G: uint8(96), B: uint8(86), A: uint8(255)}:   "ScroogeOrange",
	color.RGBA{R: uint8(228), G: uint8(110), B: uint8(104), A: uint8(255)}: "LeorOrange",

	// fluorescent proteins

	//	color.RGBA{R: uint8(224), G: uint8(120), B: uint8(240), A: uint8(255)}:  "CindylouCFP",
	color.RGBA{R: uint8(224), G: uint8(120), B: uint8(140), A: uint8(255)}: "FrostyCFP",

	// for twinkle B should = uint8(137) but this is the same colour as e.coli so changed it to uint8(138) to avoid error due to duplicate map keys
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(138), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "TwinkleCFP",
	color.RGBA{R: uint8(251), G: uint8(176), B: uint8(0), A: uint8(255)}: "YetiYFP",
	color.RGBA{R: uint8(250), G: uint8(210), B: uint8(0), A: uint8(255)}: "MarleyYFP",
	color.RGBA{R: uint8(255), G: uint8(194), B: uint8(0), A: uint8(255)}: "CratchitYFP",
	color.RGBA{R: uint8(231), G: uint8(173), B: uint8(0), A: uint8(255)}: "KringleYFP",
	//color.RGBA{R: uint8(222), G: uint8(221), B: uint8(68), A: uint8(255)}:     "CometGFP",
	color.RGBA{R: uint8(209), G: uint8(214), B: uint8(0), A: uint8(255)}:   "DasherGFP",
	color.RGBA{R: uint8(225), G: uint8(222), B: uint8(120), A: uint8(255)}: "IvyGFP",
	//color.RGBA{R: uint8(216), G: uint8(231), B: uint8(15), A: uint8(255)}:     "HollyGFP",
	color.RGBA{R: uint8(251), G: uint8(102), B: uint8(79), A: uint8(255)}: "YukonOFP",
	color.RGBA{R: uint8(215), G: uint8(72), B: uint8(76), A: uint8(255)}:  "RudolphRFP",
	color.RGBA{R: uint8(244), G: uint8(63), B: uint8(150), A: uint8(255)}: "FresnoRFP",

	// Extended fluorescent proteins
	color.RGBA{R: uint8(248), G: uint8(64), B: uint8(148), A: uint8(255)}:  "CayenneRFP",
	color.RGBA{R: uint8(241), G: uint8(84), B: uint8(152), A: uint8(255)}:  "GuajilloRFP",
	color.RGBA{R: uint8(247), G: uint8(132), B: uint8(179), A: uint8(255)}: "PaprikaRFP",
	color.RGBA{R: uint8(248), G: uint8(84), B: uint8(149), A: uint8(255)}:  "SerranoRFP",
	color.RGBA{R: uint8(254), G: uint8(253), B: uint8(252), A: uint8(255)}: "EiraCFP",
	color.RGBA{R: uint8(255), G: uint8(255), B: uint8(146), A: uint8(255)}: "BlazeYFP",
	color.RGBA{R: uint8(194), G: uint8(164), B: uint8(72), A: uint8(255)}:  "JuniperGFP",
	color.RGBA{R: uint8(243), G: uint8(138), B: uint8(112), A: uint8(255)}: "TannenGFP",

	// conventional E.coli colour
	color.RGBA{R: uint8(196), G: uint8(183), B: uint8(137), A: uint8(255)}: "E.coli",

	// lacZ expresser (e.g. pUC19) grown on S gal
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}: "veryblack",

	// plus white as a blank (or comment out to use EiraCFP)
	color.RGBA{R: uint8(242), G: uint8(243), B: uint8(242), A: uint8(255)}: "verywhite",
}

var UVProteinPaintboxmap = map[color.Color]string{
	// under UV
	//	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(255), A: uint8(255)}:  "CindylouCFP",
	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(255), A: uint8(255)}: "FrostyCFP",
	color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}: "TwinkleCFP",
	//color.RGBA{R: uint8(27), G: uint8(79), B: uint8(146), A: uint8(255)}:  "TwinkleCFP",
	color.RGBA{R: uint8(253), G: uint8(230), B: uint8(39), A: uint8(255)}: "YetiYFP",
	color.RGBA{R: uint8(236), G: uint8(255), B: uint8(0), A: uint8(255)}:  "MarleyYFP",
	color.RGBA{R: uint8(240), G: uint8(254), B: uint8(0), A: uint8(255)}:  "CratchitYFP",
	color.RGBA{R: uint8(239), G: uint8(255), B: uint8(0), A: uint8(255)}:  "KringleYFP",
	//color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}:     "CometGFP",
	color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}:   "DasherGFP",
	color.RGBA{R: uint8(0), G: uint8(232), B: uint8(216), A: uint8(255)}: "IvyGFP",
	//color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}:     "HollyGFP",
	color.RGBA{R: uint8(254), G: uint8(179), B: uint8(18), A: uint8(255)}: "YukonOFP",
	color.RGBA{R: uint8(218), G: uint8(92), B: uint8(69), A: uint8(255)}:  "RudolphRFP",
	color.RGBA{R: uint8(255), G: uint8(0), B: uint8(166), A: uint8(255)}:  "FresnoRFP",

	// Extended fluorescent proteins
	color.RGBA{R: uint8(255), G: uint8(24), B: uint8(138), A: uint8(255)}:  "CayenneRFP",
	color.RGBA{R: uint8(255), G: uint8(8), B: uint8(138), A: uint8(255)}:   "GuajilloRFP",
	color.RGBA{R: uint8(252), G: uint8(65), B: uint8(136), A: uint8(255)}:  "PaprikaRFP",
	color.RGBA{R: uint8(254), G: uint8(23), B: uint8(127), A: uint8(255)}:  "SerranoRFP",
	color.RGBA{R: uint8(173), G: uint8(253), B: uint8(218), A: uint8(255)}: "EiraCFP",
	color.RGBA{R: uint8(254), G: uint8(255), B: uint8(83), A: uint8(255)}:  "BlazeYFP",
	color.RGBA{R: uint8(0), G: uint8(231), B: uint8(162), A: uint8(255)}:   "JuniperGFP",
	color.RGBA{R: uint8(179), G: uint8(119), B: uint8(57), A: uint8(255)}:  "TannenGFP",

	// conventional E.coli colour is black under UV ??
	color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}: "E.coli",
}

func ResizeImagetoPlate(imagefilename string, plate *wtype.LHPlate, algorithm imaging.ResampleFilter) (plateimage *goimage.NRGBA) {

	// input files (just 1 in this case)
	files := []string{imagefilename}

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	img, err := imaging.Open(files[0])
	if err != nil {
		panic(err)
	}

	// have the option of changing the resize algorithm here
	plateimage = imaging.Resize(img, 0, plate.WlsY, algorithm)
	//plateimages = append(plateimages,plateimage)

	return

}

// create a map of pixel to plate position from processing a given image with a chosen colour palette.
// It's recommended to use at least 384 well plate
func ImagetoPlatelayout(imagefilename string, plate *wtype.LHPlate, chosencolourpalette color.Palette) (wellpositiontocolourmap map[string]color.Color, numberofpixels int) {

	plateimage := ResizeImagetoPlate(imagefilename, plate, imaging.CatmullRom)

	// make map of well position to colour: (array for time being)

	wellpositionarray := make([]string, 0)
	colourarray := make([]color.Color, 0)
	wellpositiontocolourmap = make(map[string]color.Color, 0)
	// need to extend for 1536 plates
	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "BB", "CC", "DD", "EE", "FF"}

	// Find out colour at each position:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)
			colourarray = append(colourarray, colour)

			// change colour to colour from a palette
			newcolour := chosencolourpalette.Convert(colour)

			plateimage.Set(x, y, newcolour)
			// equivalent well position
			wellposition := alphabet[y] + strconv.Itoa(x+1)
			fmt.Println(wellposition)
			wellpositionarray = append(wellpositionarray, wellposition)
			wellpositiontocolourmap[wellposition] = newcolour
		}
	}

	// rename file
	splitfilename := strings.Split(imagefilename, `.`)

	newname := splitfilename[0] + "_plateformat" + `.` + splitfilename[1]
	// save
	err := imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
	}

	// choose colour palette from top

	//arrayfrompalette := make([]color.Color, 0)

	for i, combo := range colourarray {
		fmt.Println("for well position ", wellpositionarray[i], ":")
		r, g, b, a := combo.RGBA()

		fmt.Println("colour (r,g,b,a)= ", r/256, g/256, b/256, a/256)

		// closest palette colour
		fmt.Println("palette colour:", chosencolourpalette.Convert(combo))
		fmt.Println("palette colour number:", chosencolourpalette.Index(combo))
	}

	numberofpixels = len(colourarray)
	fmt.Println("numberofpixels:", numberofpixels)

	return
}

func PrintFPImagePreview(imagefile string, plate *wtype.LHPlate, visiblemap, uvmap map[color.Color]string) {

	plateimage := ResizeImagetoPlate(imagefile, plate, imaging.CatmullRom)

	uvpalette := palettefromMap(uvmap)

	// Find out colour at each position under UV:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB
			colour := plateimage.At(x, y)

			// change colour to colour from a palette
			uvcolour := uvpalette.Convert(colour)

			plateimage.Set(x, y, uvcolour)

		}
	}

	// rename file
	splitfilename := strings.Split(imagefile, `.`)

	newname := splitfilename[0] + "_plateformat_UV" + `.` + splitfilename[1]
	// save
	err := imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
	}

	// repeat for visible

	// Find out colour at each position under visible light:
	for y := 0; y < plateimage.Bounds().Dy(); y++ {
		for x := 0; x < plateimage.Bounds().Dx(); x++ {
			// colour or pixel in RGB

			colour := plateimage.At(x, y)
			r, g, b, a := colour.RGBA()
			rgba := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			fmt.Println("colour", colour)
			fmt.Println("visiblemap", visiblemap)
			fmt.Println("uvmap", uvmap)
			colourstring := uvmap[rgba]
			fmt.Println("colourstring", colourstring)
			// change colour to colour of same cell + fluorescent protein under visible light
			stringkeymap, err := reversepalettemap(visiblemap)
			if err != nil {
				panic(err)
			}
			fmt.Println("stringkeymap", stringkeymap)
			viscolour, ok := stringkeymap[colourstring]
			if ok != true {
				errmessage := fmt.Sprintln("colourstring", colourstring, "not found in map", stringkeymap, "len", len(stringkeymap))
				panic(errmessage)
			}
			fmt.Println("viscolour", viscolour)
			plateimage.Set(x, y, viscolour)

		}
	}

	// rename file
	splitfilename = strings.Split(imagefile, `.`)

	newname = splitfilename[0] + "_plateformat_vis" + `.` + splitfilename[1]
	// save
	err = imaging.Save(plateimage, newname)
	if err != nil {
		panic(err)
	}
	return
}

//  Final function for user which uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
/*func PipetteImagetoPlate(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, availablecolours []string, componentlist []*wtype.LHComponent, volumeperwell wunit.Volume) (finalsolutions []*wtype.LHComponent) {

	componentmap, err := MakestringtoComponentMap(availablecolours, componentlist)
	if err != nil {
		panic(err)
	}

	solutions := make([]*wtype.LHComponent, 0)

	for locationkey, colour := range positiontocolourmap {

		component := componentmap[Colourcomponentmap[colour]]

		if component.Type == "dna" {
			component.Type = "DoNotMix"
		}

		pixelSample := mixer.Sample(component, volumeperwell)
		solution := mixer.MixTo(OutPlate, locationkey, pixelSample)
		solutions = append(solutions, solution)
	}

	finalsolutions = solutions
	return
}*/

// make a map of which colour description applies to which component assuming all components in factory are available, returns errors if either keys or components cannot be added
/*
func MakestringtoComponentMapFromFactory(colourtostringmap map[color.Color]string) (componentmap map[string]*wtype.LHComponent) {

	componentmap = make(map[string]*wtype.LHComponent, 0)
	//var previouserror error = nil

	for _, colour := range colourtostringmap {

		componentname := colourtostringmap[colour]

		componentmap[componentname] = factory.GetComponentByType(component)

	}
	return
}
*/
// else, specify colours and components to make a map of which colour description applies to which component, returns errors if either keys or components cannot be added
func MakestringtoComponentMap(keys []string, componentlist []*wtype.LHComponent) (componentmap map[string]*wtype.LHComponent, err error) {

	componentmap = make(map[string]*wtype.LHComponent, 0)
	var previouserror error = nil
	for i, key := range keys {
		for j, component := range componentlist {

			if component.CName == key {
				componentmap[key] = component
				break
			}
			if i == len(keys) {

				err = fmt.Errorf(previouserror.Error(), "+", "no key and component found for", keys, key)
			}
			if j == len(componentlist) {

				err = fmt.Errorf(previouserror.Error(), "+", "no key and component found for", component.CName)
			}
		}
	}
	return
}

//  Final function for blending colours to make and image. Uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
func PipetteImagebyBlending(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, cyan *wtype.LHComponent, magenta *wtype.LHComponent, yellow *wtype.LHComponent, black *wtype.LHComponent, volumeperfullcolour wunit.Volume) (finalsolutions []*wtype.LHComponent) {

	solutions := make([]*wtype.LHComponent, 0)

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		cmyk := ColourtoCMYK(colour)

		cyanvol := wunit.NewVolume((float64(cmyk.C) * volumeperfullcolour.SIValue()), "l")
		yellowvol := wunit.NewVolume((float64(cmyk.Y) * volumeperfullcolour.SIValue()), "l")
		magentavol := wunit.NewVolume((float64(cmyk.M) * volumeperfullcolour.SIValue()), "l")
		blackvol := wunit.NewVolume((float64(cmyk.K) * volumeperfullcolour.SIValue()), "l")

		cyanSample := mixer.Sample(cyan, cyanvol)
		components = append(components, cyanSample)
		yellowSample := mixer.Sample(yellow, yellowvol)
		components = append(components, yellowSample)
		magentaSample := mixer.Sample(magenta, magentavol)
		components = append(components, magentaSample)
		blackSample := mixer.Sample(black, blackvol)
		components = append(components, blackSample)
		solution := mixer.MixTo(OutPlate.Type, locationkey, 1, components...)
		solutions = append(solutions, solution)
	}

	finalsolutions = solutions
	return
}
