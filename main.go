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

type Opts struct {
	debug bool
	ostr  io.Writer
	estr  io.Writer
}

func parse_opts() (string, *Opts) {
	opts := new(Opts)
	flag.Parse()

	positional := make([]string, 0, len(flag.Args()))

	for i := 1; i < len(flag.Args()); i += 1 {
		arg := flag.Args()[i]
		switch arg {
			case "--debug":
				opts.debug = true
			default:
				positional = append(positional, arg)
		}
	}

	if len(positional) != 1 {
		log.Fatal("usage: boomslang <filename>")
	}

	filePath := positional[0]
	return filePath, opts
}
	
func main() {
	filePath, opts := parse_opts()
	if opts.debug {
		fmt.Printf("debug mode, good choice...\n")
	}
  execute(filePath, opts)
}

func execute(filePath string, opts *Opts) {
	// Open the file in read-only mode
	if !strings.HasSuffix(filePath, ".bs") {
		fmt.Printf("Bad file extension, '%s' does not look like a boomslang file.\n", filePath)
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Printf("Error opening file '%s': %s\n",filePath, err)
		return
	}
	defer file.Close()

	lexer := new(Lexer)
	lexer.debug = opts.debug
	lexer.buf = bufio.NewReader(file)
	tokens,err := lexer.Lex()
	if err != nil {
		log.Fatal("I am very sorry, but I could not understand this file due to: ", err)
		return
	}	

	if opts.debug {
	  log.Printf("%v\n", tokens)
  }

	parser := new(Parser)
	parser.debug = opts.debug
	parser.tokens = tokens

	ast, err := parser.Parse()
	if err != nil {
		log.Fatalf("I am sorry, but I simply could not understand the file you gave me: %v\n", err);
	}

	if opts.debug {
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
	LoadBuiltins(env)

	var val BsValue = &BsNilVal{}
  

	if opts.debug {
		fmt.Fprintf(ostr,"============================\n")
	}

	for _, node := range ast {
		val = node.Eval(env)
		if val.IsErr() {
			fmt.Fprintf(estr, "%s\n", val.PrettyPrint())
			break
		}
	}

	if opts.debug {
		fmt.Fprintf(ostr, "============================\n")
		if !val.IsErr() {
		  fmt.Fprintf(ostr, "=> %s\n", val.PrettyPrint())
	  }
	} 
}

