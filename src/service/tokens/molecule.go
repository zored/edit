package tokens

import "github.com/zored/edit/src/service/navigation"

// TokenBuffer with interval information.
type Molecule struct {
	Interval *navigation.Interval
	Tokens   Tokens
}

func NewEmptyMolecule() *Molecule {
	return NewMolecule(navigation.NewEmptyInterval(), make(Tokens, 0))
}

func NewMolecule(interval *navigation.Interval, tokens Tokens) *Molecule {
	return &Molecule{Interval: interval, Tokens: tokens}
}
