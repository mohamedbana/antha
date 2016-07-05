package wutil

import (
	"math"
)

// function for quadratics
type Quartic struct {
	a float64
	b float64
	c float64
	d float64
	e float64
}

func (c *Quartic) A() float64 {
	return c.a
}
func (c *Quartic) B() float64 {
	return c.b
}
func (c *Quartic) C() float64 {
	return c.c
}

func (c *Quartic) D() float64 {
	return c.d
}

func (c *Quartic) E() float64 {
	return c.e
}

func (c *Quartic) F(v float64) float64 {
	return c.a*math.Pow(v, 4.0) + c.b*math.Pow(v, 3.0) + c.c*math.Pow(v, 2.0) + c.d*v + c.e
}
