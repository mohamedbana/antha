// pixeltoplate.go
// image.go
package image

import (
	"github.com/antha-lang/antha/internal/github.com/disintegration/imaging"
	//	"image"
	"fmt"
	"image/color"
	"image/color/palette"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Colour palette to use // this would relate to a map of components of these available colours in factor

var AvailablePalettes = map[string]color.Palette{
	"Palette1": Chosencolourpalette,
	"WebSafe":  palette.WebSafe, //websafe,
	"plan9":    palette.Plan9,
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
	color.RGBA{R: uint8(196), G: uint8(40), B: uint8(27), A: uint8(255)}:   "black",
}

// create a map of pixel to plate position from processing a given image with a chosen colour palette.
// It's recommended to use at least 384 well plate
func ImagetoPlatelayout(imagefilename string, plate *wtype.LHPlate, chosencolourpalette color.Palette) (wellpositiontocolourmap map[string]color.Color, numberofpixels int) {

	// input files (just 1 in this case)
	files := []string{imagefilename}

	// Colour palette to use // this would relate to a map of components of these available colours in factory
	//availablecolours := chosencolourpalette //palette.WebSafe

	//var plateimages []image.Image

	for _, file := range files {
		img, err := imaging.Open(file)
		if err != nil {
			panic(err)
		}

		// have the option of changing the resize algorithm here
		plateimage := imaging.Resize(img, 0, plate.WlsY, imaging.CatmullRom)
		//plateimages = append(plateimages,plateimage)

		// make map of well position to colour: (array for time being)

		wellpositionarray := make([]string, 0)
		colourarray := make([]color.Color, 0)
		wellpositiontocolourmap = make(map[string]color.Color, 0)
		// need to extend for 1536 plates
		alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
				wellposition := string(alphabet[y]) + strconv.Itoa(x+1)
				wellpositionarray = append(wellpositionarray, wellposition)
				wellpositiontocolourmap[wellposition] = newcolour
			}
		}

		// rename file
		splitfilename := strings.Split(file, `.`)

		newname := splitfilename[0] + "_plateformat" + `.` + splitfilename[1]
		// save
		err = imaging.Save(plateimage, newname)
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

	}
	return
}

//  Final function for user which uses a given map of position to colour generated from the image processing function  along with lists of available colours, components and plate types
func PipetteImagetoPlate(OutPlate *wtype.LHPlate, positiontocolourmap map[string]color.Color, availablecolours []string, componentlist []*wtype.LHComponent, volumeperwell wunit.Volume) (finalsolutions []*wtype.LHSolution) {

	componentmap, err := MakestringtoComponentMap(availablecolours, componentlist)
	if err != nil {
		panic(err)
	}

	solutions := make([]*wtype.LHSolution, 0)

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
}

// make a map of which colour description applies to which component, returns errors if either keys or components cannot be added
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
