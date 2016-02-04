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
