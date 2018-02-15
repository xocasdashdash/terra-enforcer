package dsl

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token defines a single TFEN token which can be obtained via the Scanner
type Token struct {
	Pos   Position
	Type  TokenType
	Value string
}

//TokenType The type we're gonna be using instead of an int
type TokenType int

const (
	//Special tokens
	TokenError TokenType = iota
	TokenEOF             //

	TokenComment

	identifier_beg
	TokenResource  // resource ""
	TokenAttribute // attribute ""
	TokenHas       // has
	TokenWith      // with
	TokenOf        // Of

	literal_beg
	TokenNumber // 12345
	TokenFloat  // 123.45
	TokenBool   // true,false
	TokenString // "abc"
	//	TokenHeredoc // <<FOO\nbar\nFOO

	literal_end
	identifier_end

	operator_beg
	TokenQuote    // "
	TokenLBracket // [
	TokenLBrace   // {
	TokenComma    // ,
	TokenPeriod   // .
	TokenRBracket // ]
	TokenRBrace   // }
	TokenEqual    // =
	TokenRegex    // ~=

	operator_end
)

var keywords = map[string]TokenType{
	"with":      TokenWith,
	"has":       TokenHas,
	"attribute": TokenAttribute,
	"resource":  TokenResource,
	"of":        TokenOf,
}

func (tt TokenType) String() string {
	switch tt {
	case TokenError:
		return "error"
	case TokenQuote:
		return "\""
	case TokenEOF:
		return "eof"
	case TokenComment:
		return "comment"
	case TokenResource:
		return "resource"
	case TokenAttribute:
		return "attribute"
	case TokenHas:
		return "has"
	case TokenWith:
		return "with"
	case TokenOf:
		return "of"
	case TokenNumber:
		return "number"
	case TokenFloat:
		return "float"
	case TokenString:
		return "string"
	case TokenLBracket:
		return "["
	case TokenLBrace:
		return "{"
	case TokenComma:
		return "','"
	case TokenPeriod:
		return "."
	case TokenRBracket:
		return "']'"
	case TokenRBrace:
		return "}"
	case TokenEqual:
		return "=="
	case TokenRegex:
		return "~="
	default:
		return fmt.Sprintf("<token %d >", int(tt))
	}
}

type lexer struct {
	input string // the string being lexed

	pos   int // the current position of the input
	start int // the start of the current token
	width int // the width of the last read rune

	line int // the line number of the current token
	char int // the character number of the current token

	tokens chan Token // channel on which to emit tokens
}

func newLexer(input string) *lexer {
	return &lexer{
		input:  input,
		pos:    0,
		start:  0,
		width:  0,
		line:   1,
		char:   1,
		tokens: make(chan Token),
	}
}

//Lex Lexes an input and returns a channel where tokens are going to be emitted
func Lex(input string) <-chan Token {
	l := newLexer(input)
	go func() {
		defer close(l.tokens)
		for state := lexToken; state != nil; {
			state = state(l)
		}
	}()
	return l.tokens
}

type stateFn func(l *lexer) stateFn

const eof = -1

func (l *lexer) emit(t TokenType) {
	value := l.current()
	l.tokens <- Token{
		Pos:   l.position(),
		Type:  t,
		Value: value,
	}
	l.updatePosCounters()
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.updatePosCounters()
}

func (l *lexer) updatePosCounters() {
	value := l.current()
	// Update position counters
	l.start = l.pos

	// Count lines
	lastLine := 0
	for {
		i := strings.IndexRune(value[lastLine:], '\n')
		if i == -1 {
			break
		}
		lastLine += i + 1
		l.line++
		l.char = 1
	}
	l.char += len(value) - lastLine
}

func (l *lexer) position() Position {
	return Position{
		Line: l.line,
		Char: l.char,
	}
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

//Backup the lexer to the previous rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// error emits an error token with the err and returns the terminal state.
func (l *lexer) error(err error) stateFn {
	l.tokens <- Token{Pos: l.position(), Type: TokenError, Value: err.Error()}
	return nil
}

// errorf emits an error token with the formatted arguments and returns the terminal state.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{Pos: l.position(), Type: TokenError, Value: fmt.Sprintf(format, args...)}
	return nil
}

// ignore a contiguous block of spaces.
func (l *lexer) ignoreSpace() {
	for unicode.IsSpace(l.next()) {
		l.ignore()
	}
	l.backup()
}

/////////////////////////////
// Lex State Fns

// lexToken is the top level state
func lexToken(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsLetter(r):
			return lexWordOrKeyword
		case unicode.IsDigit(r):
			return lexNumberDigits
		case unicode.IsSpace(r):
			l.ignore()
		case r == '"':
			return lexWordOrKeyword
			//l.emit(TokenQuote)
		case r == '{':
			l.emit(TokenLBrace)
		case r == '}':
			l.emit(TokenRBrace)
		case r == ']':
			l.emit(TokenRBracket)
		case r == '[':
			l.emit(TokenLBracket)
		case r == '.':
			l.emit(TokenPeriod)
		case r == ',':
			l.emit(TokenComma)

		case r == eof:
			l.emit(TokenEOF)
			return nil
		default:
			return l.errorf("unexpected token %v", r)
		}
	}
}

func lexWordOrKeyword(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isValidIdent(r):
			// absorb
		default:
			l.backup()
			current := l.current()
			if typ, ok := keywords[current]; ok {
				l.emit(typ)
				return lexToken
			}
			l.emit(TokenString)
			return lexToken
		}
	}
}

func lexNumberDigits(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			//absorb
		default:
			l.backup()
			l.emit(TokenNumber)
			return lexToken
		}
	}
}

// isValidIdent reports whether r is either a letter or a digit
func isValidIdent(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' || r == '"' || r == '.'
}
