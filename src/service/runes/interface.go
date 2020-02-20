package runes

type IRunesMatcher interface {
	Add(r rune) bool
	Matches() int
}
