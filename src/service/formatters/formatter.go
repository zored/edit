package formatters

import (
	"github.com/zored/edit/src/service/tokens"
	"strings"
)

type formatter struct{}

func NewFormatter() IFormatter {
	return &formatter{}
}

// TODO: refactor
func (f *formatter) Format(tokens_ tokens.Tokens, options *Options) string {
	result := ""
	depth := 0
	atomNameBefore := false
	atomSeparatorBefore := false

	first := true
	for _, token := range tokens_ {
		switch options.Rule {
		case InColumn:
			tokenType := token.TokenType
			if tokenType == tokens.WrapperEnd {
				depth--
			}

			addNewLine := !first && !atomNameBefore && tokenType != tokens.AtomSeparator
			if addNewLine {
				if tokenType == tokens.WrapperEnd && !atomSeparatorBefore {
					result += string(options.TrailingSeparator)
				}
				result += "\n" + strings.Repeat(" ", depth*options.Indent)
			}

			if tokenType == tokens.WrapperStart {
				depth++
			}

			atomNameBefore = tokenType == tokens.Name
			atomSeparatorBefore = tokenType == tokens.AtomSeparator
			first = false
		}
		result += token.Value
	}

	return result
}
