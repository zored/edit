package navigation

import (
	"fmt"
)

const BadNumber = 0

type Position struct{ Column, Line int }

func NewPosition(line int, column int) *Position { return &Position{Column: column, Line: line} }
func NewZeroPosition() *Position                 { return &Position{Column: BadNumber, Line: BadNumber} }
func NewStartPosition() *Position                { return &Position{Column: 1, Line: 1} }

func (p *Position) String() string    { return fmt.Sprintf("%d:%d", p.Line, p.Column) }
func (p *Position) Wrong() bool      { return p.Column == BadNumber }
func (p *Position) LineIndex() int   { return p.Line - 1 }
func (p *Position) ColumnIndex() int { return p.Column - 1 }
