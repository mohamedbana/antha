package wtype

type PlateIterator interface {
	Rewind()
	Next() WellCoords
	Cur() WellCoords
	Valid() bool
}

type BasicPlateIterator struct {
	fst  WellCoords
	cur  WellCoords
	p    LHPlate
	rule func(WellCoords, LHPlate) WellCoords
}

func (it *BasicPlateIterator) Rewind() {
	it.cur = it.fst
	return it.cur
}
func (it *BasicPlateIterator) Cur() WellCoords {
	return it.cur
}

func (it *BasicPlateIterator) Valid() bool {
	if it.cur.X > it.p.WellsX || it.cur.X < 0 {
		return false
	}

	if it.cur.Y > it.p.WellsY || it.cur.Y < 0 {
		return false
	}

	return true
}

func (it *BasicPlateIterator) NextWell() WellCoords {
	it.cur = it.rule(it.cur)
	return it.cur
}

func NextColumn(wc WellCoords, p LHPlate) WellCoords {
	wc.X += 1
	if wc.X >= p.WellsX {
		wc.X = 0
		wc.Y += 1
	}
	if wc.Y >= p.WellsY {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}

func NextRow(wc WellCoords, p LHPlate) WellCoords {
	wc.Y += 1
	if wc.Y >= p.WellsY {
		wc.Y = 0
		wc.X += 1
	}
	if wc.X >= p.WellsX {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}

func NewColumnIterator(p LHPlate) BasicPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextColumn
	bi.p = p
}

func NewRowIterator(p LHPlate) BasicPlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextRow
	bi.p = p
}
