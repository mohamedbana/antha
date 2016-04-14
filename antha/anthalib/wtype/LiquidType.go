package wtype

import (
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type LiquidType int

const (
	LTNIL = iota
	LTWater
	LTGlycerol
	LTEthanol
	LTDetergent
	LTCulture
	LTProtein
	LTDNA
	LTload
	LTDoNotMix
	LTloadwater
	LTNeedToMix
	LTVISCOUS
	LTPAINT
	LTDISPENSEABOVE
	LTPEG
	LTProtoplasts
)

func LiquidTypeFromString(s string) LiquidType {
	if strings.Contains(s, "DOE_run") {
		fields := strings.SplitAfter(s, "DOE_run")

		runnumber, err := strconv.Atoi(fields[1])
		if err != nil {
			panic("for Liguid type " + s + err.Error())
		}
		liquid := LiquidType(100 + runnumber)

		return liquid
	} else {

		switch s {
		case "water":
		case "":
			return LTWater
		case "glycerol":
			return LTGlycerol
		case "ethanol":
			return LTEthanol
		case "detergent":
			return LTDetergent
		case "culture":
			return LTCulture
		case "protein":
			return LTProtein
		case "dna":
			return LTDNA
		case "load":
			return LTload
		case "DoNotMix":
			return LTDoNotMix
		case "loadwater":
			return LTloadwater
		case "NeedToMix":
			return LTNeedToMix
		case "viscous":
			return LTVISCOUS
		case "Paint":
			return LTPAINT
		case "DispenseAboveLiquid":
			return LTDISPENSEABOVE
		case "PEG":
			return LTPEG
		case "Protoplasts":
			return LTProtoplasts
		default:
			return LTWater
		}

	}
	return LTWater
}

func LiquidTypeName(lt LiquidType) string {

	if lt > 99 {
		return "DOE_run" + strconv.Itoa(int(lt)-100)
	}

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
	case LTDoNotMix:
		return "DoNotMix"
	case LTloadwater:
		return "loadwater"
	case LTNeedToMix:
		return "NeedToMix"
	case LTPAINT:
		return "Paint"
	case LTDISPENSEABOVE:
		return "DispenseAboveLiquid"
	case LTProtoplasts:
		return "Protoplasts"
	case LTPEG:
		return "PEG"
	default:
		return "water"
	}
}

func mergeSolubilities(c1, c2 *LHComponent) float64 {
	if c1.Smax < c2.Smax {
		return c1.Smax
	}

	return c2.Smax
}

// helper functions... will need extending eventually

func mergeTypes(c1, c2 *LHComponent) LiquidType {
	// couple of mixing rules: protein, dna etc. are basically
	// special water so we retain that characteristic whatever happens
	// ditto culture... otherwise we look for the majority
	// what we do for protein-dna mixtures I'm not sure! :)

	// nil type is overridden

	if c1.Type == LTNIL {
		return c2.Type
	} else if c2.Type == LTNIL {
		return c1.Type
	}

	if c1.Type == LTCulture || c2.Type == LTCulture {
		return LTCulture
	} else if c1.Type == LTProtoplasts || c2.Type == LTProtoplasts {
		return LTProtoplasts
	} else if c1.Type == LTDNA || c2.Type == LTDNA {
		return LTDNA
	} else if c1.Type == LTProtein || c2.Type == LTProtein {
		return LTProtein
	}
	v1 := wunit.NewVolume(c1.Vol, c1.Vunit)
	v2 := wunit.NewVolume(c2.Vol, c2.Vunit)

	if v1.LessThan(&v2) {
		return c2.Type
	}

	return c1.Type
}

// merge two names... we have a lookup function to add here
func mergeNames(a, b string) string {
	tx := strings.Split(a, "+")
	tx2 := strings.Split(b, "+")

	tx3 := mergeTox(tx, tx2)

	tx3 = Normalize(tx3)

	return strings.Join(tx3, "+")
}

// very simple, just add elements of tx2 not already in tx
func mergeTox(tx, tx2 []string) []string {
	for _, v := range tx2 {
		ix := IndexOfStringArray(v, tx)

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
	if tx[0] == "" && len(tx) > 1 {
		return tx[1:len(tx)]
	}
	return tx
}
