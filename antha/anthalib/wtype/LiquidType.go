package wtype

type LiquidType int

const (
	LTWater = iota
	LTGlycerol
	LTEthanol
	LTDetergent
	LTCulture
	LTProtein
	LTDNA
	LTload
)

func LiquidTypeName(lt LiquidType) string {
	switch lt {
	case LTWater:
		return "water"
	case LTGlycerol:
		return "glycerol"
	case LTEthanol:
		return "ethanol"
	case LTDetergent:
		return "detergent"
	case LTCulture:
		return "culture"
	case LTProtein:
		return "protein"
	case LTDNA:
		return "dna"
	case LTload:
		return "load"
	default:
		return "water"
	}
}

// helper functions... will need extending eventually

// merge two names... we have a lookup function to add here
func mergeNames(a, b string) string {
	tx := strings.Split(a, "_")
	tx2 := strings.Split(b, "_")

	tx3 := mergeTox(tx, tx2)

	tx3 = Normalize(tx3)

	return strings.Join(tx3, "_")
}

// very simple, just add elements of tx2 not already in tx
func mergeTox(tx, tx2 []string) []string {
	for _, v := range tx2 {
		ix := IndexOfStringArray(tx, v)

		if ix == -1 {
			tx = append(tx, v)
		}
	}

	return tx
}

func IndexOfStringArray(s string, a []string) int {
	ret := -1
	for i, v := range a {
		if v == s {
			ret = i
			break
		}
	}
	return ret
}

// TODO -- fill in some normalizations
// - water + salt = saline might be an eg
// but unlikely to be useful
func Normalize(tx []string) []string {
	return tx
}
