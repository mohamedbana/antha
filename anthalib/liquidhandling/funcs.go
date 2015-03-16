package liquidhandling

import "fmt"

func RaiseError(err string) {
	// TODO remove this
	fmt.Println("liquidhandling raiseError called: remove this to win")
	fmt.Println(err)
	fmt.Println("error done")
	panic("NO")
}

// looks up where a plate is mounted on a liquid handler as expressed in a request
func PlateLookup(rq LHRequest, id string) int {
	lookupmap := rq.Plate_lookup

	if len(lookupmap) == 0 {
		RaiseError("Cannot find plate lookup")
	}

	return lookupmap[id]
}
