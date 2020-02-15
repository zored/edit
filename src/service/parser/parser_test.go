package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	interval_, tokens_, err := newParser().Parse(strings.NewReader(`
{"a":"b"  , "c": { "d":"e"   }  }
`), newPosition(1, 4), newAtomWrapper("{", "}"), ",")

	newTokenString := func(s string, t tokenType) *token { return newToken([]rune(s), t) }
	assert.Nil(t, err)
	assert.Equal(t, 1, interval_.start.column)
	assert.Equal(t, 2, interval_.start.line)
	assert.Equal(t, 32, interval_.stop.column)
	assert.Equal(t, 2, interval_.stop.line)
	assert.Equal(t, []*token{
		newTokenString("{", atomWrapStart),
		newTokenString(`"a":"b"`, atom),
		newTokenString(`,`, atomSeparator),
		newTokenString(`"c":`, atomName),
		newTokenString("{", atomWrapStart),
		newTokenString(`"d":"e"`, atom),
		newTokenString("}", atomWrapEnd),
		newTokenString("}", atomWrapEnd),
	}, tokens_.all)
}
