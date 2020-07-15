package navigation

type Interval struct {
	Start, Stop *Position
}

func NewEmptyInterval() *Interval { return &Interval{Start: NewZeroPosition(), Stop: NewZeroPosition()} }
func NewInterval(start *Position, stop *Position) *Interval {
	return &Interval{Start: start, Stop: stop}
}
