package main

import (
	"errors"
	"strings"
	"fmt"
	"unicode"
	"bufio"
)


type Parser struct {
	buf *bufio.Reader
	debug bool
	indent int
}

func (p *Parser) ParseLine() (Ast, error) {
	if p.debug { fmt.Printf("in ParseLine\n") }

	line, err := p.buf.ReadString('\n')
	if err != nil { return nil, err }

	line = strings.TrimSuffix(line, "\n")
	words := strings.Fields(line)
	ast, err := p.parseStmnt(words)
	return ast, err
}

// infix parsing rules
type ruledef struct {
	opname string
	ctor func(Ast,Ast)Ast
}

var infixRules []ruledef = []ruledef {
	{ "is", func(lval,rval Ast)Ast { return AstAssign{lval,rval} }},
}

// magic infix functions
type builtindef struct {
	symbol string 
	opname string 
}

var infixBuiltins []builtindef = []builtindef {
	{ "_super-duper-secret__equals", "equals" },
	{ "_super-duper-secret__notequals", "notequals" },
	{ "_super-duper-secret__smallerthan", "smallerthan" },
	{ "_super-duper-secret__biggerthan", "biggerthan" },
	{ "_super-duper-secret__multiply", "multiply" },
	{ "_super-duper-secret__divides", "divides" },
	{ "_super-duper-secret__plus", "plus" },
	{ "_super-duper-secret__minus", "minus" },
}

func (p *Parser) parseStmnt(words []string) (Ast, error) {
	if p.debug { fmt.Printf("parseStmnt %v\n", words ) }


	if left,right,found := Partition(words, "that"); found {
		fun, err := p.parseStmnt(left)
		if err != nil { return nil,err }
		args, err := p.parseArgs(right)
		if err != nil { return nil,err }
		node := AstFunCall{
			fun: fun,
			args: args,
		}
		return node, nil
	}

	if len(words) == 0 {
		return AstLiteral{value: BsNilVal{}}, nil
	}

	// infix operators
	for _, rule := range infixRules {

		if left, right, found := Partition(words, rule.opname); found {
			lval, err := p.parseStmnt(left)
			if err != nil { return nil, err }
			rval, err := p.parseStmnt(right)
			if err != nil { return nil, err }
			node := rule.ctor(lval,rval)
			return node, nil
		}
	}
	
	// builtin infix functions

	for _, infix := range infixBuiltins {

		if left,right,found := Partition(words, infix.opname); found {
			lexpr, err := p.parseStmnt(left)
			if err != nil { return nil, err }
			rexpr, err := p.parseStmnt(right)
			if err != nil { return nil, err }
			node := AstFunCall {
				fun: AstIdent{name:infix.symbol},
				args: []Ast { lexpr, rexpr },
			}
			return node, nil
		}
	}

	// prefix operators
	switch (words[0]) {
		case "the":
			name := strings.Join(words[1:], " ")
			return AstIdent{name: name}, nil
		case "text":
			literal := strings.Join(words[1:], " ")
			value := BsStrVal{value:literal}
			node := AstLiteral{value:value}
			return node, nil
		case "if":
			cond,err := p.parseStmnt(words[1:])
			if err != nil { return nil,err }
			_,err = p.parseBlock()
			if err != nil { return nil,err }
			
			return cond,nil

		default:
			// we assume its a function call
			node, err := p.parseFunCall(words)
			if err != nil {
				return nil, err
			}
			return node, nil	
	}
}

func (p *Parser) parseBlock() (Ast, error) {
	
	return nil,nil
}

func (p *Parser) parseFunCall(words []string) (Ast, error) {
	if p.debug { fmt.Printf("parseFunCall %v\n", words ) }
	head, err := p.parseAtom(words[0])
		if err != nil { return nil, err }
	if len(words) == 1 {
		return head, nil
	}
	args, err := p.parseArgs(words[1:])
	if err != nil { return nil, err }
	node := AstFunCall{
		fun: head,
		args: args,
	}
	return node, nil
}		

func (p *Parser) parseArgs(words []string) ([]Ast, error) {
	if p.debug { fmt.Printf("parseArgs %v\n", words ) }
	node, err := p.parseStmnt(words)
	if err != nil { return nil, err }
	list := []Ast { node }
	return list, nil
}
	
func (p *Parser) parseAtom(word string) (Ast, error) {
	r := FirstRune(word)
  if unicode.IsNumber(r) {
		var value int64
		rc, err := fmt.Sscanf(word, "%d", &value)
		if err != nil {
			return nil, err
		}
		if rc == 0 {
			return nil, parseErr("not a valid integer", word)
		}
		literal := BsIntVal{value:value}
		node := AstLiteral{value:literal}
		return node, nil
	} else if word == "true" {
		node := AstLiteral{
			value: BsBooleVal{
				value: true,
			},
		}
		return node,nil
	} else if word == "false" {
		node := AstLiteral{
			value: BsBooleVal{
				value: false,
			},
		}	
		return node,nil
	} else {
		node := AstIdent{name:word}
		return node, nil
	}
}

func parseErr(msg string, src string) error {
	return errors.New("Parse Error:" + msg)
} 
