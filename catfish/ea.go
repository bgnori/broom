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
	foo  func(syn Syntax, env Enviroment, cdr Value) Value
}

func (syn *syntaxImpl) LexEnv() Enviroment {
	return syn.env
}

func (syn *syntaxImpl) Apply(env Enviroment, cdr Value) Value {
	return syn.foo(syn, env, cdr)
}


func makeBuiltinSF(f func(sync Syntax, env Enviroment, cdr Value) Value ) Syntax {
  return &syntaxImpl{foo: f}
}

func setupSpecialForms(env Enviroment) Enviroment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind("quote", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		return Car(cdr)
	}))

	//case sym("set!").Eq(car): //assignment?
	env.Bind("set!", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		panic("not implemented: set!")
		return nil
	}))

	//case sym("define").Eq(car): //definition?
	env.Bind("define", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(v, env)
		env.Bind(s.GetValue(), u)
		return s
	}))

	//case sym("if").Eq(car): //if?
	env.Bind("if", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
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
	env.Bind("lambda", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		return FromLambda(cdr, env)
	}))

	//case sym("begin").Eq(car): //begin?
	env.Bind("begin", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		var x Value
		e := NewEnvFrame()
		e.outer = env
		for _, b := range List2Arr(cdr) {
			x = Eval(b, e)
		}
		return x
	}))

	//case sym("cond").Eq(car): //cond?
	env.Bind("cond", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
		return nil
	}))

        // when macro
        // http://www.shido.info/lisp/scheme_syntax.html
	env.Bind("when", makeBuiltinSF(func(syn Syntax, env Enviroment, cdr Value) Value {
                conv := List(nil, sym("if"), Car(cdr),
                              List(Cdr(cdr), sym("begin")))
		return Eval(conv, env)
	}))
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
