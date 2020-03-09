package tokens

import (
	"github.com/zored/edit/src/service/runes"
	"strings"
)

type TokenBuffer struct {
	tokens       Tokens
	buffer       []rune
	wrapperStart *runes.Matcher
	wrapperEnd   *runes.Matcher
	separator    *runes.Matcher
	done         bool
}
type Tokens []*Token

func NewTokenBuffer(wrapper *Wrappers, separator Separator) *TokenBuffer {
	return newTokenBuffer(wrapper, separator, false)
}

func newTokenBuffer(wrapper *Wrappers, separator Separator, reverse bool) *TokenBuffer {
	if reverse {
		panic("no tokens reverse support now")
	}
	return &TokenBuffer{
		wrapperStart: runes.NewMatcher(wrapper.Start, reverse),
		wrapperEnd:   runes.NewMatcher(wrapper.End, reverse),
		separator:    runes.NewMatcher(string(separator), reverse),
		buffer:       []rune{},
		tokens:       make(Tokens, 0),
	}
}

func (t *TokenBuffer) Write(rune_ rune) {
	var token_ *Token = nil
	if t.wrapperStart.Add(rune_) {
		t.appendBufferToken(AtomName)
		token_ = NewToken(t.wrapperStart.Runes, WrapperStart)
	}
	if t.wrapperEnd.Add(rune_) {
		t.appendBufferToken(Atom)
		token_ = NewToken(t.wrapperEnd.Runes, WrapperEnd)
	}
	if t.separator.Add(rune_) {
		t.appendBufferToken(Atom)
		token_ = NewToken(t.separator.Runes, AtomSeparator)
	}
	if token_ == nil {
		t.buffer = append(t.buffer, rune_)
	} else {
		t.tokens = append(t.tokens, token_)
	}
}

func (t *TokenBuffer) Complete() Tokens {
	return t.tokens
}

func (t *TokenBuffer) appendBufferToken(tokenType TokenType) {
	if len(t.buffer) == 0 {
		return
	}
	trimmed := []rune(strings.TrimSpace(string(t.buffer)))
	if len(trimmed) != 0 {
		t.tokens = append(t.tokens, NewToken(trimmed, tokenType))
	}
	t.buffer = []rune{}
}
