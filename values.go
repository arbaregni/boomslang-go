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

type BsNameErrVal struct {
	name string
}
func (v BsNameErrVal) IsErr() bool {
	return true
}
func (v BsNameErrVal) PrettyPrint() string {
	return fmt.Sprintf("Sorry, I tried and failed to find the name '%s' in the place you requested it.", v.name)
}

type BsTypeErr struct {
	expected string
	value Ast
}
func (v BsTypeErr) IsErr() bool {
	return true
}
func (v BsTypeErr) PrettyPrint() string {
	return fmt.Sprintf("Sorry, but this is not a valid %s: %s", v.expected, v.value.ShortName())
}
