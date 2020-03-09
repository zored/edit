package tokenizers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
	"strings"
	"testing"
)

type tokenizeTestData struct {
	name             string
	input            string
	expectedMolecule *tokens.Molecule
	position         *navigation.Position
	atomWrapper      *tokens.Wrappers
	separator        tokens.Separator
}

func TestTokenizer_Parse(t *testing.T) {
	tokenizer := NewTokenizer()
	index := 0
	for _, data := range []tokenizeTestData{
		{
			name:        "flat json",
			position:    navigation.NewPosition(1, 4),
			atomWrapper: tokens.NewWrappers("{", "}"),
			input:       `{"a":"b"  , "c": { "d":"e"   }  }`,
			expectedMolecule: tokens.NewMolecule(
				navigation.NewInterval(
					navigation.NewPosition(1, 1),
					navigation.NewPosition(1, 34),
				),
				tokens.Tokens{
					tokens.NewStringToken("{", tokens.WrapperStart),
					tokens.NewStringToken(`"a":"b"`, tokens.Atom),
					tokens.NewStringToken(`,`, tokens.AtomSeparator),
					tokens.NewStringToken(`"c":`, tokens.AtomName),
					tokens.NewStringToken("{", tokens.WrapperStart),
					tokens.NewStringToken(`"d":"e"`, tokens.Atom),
					tokens.NewStringToken("}", tokens.WrapperEnd),
					tokens.NewStringToken("}", tokens.WrapperEnd),
				},
			),
			separator: ",",
		},
		{
			name:        "method",
			position:    navigation.NewPosition(3, 22),
			atomWrapper: tokens.NewWrappers("(", ")"),
			input: `
func (p *tokenizer) getLinesTillCursor(
  cursor *navigation.Position,
  linesScanner *bufio.Scanner,
) ([]string, error) {
`,
			expectedMolecule: tokens.NewMolecule(
				navigation.NewInterval(
					navigation.NewPosition(2, 39),
					navigation.NewPosition(5, 2),
				),
				tokens.Tokens{
					tokens.NewStringToken("(", tokens.WrapperStart),
					tokens.NewStringToken(`cursor *navigation.Position`, tokens.Atom),
					tokens.NewStringToken(`,`, tokens.AtomSeparator),
					tokens.NewStringToken(`linesScanner *bufio.Scanner`, tokens.Atom),
					tokens.NewStringToken(`,`, tokens.AtomSeparator),
					tokens.NewStringToken(")", tokens.WrapperEnd),
				},
			),
			separator: ",",
		},
	} {
		t.Run(fmt.Sprintf("test %s", data.name), func(t *testing.T) {
			mol, err := tokenizer.Tokenize(
				strings.NewReader(data.input),
				data.position,
				data.atomWrapper,
				data.separator,
			)
			assert.Nil(t, err)
			assert.Equal(t, data.expectedMolecule, mol)
		})
		index++
	}

}
