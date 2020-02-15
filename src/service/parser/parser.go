package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	badIndex  = -1
	badColumn = 0
)

type (
	parser   struct{}
	position struct {
		column, line int
	}
	interval struct {
		start, stop *position
	}
	atomWrapper struct {
		start, end string
	}
	runesMatcher struct {
		runes   []rune
		index   int
		count   int
		reverse bool
	}
	iScanner interface {
		Scan() bool
		Text() string
		Err() error
	}
	linesScanner struct {
		lines   []string
		index   int
		reverse bool
	}
	multiScanner struct {
		scanners []iScanner
		index    int
	}
	tokens struct {
		all          []*token
		buffer       []rune
		wrapperStart *runesMatcher
		wrapperEnd   *runesMatcher
		separator    *runesMatcher
	}
	token struct {
		value     string
		tokenType tokenType
	}
	tokenType int
	separator string
)

func newAtomWrapper(start string, end string) *atomWrapper {
	return &atomWrapper{start: start, end: end}
}

func newParser() *parser {
	return &parser{}
}

func newToken(runes []rune, tokenType tokenType) *token {
	return &token{value: string(runes), tokenType: tokenType}
}

const (
	atomWrapStart tokenType = iota
	atomWrapEnd
	atomSeparator
	atom
	atomName
)

func (t *tokens) check(rune_ rune) {
	var token_ *token = nil
	if t.wrapperStart.check(rune_) {
		t.appendBufferToken(atomName)
		token_ = t.wrapperStart.token(atomWrapStart)
	}
	if t.wrapperEnd.check(rune_) {
		t.appendBufferToken(atom)
		token_ = t.wrapperEnd.token(atomWrapEnd)
	}
	if t.separator.check(rune_) {
		t.appendBufferToken(atom)
		token_ = t.separator.token(atomSeparator)
	}
	if token_ == nil {
		t.buffer = append(t.buffer, rune_)
	} else {
		t.all = append(t.all, token_)
	}
}

func (t *tokens) appendBufferToken(tokenType tokenType) {
	if len(t.buffer) == 0 {
		return
	}
	trimmed := []rune(strings.TrimSpace(string(t.buffer)))
	if len(trimmed) != 0 {
		t.all = append(t.all, newToken(trimmed, tokenType))
	}
	t.buffer = []rune{}
}

func newTokens(wrapper *atomWrapper, separator separator, reverse bool) *tokens {
	if reverse {
		panic("no tokens reverse support now")
	}
	return &tokens{
		wrapperStart: newRunesMatcher(wrapper.start, reverse),
		wrapperEnd:   newRunesMatcher(wrapper.end, reverse),
		separator:    newRunesMatcher(string(separator), reverse),
		buffer:       []rune{},
	}
}

func newLinesScanner(lines []string, reverse bool) *linesScanner {
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

func newMultiScanner(scanners ...iScanner) *multiScanner {
	return &multiScanner{scanners: scanners}
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

func (m *multiScanner) Scan() bool {
	for ; m.index < len(m.scanners); m.index++ {
		if m.getScanner().Scan() {
			return true
		}
	}
	return false
}

func (m *multiScanner) Text() string {
	return m.getScanner().Text()
}

func (m *multiScanner) getScanner() iScanner {
	return m.scanners[m.index]
}

func (m *multiScanner) Err() error {
	return m.getScanner().Err()
}

func newPosition(line int, column int) *position {
	return &position{column: column, line: line}
}

func newRunesMatcher(runes string, reverse bool) *runesMatcher {
	matcher := &runesMatcher{
		runes:   []rune(runes),
		count:   0,
		reverse: reverse,
	}
	matcher.reset()
	return matcher
}

func (m *runesMatcher) closed(depth int) bool {
	return depth == 0 && m.count > 0
}
func (m *runesMatcher) token(t tokenType) *token {
	return newToken(m.runes, t)
}
func (m *runesMatcher) check(r rune) bool {
	if m.currentRune() != r {
		m.reset()
		return false
	}
	m.step()
	if !m.found() {
		return false
	}
	m.count++
	m.reset()
	return true
}

func (m *runesMatcher) currentRune() rune {
	return m.runes[m.index]
}

func (m *runesMatcher) found() bool {
	if m.reverse {
		return m.index < 0
	}
	return m.index >= len(m.runes)
}

func (m *runesMatcher) step() {
	delta := 1
	if m.reverse {
		delta = -1
	}
	m.index += delta
}

func (m *runesMatcher) reset() {
	if m.reverse {
		m.index = len(m.runes) - 1
		return
	}
	m.index = 0
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}

func (p *position) wrong() bool {
	return p.column == badColumn
}

func (p *position) lineYIndex() int {
	return p.line - 1
}
func (p *position) columnXIndex() int {
	return p.column - 1
}

func (p *parser) Parse(reader io.Reader, cursor *position, wrapper *atomWrapper, separator_ separator) (interval_ *interval, tokens_ *tokens, err error) {
	file := bufio.NewScanner(reader)
	linesTillCursor, err := p.getLinesTillCursor(cursor, file)
	if err != nil {
		return nil, nil, err
	}

	interval_ = &interval{}
	interval_.start, err = p.getWrapperStart(linesTillCursor, cursor, wrapper)
	if err != nil {
		return
	}

	interval_.stop, tokens_, err = p.getWrapperEnd(linesTillCursor, interval_.start, wrapper, separator_, file)
	return
}

func (p *parser) getWrapperEnd(linesTillCursor []string, startPosition *position, wrapper *atomWrapper, separator_ separator, file *bufio.Scanner) (*position, *tokens, error) {
	result := newPosition(startPosition.line, badColumn)

	lines := linesTillCursor[startPosition.lineYIndex():]
	lines[0] = lines[0][startPosition.columnXIndex():]

	start := newRunesMatcher(wrapper.start, false)
	end := newRunesMatcher(wrapper.end, false)
	depth := 0
	done := false

	allLines := newMultiScanner(newLinesScanner(lines, false), file)
	tokens := newTokens(wrapper, separator_, false)
	for allLines.Scan() {
		line := allLines.Text()
		depth, result.column, done = p.parseLine(line, depth, start, end, tokens)
		if done {
			break
		}
		result.line++
	}
	if err := allLines.Err(); err != nil {
		return nil, nil, err
	}

	if result.wrong() {
		return nil, nil, fmt.Errorf(
			"no ending wrapper '%s' for starting '%s' all, (started %d time(s), ended %d time(s)",
			wrapper.start,
			wrapper.end,
			start.count,
			end.count,
		)
	}

	return result, tokens, nil
}

func (p *parser) parseLine(line string, depth int, start *runesMatcher, end *runesMatcher, tokens *tokens) (int, int, bool) {
	runes := []rune(line)
	for runeIndex, rune_ := range runes {
		if start.check(rune_) {
			depth++
		}
		if end.check(rune_) {
			depth--
		}
		tokens.check(rune_)
		if start.closed(depth) {
			return 0, runeIndex, true
		}
	}
	return depth, badIndex, false
}

func (p *parser) getWrapperStart(linesTillCursor []string, cursor *position, wrapper *atomWrapper) (*position, error) {
	initialLineNumber := len(linesTillCursor)
	result := newPosition(initialLineNumber, badIndex)

	reverseLinesTillCursor := newLinesScanner(linesTillCursor, true)

	start := newRunesMatcher(wrapper.start, true)
	for reverseLinesTillCursor.Scan() {
		runes := []rune(reverseLinesTillCursor.Text())

		// Where to cursor search of rune:
		columnInitial := len(runes)
		if result.line == initialLineNumber {
			columnInitial = cursor.column
		}

		// Go back for a set of expected runes:
		result.column = columnInitial
		for ; result.column >= 1; result.column-- {
			_rune := runes[result.column-1]
			if start.check(_rune) {
				break
			}
		}

		if !result.wrong() {
			return result, nil
		}
	}

	return nil, fmt.Errorf(
		"no atom wrapper cursor (%s) all cursor cursor %s",
		wrapper.start,
		cursor,
	)
}

func (p *parser) getLinesTillCursor(cursor *position, linesScanner *bufio.Scanner) ([]string, error) {
	lines := make([]string, 0)

	// Get all lines start cursor line:
	line := 0
	for ; line <= cursor.line && linesScanner.Scan(); line++ {
		lines = append(lines, linesScanner.Text())
	}
	if err := linesScanner.Err(); err != nil {
		return nil, err
	}
	if line < cursor.line {
		return nil, fmt.Errorf("only %d lines (%d expected)", line, cursor.line)
	}
	return lines, nil
}
