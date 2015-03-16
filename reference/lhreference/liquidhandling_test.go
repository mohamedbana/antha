package lhreference

func ExampleOne() {
	var lhr LHRequest

	lhr = make(LHRequest, 5)

	sarr := make(map[string]LHSolution, 1)

	welltype := NewLHWell("PlateOneDW96ConicalBottom", "", "", 2000, 25, 2, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := NewLHPlate("PlateOneDW96ConicalBottom", "PlateOne", 8, 12, 44.1, "mm", welltype)

	lhr["output_platetype"] = plate
	lhr["input_platetype"] = plate

	for i := 0; i < 10; i++ {
		s := NewLHSolution()
		s["name"] = fmt.Sprintf("Solution%02d", i)
		cmp := make([]LHComponent, 0, 3)
		for j := 0; j < 3; j++ {
			c := NewLHComponent()
			c["name"] = fmt.Sprintf("Component%02d", j)
			c["type"] = "water"
			c["vol"] = 40.0
			c["vunit"] = "ul"
			cmp = append(cmp, c)
		}
		s["components"] = cmp
		sarr[s["id"].(string)] = s
	}

	lhr["output_solutions"] = sarr

	// make tips

	tipboxes := make([]LHTipbox, 2, 2)
	tip := new_tip("cybio", "250", 250.0)
	for i := 0; i < 2; i++ {
		tb := new_tipbox(8, 12, "CyBio", tip)
		tipboxes[i] = tb
	}

	lhr["tips"] = tipboxes

	// make a liquid handling structure

	lhp := NewLHProperties(12, "Felix", "CyBio", "discrete", "fixed", []string{"plate"})

	// I suspect this might need to be in the constructor
	// or at least wrapped into a factory method

	lhp["tip_preferences"] = []int{1, 5, 3}
	lhp["input_preferences"] = []int{10, 11, 12}
	lhp["output_preferences"] = []int{7, 8, 9, 2, 4}

	// need to add some configs

	hvconfig := NewLHParameter("HVconfig", 10, 250, "ul")

	cnfvol := lhp["cnfvol"].([]LHParameter)
	cnfvol[0] = hvconfig
	lhp["cnfvol"] = cnfvol
	lhp["cmnvol"] = 10.0
	lhp["cmxvol"] = 250.0
	lhp["vlunit"] = "ul"

	// these depend on the tip

	lhp["minvol"] = 10.0
	lhp["maxvol"] = 250.0

	liquidhandler := Init(lhp)
	liquidhandler.MakeSolutions(lhr)
}

func ExampleTwo() {
	names := []string{"tea", "milk", "sugar"}

	minrequired := make(map[string]float64, len(names))
	maxrequired := make(map[string]float64, len(names))
	Smax := make(map[string]float64, len(names))
	T := make(map[string]float64, len(names))
	vmin := 10.0

	for _, name := range names {
		r := rand.Float64() + 1.0
		r2 := rand.Float64() + 1.0
		r3 := rand.Float64() + 1.0

		minrequired[name] = r * r2 * 20.0
		maxrequired[name] = r * r2 * 30.0
		Smax[name] = r * r2 * r3 * 70.0
		T[name] = 100.0
	}

	cncs := choose_stock_concentrations(minrequired, maxrequired, Smax, vmin, T)

	for i, _ := range minrequired {
		var v float64
		v, ok := cncs[i]

		if !ok {
			v = -1.0
		}
		fmt.Printf("Concentration of %10s = %8.1f, volume High: %-8.1f volume Low: %-8.1f min required: %-8.1f Max required: %-8.1f Smax: %-8.1f T: %-6.1f\n", i, v, T[i]*maxrequired[i]/v, T[i]*minrequired[i]/v, minrequired[i], maxrequired[i], Smax[i], T[i])
	}

	fmt.Println()
}
