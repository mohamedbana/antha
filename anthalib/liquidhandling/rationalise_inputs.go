package liquidhandling

import "reflect"

func Rationalise_Inputs(lhr *LHRequest, lhp *LHProperties, inputs ...*LHComponent) *LHRequest {

	ip := make(map[string]*LHPlate)

	for i := 0; i < len(inputs); i++ {
		// need to check that the inputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic
		n := reflect.TypeOf(inputs[i].Container())
		if !((n.Kind() == reflect.Ptr && n.Elem().Name() == "LHWell") || n.Name() == "liquidhandling.LHWell") {
			panic("Rationalise_inputs: cannot use non-microplate formats in liquid handler... yet!")
		}

		ip[inputs[i].LContainer.Plate.ID] = inputs[i].LContainer.Plate

	}

	lhr.Input_plates = ip

	return lhr
}
