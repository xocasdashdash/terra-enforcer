package dsl

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func Parse(input string) (AST, error) {
	parser := newParser(Lex(input))
	return parser.Parse()
}

type parser struct {
	tokens <-chan Token

	lookahead [2]Token
	peekCount int
}

func newParser(tokens <-chan Token) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) Parse() (ast AST, err error) {
	// Parsing uses panics to bubble up errors
	defer p.recover(&err)

	ast = p.program()

	return
}

func (p *parser) nextToken() Token {
	return <-p.tokens
}

// next returns the next token.
func (p *parser) next() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.lookahead[0] = p.nextToken()
	}
	return p.lookahead[p.peekCount]
}

// backup backs the input stream up one token.
func (p *parser) backup() {
	p.peekCount++
}

// peek returns but does not consume the next token.
func (p *parser) peek() Token {
	if p.peekCount > 0 {
		return p.lookahead[p.peekCount-1]
	}
	p.peekCount = 1
	p.lookahead[1] = p.lookahead[0]
	p.lookahead[0] = p.nextToken()
	return p.lookahead[0]
}

// errorf formats the error and terminates processing.
func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("parser: %s", format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (p *parser) error(err error) {
	p.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (p *parser) expect(expected TokenType) Token {
	t := p.next()
	if t.Type != expected {
		debug.PrintStack()
		p.unexpected(t, expected)
	}
	return t
}

// unexpected complains about the token and terminates processing.
func (p *parser) unexpected(tok Token, expected ...TokenType) {
	expectedStrs := make([]string, len(expected))
	for i := range expected {
		expectedStrs[i] = fmt.Sprintf("%q", expected[i])
	}
	expectedStr := strings.Join(expectedStrs, ",")
	debug.PrintStack()
	p.errorf("unexpected token %q with value %q at line %d char %d, expected: %s", tok.Type, tok.Value, tok.Pos.Line, tok.Pos.Char, expectedStr)
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (p *parser) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*errp = e.(error)
	}
	return
}

var positionZero = Position{
	Line: 1,
	Char: 1,
}

func (p *parser) program() *ProgramNode {
	prog := &ProgramNode{
		Position: positionZero,
	}
	for {
		switch p.peek().Type {
		case TokenEOF:
			return prog
		case TokenResource:
			resource := p.resourceStatement()
			prog.ResourceStatements = append(prog.ResourceStatements, resource)
		default:
			p.unexpected(p.next(), TokenEOF, TokenResource)
		}
	}
}
func (p *parser) resourceStatement() ResourceNode {
	r := ResourceNode{}
	resourceToken := p.expect(TokenResource)
	r.Position = resourceToken.Pos
	idToken := p.expect(TokenString)
	r.ID = IDNode{
		ID:       idToken.Value,
		Position: idToken.Pos,
	}
	p.expect(TokenHas)
	for {
		pToken := p.peek()
		switch pToken.Type {
		case TokenAttribute:
			r.Attributes = append(r.Attributes, p.attributeStatement())
		case TokenComma:
			p.next()
		case TokenString:
			r.ID = *p.idStatement()
		case TokenLBrace:
			p.next()

		case TokenRBrace:
			p.next()
			return r
		default:
			p.unexpected(p.next(), TokenComma, TokenLBrace, TokenRBrace)
		}
	}
}
func (p *parser) idStatement() *IDNode {
	//Review
	t := p.expect(TokenString)
	return &IDNode{
		ID:       t.Value,
		Position: t.Pos,
	}
}
func (p *parser) valueStatement() []ValueNode {
	var res []ValueNode
	for {
		switch p.peek().Type {
		case TokenLBracket:
			p.next()
		case TokenRBracket:
			p.next()
			return res
		case TokenComma:
			p.next()
		case TokenString:
			tok := p.next()
			if tok.Value != "" {
				res = append(res, ValueNode{
					Position: tok.Pos,
					Value:    tok.Value,
				})
			}
		default:
			p.unexpected(p.next(), TokenLBracket, TokenString, TokenComma, TokenRBracket)
		}
	}
}
func (p *parser) attributeStatement() AttributeNode {
	p.expect(TokenAttribute)
	idNode := p.idStatement()
	p.expect(TokenWith)
	withCondition := p.expect(TokenString)
	p.expect(TokenOf)
	p.expect(TokenLBracket)
	withNode := WithNode{
		Condition: withCondition.Value,
		Position:  withCondition.Pos,
	}
	return AttributeNode{
		Position: idNode.Pos(),
		ID:       *idNode,
		With:     withNode,
		Value:    p.valueStatement(),
	}
}
