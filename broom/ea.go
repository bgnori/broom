package broom

import (
	"fmt"
)

type Environment interface {
	Bind(name string, v interface{})
	Resolve(name string) (interface{}, error)
	SetOuter(outer Environment)
	Outer() Environment
	Dump()
}

type EvalError string

func (e EvalError) Error() string {
	return string(e)
}

func Eval(env Environment, expr interface{}) interface{} {
	if v, err := env.Resolve("_debug"); err == nil && v == true {
		fmt.Println("Eval", expr)
	}
	switch {
	case expr == nil:
		return nil
	case isBoolean(expr) ||
		isNumber(expr) ||
		isString(expr) ||
		isProcedure(expr) ||
		isSyntax(expr): // self-evaluating?
		return expr
	case isRecur(expr):
		return expr
	case isSymbol(expr): // variables?
		sym, _ := expr.(Symbol)
		if v, err := env.Resolve(sym.GetValue()); err != nil {
			panic(err)
		} else {
			return v
		}
	case isPair(expr):
		car := Eval(env, Car(expr))
		r, ok := car.(*Recur)
		if ok {
			xs := make([]interface{}, 0)
			for _, v := range List2Arr(Cdr(expr)) {
				xs = append(xs, Eval(env, v))
			}
			r.Update(xs)
			return r
		}
		op, ok := car.(Closure)
		if !ok {
			panic("application error, expected SExprOperator, but got " + fmt.Sprintf("%v", car))
		}
		v := Cdr(expr)
		if v, err := env.Resolve("_debug"); err == nil && v == true {
			fmt.Println("op", car, ":", v)
		}
		return op(env, v)
	}
	return nil
}

type enviroment struct {
	variables map[string]interface{}
	outer     Environment
}

func NewEnvFrame(outer Environment) *enviroment {
	e := new(enviroment)
	e.variables = make(map[string]interface{})
	e.outer = outer
	e.Bind("_env", e)
	e.Bind("_outer", outer)
	return e
}

func NewGlobalRootFrame() *enviroment {
	e := NewEnvFrame(nil)
	setupSpecialForms(e)
	setupBuiltins(e)
	e.Bind("eval", Closure(func(env Environment, cdr Pair) interface{} {
		given := Eval(env, Car(cdr)).(Environment)
		given.Dump()
		return Eval(given, Car(Cdr(cdr)))
	}))
	e.Bind("_debug", true)
	e.Bind("_watch", false)
	return e
}

func (env *enviroment) Bind(name string, v interface{}) {
	if found, err := env.Resolve("_watch"); err == nil && found == true {
		fmt.Println("Env", env, name, "got", v)
	}
	env.variables[name] = v
}

func (env *enviroment) Resolve(name string) (interface{}, error) {
	if v, ok := env.variables[name]; ok {
		return v, nil
	}
	if env.outer != nil {
		return env.outer.Resolve(name)
	}
	return nil, EvalError(fmt.Sprintf("Unbound variable %s", name))
}

func (env *enviroment) SetOuter(outer Environment) {
	env.outer = outer
}
func (env *enviroment) Outer() Environment {
	return env.outer
}

func (env *enviroment) Dump() {
	fmt.Println("=====")
	for key, value := range env.variables {
		fmt.Println(key, ":", value)
	}
	if env.outer != nil {
		env.outer.Dump()
	}
}

func Args(p Pair) []interface{} {
	return Car(p).([]interface{})
}

func Body(p Pair) Pair {
	return Cdr(p).(Pair)
}

func NewFrameForApply(lexical Environment, dynamic Environment, args Pair, formals []interface{}) Environment {
	e := NewEnvFrame(lexical)
	as := List2Arr(args)
	for i, name := range formals {
		if len(as) <= i {
			panic("not enough argument")
		}
		if s, ok := name.(Symbol); ok {
			v := Eval(dynamic, as[i])
			e.Bind(s.GetValue(), v)
		} else {
			panic("argument name must be symbol")
		}
	}
	return e
}
