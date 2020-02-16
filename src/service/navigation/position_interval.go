package navigation

type Interval struct {
	Start, Stop *Position
}

func NewInterval(start *Position, stop *Position) *Interval {
	return &Interval{Start: start, Stop: stop}
}
