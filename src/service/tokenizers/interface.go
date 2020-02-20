package tokenizers

import (
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/tokens"
	"io"
)

type ITokenizer interface {
	// Finds and retrieves molecule starting to the left of cursor.
	Tokenize(reader io.Reader, cursor *navigation.Position, wrapper *tokens.Wrappers, separator_ tokens.Separator, ) (mol *tokens.Molecule, err error)
}
