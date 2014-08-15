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
                println("BuiltinPlus.do")
		sum := 0
                for v := range IterOverPairAsList(cdr) {
			sum += int(Eval(v, e).(Int))
                }
		return Int(sum)
	},
}
