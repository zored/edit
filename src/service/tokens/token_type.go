package tokens

type TokenType int

const (
	AtomWrapStart TokenType = iota
	AtomWrapEnd
	AtomSeparator
	Atom
	AtomName
)
