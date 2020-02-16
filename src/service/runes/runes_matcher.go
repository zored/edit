package runes

type Matcher struct {
	Runes   []rune
	index   int
	Count   int
	reverse bool
}

func NewRunesMatcher(runes string, reverse bool) *Matcher {
	matcher := &Matcher{
		Runes:   []rune(runes),
		Count:   0,
		reverse: reverse,
	}
	matcher.reset()
	return matcher
}

func (m *Matcher) Closed(depth int) bool {
	return depth == 0 && m.Count > 0
}

func (m *Matcher) Check(r rune) bool {
	if m.currentRune() != r {
		m.reset()
		return false
	}
	m.step()
	if !m.found() {
		return false
	}
	m.Count++
	m.reset()
	return true
}

func (m *Matcher) currentRune() rune {
	return m.Runes[m.index]
}

func (m *Matcher) found() bool {
	if m.reverse {
		return m.index < 0
	}
	return m.index >= len(m.Runes)
}

func (m *Matcher) step() {
	delta := 1
	if m.reverse {
		delta = -1
	}
	m.index += delta
}

func (m *Matcher) reset() {
	if m.reverse {
		m.index = len(m.Runes) - 1
		return
	}
	m.index = 0
}
