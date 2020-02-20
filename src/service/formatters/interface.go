package formatters

import "github.com/zored/edit/src/service/tokens"

type (
	IFormatter interface {
		Format(tokens_ tokens.Tokens, options *Options) string
	}
)