package broom

import (
	"fmt"
	"sync"
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
		return op(env, v)
	}
	return nil
}

type enviroment struct {
	rwm       sync.RWMutex
	variables map[string]interface{}
	outer     Environment
}

func (env *enviroment) RLock() {
	env.rwm.RLock()
}

func (env *enviroment) RUnlock() {
	env.rwm.RUnlock()
}

func (env *enviroment) Lock() {
	env.rwm.Lock()
}

func (env *enviroment) Unlock() {
	env.rwm.Unlock()
}

func NewEnvFrame(outer Environment) *enviroment {
	e := new(enviroment)
	e.variables = make(map[string]interface{})
	e.SetOuter(outer)
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
	return e
}

func (env *enviroment) Bind(name string, v interface{}) {
	env.Lock()
	defer env.Unlock()
	env.variables[name] = v
}

func (env *enviroment) Resolve(name string) (interface{}, error) {
	env.RLock()
	defer env.RUnlock()
	if v, ok := env.variables[name]; ok {
		return v, nil
	}
	outer := env.Outer()
	if outer != nil {
		return outer.Resolve(name)
	}
	return nil, EvalError(fmt.Sprintf("Unbound variable %s", name))
}

func (env *enviroment) SetOuter(outer Environment) {
	env.Lock()
	defer env.Unlock()
	env.outer = outer
}
func (env *enviroment) Outer() Environment {
	env.RLock()
	defer env.RUnlock()
	return env.outer
}

func (env *enviroment) Dump() {
	env.RLock()
	defer env.RUnlock()
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
			n := s.GetValue()
			if n == "&" {
				if rest, ok := formals[i+1].(Symbol); ok {
					n := rest.GetValue()
					ys := make([]interface{}, 0)
					for _, v := range as[i:] {
						fmt.Println(v)
						ys = append(ys, Eval(dynamic, v))
					}
					fmt.Println(n, ys)
					e.Bind(n, ys)
					break
				} else {
					panic("argument name must be symbol")
				}
			} else {
				v := Eval(dynamic, as[i])
				e.Bind(n, v)
			}
		} else {
			panic("argument name must be symbol")
		}
	}
	return e
}
