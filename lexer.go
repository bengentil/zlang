// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	LeftBlockDelim    = '{'
	RightBlockDelim   = '}'
	LeftParentDelim   = '('
	RightParentDelim  = ')'
	ParamDelim        = ','
	StringDelim       = '"'
	EscapeCar         = '\\'
	Comment           = "//"
	LeftBlockComment  = "/*"
	RightBlockComment = "*/"
)

// For error reporting
type Position struct {
	name string
	line int
	char int
}

type LexItem struct {
	Token Token
	Val   string
}

func (i LexItem) String() string {
	switch {
	case i.Token == TOK_EOF:
		return "EOF"
	case i.Token == TOK_ERROR:
		return i.Val
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

const eof = -1

// lexer holds the state of the scanner.
type Lexer struct {
	name  string // the position to be reported in case of error
	input string // the string being scanned.
	pos   int    // current position in the input.
	start int    // start position of this item.
	width int    // width of last rune read from input.
}

// Initialize the lexer
func NewLexer(name, input string) *Lexer {
	return &Lexer{
		name:  name,
		input: input,
		pos:   0,
		start: 0,
		width: 1,
	}
}

// get the line number for error reporting
func (l *Lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// Report an error
func (l *Lexer) emit_error(format string, args ...interface{}) LexItem {
	return LexItem{TOK_ERROR, fmt.Sprintf(format, args...)}
}

func (l *Lexer) emit(t Token) LexItem {
	value := l.getVal()
	l.start = l.pos
	return LexItem{t, value}
}

func (l *Lexer) getVal() string {
	return l.input[l.start:l.pos]
}

// get next rune of the input
func (l *Lexer) getRune() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// seek next rune but doesn't consume it
func (l *Lexer) peekRune() rune {
	r := l.getRune()
	l.pos -= l.width
	return r
}

// ignore the current string, used for space and tabs
func (l *Lexer) ignore() {
	l.start = l.pos
}

// Parser will call NextItem until item's token is TOK_EOF or TOK_ERROR
func (l *Lexer) NextItem() LexItem {
	lastRune := ' '

	for isSpace(lastRune) {
		l.ignore() // don't emit space token
		lastRune = l.getRune()
	}

	// identifier starts with a letter
	if isLetter(lastRune) || lastRune == POINTER_CHAR {
		for isAlphaNumeric(l.peekRune()) {
			lastRune = l.getRune()
		}
		return l.emit(resolveIdentifier(l.getVal()))
	}

	if isDigit(lastRune) {
		for isDigit(l.peekRune()) {
			lastRune = l.getRune()
		}
		return l.emit(TOK_INT)
	}

	if isEndOfLine(lastRune) {
		return l.emit(TOK_ENDL)
	}

	// parse a string
	if lastRune == StringDelim {
		l.ignore() // skip opening "
		for {
			nextRune := l.peekRune()
			if nextRune == StringDelim && lastRune != EscapeCar {
				s := l.emit(TOK_STRING)
				lastRune = l.getRune() // get the closing "
				return s
			}
			if nextRune == eof {
				return l.emit_error("unclosed string %s", l.getVal())
			}
			lastRune = l.getRune()
		}
	}

	// parse a multiline comment
	if lastRune == '/' && l.peekRune() == '*' {
		lastRune = l.getRune() // position on *
		for {
			lastRune = l.getRune()
			nextRune := l.peekRune()
			if lastRune == '*' && nextRune == '/' {
				lastRune = l.getRune() // get the closing /
				return l.emit(TOK_COMMENT)
			}
			if nextRune == eof {
				return l.emit_error("unclosed comment %s", l.getVal())
			}
		}
	}

	// parse a single line comment
	if lastRune == '/' && l.peekRune() == '/' {
		for !isEndOfLine(l.peekRune()) && lastRune != eof {
			lastRune = l.getRune()
		}
		return l.emit(TOK_COMMENT)
	}

	// TODO: replace with resolveOperator && resolveDelimiter
	switch lastRune {
	case LeftParentDelim:
		return l.emit(TOK_LPAREN)
	case RightParentDelim:
		return l.emit(TOK_RPAREN)
	case LeftBlockDelim:
		return l.emit(TOK_LBLOCK)
	case RightBlockDelim:
		return l.emit(TOK_RBLOCK)
	case ',':
		return l.emit(TOK_COMMA)
	case '*':
		return l.emit(TOK_MUL)
	case '+':
		return l.emit(TOK_PLUS)
	case '-':
		return l.emit(TOK_MINUS)
	case '/':
		return l.emit(TOK_DIV)
	case eof:
		return l.emit(TOK_EOF)
	}

	return l.emit_error("unknown token with value '%s' last rune = 0x%x", l.getVal(), lastRune)
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

// isDigit reports whether r is a digit
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
