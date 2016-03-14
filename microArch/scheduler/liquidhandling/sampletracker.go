// sample tracker service... needs to be outside this eventually but at least having one
// is a good start

package liquidhandler

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
	st.records[ID] = loc
}

func (st *SampleTracker) GetLocationOf(ID string) (string, bool) {
	return st.records[ID]
}
