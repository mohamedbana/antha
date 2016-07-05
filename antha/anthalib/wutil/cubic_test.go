package wutil

import (
	"fmt"
	"testing"
)

func TestCubic(t *testing.T) {
	c := Cubic{a: 0.00769276, b: 0.40223723, c: 7.06878347, d: 0.0}

	i := 0.0

	for k := 0; k < 20; k++ {
		v := c.F(i)
		iv := c.I(v)
		v2 := c.F(iv)
		iv2 := c.I(v2)
		fmt.Println(i, " ", v, " ", iv, " ", v2, " ", iv2)

		i += 1.0
	}
}
