package wutil

import (
	"math"
)

// function for quadratics
type Quadratic struct {
	a float64
	b float64
	c float64
}

func (c *Quadratic) A() float64 {
	return c.a
}
func (c *Quadratic) B() float64 {
	return c.b
}
func (c *Quadratic) C() float64 {
	return c.c
}

func (c *Quadratic) F(v float64) float64 {
	return c.a*math.Pow(v, 2.0) + c.b*v + c.c
}
