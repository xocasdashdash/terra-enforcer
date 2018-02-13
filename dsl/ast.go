package dsl

import "fmt"

//Node Basic node interface, besides the Pos() function you can print them
type Node interface {
	Pos() Position
}

type Position struct {
	Line int
	Char int
}

func (p Position) Pos() Position {
	return p
}
func (p Position) String() string {
	return fmt.Sprintf("Line: %d, Char: %d", p.Line, p.Char)
}

type AST Node

type IDNode struct {
	Position
	ID string
}
type IDStatement struct {
	Position
	IDs []IDNode
}
type WithNode struct {
	Position
	condition string
}
type ProgramNode struct {
	Position

	ResourceStatements []ResourceNode
}
type ValueNode struct {
	Position
	Value string
}

type AttributeNode struct {
	Position
	ID    IDNode
	With  WithNode
	Value []ValueNode
}

type ResourceNode struct {
	Position
	ID         IDNode
	Attributes []AttributeNode
}
