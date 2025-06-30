package service

import (
	"strconv"
	"strings"
	"testing"
)

func TestAskUserYN(t *testing.T) {
	input := []struct{
		name string
		input string
		excpected bool
	} {
		{name: "yes lower", input: "y\n", excpected: true},
		{name: "yes upper", input: "Y\n", excpected: true},
		{name: "no lower", input: "n\n", excpected: false},
		{name: "no upper", input: "N\n", excpected: false},
		{name: "wrong first try y lower", input: "test\ny\n", excpected: true},
		{name: "wrong first try y upper", input: "test\nY\n", excpected: true},
		{name: "wrong first try n lower", input: "test\nn\n", excpected: false},
		{name: "wrong first try n upper", input: "test\nN\n", excpected: false},
	}

	for i := range input {
		output := AskUserYN(strings.NewReader(input[i].input))
		if output != input[i].excpected {
			fatal := "excpected: " + strconv.FormatBool(input[i].excpected) + " got: " + strconv.FormatBool(output)
			t.Fatal(fatal)
		}
	}
}