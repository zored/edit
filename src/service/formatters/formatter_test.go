package formatters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zored/edit/src/service/tokens"
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	start := tokens.NewStringToken("{", tokens.WrapperStart)
	stop := tokens.NewStringToken("}", tokens.WrapperEnd)
	separator := tokens.NewStringToken(`,`, tokens.AtomSeparator)

	tokens_ := []*tokens.Token{
		start,
		tokens.NewStringToken(`"a":"b"`, tokens.Atom),
		separator,
		tokens.NewStringToken(`"c":`, tokens.AtomName),
		start,
		tokens.NewStringToken(`"d":`, tokens.AtomName),
		start,
		stop,
		stop,
		separator,
		tokens.NewStringToken(`"e":`, tokens.AtomName),
		start,
		tokens.NewStringToken(`"f":`, tokens.AtomName),
		start,
		tokens.NewStringToken(`"g":"h"`, tokens.Atom),
		separator,
		tokens.NewStringToken(`"i":"j"`, tokens.Atom),
		stop,
		stop,
		stop,
	}
	options := &Options{Indent: 2, TrailingSeparator: "", Rule: Line}
	formatter := NewFormatter()

	options.Rule = Line
	assert.Equal(t, `{"a":"b","c":{"d":{}},"e":{"f":{"g":"h","i":"j"}}}`, formatter.Format(tokens_, options))

	options.Rule = Tree
	assert.Equal(t, `{
  "a":"b",
  "c":{
    "d":{}
  },
  "e":{
    "f":{
      "g":"h",
      "i":"j"
    }
  }
}`, formatter.Format(tokens_, options))

	options.TrailingSeparator = ","
	assert.Equal(t, `{
  "a":"b",
  "c":{
    "d":{},
  },
  "e":{
    "f":{
      "g":"h",
      "i":"j",
    },
  },
}`, formatter.Format(tokens_, options))

	options.Rule = Column
	assert.Equal(t, `{
  "a":"b",
  "c":{"d":{}},
  "e":{"f":{"g":"h","i":"j"}},
}`, formatter.Format(tokens_, options))
}
