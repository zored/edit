package replacers

import (
	"github.com/zored/edit/src/service/navigation"
	"io"
)

type IReplacer interface {
	Replace(input io.Reader, interval *navigation.Interval, replacement string) (string, error)
}
