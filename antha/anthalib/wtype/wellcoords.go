package wtype

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	// only handles 96 well plates
	return WellCoords{wutil.ParseInt(a1[1:len(a1)]) - 1, AlphaToNum(string(a1[0])) - 1}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	// only handles 96 well plates
	return WellCoords{AlphaToNum(string(a1[0])) - 1, wutil.ParseInt(a1[1:len(a1)]) - 1}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	return WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
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
