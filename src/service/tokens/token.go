package tokens

type (
	Token struct {
		Value     string
		TokenType TokenType
	}
	TokenType int
)

const (
	WrapperStart TokenType = iota
	WrapperEnd
	AtomSeparator
	Atom
	AtomName
)

func NewToken(runes []rune, tokenType TokenType) *Token {
	return &Token{Value: string(runes), TokenType: tokenType}
}

func NewStringToken(runes string, tokenType TokenType) *Token {
	return &Token{Value: runes, TokenType: tokenType}
}
