package formatters

import (
	"github.com/zored/edit/src/service/tokens"
	"strings"
)

type (
	formatter  struct{}
	IFormatter interface {
		Format(tokens_ []*tokens.Token, formatRule FormatRule, ident int, separator tokens.Separator) string
	}
)

func NewFormatter() IFormatter {
	return &formatter{}
}

func (f *formatter) Format(tokens_ []*tokens.Token, formatRule FormatRule, ident int, separator tokens.Separator) string {
	result := ""
	depth := 0
	atomNameBefore := false
	atomSeparatorBefore := false

	first := true
	for _, token := range tokens_ {
		switch formatRule {
		case InColumn:
			tokenType := token.TokenType
			if tokenType == tokens.AtomWrapEnd {
				depth--
			}

			addNewLine := !first && !atomNameBefore && tokenType != tokens.AtomSeparator
			if addNewLine {
				if tokenType == tokens.AtomWrapEnd && !atomSeparatorBefore {
					result += string(separator)
				}
				result += "\n" + strings.Repeat(" ", depth*ident)
			}

			if tokenType == tokens.AtomWrapStart {
				depth++
			}

			atomNameBefore = tokenType == tokens.AtomName
			atomSeparatorBefore = tokenType == tokens.AtomSeparator
			first = false
		}
		result += token.Value
	}

	return result
}
