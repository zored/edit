package scan

type linesScanner struct {
	lines   []string
	index   int
	reverse bool
}

func NewLinesScanner(lines []string, reverse bool) *linesScanner {
	index := 0
	if reverse {
		index = len(lines) - 1
	}

	return &linesScanner{
		lines:   lines,
		reverse: reverse,
		index:   index,
	}
}

func (l *linesScanner) Scan() bool {
	if l.reverse {
		return l.index >= 0
	}
	return l.index < len(l.lines)
}

func (l *linesScanner) Text() string {
	i := l.index
	delta := 1
	if l.reverse {
		delta = -1
	}
	l.index += delta
	return l.lines[i]
}

func (l *linesScanner) Err() error {
	return nil
}
