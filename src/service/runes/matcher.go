package runes

type Matcher struct {
	Runes   []rune
	matches int
	index   int
	reverse bool
}

func NewMatcher(runes string, reverse bool) *Matcher {
	matcher := &Matcher{
		Runes:   []rune(runes),
		matches: 0,
		reverse: reverse,
	}
	matcher.reset()
	return matcher
}

func (m *Matcher) Add(r rune) bool {
	if m.currentRune() != r {
		m.reset()
		return false
	}
	m.step()
	if !m.found() {
		return false
	}
	m.matches++
	m.reset()
	return true
}

func (m *Matcher) Matches() int {
	return m.matches
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
