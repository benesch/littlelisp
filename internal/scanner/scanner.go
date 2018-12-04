package scanner

import (
	"errors"
)

type Token int

const (
	Lparen Token = iota
	Rparen
	Symbol
	Number
	String
	EOF
)

type Scanner struct {
	src []byte
	pos int
}

func NewScanner(src []byte) *Scanner {
	return &Scanner{
		src: src,
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n'
}

func isDelim(b byte) bool {
	return isWhitespace(b) || b == '(' || b == ')'
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9' || b == '.'
}

func (s *Scanner) Scan() (Token, string, error) {
	for s.pos < len(s.src) && isWhitespace(s.src[s.pos]) {
		s.pos++
	}

	if s.pos >= len(s.src) {
		return EOF, "", nil
	}

	switch s.src[s.pos] {
	case '(':
		s.pos++
		return Lparen, "(", nil
	case ')':
		s.pos++
		return Rparen, ")", nil
	case '"':
		start := s.pos
		s.pos++
		for s.pos < len(s.src) && (s.src[s.pos] != '"') {
			s.pos++
		}
		if s.pos >= len(s.src) {
			return EOF, "", errors.New("unterminated string literal")
		}
		s.pos++
		return String, string(s.src[start:s.pos]), nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		start := s.pos
		for s.pos < len(s.src) && isNum(s.src[s.pos]) {
			s.pos++
		}
		return Number, string(s.src[start:s.pos]), nil
	default:
		start := s.pos
		for s.pos < len(s.src) && !isDelim(s.src[s.pos]) {
			s.pos++
		}
		return Symbol, string(s.src[start:s.pos]), nil
	}
}
