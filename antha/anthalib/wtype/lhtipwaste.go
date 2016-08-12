package wtype

// defines a tip waste
import "fmt"

// tip waste

type LHTipwaste struct {
	Name       string
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
	bounds     BBox
	parent     LHObject
}

func (tw LHTipwaste) SpaceLeft() int {
	return tw.Contents - tw.Capacity
}

func (te LHTipwaste) String() string {
	return fmt.Sprintf(
		`LHTipwaste {
	ID: %s,
	Type: %s,
    Name: %s,
	Mnfr: %s,
	Capacity: %d,
	Contents: %d,
    Length: %f,
    Width: %f,
	Height: %f,
	WellXStart: %f,
	WellYStart: %f,
	WellZStart: %f,
	AsWell: %p,
}
`,
		te.ID,
		te.Type,
		te.Name,
		te.Mnfr,
		te.Capacity,
		te.Contents,
		te.bounds.GetSize().X,
		te.bounds.GetSize().Y,
		te.bounds.GetSize().Z,
		te.WellXStart,
		te.WellYStart,
		te.WellZStart,
		te.AsWell, //AsWell is printed as pointer to keep things short
	)
}

func (tw *LHTipwaste) Dup() *LHTipwaste {
	return NewLHTipwaste(tw.Capacity, tw.Type, tw.Mnfr, tw.bounds.GetSize(), tw.AsWell, tw.WellXStart, tw.WellYStart, tw.WellZStart)
}

func (tw *LHTipwaste) GetName() string {
	return tw.Name
}

func (tw *LHTipwaste) GetType() string {
	return tw.Type
}

func NewLHTipwaste(capacity int, typ, mfr string, size Coordinates, w *LHWell, wellxstart, wellystart, wellzstart float64) *LHTipwaste {
	var lht LHTipwaste
	lht.ID = GetUUID()
	lht.Type = typ
	lht.Name = fmt.Sprintf("%s_%s", typ, lht.ID[1:len(lht.ID)-2])
	lht.Mnfr = mfr
	lht.Capacity = capacity
	lht.bounds.SetSize(size)
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

//##############################################
//@implement LHObject
//##############################################

func (self *LHTipwaste) GetPosition() Coordinates {
	return self.bounds.GetPosition()
}

func (self *LHTipwaste) GetSize() Coordinates {
	return self.bounds.GetSize()
}

func (self *LHTipwaste) GetBoxIntersections(box BBox) []LHObject {
	ret := []LHObject{}
	//todo, test well
	if self.bounds.IntersectsBox(box) {
		ret = append(ret, self)
	}
	return ret
}

func (self *LHTipwaste) GetPointIntersections(point Coordinates) []LHObject {
	ret := []LHObject{}
	//Todo, test well
	if self.bounds.IntersectsPoint(point) {
		ret = append(ret, self)
	}
	return ret
}

func (self *LHTipwaste) SetOffset(o Coordinates) {
	if self.parent != nil {
		o = o.Add(self.parent.GetSize())
	}
	self.bounds.SetPosition(o)
}

func (self *LHTipwaste) SetParent(p LHObject) {
	self.parent = p
}

func (self *LHTipwaste) GetParent() LHObject {
	return self.parent
}

//##############################################
//@implement Addressable
//##############################################

func (self *LHTipwaste) AddressExists(c WellCoords) bool {
	return c.X == 0 && c.Y == 0
}

func (self *LHTipwaste) NRows() int {
	return 1
}

func (self *LHTipwaste) NCols() int {
	return 1
}

func (self *LHTipwaste) GetChildByAddress(c WellCoords) LHObject {
	if !self.AddressExists(c) {
		return nil
	}
	//LHWells arent LHObjects yet
	//return self.AsWell
	return nil
}

func (self *LHTipwaste) CoordsToWellCoords(r Coordinates) (WellCoords, Coordinates) {
	wc := WellCoords{0, 0}

	c, _ := self.WellCoordsToCoords(wc, TopReference)

	return wc, r.Subtract(c)
}

func (self *LHTipwaste) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates, bool) {
	if !self.AddressExists(wc) {
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
		self.WellXStart + 0.5*self.AsWell.Xdim,
		self.WellYStart + 0.5*self.AsWell.Ydim,
		z}, true
}
