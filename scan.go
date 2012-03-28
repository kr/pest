package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type position struct {
	file string
	line int
	col  int
}

type scannerError struct {
	off int
	msg string
}

type scanner struct {
	src     []byte
	ch      rune
	off     int
	offNext int
	err     []scannerError
}

func (s *scanner) init(src []byte) {
	s.src = src
	s.next()
}

func (s *scanner) scan() (off int, tok token, lit string) {
	s.skipWhitespace()

	off = s.off
	switch ch := s.ch; {
	case isLetter(s.ch):
		tok = IDENT
		lit = s.scanIdent()
	case isDigit(s.ch):
		tok = INT
		lit = s.scanInt()
	default:
		s.next()
		switch ch {
		case -1:
			tok = EOF
		case '\n':
			tok = NL
		case ':':
			tok = TSTART
			lit = s.scanTPart()
		case '}':
			tok = TCONT
			lit = s.scanTPart()
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN
		case '"':
			tok = STRING
			lit = s.scanString()
		case '.':
			tok = DOT
		case '=':
			if s.ch == '=' {
				s.next()
				tok = EQ
			} else {
				tok = ASSIGN
			}
		case '~':
			tok = TEST
		case '+':
			tok = ADD
		case '-':
			tok = SUB
		case '*':
			tok = MUL
		case '/':
			tok = QUO
		case '%':
			tok = REM
		case '!':
			if s.ch == '=' {
				s.next()
				tok = NE
			} else {
				tok = NOT
			}
		case '&':
			if s.ch == '&' {
				s.next()
				tok = AND
				break
			}
			s.error(off, "expected '&&'")
		case '|':
			if s.ch == '|' {
				s.next()
				tok = OR
				break
			}
			s.error(off, "expected '||'")
		case '<':
			if s.ch == '=' {
				s.next()
				tok = LE
			} else {
				tok = LT
			}
		case '>':
			if s.ch == '=' {
				s.next()
				tok = GE
			} else {
				tok = GT
			}
		default:
			tok = INVALID
			s.errorf(off, "invalid character %#v", string(ch))
		}
	}
	return
}

func (s *scanner) scanTPart() string {
	// opening ':' or '}' already consumed
	off := s.off - 1

	for s.ch != '\n' && s.ch != '{' {
		s.next()
	}

	if s.ch == '{' {
		s.next()
		return string(s.src[off:s.off])
	}

	return string(s.src[off:s.off]) + "\n"
}

func (s *scanner) scanIdent() string {
	off := s.off
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[off:s.off])
}

func (s *scanner) scanInt() string {
	off := s.off
	if s.ch == '0' {
		s.next()
		if s.ch == 'x' || s.ch == 'X' {
			s.next()
			s.scanIntDigits(16)
		} else {
			s.scanIntDigits(8)
		}
	} else {
		s.scanIntDigits(10)
	}
	return string(s.src[off:s.off])
}

func (s *scanner) scanIntDigits(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func (s *scanner) scanString() string {
	// opening '"' already consumed
	off := s.off - 1

	for s.ch != '"' {
		ch := s.ch
		s.next()
		if ch == '\n' || ch < 0 {
			s.error(off, "string not terminated")
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}
	s.next()
	return string(s.src[off:s.off])
}

func (s *scanner) scanEscape(quote rune) {
	off := s.off

	var i, base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return
	case '0', '1', '2', '3', '4', '5', '6', '7':
		i, base, max = 3, 8, 255
	case 'x':
		s.next()
		i, base, max = 2, 16, 255
	case 'u':
		s.next()
		i, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		i, base, max = 8, 16, unicode.MaxRune
	default:
		s.next() // always make progress
		s.error(off, "unknown escape sequence")
		return
	}

	var x uint32
	for ; i > 0 && s.ch != quote && s.ch >= 0; i-- {
		d := uint32(digitVal(s.ch))
		if d >= base {
			s.error(s.off, "illegal character in escape sequence")
			break
		}
		x = x*base + d
		s.next()
	}
	// in case of an error, consume remaining chars
	for ; i > 0 && s.ch != quote && s.ch >= 0; i-- {
		s.next()
	}
	if x > max || 0xd800 <= x && x < 0xe000 {
		s.error(off, "escape sequence is invalid Unicode code point")
	}
}

func (s *scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\r' {
		s.next()
	}
}

// s.ch < 0 means eof.
func (s *scanner) next() {
	if s.offNext < len(s.src) {
		s.off = s.offNext
		r, n := utf8.DecodeRune(s.src[s.off:])
		s.offNext += n
		s.ch = r
	} else {
		s.off = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *scanner) error(off int, msg string) {
	s.err = append(s.err, scannerError{off, msg})
}

func (s *scanner) errorf(off int, f string, a ...interface{}) {
	s.error(off, fmt.Sprintf(f, a...))
}

func isLetter(c rune) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func digitVal(c rune) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}
