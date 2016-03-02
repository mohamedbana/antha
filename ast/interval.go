package ast

// TODO(ddn): Replace with a more efficient data structure (interval tree)

// An interval or union thereof
type Interval struct {
	values []struct{ A, B float64 }
}

func (a Interval) Extrema() (min, max float64) {
	if len(a.values) == 0 {
		return
	}
	min = a.values[0].A
	max = a.values[0].B
	for i, n := 1, len(a.values); i < n; i += 1 {
		if a.values[i].A < min {
			min = a.values[i].A
		}

		if a.values[i].B > max {
			max = a.values[i].B
		}
	}
	return
}

// The nil interval does not contain any points
func (a Interval) Nil() bool {
	return len(a.values) == 0
}

func (a Interval) Contains(x, y float64) bool {
	for _, v := range a.values {
		if v.A <= x && y <= v.B {
			return true
		}
	}
	return false
}

func (a Interval) Add(x Interval) *Interval {
	var values []struct{ A, B float64 }
	for _, v := range a.values {
		values = append(values, v)
	}
	for _, v := range x.values {
		values = append(values, v)
	}
	return &Interval{values: values}
}

// Create the interval [a, b]
func NewInterval(a, b float64) *Interval {
	return &Interval{
		values: []struct{ A, B float64 }{struct{ A, B float64 }{A: a, B: b}},
	}
}

// Create the interval [a, a]
func NewPoint(a float64) *Interval {
	return &Interval{
		values: []struct{ A, B float64 }{struct{ A, B float64 }{A: a, B: a}},
	}
}
