package formatters

import (
	"github.com/magiconair/properties/assert"
	"github.com/zored/edit/src/service/tokens"
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	tokens_ := []*tokens.Token{
		tokens.NewStringToken("{", tokens.WrapperStart),
		tokens.NewStringToken(`"a":"b"`, tokens.Atom),
		tokens.NewStringToken(`,`, tokens.AtomSeparator),
		tokens.NewStringToken(`"c":`, tokens.Name),
		tokens.NewStringToken("{", tokens.WrapperStart),
		tokens.NewStringToken(`"d":"e"`, tokens.Atom),
		tokens.NewStringToken("}", tokens.WrapperEnd),
		tokens.NewStringToken("}", tokens.WrapperEnd),
	}
	options := &Options{Indent: 2, TrailingSeparator: "", Rule: InLine}
	formatter := NewFormatter()

	assert.Equal(t, formatter.Format(tokens_, options), `{"a":"b","c":{"d":"e"}}`)

	options.Rule = InColumn
	assert.Equal(t, formatter.Format(tokens_, options), `{
  "a":"b",
  "c":{
    "d":"e"
  }
}`)

	options.TrailingSeparator = ","
	assert.Equal(t, formatter.Format(tokens_, options), `{
  "a":"b",
  "c":{
    "d":"e",
  },
}`)
}
