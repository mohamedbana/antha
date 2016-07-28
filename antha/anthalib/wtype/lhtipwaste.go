package wtype

// defines a tip waste
import "fmt"

// tip waste

type LHTipwaste struct {
	ID         string
	Type       string
	Mnfr       string
	Capacity   int
	Contents   int
	Height     float64
	WellXStart float64
	WellYStart float64
	WellZStart float64
	AsWell     *LHWell
}

func (tw LHTipwaste) SpaceLeft() int {
	return tw.Contents - tw.Capacity
}

func (te LHTipwaste) String() string {
	return fmt.Sprintf(
		`LHTipwaste {
	ID: %s,
	Type: %s,
	Mnfr: %s,
	Capacity: %d,
	Contents: %d,
	Height: %f,
	WellXStart: %f,
	WellYStart: %f,
	WellZStart: %f,
	AsWell: %p,
}
`,
		te.ID,
		te.Type,
		te.Mnfr,
		te.Capacity,
		te.Contents,
		te.Height,
		te.WellXStart,
		te.WellYStart,
		te.WellZStart,
		te.AsWell, //AsWell is printed as pointer to kepp things short
	)
}

func (tw *LHTipwaste) Dup() *LHTipwaste {
	return NewLHTipwaste(tw.Capacity, tw.Type, tw.Mnfr, tw.Height, tw.AsWell, tw.WellXStart, tw.WellYStart, tw.WellZStart)
}

func (tw *LHTipwaste) GetName() string {
	return tw.Type
}

func (tw *LHTipwaste) GetType() string {
    return tw.Type
}

func NewLHTipwaste(capacity int, typ, mfr string, height float64, w *LHWell, wellxstart, wellystart, wellzstart float64) *LHTipwaste {
	var lht LHTipwaste
	lht.ID = GetUUID()
	lht.Type = typ
	lht.Mnfr = mfr
	lht.Capacity = capacity
	lht.Height = height
	lht.AsWell = w
	lht.WellXStart = wellxstart
	lht.WellYStart = wellystart
	lht.WellZStart = wellzstart
	return &lht
}

func (lht *LHTipwaste) Empty() {
	lht.Contents = 0
}

func (lht *LHTipwaste) Dispose(n int) bool {
	if lht.Capacity-lht.Contents < n {
		return false
	}

	lht.Contents += n
	return true
}

//@implement LHDeckObject

func (self *LHTipwaste) GetSize() Coordinates {
    //Assume that TipX/YStart is repeated the other side
    return Coordinates{
        2 * self.WellXStart + self.AsWell.Xdim,
        2 * self.WellYStart + self.AsWell.Ydim,
        self.Height,
    }
}

func (self *LHTipwaste) HasCoords(c WellCoords) bool {
    return c.X == 0 &&
           c.Y == 0 
}

func (self *LHTipwaste) GetCoords(c WellCoords) (interface{}, bool) {
    if !self.HasCoords(c) {
        return nil, false
    }
    return self.AsWell, true
}

func (self *LHTipwaste) CoordsToWellCoords(r Coordinates) (WellCoords, Coordinates) {
    wc := WellCoords{0,0}

    c, _ := self.WellCoordsToCoords(wc, TopReference)

    return wc, r.Subtract(c)
}

func (self *LHTipwaste) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
    if !self.HasCoords(wc) {
        return Coordinates{}, false
    }

    var z float64
    if r == BottomReference {
        z = self.WellZStart
    } else if r == TopReference {
        z = self.Height
    } else {
        return Coordinates{}, false
    }

    return Coordinates{
        self.WellXStart + 0.5 * self.AsWell.Xdim,
        self.WellYStart + 0.5 * self.AsWell.Ydim,
        z}, true
}

