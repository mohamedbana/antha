package liquidhandling

import "reflect"

func Rationalise_Outputs(lhr *LHRequest, lhp *LHProperties, outputs ...*LHComponent) *LHRequest {

	op := make(map[string]*LHPlate)

	for i := 0; i < len(outputs); i++ {
		// need to check that the outputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic
		n := reflect.TypeOf(outputs[i].Container())
		if !((n.Kind() == reflect.Ptr && n.Elem().Name() == "LHWell") || n.Name() == "liquidhandling.LHWell") {
			panic("Rationalise_outputs: cannot use non-microplate formats in liquid handler... yet!")
		}

		op[outputs[i].LContainer.Plate.ID] = outputs[i].LContainer.Plate

	}

	lhr.Output_plates = op

	return lhr
}
