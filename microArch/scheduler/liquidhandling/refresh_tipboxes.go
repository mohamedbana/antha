package liquidhandling

// no longer need to supply tipboxes after the fact

func (lh *Liquidhandler) Refresh_tipboxes(rq *LHRequest) {

	// dead simple

	for pos, _ := range lh.Properties.PosLookup {
		tb, ok := lh.Properties.Tipboxes[pos]

		if !ok {
			continue
		}

		lh.FinalProperties.AddTipBoxTo(pos, tb.Dup())
		tb.Refresh()
	}
}
