package broom

import (
	"fmt"
)

type Environment interface {
	Bind(name string, v Value)
	Resolve(name string) Value
    SetOuter(outer Environment)
	Dump()
}

func Args(p Pair) Pair {
    return Car(p).(Pair)
}

func Body(p Pair) Pair {
    return Cdr(p).(Pair)
}

func FromLambda(cdr Pair, lexical Environment) Closure {
	return func(dynamic Environment, args Pair) Value {
        e := NewFrameForApply(lexical, dynamic, args, Args(cdr))
		var x Value
		for _, b := range List2Arr(Body(cdr)) {
			x = Eval(e, b)
		}
		return x
	}
}

func setupSpecialForms(env Environment) Environment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind("quote", Closure(func(env Environment, cdr Pair) Value {
		return Car(cdr)
	}))

	//case sym("set!").Eq(car): //assignment?
	env.Bind("set!", Closure(func(env Environment, cdr Pair) Value {
		panic("not implemented: set!")
		return nil
	}))

	//case sym("define").Eq(car): //definition?
	env.Bind("define", Closure(func(env Environment, cdr Pair) Value {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(env, v)
		env.Bind(s.GetValue(), u)
		return s
	}))

	//case sym("if").Eq(car): //if?
	env.Bind("if", Closure(func(env Environment, cdr Pair) Value {
		cond := Car(cdr)
		if Eval(env, cond) == true {
			clauseThen := Car(Cdr(cdr))
			return Eval(env, clauseThen)
		} else {
			clauseElse := Car(Cdr(Cdr(cdr)))
			return Eval(env, clauseElse)
		}
		return nil
	}))

	//case sym("lambda").Eq(car): //lambda?
	env.Bind("lambda", Closure(func(env Environment, cdr Pair) Value {
		return FromLambda(cdr, env)
	}))

	//case sym("begin").Eq(car): //begin?
	env.Bind("begin", Closure(func(env Environment, cdr Pair) Value {
		var x Value
		e := NewEnvFrame(env)
		for _, b := range List2Arr(cdr) {
			x = Eval(e, b)
		}
		return x
	}))

	//case sym("cond").Eq(car): //cond?
	env.Bind("cond", Closure(func(env Environment, cdr Pair) Value {
		return nil
	}))

	// when macro
	// http://www.shido.info/lisp/scheme_syntax.html
	env.Bind("when", Closure(func(env Environment, cdr Pair) Value {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
	}))

	// to be implemented
	env.Bind("macroexpand", Closure(func(env Environment, cdr Pair) Value {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
	}))

	//
	env.Bind("and", Closure(and))
	env.Bind("or", Closure(or))
	return env
}

func and(env Environment, cdr Pair) Value {
	v := Eval(env, Car(cdr)).(bool)
	if v {
		next := Cdr(cdr)
		if next == nil {
			return true
		} else {
			return and(env, next)
		}
	}
	return false
}

func or(env Environment, cdr Pair) Value {
	v := Eval(env, Car(cdr)).(bool)
	if !v {
		next := Cdr(cdr)
		if next == nil {
			return false
		} else {
			return or(env, next)
		}
	}
	return true
}

func Eval(env Environment, expr Value) Value {
	/*
	   (define (eval exp env)
	     (cond ((self-evaluating? exp) exp)
	           ((variable? exp) (lookup-variable-value exp env))
	           ((quoted? exp) (text-of-quotation exp))
	           ((assignment? exp) (eval-assignment exp env))
	           ((definition? exp) (eval-definition exp env))
	           ((if? exp) (eval-if exp env))
	           ((lambda? exp)
	            (make-procedure (lambda-parameters exp)
	                            (lambda-body exp)
	                            env))
	           ((begin? exp)
	            (eval-sequence (begin-actions exp) env))
	           ((cond? exp) (eval (cond->if exp) env))
	           ((application? exp)
	            (apply (eval (operator exp) env)
	                   (list-of-values (operands exp) env)))
	           (else
	            (error "Unknown expression type -- EVAL" exp))))
	*/
	switch {
	case isBoolean(expr) || isNumber(expr) || isString(expr) || isProcedure(expr) || isSyntax(expr): // self-evaluating?
		return expr
	case isSymbol(expr): // variables?
		sym, _ := expr.(Symbol)
		return env.Resolve(sym.GetValue())
	case isPair(expr):
		car := Eval(env, Car(expr))
		op, ok := car.(Closure)
		if !ok {
			panic("application error, expected SExprOperator, but got " + fmt.Sprintf("%v", car))
		}
		return op(env, Cdr(expr))
	}
	return nil
}

type enviroment struct {
	variables map[string]Value
	outer     Environment
}

func NewEnvFrame(outer Environment) *enviroment {
	e := new(enviroment)
	e.variables = make(map[string]Value)
	e.outer = outer
	e.Bind("_env", e)
	e.Bind("_outer", outer)
	return e
}

func NewGlobalRootFrame() *enviroment {
	e := NewEnvFrame(nil)
	setupSpecialForms(e)
	setupBuiltins(e)
	e.Bind("eval", Closure(func (env Environment, cdr Pair) Value {
        given := Eval(env, Car(cdr)).(Environment)
        given.Dump()
        return Eval(given, Car(Cdr(cdr)))
    }))
	return e
}

func (env *enviroment) Bind(name string, v Value) {
	env.variables[name] = v
}

func (env *enviroment) Resolve(name string) Value {
	if v, ok := env.variables[name]; ok {
		return v
	}
	if env.outer != nil {
		return env.outer.Resolve(name)
	}
	panic(fmt.Sprintf("Unbound variable %s", name))
	return nil
}

func (env *enviroment) SetOuter (outer Environment) {
    env.outer = outer
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

func NewFrameForApply(lexical Environment, dynamic Environment, args Pair, formals Pair) Environment {
	e := NewEnvFrame(lexical)
    as := List2Arr(args)
    for i, name := range List2Arr(formals) {
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

