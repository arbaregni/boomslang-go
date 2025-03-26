package main

import (
	"fmt"
	"io"
	"log"
)

type BsEnv struct {
	symbols map[string]BsValue
	debug   bool
	ostr    io.Writer
	estr    io.Writer
}

func MakeEnv(ostr, estr io.Writer) *BsEnv {
	env := new(BsEnv)
	env.symbols = make(map[string]BsValue, 50)
	env.ostr = ostr
	env.estr = estr
	return env
}
func (env *BsEnv) AssignName(name string, value BsValue) {
	env.symbols[name] = value
}
func (env *BsEnv) Lookup(name string) BsValue {
	val, ok := env.symbols[name]
	if !ok {
		return BsNameErr{name: name}
	}
	return val
}

// implement Eval for all Ast nodes

func (node AstFunCall) Eval(env *BsEnv) BsValue {
	if env.debug {
		log.Printf(" Eval AstFunCall\n")
	}

	fun := node.fun.Eval(env)
	if fun.IsErr() {
		fmt.Fprintf(env.estr, "[[Encountered failure while evaluating head expression]]\n")
		return fun
	}

	args := make([]BsValue, len(node.args))

	for i, _ := range node.args {
		args[i] = node.args[i].Eval(env)
		if args[i].IsErr() {
			fmt.Fprintf(env.estr, "[[Encountered dailure getting the value on the %dth thing]]\n", i+1)
			return args[i]
		}
	}

	// if its a function, it can now be called
	var out BsValue = BsNilVal{}
	funVal, ok := fun.(BsFunVal)
	if !ok {
		return BsUnpackErr{expected: "procedure", value: node.fun}
	}
	out = funVal.thunk.Call(env, args)
	if out.IsErr() {
		fmt.Fprintf(env.estr, "[[Encountered a failure while invoking a function]]\n")
		return out
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
		return rvalue
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
		return cond
	}
	// truthyness evaluation
	cond_b, ok := cond.(BsBooleVal)
	if !ok {
		return BsTypeErr{}
	}

	if env.debug {
		log.Printf(" Eval AstIfStmnt: cond_b is %v\n", cond_b)
	}

	if cond_b.value {
		return EvalAll(env, node.if_block)
	} else {
		return EvalAll(env, node.else_block)
	}
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
