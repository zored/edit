package parser

import (
	"bufio"
	"fmt"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/runes"
	"github.com/zored/edit/src/service/scan"
	"github.com/zored/edit/src/service/tokens"
	"io"
)

const (
	badIndex  = -1
	badColumn = 0
)

type (
	parser struct{}
)

func newParser() *parser {
	return &parser{}
}

func (p *parser) Parse(reader io.Reader, cursor *navigation.Position, wrapper *tokens.AtomWrapper, separator_ tokens.Separator) (interval_ *navigation.Interval, tokens_ *tokens.Tokens, err error) {
	file := bufio.NewScanner(reader)
	linesTillCursor, err := p.getLinesTillCursor(cursor, file)
	if err != nil {
		return nil, nil, err
	}

	interval_ = &navigation.Interval{}
	interval_.Start, err = p.getWrapperStart(linesTillCursor, cursor, wrapper.Start)
	if err != nil {
		return
	}

	interval_.Stop, tokens_, err = p.getWrapperEnd(linesTillCursor, interval_.Start, wrapper, separator_, file)
	return
}
func (p *parser) getWrapperStart(linesTillCursor []string, cursor *navigation.Position, needle string) (*navigation.Position, error) {
	initialLineNumber := len(linesTillCursor)
	result := navigation.NewPosition(initialLineNumber, badIndex)

	reverseLinesTillCursor := scan.NewLinesScanner(linesTillCursor, true)

	start := runes.NewRunesMatcher(needle, true)
	for reverseLinesTillCursor.Scan() {
		lineRunes := []rune(reverseLinesTillCursor.Text())

		// Where to cursor search of rune:
		columnInitial := len(lineRunes)
		if result.Line == initialLineNumber {
			columnInitial = cursor.Column
		}

		// Go back for a set of expected lineRunes:
		result.Column = columnInitial
		for ; result.Column >= 1; result.Column-- {
			_rune := lineRunes[result.Column-1]
			if start.Check(_rune) {
				break
			}
		}

		if !result.Wrong() {
			return result, nil
		}
	}

	return nil, fmt.Errorf(
		"no atom wrapper cursor (%s) all cursor cursor %s",
		needle,
		cursor,
	)
}

func (p *parser) getWrapperEnd(linesTillCursor []string, startPosition *navigation.Position, wrapper *tokens.AtomWrapper, separator_ tokens.Separator, file *bufio.Scanner) (*navigation.Position, *tokens.Tokens, error) {
	result := navigation.NewPosition(startPosition.Line, badColumn)

	lines := linesTillCursor[startPosition.LineIndex():]
	lines[0] = lines[0][startPosition.ColumnIndex():]

	start := runes.NewRunesMatcher(wrapper.Start, false)
	end := runes.NewRunesMatcher(wrapper.End, false)
	depth := 0
	done := false

	allLines := scan.NewMultiScanner(scan.NewLinesScanner(lines, false), file)
	tokens_ := tokens.NewTokens(wrapper, separator_, false)
	for allLines.Scan() {
		line := allLines.Text()
		depth, result.Column, done = p.parseLine(line, depth, start, end, tokens_)
		if done {
			break
		}
		result.Line++
	}
	if err := allLines.Err(); err != nil {
		return nil, nil, err
	}

	if result.Wrong() {
		return nil, nil, fmt.Errorf(
			"no ending wrapper '%s' for starting '%s' all, (started %d time(s), ended %d time(s)",
			wrapper.Start,
			wrapper.End,
			start.Count,
			end.Count,
		)
	}

	return result, tokens_, nil
}

func (p *parser) parseLine(line string, depth int, start *runes.Matcher, end *runes.Matcher, tokens *tokens.Tokens) (int, int, bool) {
	lineRunes := []rune(line)
	for runeIndex, rune_ := range lineRunes {
		if start.Check(rune_) {
			depth++
		}
		if end.Check(rune_) {
			depth--
		}
		tokens.Check(rune_)
		if start.Closed(depth) {
			return 0, runeIndex, true
		}
	}
	return depth, badIndex, false
}

func (p *parser) getLinesTillCursor(cursor *navigation.Position, linesScanner *bufio.Scanner) ([]string, error) {
	lines := make([]string, 0)

	// Get all lines start cursor line:
	line := 0
	for ; line <= cursor.Line && linesScanner.Scan(); line++ {
		lines = append(lines, linesScanner.Text())
	}
	if err := linesScanner.Err(); err != nil {
		return nil, err
	}
	if line < cursor.Line {
		return nil, fmt.Errorf("only %d lines (%d expected)", line, cursor.Line)
	}
	return lines, nil
}
