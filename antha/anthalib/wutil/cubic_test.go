package wutil

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCubic(t *testing.T) {
	c := Cubic{A: 0.00769276, B: 0.40223723, C: 7.06878347, D: 0.0}
	q := Quadratic{A: -0.00015253, B: 0.0992489368, C: 0.46400556}
	qt := Quartic{A: -3.3317851312e-09, B: 0.00000225834467, C: -0.0006305492472, D: 0.1328156706978, E: 0}
	i := 0.0

	for k := 0; k < 15; k++ {
		v := c.F(i)
		iv := c.I(v)
		i += 1.0
	}
}

func TestSerialize(t *testing.T) {
	c := Cubic{A: 0.00769276, B: 0.40223723, C: 7.06878347, D: 0.0}

	m, _ := json.Marshal(c)

	fmt.Println(string(m))
}
