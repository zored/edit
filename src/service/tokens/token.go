package tokens

type Token struct {
	Value     string
	TokenType TokenType
}

func NewToken(runes []rune, tokenType TokenType) *Token {
	return &Token{Value: string(runes), TokenType: tokenType}
}

func NewStringToken(runes string, tokenType TokenType) *Token {
	return &Token{Value: runes, TokenType: tokenType}
}
