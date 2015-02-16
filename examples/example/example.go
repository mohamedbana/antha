// Simple do nothing protocol for hello world example
package Example
	
import "github.com/antha-lang/antha/execute"
import "github.com/antha-lang/goflow"
import "sync"
import "log"
import "bytes"
import "encoding/json"
import "io"


//import "github.com/antha-lang/antha/examples/example"

import (
	"time"
)

// no physical inputs

// none

// no physical outputs

// none

func (e *Example) setup(p ParamBlock) {
	// None
}

func (e *Example) steps(p ParamBlock) {
	time.Sleep(p.SleepTime)
	OutColor = p.Color
}
// AsyncBag functions
func (e *Example) Complete(params interface{}) {
	p := params.(ParamBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p)
	
}

// empty function for interface support
func (e *Example) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *Example) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func New() *Example {
	e := new(Example)
	e.init()
	return e
}

// Mapper function
func (e *Example) Map(m map[string]interface{}) interface{} {
	var res ParamBlock

	res.Color = m["Color"].(execute.ThreadParam).Value.(string)	

	res.SleepTime = m["SleepTime"].(execute.ThreadParam).Value.(time.Duration)	

	return res
}


type Example struct {
	flow.Component                    // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once	
	params         map[execute.ThreadID]*execute.AsyncBag
	Color          <-chan execute.ThreadParam
	SleepTime          <-chan execute.ThreadParam
	WellColor      chan<- execute.ThreadParam
}

type ParamBlock struct {
	ID        execute.ThreadID
	Color string
	SleepTime time.Duration
}

type ResultBlock struct {
	ID        execute.ThreadID
	WellColor string
}

type JSONBlock struct {
	ID        *execute.ThreadID
	WellColor *string
}

