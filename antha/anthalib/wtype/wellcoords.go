package wtype

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

func MakeWellCoords(wc string) WellCoords {
	// try each one in turn

	r := MakeWellCoordsA1(wc)

	if r != nil {
		return r
	}

	r = MakeWellCoords1A(wc)

	if r != nil {
		return r
	}

	r = MakeWellCoordsXY(wc)

	return r
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	// only handles 96 well plates
	if !regexp.MatchString("[A-Z][0-9]{1,2}", a1) {
		return nil
	}
	return WellCoords{wutil.ParseInt(a1[1:len(a1)]) - 1, AlphaToNum(string(a1[0])) - 1}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	// only handles 96 well plates

	if !regexp.MatchString("[0-9]{1,2}[A-Z]", a1) {
		return nil
	}
	return WellCoords{AlphaToNum(string(a1[0])) - 1, wutil.ParseInt(a1[1:len(a1)]) - 1}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	return WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	if tx == nil || len(tx) != 2 || len(tx[0]) == 0 || len(tx[1]) == 0 {
		return nil
	}
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc *WellCoords) FormatXY() string {
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}
func (wc *WellCoords) Format1A() string {
	return strconv.Itoa(wc.X+1) + NumToAlpha(wc.Y+1)
}
func (wc *WellCoords) FormatA1() string {
	return NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}
func (wc *WellCoords) WellNumber() int {
	return (8*(wc.X-1) + wc.Y)
}

func (wc *WellCoords) ColNumString() string {
	return strconv.Itoa(wc.X + 1)
}
func (wc *WellCoords) RowLettString() string {
	return NumToAlpha(wc.Y + 1)
}
