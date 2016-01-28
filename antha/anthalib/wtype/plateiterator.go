package wtype

type PlateIterator interface {
	Rewind() WellCoords
	Next() WellCoords
	Curr() WellCoords
	Valid() bool
}

type BasicPlateIterator struct {
	fst  WellCoords
	cur  WellCoords
	p    *LHPlate
	rule func(WellCoords, *LHPlate) WellCoords
}

func (it *BasicPlateIterator) Rewind() WellCoords {
	it.cur = it.fst
	return it.cur
}
func (it *BasicPlateIterator) Curr() WellCoords {
	return it.cur
}

func (it *BasicPlateIterator) Valid() bool {
	if it.cur.X > it.p.WellsX() || it.cur.X < 0 {
		return false
	}

	if it.cur.Y > it.p.WellsY() || it.cur.Y < 0 {
		return false
	}

	return true
}

func (it *BasicPlateIterator) Next() WellCoords {
	it.cur = it.rule(it.cur, it.p)
	return it.cur
}

func NextColumnOnce(wc WellCoords, p *LHPlate) WellCoords {
	wc.X += 1
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y += 1
	}
	return wc
}
func NextColumn(wc WellCoords, p *LHPlate) WellCoords {
	wc.X += 1
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y += 1
	}
	if wc.Y >= p.WellsY() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}

func NextRow(wc WellCoords, p *LHPlate) WellCoords {
	wc.Y += 1
	if wc.Y >= p.WellsY() {
		wc.Y = 0
		wc.X += 1
	}
	if wc.X >= p.WellsX() {
		wc.X = 0
		wc.Y = 0
	}
	return wc
}
func NextRowOnce(wc WellCoords, p *LHPlate) WellCoords {
	wc.Y += 1
	if wc.Y >= p.WellsY() {
		wc.Y = 0
		wc.X += 1
	}
	return wc
}

func NewColumnIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextColumn
	bi.p = p
	return &bi
}
func NewOneTimeColumnIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextColumnOnce
	bi.p = p
	return &bi
}

func NewRowIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextRow
	bi.p = p
	return &bi
}
func NewOneTimeRowIterator(p *LHPlate) PlateIterator {
	var bi BasicPlateIterator
	bi.fst = WellCoords{0, 0}
	bi.cur = WellCoords{0, 0}
	bi.rule = NextRowOnce
	bi.p = p
	return &bi
}
