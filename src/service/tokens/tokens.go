package tokens

import (
	"github.com/zored/edit/src/service/runes"
	"strings"
)

type Tokens struct {
	All          []*Token
	buffer       []rune
	wrapperStart *runes.Matcher
	wrapperEnd   *runes.Matcher
	separator    *runes.Matcher
}

func NewTokens(wrapper *AtomWrapper, separator Separator, reverse bool) *Tokens {
	if reverse {
		panic("no tokens reverse support now")
	}
	return &Tokens{
		wrapperStart: runes.NewRunesMatcher(wrapper.Start, reverse),
		wrapperEnd:   runes.NewRunesMatcher(wrapper.End, reverse),
		separator:    runes.NewRunesMatcher(string(separator), reverse),
		buffer:       []rune{},
	}
}

func (t *Tokens) Check(rune_ rune) {
	var token_ *Token = nil
	if t.wrapperStart.Check(rune_) {
		t.appendBufferToken(AtomName)
		token_ = NewToken(t.wrapperStart.Runes, AtomWrapStart)
	}
	if t.wrapperEnd.Check(rune_) {
		t.appendBufferToken(Atom)
		token_ = NewToken(t.wrapperEnd.Runes, AtomWrapEnd)
	}
	if t.separator.Check(rune_) {
		t.appendBufferToken(Atom)
		token_ = NewToken(t.separator.Runes, AtomSeparator)
	}
	if token_ == nil {
		t.buffer = append(t.buffer, rune_)
	} else {
		t.All = append(t.All, token_)
	}
}

func (t *Tokens) appendBufferToken(tokenType TokenType) {
	if len(t.buffer) == 0 {
		return
	}
	trimmed := []rune(strings.TrimSpace(string(t.buffer)))
	if len(trimmed) != 0 {
		t.All = append(t.All, NewToken(trimmed, tokenType))
	}
	t.buffer = []rune{}
}
