// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"fmt"
	"os"
	//"path"
	//"runtime"
	"image"
	"image/draw"
	"math"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/internal/code.google.com/p/draw2d/draw2d"
	//"code.google.com/p/draw2d/draw2d"
	"github.com/antha-lang/antha/internal/github.com/disintegration/imaging"
	//"github.com/antha-lang/antha/internal/github.com/hybridgroup/go-opencv/opencv"
	"github.com/antha-lang/antha/internal/github.com/lazywei/go-opencv/opencv"
	//"../opencv" // can be used in forks, comment in real application
)

func Pick(imagefile string, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	var newname string
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname = splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			panic("img == nil")
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname = splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)

					if err != nil {
						panic(err)
					}
					/*
						wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
						return*/
					/*img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")

					*/
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					//os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					//os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			//	os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	//os.Exit(0)
	/*if rotated {
		wells = PickAgain(newname, numbertopick, setplateperimeterfirst, false)
	}*/
	return
}

/*
func PickAgain(imagefile string, numbertopick int, setplateperimeterfirst bool, rotate bool) (wells []string) {
	var topleft opencv.Point
	var bottomright opencv.Point
	var topright opencv.Point
	var well string
	var rotated bool
	var rotatedimg draw.Image
	counter := 0
	//_, currentfile, _, _ := runtime.Caller(0)
	filename := imagefile //:= path.Join(path.Dir(currentfile), imagefile)
	wells = make([]string, 0)
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}

	img0 := opencv.LoadImage(filename)
	if img0 == nil {
		panic("LoadImage fail")
	}

	// resize image if bigger than full screen
	if img0.Height() > 801 {
		img0.Release()
		imgtoresize, err := imaging.Open(filename)
		if err != nil {
			panic(err)
		}
		resizedimage := imaging.Resize(imgtoresize, 0, 800, imaging.CatmullRom)
		splitfilename := strings.Split(filename, `.`)

		newname := splitfilename[0] + "_resized" + `.` + splitfilename[1]

		err = imaging.Save(resizedimage, newname)
		if err != nil {
			panic(err)
		}
		img0 = opencv.LoadImage(newname)
		filename = newname
		if img0 == nil {
			panic("LoadImage fail")
		}
	}

	defer img0.Release()

	fmt.Print("Hot keys: \n",
		"\tESC - quit the program\n",
		"\tr - restore the original image\n",
		"\ti or ENTER - run inpainting algorithm\n",
		"\t\t(before running it, paint something on the image)\n",
	)

	img := img0.Clone()
	inpainted := img0.Clone()
	inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

	opencv.Zero(inpaint_mask)
	//opencv.Zero( inpainted )

	win := opencv.NewWindow("image")
	defer win.Destroy()

	prev_pt := opencv.Point{-1, -1}
	win.SetMouseCallback(func(event, x, y, flags int, param ...interface{}) {
		if img == nil {
			os.Exit(0)
		}

		if event == opencv.CV_EVENT_LBUTTONUP ||
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) == 0 {
			prev_pt = opencv.Point{-1, -1}
		} else if event == opencv.CV_EVENT_LBUTTONDOWN {

			rgb := opencv.ScalarAll(255.0)
			counter = counter + 1

			if rotate && rotated == false {
				imgtorotate, err := imaging.Open(filename)

				if err != nil {
					panic(err)
				}
				if counter == 1 {
					fmt.Println("Rotating image")
					fmt.Println("topleft point:", x, y)
					topleft = opencv.Point{x, y}
				} else if counter == 2 {
					topright = opencv.Point{x, y}
					fmt.Println("topright point:", x, y)
					opposite := float64(topleft.Y) - float64(topright.Y)
					adjacent := float64(topright.X) - float64(topleft.X)
					tantheta := opposite / adjacent
					fmt.Println("adjacent:", adjacent)
					fmt.Println("opposite:", opposite)
					fmt.Println("costheta:", tantheta)
					thetainrad := math.Atan(tantheta)
					fmt.Println("thetainrad:", thetainrad)
					degrees := (180 / math.Pi) * thetainrad
					fmt.Println("degrees:", degrees)
					tr := draw2d.NewRotationMatrix(thetainrad)

					//
					ar := imgtorotate.Bounds()
					w, h, _ := ar.Dx(), ar.Dy(), 30.0
					rotatedimg = image.NewRGBA(image.Rect(0, 0, w, h))
					draw.Draw(rotatedimg, ar, imgtorotate, ar.Min, draw.Src)

					draw2d.DrawImage(imgtorotate, rotatedimg, tr, draw.Src, draw2d.LinearFilter)

					// open new window
					splitfilename := strings.Split(filename, `.`)
					newname := splitfilename[0] + "_rotated" + `.` + splitfilename[1]

					//err = imaging.Save(rotatedimg, newname)
					err = imaging.Save(rotatedimg, newname)
					if err != nil {
						panic(err)
					}
					img0 = opencv.LoadImage(newname)
					filename = newname
					if img0 == nil {
						panic("LoadImage fail")
					}
					defer img0.Release()

					fmt.Print("Hot keys: \n",
						"\tESC - quit the program\n",
						"\tr - restore the original image\n",
						"\ti or ENTER - run inpainting algorithm\n",
						"\t\t(before running it, paint something on the image)\n",
					)

					//img := img0.Clone()
					//inpainted := img0.Clone()
					//inpaint_mask := opencv.CreateImage(img0.Width(), img0.Height(), 8, 1)

					//opencv.Zero(inpaint_mask)
					//opencv.Zero( inpainted )

					win := opencv.NewWindow("rotated")
					defer win.Destroy()
					// change status
					rotated = true
					fmt.Println("Now choose top left of plate")
				}
			} else if rotate && setplateperimeterfirst && counter == 3 && rotated || setplateperimeterfirst && counter == 1 && rotate == false {
				topleft = opencv.Point{x, y}
				fmt.Println("Topleft chosen")
				fmt.Println("Now choose bottom right of plate")

			} else if rotate && setplateperimeterfirst && counter == 4 && rotated || setplateperimeterfirst && counter == 2 && rotate == false {
				bottomright = opencv.Point{x, y}
				fmt.Println("plate boundaries", x, y)
				opencv.Rectangle(img, topleft, bottomright, rgb, 1, 8, 0)
				win.ShowImage(img)
			} else {

				prev_pt = opencv.Point{x, y}
				fmt.Println("actual pixels:", x, y)
				fmt.Println("imagesize:", img.ImageSize())
				fmt.Println("width:", img.Width())
				fmt.Println("height:", img.Height())
				fmt.Println("Well:", PixelstoWellPosition(x, y, img))
				fmt.Println("colony count:", counter)

				if setplateperimeterfirst {
					well = PixelstoWellPositionFromRectangle(x, y, topleft, bottomright)
				} else if setplateperimeterfirst == false {
					well = PixelstoWellPosition(x, y, img)
				}

				wells = append(wells, well)
				fmt.Println(wells)

				// draw circle on click

				rgb = opencv.ScalarAll(255.0)
				opencv.Circle(img, prev_pt, 1, rgb, 5, 8, 0)

				win.ShowImage(img)

				if rotate && setplateperimeterfirst && counter == numbertopick+4 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst && counter == numbertopick+2 {
					//return
					splitfilename := strings.Split(filename, `.`)

					newname := splitfilename[0] + "_picked" + `.` + splitfilename[1]
					opencv.SaveImage(newname, img, 0)
					os.Exit(0)
					return
				}
				if rotate == false && setplateperimeterfirst == false && counter == numbertopick {
					//return
					os.Exit(0)
					return
				}
			}
		} else if event == opencv.CV_EVENT_MOUSEMOVE &&
			(flags&opencv.CV_EVENT_FLAG_LBUTTON) != 0 {
			pt := opencv.Point{x, y}
			if prev_pt.X < 0 {
				prev_pt = pt
			}

			rgb := opencv.ScalarAll(255.0)

			opencv.Rectangle(img, prev_pt, pt, rgb, 1, 8, 0)
			//opencv.Line(inpaint_mask, prev_pt, pt, rgb, 5, 8, 0)
			//opencv.Line(img, prev_pt, pt, rgb, 5, 8, 0)
			prev_pt = pt

			win.ShowImage(img)
		}
	})
	win.ShowImage(img)
	opencv.WaitKey(0)

	win2 := opencv.NewWindow("inpainted image")

	defer win2.Destroy()
	win2.ShowImage(inpainted)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			os.Exit(0)
		} else if key == 'r' {
			opencv.Zero(inpaint_mask)
			opencv.Copy(img0, img, nil)
			win.ShowImage(img)
		} else if key == 'i' || key == '\n' {
			opencv.Inpaint(img, inpaint_mask, inpainted, 3,
				opencv.CV_INPAINT_TELEA,
			)
			win2.ShowImage(inpainted)
		}
	}
	os.Exit(0)
	if rotated {

	}
	return
}
*/
var alphabet []string = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z", "AA", "BB", "CC", "DD", "EE", "FF", "GG", "HH", "II", "JJ"}

func PixelstoWellPosition(x, y int, image *opencv.IplImage) (a1 string) {

	fmt.Println("func pixels:", x, y)
	fmt.Println("func width:", image.Width())
	fmt.Println("func height:", image.Height())

	col1536 := (float64(x) / float64(image.Width())) * float64(48)
	row1536 := (float64(y) / float64(image.Height())) * float64(36)

	colint, err := wutil.RoundDown(col1536)
	if err != nil {
		return
	}
	rowint, err := wutil.RoundDown(row1536)
	if err != nil {
		return
	}
	a1 = alphabet[rowint] + strconv.Itoa(colint+1)
	return
}
func PixelstoWellPositionFromRectangle(x, y int, topleft, bottomright opencv.Point) (a1 string) {

	fmt.Println("func pixels:", x, y)
	fmt.Println("rectangle width:", bottomright.X-topleft.X)
	fmt.Println("rectangle height:", bottomright.Y-topleft.Y)

	col1536 := (float64(x-topleft.X) / float64(bottomright.X-topleft.X)) * float64(48)
	row1536 := (float64(y-topleft.Y) / float64(bottomright.Y-topleft.Y)) * float64(36)

	colint, err := wutil.RoundDown(col1536)
	if err != nil {
		return
	}
	rowint, err := wutil.RoundDown(row1536)
	if err != nil {
		return
	}
	a1 = alphabet[rowint] + strconv.Itoa(colint+1)
	return
}
