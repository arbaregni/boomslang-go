package main

import (
	"strings"
	"testing"
)

type testcase struct {
	filePath string
	expected string
}

var cases = []testcase{
	{"examples/hello.bs", "hello world!\n"},
	{"examples/assign.bs", "♤♡◇♧♧\n"},
	{"examples/sums.bs", "3\n"},
	{"examples/ifstmnt.bs", "its here!\n"},
	{"examples/nestedif.bs", "true and true is true!\nall done\n"},
	{"examples/elses.bs", "it works!\nbetween 10 and 20\n"},
	{"examples/simpleloop.bs", "0\n1\n2\n3\n4\n"},
	{"examples/loopwithbreak.bs", "done\n"},
	{"examples/loopwithbreak2.bs", "0\nand again\n1\nand again\n2\n"},
}

func TestExamples(t *testing.T) {

	for _, tc := range cases {
		expected := tc.expected

		buf := new(strings.Builder)
		opts := new(Opts)
		opts.ostr = buf
		opts.estr = buf
		opts.filePath = tc.filePath

		execute(opts)

		actual := buf.String()

		if expected != actual {
			t.Errorf("ran '%s' expected '%s', got '%s'", tc.filePath, expected, actual)
		}

	}
}
