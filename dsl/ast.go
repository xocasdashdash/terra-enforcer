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

type IDStatement struct {
	Position
	ID []string
}

type WithStatement struct {
	Position
	condition string
}
type ProgramNode struct {
	Position

	ResourceStatements []ResourceNode
}
type Value struct {
	Position
}

type ValueStatement struct {
	Position
	Value string
}

type AttributeNode struct {
	Position
	ID    IDStatement
	With  WithStatement
	Value []ValueStatement
}

type ResourceNode struct {
	Position
	ID             IDStatement
	BlockStatement []AttributeNode
}
