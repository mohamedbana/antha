package mixer

import (
	"bytes"
	"testing"
)

var testPlateCsv = `pcrplate_with_cooler,
A1,water,water,50.0,ul,
A4,tea,water,50.0,ul,
A5,milk,water,100.0,ul,
`

func TestParsePlate(t *testing.T) {
	_, err := parseInputPlateData(bytes.NewBuffer([]byte(testPlateCsv)))
	if err != nil {
		t.Fatal(err)
	}
}
