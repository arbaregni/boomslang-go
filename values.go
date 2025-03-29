package main

import (
	"fmt"
)

type BsValue interface {
	IsErr() bool
	PrettyPrint() string
}
type BsNilVal struct{}

func (v BsNilVal) IsErr() bool {
	return false
}
func (v BsNilVal) PrettyPrint() string {
	return fmt.Sprintf("nothing")
}

type BsStrVal struct {
	value string
}

func (v BsStrVal) IsErr() bool {
	return false
}
func (v BsStrVal) PrettyPrint() string {
	return fmt.Sprintf("%s", v.value)
}

type BsBooleVal struct {
	value bool
}

func (v BsBooleVal) IsErr() bool {
	return false
}
func (v BsBooleVal) PrettyPrint() string {
	if v.value {
		return "true"
	} else {
		return "false"
	}
}

type BsIntVal struct {
	value int64
}

func (v BsIntVal) IsErr() bool {
	return false
}
func (v BsIntVal) PrettyPrint() string {
	return fmt.Sprintf("%d", v.value)
}

type BsFunVal struct {
	thunk BsFunThunk
}

func (v BsFunVal) IsErr() bool {
	return false
}
func (v BsFunVal) PrettyPrint() string {
	return v.thunk.PrettyPrint()
}

type BsFunThunk interface {
	Call(env *BsEnv, args []BsValue) BsValue
	PrettyPrint() string
}

type BsRuntimeFunc struct {
	name   *AstIdent
	env    *BsEnv     // functions use the scope where they were defined
	params []AstIdent // to be instantiated
	body   []Ast
	// todo: python has much worse semantics
}

func (v BsRuntimeFunc) Call(callerEnv *BsEnv, args []BsValue) BsValue {
	if len(args) != len(v.params) {
		return BsMethodErr{expected: fmt.Sprintf("need %d arguments, got %d", len(v.params), len(args))}
	}
	// each invocarion gets a fresh state
	invocationEnv := v.env.NewChild()
	// instantiate the parameters
	for i, param := range v.params {
		arg := args[i]
		invocationEnv.AssignName(param.name, arg)
	}
	out := EvalAll(invocationEnv, v.body)
	if out.IsErr() {
		// todo: this is actually a call frame
		return invocationEnv.addFrame(out, v.name, "during function invocation")
	}
	// todo: catch returns here
	return BsNilVal{}
}

func (v BsRuntimeFunc) PrettyPrint() string {
	if v.name == nil {
		return "<unnamed procedure>"
	}
	return "<procedure '" + v.name.name + "'>"
}

// ====================================
//
//	name errors
type BsNameErr struct {
	name string
}

func (v BsNameErr) IsErr() bool {
	return true
}
func (v BsNameErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, I tried and failed to find the name '%s' in the place you requested it.", v.name)
}

// ====================================
//
//	type errors
type BsTypeErr struct {
	expected string
	value    BsValue
}

func (v BsTypeErr) IsErr() bool {
	return true
}
func (v BsTypeErr) PrettyPrint() string {
	return fmt.Sprintf("(TypeError) Sorry, but this is not a valid %s: %s", v.expected, v.value.PrettyPrint())
}

// ====================================
//  Method errors

type BsMethodErr struct {
	expected string
}

func (v BsMethodErr) IsErr() bool {
	return true
}
func (v BsMethodErr) PrettyPrint() string {
	return fmt.Sprintf("(MethodError) Sorry, but you invoked a procedure with a bad set of arguments: %s", v.expected)
}

// ====================================
//
//	unpack errors
type BsUnpackErr struct {
	expected string
	value    Ast
}

func (v BsUnpackErr) IsErr() bool {
	return true
}
func (v BsUnpackErr) PrettyPrint() string {
	return fmt.Sprintf("(UnpackError) Sorry, but I can not assign to this %s: %s", v.expected, v.value.ShortName())
}

// ====================================
//  io errors

type BsIoErr struct {
	msg string
}

func (v BsIoErr) IsErr() bool {
	return true
}
func (v BsIoErr) PrettyPrint() string {
	return fmt.Sprintf("(IoError) Sorry, but something happened with the file system: %s", v.msg)
}

// ====================================
//  break exception - used for breaking out loops

type BsBreakExc struct {
}

func (v BsBreakExc) IsErr() bool {
	return true
}
func (v BsBreakExc) PrettyPrint() string {
	return fmt.Sprintf("(BreakException) break out of loop")
}
