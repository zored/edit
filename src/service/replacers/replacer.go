package replacers

import (
	"bufio"
	"github.com/zored/edit/src/service/navigation"
	"io"
)

type (
	replacer  struct{}
	IReplacer interface {
		Replace(input io.Reader, interval *navigation.Interval, replacement string) (string, error)
	}
)

func NewReplacer() IReplacer {
	return &replacer{}
}

func (r *replacer) Replace(input io.Reader, interval *navigation.Interval, replacement string) (string, error) {
	lines := bufio.NewScanner(input)
	// TODO: use more efficient way:
	position := navigation.NewPosition(1, 1)
	start, stop := interval.Start, interval.Stop
	result := ""
	for lines.Scan() {
		lineText := lines.Text()

		if position.Line < start.Line || position.Line > stop.Line {
			result += lineText + "\n"
		} else {
			if position.Line == start.Line {
				result += lineText[:start.Column-1] + replacement
			}
			if position.Line == stop.Line {
				result += lineText[stop.Column-1:] + "\n"
			}
		}
		position.Line++
	}

	if err := lines.Err(); err != nil {
		return "", err
	}

	return result, nil
}
