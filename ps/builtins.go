package ps

import (
	"fmt"
)

type Builtin struct {
	do func(cdr Value, e *EnvFrame) Value
}

func (x *Builtin) Pair() bool {
	return false
}

func (x *Builtin) Do(cdr Value, e *EnvFrame) Value {
	return x.do(cdr, e)
}

func (x *Builtin) String() string {
	return fmt.Sprintf("%x", x)
}

var BuiltinPlus = &Builtin{
	do: func(cdr Value, e *EnvFrame) Value {
		sum := 0
		for v := range IterOverPairAsList(cdr) {
			sum += int(Eval(v, e).(Int))
		}
		return Int(sum)
	},
}

var BuiltinMinus = &Builtin{
	do: func(cdr Value, e *EnvFrame) Value {
		sum := 0
		for v := range IterOverPairAsList(cdr) {
			sum -= int(Eval(v, e).(Int))
		}
		return Int(sum)
	},
}

var BuiltinMul = &Builtin{
	do: func(cdr Value, e *EnvFrame) Value {
		sum := 0
		for v := range IterOverPairAsList(cdr) {
			sum *= int(Eval(v, e).(Int))
		}
		return Int(sum)
	},
}

var BuiltinDiv = &Builtin{
	do: func(cdr Value, e *EnvFrame) Value {
		sum := 0
		for v := range IterOverPairAsList(cdr) {
			sum /= int(Eval(v, e).(Int))
		}
		return Int(sum)
	},
}

var BuiltinDefine = &Builtin{
	do: func(cdr Value, e *EnvFrame) Value {
		name, ok  := Car(cdr).(Name)
                if !ok {
                      panic("first must be name")
                }
		value := Car(Cdr(cdr))
		e.Bind(string(name), value)
		return nil
	},
}

func GlobalEnv() *EnvFrame {
	env := MakeEnv()
	env.Bind("+", BuiltinPlus)
	env.Bind("-", BuiltinMinus)
	env.Bind("*", BuiltinMul)
	env.Bind("/", BuiltinDiv)
	env.Bind("define", BuiltinDefine)
	return env
}
