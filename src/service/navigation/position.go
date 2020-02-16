package navigation

import (
	"fmt"
)

type Position struct{ Column, Line int }

func NewPosition(line int, column int) *Position { return &Position{Column: column, Line: line} }

func (p Position) String() string    { return fmt.Sprintf("%d:%d", p.Line, p.Column) }
func (p *Position) Wrong() bool      { return p.Column == 0 }
func (p *Position) LineIndex() int   { return p.Line - 1 }
func (p *Position) ColumnIndex() int { return p.Column - 1 }
