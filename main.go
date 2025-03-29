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

	if len(positional) > 1 {
		fmt.Printf("usage: boomslang <filename> [flags], got: %v\n", positional)
		os.Exit(EXIT_BAD_OPTS)
	}
	if len(positional) == 1 {
		opts.filePath = positional[0]
	}

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
		log.Printf("debug mode, good choice...\n")
	}
	if opts.filePath != "" {
		rc := execute(opts, opts.filePath)
		os.Exit(rc)
	}
	repl(opts)
}

func execute(opts *Opts, filePath string) int {
	// Open the file in read-only mode
	if !strings.HasSuffix(filePath, ".bs") {
		fmt.Fprintf(opts.estr, "Bad file extension, '%s' does not look like a boomslang file.\n", filePath)
		return (EXIT_BAD_FILE)
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Fprintf(opts.estr, "Error opening file '%s': %s\n", filePath, err)
		return (EXIT_BAD_FILE)
	}
	defer file.Close()
	buf := bufio.NewReader(file)

	// evaluate the program
	env := MakeEnv(opts)
	LoadBuiltins(env)

	rc, _ := run(opts, opts.filePath, buf, env)
	return rc
}

type replsource struct{}

func repl(opts *Opts) {
	env := MakeEnv(opts)
	LoadBuiltins(env)

	source := bufio.NewReader(os.Stdin)

	for {
		_, val := run(opts, "<repl>", source, env)
		if val != nil {
			fmt.Fprintf(opts.ostr, " => %s\n", val.PrettyPrint())
		}
	}
}

func run(opts *Opts, sourceName string, source *bufio.Reader, env *BsEnv) (int, BsValue) {
	lexer := MakeLexer(opts, sourceName, source)
	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Fprintf(opts.estr, "\033[0;31m I am very sorry, but I could not understand this file due to: %v\n\033[0m ", err)
		return EXIT_LEX_FAILURE, nil
	}
	parser := MakeParser(opts, tokens)

	ast, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(opts.estr, "\033[0;31m I am sorry, but I simply could not understand the file you gave me: %v\n\033[0m ", err)
		return EXIT_PARSE_FAILURE, nil
	}

	if opts.debug != 0 {
		fmt.Fprintf(opts.ostr, "============================ BEGIN EVAL ===========================\n")
	}

	val := EvalAll(env, ast)

	if val.ShouldUnwind() {
		fmt.Fprintf(opts.estr, "\033[0;31m Failure occured during runtime:\n%v\033[0m\n", val.PrettyPrint())
		return EXIT_RUNTIME_FAILURE, nil
	}

	return 0, val
}
