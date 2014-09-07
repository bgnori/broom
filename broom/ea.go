package broom

import (
	"fmt"
)

type Enviroment interface {
	Bind(name string, v Value)
	Resolve(name string) Value
    Dump()
}

func FromLambda(cdr Pair, lexical Enviroment) Closure {
	return func(dynamic Enviroment, args Pair) Value {
		e := NewEnvFrame(lexical)
		formals := List2Arr(Car(cdr))
		for i, a := range List2Arr(args) {
			fmt.Println(formals[i], "Eval", a)
			v := Eval(a, dynamic)
			fmt.Println(formals[i], "Bind", v)
			s, _ := formals[i].(Symbol)
			e.Bind(s.GetValue(), v)
		}
		var x Value
		for _, b := range List2Arr(Cdr(cdr)) {
			x = Eval(b, e)
		}
		return x
	}
}

func setupSpecialForms(env Enviroment) Enviroment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind("quote", Closure(func(env Enviroment, cdr Pair) Value {
		return Car(cdr)
	}))

	//case sym("set!").Eq(car): //assignment?
	env.Bind("set!", Closure(func(env Enviroment, cdr Pair) Value {
		panic("not implemented: set!")
		return nil
	}))

	//case sym("define").Eq(car): //definition?
	env.Bind("define", Closure(func(env Enviroment, cdr Pair) Value {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(v, env)
		env.Bind(s.GetValue(), u)
		return s
	}))

	//case sym("if").Eq(car): //if?
	env.Bind("if", Closure(func(env Enviroment, cdr Pair) Value {
		cond := Car(cdr)
		fmt.Println(cond)
		if Eval(cond, env) == true {
			clauseThen := Car(Cdr(cdr))
			fmt.Println(clauseThen)
			return Eval(clauseThen, env)
		} else {
			clauseElse := Car(Cdr(Cdr(cdr)))
			return Eval(clauseElse, env)
		}
		return nil
	}))

	//case sym("lambda").Eq(car): //lambda?
	env.Bind("lambda", Closure(func(env Enviroment, cdr Pair) Value {
		return FromLambda(cdr, env)
	}))

	//case sym("begin").Eq(car): //begin?
	env.Bind("begin", Closure(func(env Enviroment, cdr Pair) Value {
		var x Value
		e := NewEnvFrame(env)
		for _, b := range List2Arr(cdr) {
			x = Eval(b, e)
		}
		return x
	}))

	//case sym("cond").Eq(car): //cond?
	env.Bind("cond", Closure(func(env Enviroment, cdr Pair) Value {
		return nil
	}))

	// when macro
	// http://www.shido.info/lisp/scheme_syntax.html
	env.Bind("when", Closure(func(env Enviroment, cdr Pair) Value {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(conv, env)
	}))

	// to be implemented
	env.Bind("macroexpand", Closure(func(env Enviroment, cdr Pair) Value {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(conv, env)
	}))
	return env
}

func Eval(expr Value, env Enviroment) Value {
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
		car := Eval(Car(expr), env)
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
	outer     Enviroment
}

func NewEnvFrame(outer Enviroment) *enviroment {
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
	return nil
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
