package files

import (
	"github.com/zored/edit/src/service/formatters"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
)

type FileFormatConfig struct {
	cursor     *navigation.Position
	wrapper    *tokens.AtomWrapper
	separator  tokens.Separator
	formatRule formatters.FormatRule
	file       string
	indent     int
}

func NewFileFormatConfig(file string, cursor *navigation.Position) *FileFormatConfig {
	return &FileFormatConfig{
		cursor:     cursor,
		wrapper:    tokens.NewAtomWrapper("(", ")"),
		separator:  ",",
		formatRule: formatters.InColumn,
		file:       file,
		indent:     4,
	}
}
