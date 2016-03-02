package ast

// TODO make more specific
type Location interface{}

type Movement struct {
	From Location
	To   Location
}

type Request struct {
	MixVol *Interval
	Temp   *Interval
	Time   *Interval
	Move   *Movement
}
