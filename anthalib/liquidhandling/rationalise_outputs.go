package liquidhandling

import "reflect"

func Rationalise_Outputs(lhr *LHRequest, lhp *LHProperties, outputs ...*LHComponent) *LHRequest {

	op := make(map[string]*LHPlate)

	for i := 0; i < len(outputs); i++ {
		// need to check that the outputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic

		if reflect.TypeOf(outputs[i].Container).Name() != "*liquidhandling.LHWell" {
			panic("Cannot use non-microplate formats in liquid handler... yet!")
		}

		op[outputs[i].Container.Plate.ID] = outputs[i].Container.Plate

	}

	lhr.Output_plates = op

	return lhr
}
