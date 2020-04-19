package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParser(t *testing.T) {
	//var edef metadata.ExperimentDefinition
	parser := &Parser{}
	data, err := ioutil.ReadFile("../cue.mod/p2plab.cue")
	if err != nil {
		t.Fatal(err)
	}
	inst, err := parser.Compile("p2plab", data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", inst)
}
