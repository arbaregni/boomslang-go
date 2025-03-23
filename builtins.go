package main

import (
	"fmt"
)

func LoadBuiltins(env *BsEnv) {
	env.AssignName("show", BsFunVal { thunk: BsBuiltinShow{} })
}

type BsBuiltinShow struct {
}

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

