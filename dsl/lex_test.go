package dsl_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/xocasdashdash/terra-enforcer/dsl"
)

func TestLexer(t *testing.T) {
	basicString, _ := ioutil.ReadFile("fixtures/test_lexer_01.tfen")
	for t := range dsl.Lex(string(basicString)) {
		fmt.Printf("Token: %v\n", t)
	}
}
