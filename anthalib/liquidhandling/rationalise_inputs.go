package liquidhandling

func Rationalise_Inputs(lhr *LHRequest, lhp *LHProperties, inputs ...*LHSolution) *LHRequest {

	for i := 0; i < len(inputs); i++ {
		// need to check that the inputs are in the right format
		// at the moment it's just a case of figuring out if they're in plates
		// already... if not we panic

	}

	return lhr
}
