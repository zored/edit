package tokens

type Wrappers struct {
	Start, End string
}

func NewWrappers(start string, end string) *Wrappers {
	return &Wrappers{Start: start, End: end}
}
