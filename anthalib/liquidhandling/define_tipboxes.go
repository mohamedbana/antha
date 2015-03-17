package liquidhandling

func Define_Tipboxes(lhr *LHRequest, lhp *LHProperties) *LHRequest {
	// make tips
	// this needs to be refactored elsewhere
	tipboxes := make([]*LHTipbox, 2, 2)
	tip := NewLHTip("ACMEliquidhandlers", "ACMEliquidhandlers250", 20.0, 250.0)
	for i := 0; i < 2; i++ {
		tb := NewLHTipbox(8, 12, "ACMEliquidhandlers", tip)
		tipboxes[i] = tb
	}
	lhr.Tips = tipboxes

	return lhr
}
