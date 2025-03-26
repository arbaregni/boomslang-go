package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"flag"
	"log"
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
	debug DebugTarget
	ostr  io.Writer
	estr  io.Writer
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
		fmt.Printf("usage: boomslang <filename>\n")
		os.Exit(EXIT_BAD_OPTS)
	}

	opts.filePath = positional[0]
	return opts
}
	
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

  opts := parse_opts()
	if opts.debug > 0 {
		fmt.Printf("debug mode, good choice...\n")
	}
  execute(opts)
}

func execute(opts *Opts) {
	// Open the file in read-only mode
	filePath := opts.filePath
	if !strings.HasSuffix(filePath, ".bs") {
		fmt.Printf("Bad file extension, '%s' does not look like a boomslang file.\n", filePath)
		os.Exit(EXIT_BAD_FILE)
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Printf("Error opening file '%s': %s\n",filePath, err)
		os.Exit(EXIT_BAD_FILE)
	}
	defer file.Close()

	lexer := new(Lexer)
	lexer.debug = opts.debug & DBG_LEX > 0
	lexer.buf = bufio.NewReader(file)
	tokens,err := lexer.Lex()
	if err != nil {
		fmt.Printf("I am very sorry, but I could not understand this file due to: %v\n", err)
		os.Exit(EXIT_LEX_FAILURE)
	}	

	if lexer.debug {
	  log.Printf("%v\n", tokens)
  }

	parser := new(Parser)
	parser.debug = opts.debug & DBG_PARSE > 0
	parser.tokens = tokens

	ast, err := parser.Parse()
	if err != nil {
		fmt.Printf("I am sorry, but I simply could not understand the file you gave me: %v\n", err)
		os.Exit(EXIT_PARSE_FAILURE)
	}

	if parser.debug {
		for _, a := range ast {
			log.Printf("%#v\n", a)
		}
	}

	ostr := opts.ostr
	if ostr == nil {
		ostr = os.Stdout
	}

	estr := opts.estr
	if estr == nil {
		estr = os.Stderr
	}

	// evaluate the program
	env := MakeEnv(ostr, estr)
	env.debug = opts.debug & DBG_EVAL > 0
	LoadBuiltins(env)

	if env.debug {
		fmt.Fprintf(ostr,"============================\n")
	}

	val := EvalAll(env, ast)

	if env.debug {
		fmt.Fprintf(ostr, "============================\n")
	}

	if val.IsErr() {
		fmt.Fprintf(estr, "Failure occured during runtime:\n%v\n", val.PrettyPrint())
		os.Exit(EXIT_RUNTIME_FAILURE)
	} else {
		fmt.Fprintf(ostr, "=> %s\n", val.PrettyPrint())
	}

}

