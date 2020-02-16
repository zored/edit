package tokens

type AtomWrapper struct {
	Start, End string
}

func NewAtomWrapper(start string, end string) *AtomWrapper {
	return &AtomWrapper{Start: start, End: end}
}
