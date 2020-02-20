package formatters

import "github.com/zored/edit/src/service/tokens"

type Options struct {
	Indent            int
	TrailingSeparator tokens.Separator
	Rule              Rule
}
