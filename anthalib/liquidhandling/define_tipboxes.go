package liquidhandling

func Define_Tipboxes(lhr *LHRequest, lhp *LHProperties) {
	// make tips
	// this needs to be refactored elsewhere
	tipboxes := make([]*liquidhandling.LHTipbox, 2, 2)
	tip := liquidhandling.NewLHTip("ACMEliquidhandlers", "ACMEliquidhandlers250", 20.0, 250.0)
	for i := 0; i < 2; i++ {
		tb := liquidhandling.NewLHTipbox(8, 12, "ACMEliquidhandlers", tip)
		tipboxes[i] = tb
	}
	lhr.Tips = tipboxes

}
