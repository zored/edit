package tokenizers

import (
	"bufio"
	"fmt"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/runes"
	"github.com/zored/edit/src/service/scanners"
	"github.com/zored/edit/src/service/tokens"
	"io"
)

const badIndex = -1

type (
	tokenizer struct{}
)

func NewTokenizer() ITokenizer {
	return &tokenizer{}
}

func (p *tokenizer) Tokenize(
	reader io.Reader,
	cursor *navigation.Position,
	wrapper *tokens.Wrappers,
	separator_ tokens.Separator,
) (mol *tokens.Molecule, err error) {
	file := bufio.NewScanner(reader)
	linesTillCursor, err := p.getLinesTillCursor(cursor, file)
	if err != nil {
		return mol, err
	}

	mol = tokens.NewEmptyMolecule()
	mol.Interval.Start, err = p.getWrapperStart(linesTillCursor, cursor, wrapper.Start)
	if err != nil {
		return
	}

	mol.Interval.Stop, mol.Tokens, err = p.getWrapperEnd(linesTillCursor, mol.Interval.Start, wrapper, separator_, file)
	return
}

func (p *tokenizer) getWrapperStart(linesTillCursor []string, cursor *navigation.Position, start string) (*navigation.Position, error) {
	initialLineNumber := len(linesTillCursor)
	result := navigation.NewZeroPosition()
	result.Line = initialLineNumber

	startMatcher := runes.NewMatcher(start, true)
	err := scanners.ScanAll(
		scanners.NewReverseLinesScanner(linesTillCursor),
		func(line string) error {
			lineRunes := []rune(line)

			// Where to cursor search of rune:
			columnInitial := len(lineRunes)
			if result.Line == initialLineNumber {
				columnInitial = cursor.Column
			}

			// Go back for a set of expected lineRunes:
			result.Column = columnInitial
			for ; result.Column >= 1; result.Column-- {
				_rune := lineRunes[result.Column-1]
				if startMatcher.Add(_rune) {
					break
				}
			}

			if !result.Wrong() {
				return scanners.BreakErr
			}

			result.Line--
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !result.Wrong() {
		return result, nil
	}

	return nil, fmt.Errorf(
		"no atom wrapper cursor (%s) all cursor cursor %s",
		start,
		cursor,
	)
}

func (p *tokenizer) getWrapperEnd(
	linesTillCursor []string,
	startPosition *navigation.Position,
	wrapper *tokens.Wrappers,
	separator_ tokens.Separator,
	fileScanner *bufio.Scanner,
) (*navigation.Position, tokens.Tokens, error) {
	result := navigation.NewPosition(startPosition.Line, startPosition.Column)
	resultTokensBuffer := tokens.NewTokenBuffer(wrapper, separator_)

	lines := linesTillCursor[startPosition.LineIndex():]
	lines[0] = lines[0][startPosition.ColumnIndex():]

	startMatcher := runes.NewMatcher(wrapper.Start, false)
	endMatcher := runes.NewMatcher(wrapper.End, false)
	depth := 0
	done := false

	err := scanners.ScanAll(
		scanners.NewMultiScanner(
			scanners.NewLinesScanner(lines),
			fileScanner,
		),
		func(line string) error {
			var runeIndex int
			depth, runeIndex, done = p.parseLine(line, depth, startMatcher, endMatcher, resultTokensBuffer)
			if done {
				result.Column = runeIndex + 2
				if startPosition.Line == result.Line {
					result.Column += startPosition.Column - 1
				}
				return scanners.BreakErr
			}
			result.Line++
			return nil
		},
	)

	if err != nil {
		return nil, nil, err
	}

	if done {
		return result, resultTokensBuffer.Complete(), nil
	}

	return nil, nil, fmt.Errorf(
		"no ending wrapper '%s' for starting '%s' all, (started %d time(s), ended %d time(s)",
		wrapper.Start,
		wrapper.End,
		startMatcher.Matches(),
		endMatcher.Matches(),
	)
}

func (p *tokenizer) parseLine(line string, depth int, start *runes.Matcher, end *runes.Matcher, tokensBuffer *tokens.TokenBuffer) (int, int, bool) {
	lineRunes := []rune(line)
	for runeIndex, rune_ := range lineRunes {
		if start.Add(rune_) {
			depth++
		}
		if end.Add(rune_) {
			depth--
		}
		tokensBuffer.Write(rune_)
		if depth == 0 && start.Matches() > 0 {
			return 0, runeIndex, true
		}
	}
	return depth, badIndex, false
}

func (p *tokenizer) getLinesTillCursor(cursor *navigation.Position, file *bufio.Scanner) ([]string, error) {
	lines := make([]string, 0)

	// Get all lines start cursor lineNumber:
	lineNumber := 1
	err := scanners.ScanAll(file, func(line string) error {
		lineNumber++
		lines = append(lines, line)
		if lineNumber > cursor.Line {
			return scanners.BreakErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if lineNumber < cursor.Line {
		return nil, fmt.Errorf("only %d lines (%d expected)", lineNumber, cursor.Line)
	}
	return lines, nil
}
