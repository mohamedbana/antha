package wutil

func StrInStrArray(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}

	return false
}
