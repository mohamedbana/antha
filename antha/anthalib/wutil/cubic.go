package wutil

import (
	"math"
	"math/cmplx"
)

// function for cubics
type Cubic struct {
	a float64
	b float64
	c float64
	d float64
	p float64
	r float64
}

func (c *Cubic) A() float64 {
	return c.a
}
func (c *Cubic) B() float64 {
	return c.b
}
func (c *Cubic) C() float64 {
	return c.c
}
func (c *Cubic) D() float64 {
	return c.d
}

func (c *Cubic) F(v float64) float64 {
	return c.a*math.Pow(v, 3.0) + c.b*math.Pow(v, 2.0) + c.c*v + c.d
}

func (c *Cubic) P() float64 {
	if c.p == 0.0 && c.a != 0.0 {
		c.p = (-1.0 * c.b) / (3.0 * c.a)
	}

	return c.p
}

func (c *Cubic) Q(v float64) float64 {
	q := math.Pow(c.P(), 3.0) + ((c.b*c.c)-(3.0*c.a*(c.d-v)))/(6.0*math.Pow(c.a, 2.0))
	return q
}

func (c *Cubic) R() float64 {
	if c.r == 0.0 {
		c.r = c.c / (3.0 * c.a)
	}
	return c.r
}

// should decide when this is safe
func (c *Cubic) I(v float64) float64 {
	s := complex(c.P(), 0.0)
	vv := cmplx.Pow(cmplx.Pow(complex(c.Q(v), 0.0), 2.0)+cmplx.Pow(complex(c.R(), 0.0)-cmplx.Pow(complex(c.P(), 0.0), 2.0), 3.0), 0.5)
	vv1 := cmplx.Pow(complex(c.Q(v), 0.0)+vv, 1.0/3.0)
	vv2 := cmplx.Pow(complex(c.Q(v), 0.0)-vv, 1.0/3.0)

	// this should be real anyway, but Abs should be a bit better than real(...)
	return cmplx.Abs(s + vv1 + vv2)
}
