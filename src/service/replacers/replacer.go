package replacers

import (
	"bufio"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/scanners"
	"io"
)

type replacer struct{}

func NewReplacer() IReplacer {
	return &replacer{}
}

func (r *replacer) Replace(input io.Reader, interval *navigation.Interval, replacement string) (string, error) {
	position := navigation.NewStartPosition()
	result := ""
	err := scanners.ScanAll(bufio.NewScanner(input), func(lineText string) error {
		defer func() { position.Line++ }()

		if position.Line < interval.Start.Line || position.Line > interval.Stop.Line {
			result += lineText + "\n"
			return nil
		}

		if position.Line == interval.Start.Line {
			result += lineText[:interval.Start.Column-1] + replacement
		}

		if position.Line == interval.Stop.Line {
			result += lineText[interval.Stop.Column-1:] + "\n"
		}

		return nil
	})
	return result, err
}
