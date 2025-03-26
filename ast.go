package main

type Ast interface {
	Eval(env *BsEnv) BsValue
	ShortName() string
}

type AstFunCall struct {
	fun  Ast
	args []Ast
}

func (node AstFunCall) ShortName() string { return "procedure" }

type AstIdent struct {
	name string
}

func (node AstIdent) ShortName() string { return "name" }

type AstLiteral struct {
	value BsValue
}

func (node AstLiteral) ShortName() string { return node.value.PrettyPrint() }

type AstAssign struct {
	lvalue Ast
	rvalue Ast
}

func (node AstAssign) ShortName() string { return "assignment" }

type AstIfStmnt struct {
	cond       Ast
	if_block   []Ast
	else_block []Ast
}

func (node AstIfStmnt) ShortName() string { return "if statement" }
