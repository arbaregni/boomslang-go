package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type BsEnv struct {
	symbols map[string]BsValue
	debug   bool
	istr    io.Reader
	ostr    io.Writer
	estr    io.Writer
}

func MakeEnv(opts *Opts) *BsEnv {
	env := new(BsEnv)
	env.symbols = make(map[string]BsValue, 50)
	env.istr = opts.istr
	env.ostr = opts.ostr
	env.estr = opts.estr
	return env
}
func (env *BsEnv) AssignName(name string, value BsValue) {
	if env.debug {
		log.Printf("assigning symbol '%s' to %v\n", name, value)
	}
	env.symbols[name] = value
}
func (env *BsEnv) Lookup(name string) BsValue {
	val, ok := env.symbols[name]
	if !ok {
		return BsNameErr{name: name}
	}
	return val
}

// For collecting context on the way up the stack
type BsUnwindCtx struct {
	init BsValue // the value that was initially thrown
	frames []BsEvalFrame
}

func (v BsUnwindCtx) IsErr() bool {
	return true
}
func (v BsUnwindCtx) PrettyPrint() string {
	b := new(strings.Builder)
	b.WriteString(v.init.PrettyPrint())
	b.WriteString("\n")
	for i, frame := range v.frames {
		b.WriteString(fmt.Sprintf("  [%d] : %s\n", i,frame.msg))
	}
	return b.String()
}

type BsEvalFrame struct {
	node Ast
	msg string
}

func (env *BsEnv) addFrame(throw BsValue, node Ast, format string, args ...any) BsUnwindCtx {
	msg := fmt.Sprintf(format, args...)
	frame := BsEvalFrame{node:node,msg:msg}
	if env.debug {
		log.Printf("unwinding. New frame: msg = %s, node = %s\n",msg,node.ShortName())
	}
	if ctx, ok := throw.(BsUnwindCtx); ok {
		ctx.frames = append(ctx.frames, frame)
		return ctx
	}
	return BsUnwindCtx{init:throw,frames:[]BsEvalFrame{frame}}
}

// implement Eval for all Ast nodes

func (node AstFunCall) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstFunCall\n")
	}

	fun := node.fun.Eval(env)
	if fun.IsErr() {
		return env.addFrame(fun, node, "while evaluating head expression of procedure")
	}

	args := make([]BsValue, len(node.args))

	for i, _ := range node.args {
		args[i] = node.args[i].Eval(env)
		if args[i].IsErr() {
			return env.addFrame(args[i], node, "Encountered failure evaluating the %dth argument", i+1)
		}
	}

	// if its a function, it can now be called
	var out BsValue = BsNilVal{}
	funVal, ok := fun.(BsFunVal)
	if !ok {
		return BsMethodErr{expected: "can not invoke '" + fun.PrettyPrint() + "'"}
	}
	out = funVal.thunk.Call(env, args)
	if out.IsErr() {
		 return env.addFrame(out, node, "Encountered a failure while invoking a function")
	}
	return out
}
func (node AstIdent) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstIdent\n")
	}

	// check the name in lookup table
	return env.Lookup(node.name)
}
func (node AstLiteral) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstLiteral\n")
	}

	return node.value
}
func (node AstAssign) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstAssign\n")
	}

	lvalue, ok := node.lvalue.(AstIdent)
	if !ok {
		return BsUnpackErr{expected: "lvalue", value: node.lvalue}
	}
	rvalue := node.rvalue.Eval(env)
	if rvalue.IsErr() {
		return env.addFrame(rvalue, node, "Encountered failure evaluating right side of expression")
	}
	env.AssignName(lvalue.name, rvalue)
	return BsNilVal{}
}
func (node AstIfStmnt) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstIfStmnt\n")
	}

	cond := node.cond.Eval(env)
	if cond.IsErr() {
		return env.addFrame(cond, node, "Encountered failure evaluating condition of if statement")
	}
	cond_b := bsTruthy(cond)

	if env.debug {
		log.Printf(" Eval AstIfStmnt: cond_b is %v\n", cond_b)
	}

	if cond_b {
		return EvalAll(env, node.if_block)
	} else {
		return EvalAll(env, node.else_block)
	}
}
func (node AstLoop) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstLoop\n")
	}

	for {
		cond := node.cond.Eval(env)
		if cond.IsErr() {
			return env.addFrame(cond, node, "Encountered failure evaluating condition of loop")
		}
		cond_b := bsTruthy(cond)
		if env.debug {
			log.Printf(" Eval AstIfLoop: cond_b is %v\n", cond_b)
		}
		if !cond_b {
			break
		}
		expr := EvalAll(env, node.block)
		// todo: share some of the try/catch stuff
		if _,ok := expr.(BsBreakExc); ok {
			if env.debug {
				log.Printf(" Eval AstLoop: caught break \n")
			}
			break
		}
		if expr.IsErr() {
			return env.addFrame(expr, node, "encountered error in loop body")
		}
	}

	return BsNilVal{}
}
func (node AstBreak) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstBreak\n")
	}
	return BsBreakExc{}
}

	
// utilities for multiple ast nodes
func EvalAll(env *BsEnv, ast []Ast) BsValue {
	if env.debug {
		log.Printf(" Eval [%d]Ast\n", len(ast))
	}
	var out BsValue = BsNilVal{}
	for _, node := range ast {
		out = node.Eval(env)
		if out.IsErr() {
			return out
		}
	}
	return out
}
// casting utilities


// truthyness evaluation
func bsTruthy(value BsValue) bool {
	switch v := value.(type) {
	case BsBooleVal:
		return v.value
	case BsIntVal:
		return v.value > 0 // note: purposefully annoying, negatives are falsey
	case BsStrVal:
		return len(v.value) > 0
	case BsNilVal:
		return false
	}
	if value.IsErr() {
		return false 
	}
	return true
}
