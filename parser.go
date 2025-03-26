package main

import (
	"errors"
	"strings"
	"fmt"
	"log"
	"bufio"
)


type Parser struct {
	buf *bufio.Reader
	debug bool
	tokens []Token
	pos int
}
func eof() Token {
	return Token{
		Lex:"",
		Ty:TOKEN_EOF}
}
func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return eof()
	}
	return p.tokens[p.pos]
}
func (p *Parser) consumeOrFail(ty TokenType) (Token, error) {
	tok := p.peek()
	if tok.Ty != ty {
		return eof(), errors.New(fmt.Sprintf("expected %v, found %v\n", tok.Ty, ty))
	}
	if ty != TOKEN_EOF {
		p.pos += 1
	}
	return tok, nil
}
func (p *Parser) hasTokens() bool {
	return p.peek().Ty != TOKEN_EOF
}

func (p *Parser) consumeLine() []Token {
begin := p.pos
		for p.peek().Ty != TOKEN_NEWLINE && p.peek().Ty != TOKEN_EOF {
			p.pos += 1
		}
		end := p.pos
		p.pos += 1
		return p.tokens[begin:end]
}


func (p *Parser) Parse() ([]Ast, error) {
	ast := make([]Ast, 0, 50)
	for p.hasTokens() {
		line := p.consumeLine()
		node, err := p.parseStmnt(line)
		if err != nil { return nil,err }
		ast = append(ast, node)
	}
	return ast,nil
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

func (p *Parser) parseStmnt(words []Token) (Ast, error) {
	if p.debug { fmt.Printf("parseStmnt %v\n", words ) }

	if len(words) == 0 { return AstLiteral{value: BsNilVal{}}, nil }

	if left, right, found := Partition(words, TOKEN_KW_IS); found {
		lval, err := p.parseIdent(left)
		if err != nil { return nil, err }
		rval, err := p.parseExpr(right)
		if err != nil { return nil, err }
		node := AstAssign{lval,rval}
		return node, nil
	}

  if words[0].Ty == TOKEN_KW_IF {
		cond,err := p.parseExpr(words[1:])
		if err != nil { return nil,err }
		_,err = p.parseBlock()
		if err != nil { return nil,err }
		
		return cond,nil
	}

	return p.parseExpr(words)
}

func (p *Parser) parseExpr(words []Token) (Ast, error) {
  if p.debug { log.Printf("parseExpr %v\n", words) }

	if left,right,found := Partition(words, TOKEN_KW_OF); found {
		fun, err := p.parseExpr(left)
		if err != nil { return nil,err }
		args, err := p.parseArgs(right)
		if err != nil { return nil,err }
		node := AstFunCall{
			fun: fun,
			args: args,
		}
		return node, nil
	}

	// builtin infix functions

	for _, infix := range infixBuiltins {

		if left,right,found := PartitionByLexeme(words, infix.opname); found {
			lexpr, err := p.parseExpr(left)
			if err != nil { return nil, err }
			rexpr, err := p.parseExpr(right)
			if err != nil { return nil, err }
			node := AstFunCall {
				fun: AstIdent{name:infix.symbol},
				args: []Ast { lexpr, rexpr },
			}
			return node, nil
		}
	}
	
	// prefix operators
	if words[0].Ty == TOKEN_KW_THE {
		return p.parseIdent(words)
	}

	// we assume its a function call
	return p.parseFunCall(words)
}

func (p *Parser) parseBlock() ([]Ast, error) {
	if p.debug { log.Printf("parseBlock \n") }
	p.consumeOrFail(TOKEN_BEGIN_INDENT)

	ast := make([]Ast, 0, 5)
	for p.peek().Ty != TOKEN_END_INDENT {
		line := p.consumeLine()
		n, err := p.parseStmnt(line)
		if err != nil { return nil, err }
		ast = append(ast, n)
	}
	p.consumeOrFail(TOKEN_END_INDENT)
	
	return ast, nil
}

func (p *Parser) parseFunCall(words []Token) (Ast, error) {
	if p.debug { log.Printf("parseFunCall %v\n", words ) }

	if len(words) == 0 {
		return nil,errors.New("need tokens inside function call")
	}

	if words[0].Ty != TOKEN_WORD {
		// its an atomic node, only a single token ks allowed
		if len(words) > 1 {
			return nil,parseErr("extra tokens following atom. not expected", words[1])
		}

		atom, err := p.parseAtom(words[0])
		if err != nil { return nil, err }
		return atom, nil
	}

	// only case where a single bare word can become an identifier

	head := AstIdent{
		name:words[0].Lex,
	}
	args, err := p.parseArgs(words[1:])
	if err != nil { return nil, err }
	node := AstFunCall{
		fun: head,
		args: args,
	}
	return node, nil
}		

func (p *Parser) parseArgs(words []Token) ([]Ast, error) {
	if p.debug { log.Printf("parseArgs %v\n", words ) }

	// todo: handle empty list  multiple args better
	node, err := p.parseExpr(words)
	if err != nil { return nil, err }
	list := []Ast { node }
	return list, nil
}

func (p *Parser) parseIdent(words []Token) (Ast, error) {
	if p.debug { log.Printf("parseIdent %v\n", words) }

	if words[0].Ty != TOKEN_KW_THE {
		return nil,parseErr(fmt.Sprintf("expected 'the' keyword to begin identifier, found %v", words[0]), words[0])
	}
	b := new(strings.Builder)
	for i, w := range words[1:] {
		if i != 0 {
			b.WriteString(" ")
		}
		if w.Ty != TOKEN_WORD {
			return nil,parseErr(fmt.Sprintf("expected TOKEN_WORD after 'the' keyword, found %v", w.Ty), w)
		}
		b.WriteString(w.Lex)
	}
	name := b.String()
	return AstIdent{name: name}, nil 
}

func (p *Parser) parseAtom(word Token) (Ast, error) {
	if p.debug { log.Printf("parseAtom %v\n", word) }

  if word.Ty == TOKEN_NUMBER {
		var value int64
		rc, err := fmt.Sscanf(word.Lex, "%d", &value)
		if err != nil {
			return nil, err
		}
		if rc == 0 {
			return nil, parseErr("not a valid number", word)
		}
		literal := BsIntVal{value:value}
		node := AstLiteral{value:literal}
		return node, nil
	} else if word.Ty == TOKEN_KW_TRUE {
		node := AstLiteral{
			value: BsBooleVal{
				value: true,
			},
		}
		return node,nil
	} else if word.Ty == TOKEN_KW_FALSE {
		node := AstLiteral{
			value: BsBooleVal{
				value: false,
			},
		}	
		return node,nil
	} else if word.Ty == TOKEN_TEXT {
			value := BsStrVal{value:word.Lex}
			node := AstLiteral{value:value}
			return node, nil
	} else {
		return nil, parseErr(fmt.Sprintf("unexpected token type %v inside atom", word), word)
	}
}

func parseErr(msg string, src Token) error {
	return errors.New(" Parse Error:" + msg + " at : " )
} 

// utilities

func PartitionByLexeme(words []Token, split string) ([]Token, []Token, bool) {
	idx := FindFirst(words, func(tok Token)bool { return tok.Lex == split})
	if idx == -1 {
		return nil,nil,false
	}
	left := words[:idx]
	right := words[idx+1:]
	return left,right,true
}
	

func Partition(words []Token, split TokenType) ([]Token, []Token, bool) {
	idx := FindFirst(words, func(tok Token)bool { return tok.Ty == split})
	if idx == -1 {
		return nil,nil,false
	}
	left := words[:idx]
	right := words[idx+1:]
	return left,right,true
}
	
func FindFirst(words []Token, pred func(Token)bool) int {
	for i := range words {
		if pred(words[i]) { return i }
	}
	return -1
}

