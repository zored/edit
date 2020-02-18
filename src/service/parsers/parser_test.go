package parsers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
	"strings"
	"testing"
)

type parseTestData struct {
	name        string
	input       string
	interval    *navigation.Interval
	tokens      []*tokens.Token
	position    *navigation.Position
	atomWrapper *tokens.AtomWrapper
	separator   tokens.Separator
}

func TestParser_Parse(t *testing.T) {
	index := 0
	for _, data := range []parseTestData{
		{
			name:        "flat json",
			position:    navigation.NewPosition(1, 4),
			atomWrapper: tokens.NewAtomWrapper("{", "}"),
			input:       `{"a":"b"  , "c": { "d":"e"   }  }`,
			interval: navigation.NewInterval(
				navigation.NewPosition(1, 1),
				navigation.NewPosition(1, 34),
			),
			tokens: []*tokens.Token{
				tokens.NewStringToken("{", tokens.AtomWrapStart),
				tokens.NewStringToken(`"a":"b"`, tokens.Atom),
				tokens.NewStringToken(`,`, tokens.AtomSeparator),
				tokens.NewStringToken(`"c":`, tokens.AtomName),
				tokens.NewStringToken("{", tokens.AtomWrapStart),
				tokens.NewStringToken(`"d":"e"`, tokens.Atom),
				tokens.NewStringToken("}", tokens.AtomWrapEnd),
				tokens.NewStringToken("}", tokens.AtomWrapEnd),
			},
			separator: ",",
		},
		{
			name:        "method",
			position:    navigation.NewPosition(2, 49),
			atomWrapper: tokens.NewAtomWrapper("(", ")"),
			input: `
func (p *parser) getLinesTillCursor(cursor *navigation.Position, linesScanner *bufio.Scanner) ([]string, error) {
`,
			interval: navigation.NewInterval(
				navigation.NewPosition(2, 36),
				navigation.NewPosition(2, 94),
			),
			tokens: []*tokens.Token{
				tokens.NewStringToken("(", tokens.AtomWrapStart),
				tokens.NewStringToken(`cursor *navigation.Position`, tokens.Atom),
				tokens.NewStringToken(`,`, tokens.AtomSeparator),
				tokens.NewStringToken(`linesScanner *bufio.Scanner`, tokens.Atom),
				tokens.NewStringToken(")", tokens.AtomWrapEnd),
			},
			separator: ",",
		},
	} {
		t.Run(fmt.Sprintf("test %s", data.name), func(t *testing.T) {
			interval, tokens_, err := NewParser().Parse(
				strings.NewReader(data.input),
				data.position,
				data.atomWrapper,
				data.separator,
			)
			assert.Nil(t, err)
			assert.Equal(t, data.interval, interval)
			assert.Equal(t, data.tokens, tokens_.All)
		})
		index++
	}

}
