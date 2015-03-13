// Example pancake protocol(recipe).
// Provides instructions on how to make a pancake

package pancake
	
import "github.com/antha-lang/antha/execute"
import "github.com/antha-lang/goflow"
import "sync"
import "log"
import "bytes"
import "encoding/json"
import "io"


// import nuttin'
import ()

// Input parameters for this protocol (data)
// recipe from http://www.bbc.co.uk/food/recipes/fluffyamericanpancak_74828

//KASM = KitchenAid Stand Mixer, representing the numbered speeds available
//seconds

// tablespoons

// teaspoon

// HEN = Heat Element Number, arbitrarily set at 5

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// No special requirements on inputs
func (e *Pancake) requirements() {
	// Must be delicious
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func (e *Pancake) setup(p ParamBlock) {
	e.NumberPancakes <- execute.ThreadParam{0, p.ID}
}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Pancake) steps(p ParamBlock) {
	var batter = mix(p.AllPurposeFlour(p.FlourMass)+p.Sugar(p.SugarVolume)+p.BakingPowder(p.BakingPowderVolume)+p.Eggs(p.EggsEach)+p.Salt(p.SaltVolume)+p.Butter(p.ButterVolume)+p.Milk(p.MilkVolume), p.MixTime, p.KitchenAid(p.MixSpeed))
	e.BatterVolume <- execute.ThreadParam{human.observe(batter), p.ID}
	for e.BatterVolume > 0 {
		e.BatterVolume = e.BatterVolume - p.OneCakeBatterVolume
		e.ResultantPancake = fry(batter, p.OneCakeBatterVolume, p.FryTime, p.FryHeat)
		e.NumberPancakes++
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Pancake) analysis(p ParamBlock, r ResultBlock) {
	// None. Pancakes are delicious, everyone knows this.
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Pancake) validation(p ParamBlock, r ResultBlock) {
	if e.NumberPancakes < 0 {
		panic("Guys, did you forget to make the pancakes?")
	}
	if e.NumberPancakes > 100 {
		warn("Woah baby, that's a lot of pancakes!")
	}
}
// AsyncBag functions
func (e *Pancake) Complete(params interface{}) {
fmt.Println("generator.go:Complete")
	p := params.(ParamBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p)
	
}

// empty function for interface support
func (e *Pancake) anthaElement() {}
fmt.Println("generator.go:anthaElement")

// init function, read characterization info from seperate file to validate ranges?
func (e *Pancake) init() {
fmt.Println("generator.go:init")
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func New() *Pancake {
fmt.Println("generator.go:New")
	e := new(Pancake)
	e.init()
	return e
}

// Mapper function
func (e *Pancake) Map(m map[string]interface{}) interface{} {
fmt.Println("generator.go:Map")
	var res ParamBlock

	res.BowlSize = m["BowlSize"].(execute.ThreadParam).Value.(Volume)	

	res.MixSpeed = m["MixSpeed"].(execute.ThreadParam).Value.(KASMUnit)	

	res.MixTime = m["MixTime"].(execute.ThreadParam).Value.(Duration)	

	res.FlourMass = m["FlourMass"].(execute.ThreadParam).Value.(Mass)	

	res.SugarVolume = m["SugarVolume"].(execute.ThreadParam).Value.(Volume)	

	res.BakingPowderVolume = m["BakingPowderVolume"].(execute.ThreadParam).Value.(Volume)	

	res.EggsEach = m["EggsEach"].(execute.ThreadParam).Value.(Each)	

	res.SaltVolume = m["SaltVolume"].(execute.ThreadParam).Value.(Volume)	

	res.ButterVolume = m["ButterVolume"].(execute.ThreadParam).Value.(Volume)	

	res.MilkVolume = m["MilkVolume"].(execute.ThreadParam).Value.(Volume)	

	res.FryTime = m["FryTime"].(execute.ThreadParam).Value.(Duration)	

	res.FryHeat = m["FryHeat"].(execute.ThreadParam).Value.(HeatElementNumber)	

	res.OneCakeBatterVolume = m["OneCakeBatterVolume"].(execute.ThreadParam).Value.(Volume)	

	res.KitchenAid = m["KitchenAid"].(execute.ThreadParam).Value.(Mixer)	

	res.AllPurposeFlour = m["AllPurposeFlour"].(execute.ThreadParam).Value.(Flour)	

	res.Sugar = m["Sugar"].(execute.ThreadParam).Value.(Sweetener)	

	res.BakingPowder = m["BakingPowder"].(execute.ThreadParam).Value.(Leavener)	

	res.Eggs = m["Eggs"].(execute.ThreadParam).Value.(Leavener)	

	res.Salt = m["Salt"].(execute.ThreadParam).Value.(Seasoning)	

	res.Butter = m["Butter"].(execute.ThreadParam).Value.(Fat)	

	res.Milk = m["Milk"].(execute.ThreadParam).Value.(Liquid)	

	return res
}


type Pancake struct {
	flow.Component                    // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once	
	params         map[execute.ThreadID]*execute.AsyncBag
	BowlSize          <-chan execute.ThreadParam
	MixSpeed          <-chan execute.ThreadParam
	MixTime          <-chan execute.ThreadParam
	FlourMass          <-chan execute.ThreadParam
	SugarVolume          <-chan execute.ThreadParam
	BakingPowderVolume          <-chan execute.ThreadParam
	EggsEach          <-chan execute.ThreadParam
	SaltVolume          <-chan execute.ThreadParam
	ButterVolume          <-chan execute.ThreadParam
	MilkVolume          <-chan execute.ThreadParam
	FryTime          <-chan execute.ThreadParam
	FryHeat          <-chan execute.ThreadParam
	OneCakeBatterVolume          <-chan execute.ThreadParam
	KitchenAid          <-chan execute.ThreadParam
	AllPurposeFlour          <-chan execute.ThreadParam
	Sugar          <-chan execute.ThreadParam
	BakingPowder          <-chan execute.ThreadParam
	Eggs          <-chan execute.ThreadParam
	Salt          <-chan execute.ThreadParam
	Butter          <-chan execute.ThreadParam
	Milk          <-chan execute.ThreadParam
	NumberPancakes      chan<- execute.ThreadParam
	BatterVolume      chan<- execute.ThreadParam
	ResultantPancake      chan<- execute.ThreadParam
}

type ParamBlock struct {
	ID        execute.ThreadID
	BowlSize Volume
	MixSpeed KASMUnit
	MixTime Duration
	FlourMass Mass
	SugarVolume Volume
	BakingPowderVolume Volume
	EggsEach Each
	SaltVolume Volume
	ButterVolume Volume
	MilkVolume Volume
	FryTime Duration
	FryHeat HeatElementNumber
	OneCakeBatterVolume Volume
	KitchenAid Mixer
	AllPurposeFlour Flour
	Sugar Sweetener
	BakingPowder Leavener
	Eggs Leavener
	Salt Seasoning
	Butter Fat
	Milk Liquid
}

type ResultBlock struct {
	ID        execute.ThreadID
	NumberPancakes Number
	BatterVolume Volume
	ResultantPancake Pancake
}

type JSONBlock struct {
	ID        *execute.ThreadID
	NumberPancakes *Number
	BatterVolume *Volume
	ResultantPancake *Pancake
}

