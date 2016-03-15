// sample tracker service... needs to be outside this eventually but at least having one
// is a good start

package liquidhandling

import "fmt"

type SampleTracker struct {
	records map[string]string
}

func newSampleTracker() *SampleTracker {
	r := make(map[string]string)
	st := SampleTracker{r}
	return &st
}

func GetSampleTracker() *SampleTracker {
	return newSampleTracker()
}

func (st *SampleTracker) SetLocationOf(ID string, loc string) {
	fmt.Println("LOCATION OF ", ID, " SET TO ", loc)
	st.records[ID] = loc
}

func (st *SampleTracker) GetLocationOf(ID string) (string, bool) {
	fmt.Println("ASKING FOR LOCATION OF ", ID)

	fmt.Println("HERE'S WHAT I KNOW:")
	for k, v := range st.records {
		fmt.Println(k, v)
	}

	s, ok := st.records[ID]
	return s, ok
}
