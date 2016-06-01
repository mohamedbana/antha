package sampletracker

var st *SampleTracker

type SampleTracker struct {
	records  map[string]string
	forwards map[string]string
}

func newSampleTracker() *SampleTracker {
	r := make(map[string]string)
	f := make(map[string]string)
	st := SampleTracker{r, f}
	return &st
}

func GetSampleTracker() *SampleTracker {
	if st == nil {
		st = newSampleTracker()
	}

	return st
}

func (st *SampleTracker) SetLocationOf(ID string, loc string) {
	//fmt.Println("LOCATION OF ", ID, " SET TO ", loc)
	st.records[ID] = loc
}

func (st *SampleTracker) GetLocationOf(ID string) (string, bool) {
	//fmt.Println("GET LOCATION OF :", ID)
	if ID == "" {
		return "", false
	}

	s, ok := st.records[ID]

	// look to see if there's a forwarding address
	// can this lead to an out of date location???

	if !ok {
		return st.GetLocationOf(st.forwards[ID])
	}

	return s, ok
}

func (st *SampleTracker) UpdateIDOf(ID string, newID string) {
	_, ok := st.records[ID]
	if ok {
		st.records[newID] = st.records[ID]
	} else {
		// set up a forward
		// actually a backward...
		st.forwards[newID] = ID
	}
}
