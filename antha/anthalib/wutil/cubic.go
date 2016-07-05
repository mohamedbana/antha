package wutil

import (
	"math"
	"math/cmplx"
)

// inverse function for cubic polynomials

type Cubic struct {
	a float64
	b float64
	c float64
	d float64
	p float64
	r float64
}

func (c *Cubic) F(v float64) float64 {
	return a*math.Pow(v, 3.0) + b*math.Pow(v, 2.0) + c*v + d
}

func (c *Cubic) P() float64 {
	if c.p == 0.0 {
		c.p = (-1.0 * c.b) / (3.0 * c.a)
	}

	return c.p
}

func (c *Cubic) Q(v float64) float64 {
	q := math.Pow(c.p, 3.0) + (c.b*c.c-3.0*c.a*(c.d-v))/(6.0*math.Pow(c.a, 2.0))
	return q
}

func (c *Cubic) R() float64 {
	if c.r == 0.0 {
		c.r = c.c / (3.0 * c.a)
	}
	return c.r
}
func (c *Cubic) I(v float64) float64 {
	s := complex128(c.p)
	vv := cmplx.Pow(c.q, 2.0) + cmplx.Pow(c.r-math.Pow(c.p, 2.0), 0.5)
	vv1 := cmplx.Pow(complex128(c.q)+vv, 1.0/3.0)
	vv2 := cmplx.Pow(complex128(c.q)-vv, 1.0/3.0)

	// this should be real anyway
	return cmplx.Abs(s + vv1 + vv2)
}
