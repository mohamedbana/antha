package wtype

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

func MatchString(s1, s2 string) bool {
	m, _ := regexp.MatchString(s1, s2)
	return m
}

func (wc WellCoords) Equals(w2 WellCoords) bool {
	if wc.X == w2.X && wc.Y == w2.Y {
		return true
	}

	return false
}

func MakeWellCoords(wc string) WellCoords {
	// try each one in turn

	r := MakeWellCoordsA1(wc)

	zero := WellCoords{-1, -1}

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoords1A(wc)

	if !r.Equals(zero) {
		return r
	}

	r = MakeWellCoordsXY(wc)

	return r
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	// only handles 96 well plates
	if !MatchString("[A-Z][0-9]{1,2}", a1) {
		return WellCoords{-1, -1}
	}
	return WellCoords{wutil.ParseInt(a1[1:len(a1)]) - 1, AlphaToNum(string(a1[0])) - 1}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	// only handles 96 well plates

	if !MatchString("[0-9]{1,2}[A-Z]", a1) {
		return WellCoords{-1, -1}
	}
	return WellCoords{AlphaToNum(string(a1[0])) - 1, wutil.ParseInt(a1[1:len(a1)]) - 1}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	r := WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}

	if r.X < 0 || r.Y < 0 {
		return WellCoords{-1, -1}
	}

	return r
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	if tx == nil || len(tx) != 2 || len(tx[0]) == 0 || len(tx[1]) == 0 {
		return WellCoords{-1, -1}
	}
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc WellCoords) FormatXY() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}
func (wc WellCoords) Format1A() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X+1) + NumToAlpha(wc.Y+1)
}
func (wc WellCoords) FormatA1() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}
func (wc WellCoords) WellNumber() int {
	if wc.X < 0 || wc.Y < 0 {
		return -1
	}
	return (8*(wc.X-1) + wc.Y)
}

func (wc WellCoords) ColNumString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return strconv.Itoa(wc.X + 1)
}
func (wc WellCoords) RowLettString() string {
	if wc.X < 0 || wc.Y < 0 {
		return ""
	}
	return NumToAlpha(wc.Y + 1)
}
