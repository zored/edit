package saver

import (
	"github.com/zored/edit/src/service/formatters"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
)

type FileOptions struct {
	file string

	cursor                  *navigation.Position
	wrappers                *tokens.Wrappers
	separator               tokens.Separator
	formatRule              formatters.Rule
	toggleTrailingSeparator bool
	indent                  int
}

func (c FileOptions) GetFormatterOptions() *formatters.Options {
	trailingSeparator := tokens.Separator("")
	if c.toggleTrailingSeparator {
		trailingSeparator = c.separator
	}

	return &formatters.Options{
		Indent:            c.indent,
		TrailingSeparator: trailingSeparator,
		Rule:              c.formatRule,
	}
}

func NewFileOptions(file string, cursor *navigation.Position, wrappers *tokens.Wrappers, formatRule formatters.Rule, trailingSeparator bool) *FileOptions {
	return &FileOptions{
		cursor:                  cursor,
		wrappers:                wrappers,
		separator:               ",",
		formatRule:              formatRule,
		file:                    file,
		indent:                  4,
		toggleTrailingSeparator: trailingSeparator,
	}
}
