package catfish

import (
	"fmt"
)

type Enviroment interface {
	Bind(name string, v Value)
	Resolve(name string) Value
}

type clousure struct {
	body Value
	args Value
	env  Enviroment
}

func FromLambda(cdr Value, env Enviroment) *clousure {
	c := new(clousure)
	c.args = Car(cdr)
	c.body = Cdr(cdr)
	c.env = env
	fmt.Println("FromLambda", c)
	return c
}

func (c *clousure) Apply(env Enviroment, cdr Value) Value {
	e := NewEnvFrame()
	e.outer = c.LexEnv()
	formals := List2Arr(c.args)
	for i, a := range List2Arr(cdr) {
		fmt.Println(formals[i], "Eval", a)
		v := Eval(a, e)
		fmt.Println(formals[i], "Bind", v)
		s, _ := formals[i].(Symbol)
		e.Bind(s.GetValue(), v)
	}
	var x Value
	for _, b := range List2Arr(c.body) {
		x = Eval(b, e)
	}
	fmt.Println("Apply got", x)
	return x
}

func (c *clousure) LexEnv() Enviroment {
	return c.env
}

type syntaxImpl struct {
	body Value
	args Value
	env  Enviroment
	foo  func(syn *syntaxImpl, env Enviroment, cdr Value) Value
}

func (s *syntaxImpl) LexEnv() Enviroment {
	return s.env
}

func (syn *syntaxImpl) Apply(env Enviroment, cdr Value) Value {
	return syn.foo(syn, env, cdr)
}

func setupSpecialForms(env Enviroment) Enviroment {
	//case sym("quote").Eq(car): //quoted?
	quote := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		return Car(cdr)
	}}
	env.Bind("quote", quote)

	//case sym("set!").Eq(car): //assignment?
	setbang := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		panic("not implemented: set!")
		return nil
	}}
	env.Bind("set!", setbang)

	//case sym("define").Eq(car): //definition?
	define := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(v, env)
		env.Bind(s.GetValue(), u)
		return s
	}}
	env.Bind("define", define)

	//case sym("if").Eq(car): //if?
	_if := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
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
	}}
	env.Bind("if", _if)

	//case sym("lambda").Eq(car): //lambda?
	lambda := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		return FromLambda(cdr, env)
	}}
	env.Bind("lambda", lambda)

	//case sym("begin").Eq(car): //begin?
	begin := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		var x Value
		e := NewEnvFrame()
		e.outer = env
		for _, b := range List2Arr(cdr) {
			x = Eval(b, e)
		}
		return x
	}}
	env.Bind("begin", begin)

	//case sym("cond").Eq(car): //cond?
	cond := &syntaxImpl{foo: func(syn *syntaxImpl, env Enviroment, cdr Value) Value {
		return nil
	}}
	env.Bind("cond", cond)
	return env
}

func NullEnviroment(version int) Enviroment {
	return nil
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
		op, ok := car.(SExprOperator)
		if !ok {
			panic("application error, expected SExprOperator")
		}
		return op.Apply(env, Cdr(expr))
	}
	return nil
}

type enviroment struct {
	variables map[string]Value
	outer     Enviroment
}

func NewEnvFrame() *enviroment {
	e := new(enviroment)
	e.variables = make(map[string]Value)
	e.outer = nil
	return e
}

func NewGlobalRootFrame() *enviroment {
	e := NewEnvFrame()
	setupSpecialForms(e)
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
