package wtype

type LiquidType int

const (
	Water = iota
	Glycerol
	Ethanol
)

func LiquidTypeName(lt LiquidType) string {
	switch lt {
	case Water:
		return "water"
	case Glycerol:
		return "glycerol"
	case Ethanol:
		return "ethanol"
	default:
		return "water"
	}
}
