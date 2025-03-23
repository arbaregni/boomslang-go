package main

import (
	"fmt"
)

func LoadBuiltins(env *BsEnv) {
	// top level functions
	env.AssignName("show", BsFunVal { thunk: BsBuiltinShow{} })
	// magic names
	env.AssignName("_super-duper-secret__plus", makeIntBinOp("plus", func(x,y int64)int64 { return x + y}))
	env.AssignName("_super-duper-secret__minus", makeIntBinOp("minus", func(x,y int64)int64 { return x - y}))
	env.AssignName("_super-duper-secret__multiply", makeIntBinOp("multiply", func(x,y int64)int64 { return x * y}))
	env.AssignName("_super-duper-secret__divide", makeIntBinOp("divide", func(x,y int64)int64 { return x / y}))
	env.AssignName("_super-duper-secret__smallerthan", makeIntBinPred("smallerthan", func(x,y int64)bool { return x < y}))
	env.AssignName("_super-duper-secret__biggerthan", makeIntBinPred("biggerthan", func(x,y int64)bool { return x > y}))
	env.AssignName("_super-duper-secret__equals", makeIntBinPred("equls", func(x,y int64)bool { return x == y}))
}

type BsBuiltinShow struct {}
func (this BsBuiltinShow) PrettyPrint() string {
	return fmt.Sprintf("<builtin procedure 'show'>")
}
func (this BsBuiltinShow) Call(env *BsEnv, args []BsValue) BsValue {
	for i, arg := range args {
		if i != 0 {
			fmt.Fprintf(env.ostr, " ")
		}
		fmt.Fprintf(env.ostr,"%s", arg.PrettyPrint());
	}
	fmt.Fprintf(env.ostr,"\n")
	return new(BsNilVal)
}

type BsBuiltinIntBinOp struct {
	name string
	op func(int64,int64)int64
}
func makeIntBinOp(name string, op func(int64,int64)int64) BsFunVal {
	thunk := BsBuiltinIntBinOp {
		name:name,
		op:op,
	}
	fun := BsFunVal { thunk: thunk }
	return fun
}
func (this BsBuiltinIntBinOp) PrettyPrint() string {
	return fmt.Sprintf("<builtin procedure '%s'>", this.name)
}
func (this BsBuiltinIntBinOp) Call(env *BsEnv, args []BsValue) BsValue {
	if len(args) != 2 {
		return BsMethodErr{expected:"2 parameters"}
	}
	left,ok := args[0].(BsIntVal)
	if !ok { return BsTypeErr{expected:"number", value:args[0]} }
	right,ok := args[1].(BsIntVal)
	if !ok { return BsTypeErr{expected:"number", value:args[1]} }
	value := this.op(left.value, right.value)
	return BsIntVal{value:value}
}

type BsBuiltinIntBinPred struct {
	name string
	op func(int64,int64)bool
}
func makeIntBinPred(name string, op func(int64,int64)bool) BsFunVal {
	thunk := BsBuiltinIntBinPred {
		name:name,
		op:op,
	}
	fun := BsFunVal { thunk: thunk }
	return fun
}
func (this BsBuiltinIntBinPred) PrettyPrint() string {
	return fmt.Sprintf("<builtin procedure '%s'>", this.name)
}
func (this BsBuiltinIntBinPred) Call(env *BsEnv, args []BsValue) BsValue {
	if len(args) != 2 {
		return BsMethodErr{expected:"2 parameters"}
	}
	left,ok := args[0].(BsIntVal)
	if !ok { return BsTypeErr{expected:"number", value:args[0]} }
	right,ok := args[1].(BsIntVal)
	if !ok { return BsTypeErr{expected:"number", value:args[1]} }
	value := this.op(left.value, right.value)
	return BsBooleVal{value:value}
}
