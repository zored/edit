package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	interval_, tokens_, err := newParser().Parse(strings.NewReader(`
{"a":"b"  , "c": { "d":"e"   }  }
`), navigation.NewPosition(1, 4), tokens.NewAtomWrapper("{", "}"), ",")

	newTokenString := func(s string, t tokens.TokenType) *tokens.Token { return tokens.NewToken([]rune(s), t) }
	assert.Nil(t, err)
	assert.Equal(t, 1, interval_.Start.Column)
	assert.Equal(t, 2, interval_.Start.Line)
	assert.Equal(t, 32, interval_.Stop.Column)
	assert.Equal(t, 2, interval_.Stop.Line)
	assert.Equal(t, []*tokens.Token{
		newTokenString("{", tokens.AtomWrapStart),
		newTokenString(`"a":"b"`, tokens.Atom),
		newTokenString(`,`, tokens.AtomSeparator),
		newTokenString(`"c":`, tokens.AtomName),
		newTokenString("{", tokens.AtomWrapStart),
		newTokenString(`"d":"e"`, tokens.Atom),
		newTokenString("}", tokens.AtomWrapEnd),
		newTokenString("}", tokens.AtomWrapEnd),
	}, tokens_.All)
}
