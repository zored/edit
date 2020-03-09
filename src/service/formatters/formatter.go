package formatters

import (
	"github.com/zored/edit/src/service/tokens"
	"math"
	"strings"
)

type formatter struct{}

func NewFormatter() IFormatter {
	return &formatter{}
}

func (f *formatter) Format(tokens_ tokens.Tokens, options *Options) string {
	result := ""
	depthNow := 0
	depthNext := 0
	var typeBefore tokens.TokenType

	treeDepth := math.MaxInt32
	switch options.Rule {
	case Line:
		treeDepth = 0
	case Column:
		treeDepth = 2
	}

	first := true
	for _, token := range tokens_ {
		typeNow := token.TokenType

		switch typeNow {
		case tokens.WrapperEnd:
			depthNext--
		case tokens.WrapperStart:
			depthNext++
		}

		sameLine := first ||
			typeBefore == tokens.AtomName ||
			typeNow == tokens.AtomSeparator ||
			(typeBefore == tokens.WrapperStart && typeNow == tokens.WrapperEnd) ||
			(depthNow >= treeDepth)
		if !sameLine {
			result += nextLine(typeNow, typeBefore, options, depthNext)
		}

		first = false
		typeBefore = typeNow
		depthNow = depthNext
		result += token.Value
	}

	return result
}

func nextLine(typeNow, typeBefore tokens.TokenType, options *Options, depthNext int) string {
	// Separator:
	separator := ""
	if typeNow == tokens.WrapperEnd && typeBefore != tokens.AtomSeparator {
		separator = string(options.TrailingSeparator)
	}

	return separator + "\n" + strings.Repeat(" ", depthNext*options.Indent)
}
