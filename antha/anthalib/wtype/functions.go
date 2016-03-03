package wtype

func CopyComponentArray(arin []*LHComponent) []*LHComponent {
	r := make([]*LHComponent, len(arin))

	for i, v := range arin {
		r[i] = v.Dup()
	}

	return r
}
