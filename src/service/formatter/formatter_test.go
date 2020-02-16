package formatter

import (
	"github.com/magiconair/properties/assert"
	"github.com/zored/edit/src/service/tokens"
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	tokens_ := []*tokens.Token{
		tokens.NewStringToken("{", tokens.AtomWrapStart),
		tokens.NewStringToken(`"a":"b"`, tokens.Atom),
		tokens.NewStringToken(`,`, tokens.AtomSeparator),
		tokens.NewStringToken(`"c":`, tokens.AtomName),
		tokens.NewStringToken("{", tokens.AtomWrapStart),
		tokens.NewStringToken(`"d":"e"`, tokens.Atom),
		tokens.NewStringToken("}", tokens.AtomWrapEnd),
		tokens.NewStringToken("}", tokens.AtomWrapEnd),
	}
	code := NewFormatter().Format(tokens_, in_line, 0)
	assert.Equal(t, code, `{"a":"b","c":{"d":"e"}}`)
	/*
	code = NewFormatter().Format(tokens_, in_column, 2)
	assert.Equal(t, code, `{
  "a":"b",
  "c":{
    "d":"e"
  }
}`)
	 */
}