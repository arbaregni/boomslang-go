package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)
import "github.com/davecgh/go-spew/spew"

const (
	EXIT_BAD_OPTS int = 10 + iota
	EXIT_BAD_FILE
	EXIT_LEX_FAILURE
	EXIT_PARSE_FAILURE
	EXIT_RUNTIME_FAILURE
)

type DebugTarget int

const (
	DBG_PRE DebugTarget = 1 << iota
	DBG_LEX
	DBG_PARSE
	DBG_EVAL
)
const DBG_ALL DebugTarget = ^0

type Opts struct {
	debug    DebugTarget
	istr     io.Reader
	ostr     io.Writer
	estr     io.Writer
	filePath string
}

func parse_opts() *Opts {
	opts := new(Opts)
	flag.Parse()

	positional := make([]string, 0, len(flag.Args()))

	for _, arg := range flag.Args()[1:] {
		if !strings.HasPrefix(arg, "--") {
			positional = append(positional, arg)
			continue
		}

		// start: parsing the arg
		if arg == "--debug" {
			opts.debug = DBG_ALL
		} else if strings.HasPrefix(arg, "--debug=") {
			elems := strings.Split(strings.TrimPrefix(arg, "--debug="), ",")
			for _, elem := range elems {
				if elem == "lex" {
					opts.debug |= DBG_LEX
				} else if elem == "parse" {
					opts.debug |= DBG_PARSE
				} else if elem == "eval" {
					opts.debug |= DBG_EVAL
				} else {
					fmt.Printf("Bad choice for --debug, '%s' not supported\n", elem)
					os.Exit(EXIT_BAD_OPTS)
				}
			}
		} else {
			fmt.Printf("Bad flag: I do not recognize %s\n", arg)
			os.Exit(EXIT_BAD_OPTS)
		}

		//end: parsing the flag args

	}

	// use the postional args

	if len(positional) != 1 {
		fmt.Printf("usage: boomslang <filename> [flags], got: %v\n", positional)
		os.Exit(EXIT_BAD_OPTS)
	}

	opts.filePath = positional[0]

	// set good defaults for other args
	opts.istr = os.Stdin
	opts.ostr = os.Stdout
	opts.estr = os.Stderr

	return opts
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	opts := parse_opts()
	if opts.debug > 0 {
		fmt.Printf("debug mode, good choice...\n")
	}
	rc := execute(opts)
	os.Exit(rc)
}

func execute(opts *Opts) int {	

	// Open the file in read-only mode
	filePath := opts.filePath
	if !strings.HasSuffix(filePath, ".bs") {
		fmt.Fprintf(opts.estr,"Bad file extension, '%s' does not look like a boomslang file.\n", filePath)
		return (EXIT_BAD_FILE)
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Fprintf(opts.estr,"Error opening file '%s': %s\n", filePath, err)
		return (EXIT_BAD_FILE)
	}
	defer file.Close()

	lexer := new(Lexer)
	lexer.filePath = opts.filePath
	lexer.debug = opts.debug&DBG_LEX > 0
	lexer.buf = bufio.NewReader(file)
	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Fprintf(opts.estr,"I am very sorry, but I could not understand this file due to: %v\n", err)
		return (EXIT_LEX_FAILURE)
	}

	if lexer.debug {
		log.Printf("%v\n", tokens)
	}

	parser := new(Parser)
	parser.debug = opts.debug&DBG_PARSE > 0
	parser.tokens = tokens

	ast, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(opts.estr, "I am sorry, but I simply could not understand the file you gave me: %v\n", err)
		return (EXIT_PARSE_FAILURE)
	}

	if parser.debug {
		log.Printf("ast = %v\n", spew.Sdump(ast))
	}

	// evaluate the program
	env := MakeEnv(opts)
	env.debug = opts.debug&DBG_EVAL > 0
	LoadBuiltins(env)

	if opts.debug != 0 {
		fmt.Fprintf(opts.ostr, "============================\n")
	}

	val := EvalAll(env, ast)

	if val.IsErr() {
		fmt.Fprintf(opts.estr, "Failure occured during runtime:\n%v\n", val.PrettyPrint())
		return (EXIT_RUNTIME_FAILURE)
	}

	return 0
}
