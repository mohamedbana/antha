package mixer

import (
	"testing"
)

func TestParsePlate(t *testing.T) {
	_, err := parseInputPlateFile("test.csv")
	if err != nil {
		t.Fatal(err)
	}

}
