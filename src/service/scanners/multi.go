package scanners

type multiScanner struct {
	scanners []IScanner
	index    int
}

func NewMultiScanner(scanners ...IScanner) *multiScanner {
	return &multiScanner{scanners: scanners}
}

func (m *multiScanner) Scan() bool {
	for ; m.index < len(m.scanners); m.index++ {
		if m.getScanner().Scan() {
			return true
		}
	}
	m.index--
	return false
}

func (m *multiScanner) Text() string {
	return m.getScanner().Text()
}

func (m *multiScanner) getScanner() IScanner {
	return m.scanners[m.index]
}

func (m *multiScanner) Err() error {
	return m.getScanner().Err()
}
