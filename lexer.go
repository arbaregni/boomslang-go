package main

import (
	"bufio"
	"strings"
	"unicode"
	"errors"
	"io"
	"fmt"
	"log"
)


type TokenType string
const (
	TOKEN_NEWLINE TokenType = "TOKEN_NEWLINE"
	TOKEN_EOF = "TOKEN_EOF"
	TOKEN_BEGIN_INDENT = "TOKEN_BEGIN_INDENT"
	TOKEN_END_INDENT = "TOKEN_END_INDENT"
	TOKEN_NUMBER = "TOKEN_NUMBER"
	TOKEN_WORD = "TOKEN_WORD"
	TOKEN_TEXT = "TOKEN_TEXT"
	TOKEN_KW_IS = "TOKEN_KW_IS"
	TOKEN_KW_THE = "TOKEN_KW_THE"
	TOKEN_KW_FALSE = "TOKEN_KW_FALSE"
	TOKEN_KW_TRUE = "TOKEN_KW_TRUE"
	TOKEN_KW_THAT = "TOKEN_KW_THAT"
	TOKEN_KW_OF = "TOKEN_KW_OF"
	TOKEN_KW_IF = "TOKEN_KW_IF"
	TOKEN_KW_FOR = "TOKEN_KW_FOR"
	TOKEN_KW_WHILE = "TOKEN_KW_WHILE"
)

type Span struct {
	Filename string
	Lineno int
	Begin int
	End int
}

type Token struct {
	Ty TokenType
	Lex string
}

type Lexer struct {
	debug bool
	filename string
	lineno int
	buf *bufio.Reader
	tokens []Token
	indent int
	shiftWidth indentlevel
}

func (l *Lexer) emit(lex string, ty TokenType) {
	tok := Token{ty,lex}
	l.tokens = append(l.tokens,tok)
}

func (l *Lexer) Lex() ([]Token, error) {
	for {
		err := l.lexLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil,err
		}

	}
	// emit dedents down to zero
	if err := l.handleIndent(indentlevel{}); err != nil {
		return nil,err
	}

	return l.tokens, nil
}
func (l *Lexer) lexLine() error {
	if l.debug { log.Printf("entering lexLine\n") }
	line, err := l.buf.ReadString('\n')
	if err != nil { return err }
	l.lineno += 1

	// emit indents
	line, indent := TrimIndent(line)
	if err := l.handleIndent(indent); err != nil {
		return err
	}

	if err != nil { return err }
		
	// word to token
	line = strings.TrimSpace(line)
	words := strings.Fields(line)
	for i, word := range words {
		if word == "is" {
			l.emit(word, TOKEN_KW_IS)
		} else if word == "if" {
			l.emit(word, TOKEN_KW_IF)
		} else if word == "of" {
			l.emit(word, TOKEN_KW_OF)
		} else if word == "the" {
			l.emit(word, TOKEN_KW_THE)
		} else if word == "that" {
			l.emit(word, TOKEN_KW_THAT)
		} else if word == "for" {
			l.emit(word, TOKEN_KW_FOR)
		} else if word == "while" {
			l.emit(word, TOKEN_KW_WHILE)
		} else if word == "true" {
			l.emit(word, TOKEN_KW_TRUE)
		} else if word == "false" {
			l.emit(word, TOKEN_KW_FALSE)
		} else if word == "text" {
			text := strings.Join(words[i+1:], " ")
			l.emit(text, TOKEN_TEXT)
			break // break out of for loop
		} else if unicode.IsNumber(FirstRune(word)) {
			l.emit(word, TOKEN_NUMBER)
		} else {
			l.emit(word, TOKEN_WORD)
		}	
		
	}

	// emit newline
	l.emit("\n", TOKEN_NEWLINE)

	return nil
}

func (l *Lexer) handleIndent(newIndent indentlevel) error {

	if l.debug { log.Printf("entering handleIndrnt, newIdent = %v, curr = %v\n", newIndent, l.indent) }
	curr := l.indent
	if l.shiftWidth.spaces == 0 && l.shiftWidth.tabs == 0 {
		// if its zero, set the shift width for the first time
		l.shiftWidth = newIndent
		if l.debug { log.Printf("setting shift widths to %v\n", l.shiftWidth) }
	}

	newLevel, err := translateIndent(newIndent, l.shiftWidth)
	if err != nil { return err }

	diff := newLevel - curr
	
	if diff > 0 {
		if diff != 1 {
			return errors.New(fmt.Sprintf("can not indent multiple at a time: you tried to indent %d levels", diff))
		}
		l.emit("  ", TOKEN_BEGIN_INDENT)
		l.indent = newLevel
	  return nil
	}

	diff = -diff
	if l.debug { log.Printf("dedenting from curr = %d to newLevel = %d\n", curr, newLevel) }

	for i := 0; i < diff; i += 1 {
		l.emit("  ", TOKEN_END_INDENT)
	}
	l.indent = newLevel
	return nil
}

func translateIndent(newIndent indentlevel, shiftWidth indentlevel) (int, error) {
	if shiftWidth.tabs != 0 && shiftWidth.spaces != 0 {
		return 0, errors.New("can not mix tabs and spaces")
	}

	if shiftWidth.tabs > 0 {
		if newIndent.spaces > 0 {
			return 0, errors.New("can not mix tabs and spaces")
		}
		if shiftWidth.tabs != 1 {
			return 0, errors.New("can not have multipletabs in shift width")
		}
		return newIndent.tabs, nil
	}

	if shiftWidth.spaces > 0 {
		if newIndent.spaces % shiftWidth.spaces != 0 {
			return 0,errors.New(fmt.Sprintf("wrong number of spaces in indentation: %d. You are using %d", newIndent.spaces, shiftWidth.spaces))
		}
		return newIndent.spaces / shiftWidth.spaces, nil
	}

	if newIndent.spaces == 0 && newIndent.tabs == 0 {
		return 0, nil
	}
	
	return 0, errors.New("can have zero shift width (this is likely a compiler bug)")
}





