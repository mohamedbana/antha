package liquidhandling

import "reflect"

func Rationalise_Inputs(lhr *LHRequest, lhp *LHProperties, inputs ...*LHComponent) *LHRequest {

	ip := make(map[string]*LHPlate)

	for i := 0; i < len(inputs); i++ {
		// need to check that the inputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic

		if reflect.TypeOf(inputs[i].Container).Name() != "*liquidhandling.LHWell" {
			panic("Cannot use non-microplate formats in liquid handler... yet!")
		}

		ip[inputs[i].Container.Plate.ID] = inputs[i].Container.Plate

	}

	lhr.Input_plates = ip

	return lhr
}
