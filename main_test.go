package main

import (
	"testing"
	"strings"
)

type testcase struct {
	filePath string
	expected string
}

var cases = []testcase {
	{ "examples/hello.bs", "hello world!\n" },
	{ "examples/assign.bs", "♤♡◇♧♧\n" },
	{ "examples/sums.bs", "3\n" },
	{ "examples/ifstmnt.bs", "its here!\n" },
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
			t.Errorf("ran '%s' expected '%s', got '%s'", filePath, expected, actual)
		}

	}
}
