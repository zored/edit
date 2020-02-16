package formatter

import (
	"github.com/zored/edit/src/service/tokens"
	"strings"
)

type formatter struct{}

func NewFormatter() *formatter {
	return &formatter{}
}

const (
	in_column = iota
	in_line   = 1
)

func (f *formatter) Format(tokens_ []*tokens.Token, t, ident int) string {
	result := ""
	depth := 0
	atomNameBefore := false
	first := true
	for _, token := range tokens_ {
		switch t {
		case in_column:
			tokenType := token.TokenType
			if tokenType == tokens.AtomWrapEnd {
				depth--
			}

			if first || !atomNameBefore {
				switch tokenType {
				case tokens.AtomSeparator:
				default:
					result += "\n" + strings.Repeat(" ", depth*ident)
				}
			}

			if tokenType == tokens.AtomWrapStart {
				depth++
			}

			atomNameBefore = tokenType == tokens.AtomName
			first = false
		}
		result += token.Value
	}

	return result
}
