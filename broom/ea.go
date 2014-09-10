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
	e.Bind("eval", Closure(func(env Environment, cdr Pair) Value {
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

func (env *enviroment) SetOuter(outer Environment) {
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

func Args(p Pair) []Value {
	return Car(p).([]Value)
}

func Body(p Pair) Pair {
	return Cdr(p).(Pair)
}

func NewFrameForApply(lexical Environment, dynamic Environment, args Pair, formals []Value) Environment {
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
