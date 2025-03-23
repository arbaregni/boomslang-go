package main


import (
	"fmt"
)

type BsValue interface {
	IsErr() bool
	PrettyPrint() string
}
type BsNilVal struct {}
func (v BsNilVal) IsErr() bool {
	return false;
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

type BsNameErr struct {
	name string
}
func (v BsNameErr) IsErr() bool {
	return true
}
func (v BsNameErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, I tried and failed to find the name '%s' in the place you requested it.", v.name)
}

type BsTypeErr struct {
	expected string
	value BsValue
}
func (v BsTypeErr) IsErr() bool {
	return true
}
func (v BsTypeErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, but this is not a valid %s: %s", v.expected, v.value.PrettyPrint())
}

type BsMethodErr struct {
	expected string
}
func (v BsMethodErr) IsErr() bool {
	return true
}
func (v BsMethodErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, but you invoked a procedure with a bad set of arguments: %s", v.expected)
}

type BsUnpackErr struct {
	expected string
	value Ast
}
func (v BsUnpackErr) IsErr() bool {
	return true
}
func (v BsUnpackErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, but I can not assign to this %s: %s", v.expected, v.value.ShortName())
}
