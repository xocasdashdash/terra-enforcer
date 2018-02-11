package dsl

// Token defines a single TFEN token which can be obtained via the Scanner
type Token struct {
	Type Type
	Pos  Pos
	Text string
	JSON bool
}
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	identifier_beg
	IDENT // literals
	RESOURCE // resource ""
	ATTRIBUTE // attribute ""
	HAS // has 
	WITH // with
	OF // OF []

	literal_beg
	NUMBER  // 12345
	FLOAT   // 123.45
	BOOL    // true,false
	STRING  // "abc"
	HEREDOC // <<FOO\nbar\nFOO
	literal_end
	identifier_end

	operator_beg
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RBRACK // ]
	RBRACE // }

	ASSIGN // =
	ADD    // +
	SUB    // -
	REGEX  // ~=

	operator_end
)
