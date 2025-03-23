package main

import (
	"errors"
	"strings"
)


type Parser struct {

}

func MakeParser() Parser {
	return Parser{}
}

func (p *Parser) ParseLine(line string) (Ast, error) {

	line = line[:len(line)-1]
	words := strings.Split(line, " ")
	ast, err := p.parseExpr(words)
	return ast, err
}
func (p *Parser) parseExpr(words []string) (Ast, error) {

	if len(words) == 0 {
		return AstLiteral{value: BsNilVal{}}, nil
	}
	// infix operators
	if contains(words, "is") {
		idx := indexof(words, "is")
		lval, err := p.parseExpr(words[:idx])
		if err != nil {
			return nil, err
		}
		rval, err := p.parseExpr(words[idx+1:])
		node := AstAssign {
			lvalue: lval,
			rvalue: rval}
		return node, nil
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
		default:
			// we assume its a function call
			node, err := p.parseFunCall(words)
			if err != nil {
				return nil, err
			}
			return node, nil	
	}
}


func (p *Parser) parseFunCall(words []string) (Ast, error) {
	head := AstIdent{name:words[0]}
	rest, err := p.parseExpr(words[1:])
	if err != nil {
		return nil, err
	}
	node := AstFunCall{fun:head, args:[]Ast{ rest }}
	return node, nil
}		

func parseErr(msg string, src string) error {
	return errors.New("Parse Error:" + msg)
}

func contains(words []string, w string) bool {
    return indexof(words, w) != -1
}
func indexof(words []string, w string) int {
	for i := range words {
		if words[i] == w { return i }
	}
	return -1
}
