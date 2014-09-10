package broom

import (
	"fmt"
)

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

	env.Bind("lambda", Closure(func(env Environment, cdr Pair) Value {
		//sort of macro
		xs := List2Arr(Car(cdr).(Pair))
		body := Cdr(cdr).(Pair)
		transformed := Cons(sym("fn"), Cons(xs, body))
		fmt.Println("transformed:", transformed)
		return Eval(env, transformed) //fn
	}))

	//http://clojuredocs.org/clojure_core/clojure.core/fn
	env.Bind("fn", Closure(func(lexical Environment, cdr Pair) Value {
		return Closure(func(dynamic Environment, args Pair) Value {
			e := NewFrameForApply(lexical, dynamic, args, Args(cdr))
			var x Value
			for _, b := range List2Arr(Body(cdr)) {
				x = Eval(e, b)
			}
			return x
		})
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
