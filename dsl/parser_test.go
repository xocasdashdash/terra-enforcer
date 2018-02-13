package dsl_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/xocasdashdash/terra-enforcer/dsl"
)

func TestParser(t *testing.T) {
	basicString, _ := ioutil.ReadFile("fixtures/test_parser_01.tfen")
	ast, err := dsl.Parse(string(basicString))
	if err != nil {
		t.Errorf("error should be nil %#v", err)
		t.FailNow()
	}
	pn, ok := ast.(dsl.ProgramNode)
	if !ok {
		t.Errorf("error casting the ast to a program")
	}
	for _, r := range pn.ResourceStatements {
		fmt.Printf("Resource: %s", r.ID)
	}
}
