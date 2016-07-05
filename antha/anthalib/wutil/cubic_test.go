package wutil

import (
	"fmt"
	"testing"
)

func TestCubic(t *testing.T) {
	c := Cubic{a: 0.00769276, b: 0.40223723, c: 7.06878347, d: 0.0}
	q := Quadratic{a: -0.00015253, b: 0.0992489368, c: 0.46400556}
	qt := Quartic{a: -3.3317851312e-09, b: 0.00000225834467, c: -0.0006305492472, d: 0.1328156706978, e: 0}
	i := 0.0

	for k := 0; k < 15; k++ {
		v := c.F(i)
		iv := c.I(v)
		fmt.Println(i, " ", v, " ", iv, " ", q.F(v), qt.F(v))

		i += 1.0
	}
}
