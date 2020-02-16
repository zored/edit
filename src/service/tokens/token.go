package tokens

type Token struct {
	value     string
	tokenType TokenType
}

func NewToken(runes []rune, tokenType TokenType) *Token {
	return &Token{value: string(runes), tokenType: tokenType}
}
