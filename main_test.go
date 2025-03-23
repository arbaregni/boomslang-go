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
}


func TestExamples(t *testing.T) {

	for _, tc := range cases {
	  filePath := tc.filePath
	  expected := tc.expected
		
		buf := new(strings.Builder)
		opts := new(Opts)
		opts.ostr = buf
		opts.estr = buf

		execute(filePath, opts)

		actual := buf.String()

		if expected != actual {
			t.Errorf("ran '%s' expected '%s', got '%s'", filePath, expected, actual)
		}

	}
}
