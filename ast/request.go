package ast

type Location struct{}

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
