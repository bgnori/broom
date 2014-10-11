package broom

import (
	"fmt"
)

func setupSpecialForms(env Environment) Environment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind("quote", func(env Environment, cdr Pair) interface{} {
		if cdr == nil {
			println("got nil")
			return nil
		}
		return Car(cdr)
	})

	//case sym("set!").Eq(car): //assignment?
	env.Bind("set!", func(env Environment, cdr Pair) interface{} {
		panic("not implemented: set!")
		return nil
	})

	//case sym("def").Eq(car): //definition?
	env.Bind("def", func(env Environment, cdr Pair) interface{} {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(env, v)
		env.Bind(s.GetValue(), u)
		return s
	})

	//case sym("if").Eq(car): //if?
	env.Bind("if", func(env Environment, cdr Pair) interface{} {
		cond := Car(cdr)
		if Eval(env, cond) == true {
			clauseThen := Car(Cdr(cdr))
			return Eval(env, clauseThen)
		} else {
			clauseElse := Car(Cdr(Cdr(cdr)))
			return Eval(env, clauseElse)
		}
		return nil
	})

	//idea from http://clojuredocs.org/clojure_core/clojure.core/fn
	env.Bind("fn", func(lexical Environment, cdr Pair) interface{} {
		r := func(dynamic Environment, args Pair) interface{} {
			eb := NewEnvBuilder(Car(cdr).([]interface{}))
			e := eb.EvalAndBindAll(List2Arr(args), NewEnvFrame(lexical), dynamic)
			return EvalExprs(e, List2Arr(Body(cdr)))
		}
		return r
	})

	env.Bind("let", func(env Environment, cdr Pair) interface{} {
		//(let [x 1] ,body)
		xs := Car(cdr).([]interface{})

		e := NewEnvFrame(env)
		eb := NewEnvBuilder(Evens(xs))
		x := eb.EvalAndBindAll(Odds(xs), e, e)
		return EvalExprs(x, List2Arr(Body(cdr)))
	})

	// (loop [i 1] (recur (+ i 1))) ... never stops
	env.Bind("loop", func(dynamic Environment, cdr Pair) interface{} {
		r := NewRecur(dynamic, Car(cdr).([]interface{}))
		for {
			x := EvalExprs(r.Env(), List2Arr(Body(cdr)))
			_, recuring := x.(*Recur)
			if !recuring {
				return x
			}
		}
		panic("never reach")
	})

	//case sym("begin").Eq(car): //begin?
	env.Bind("begin", func(env Environment, cdr Pair) interface{} {
		e := NewEnvFrame(env)
		return EvalExprs(e, List2Arr(cdr))
	})

	//macro
	env.Bind("macro", func(lexical Environment, cdr Pair) interface{} {
		m := func(dynamic Environment, args Pair) interface{} {
			env = NewEnvFrame(lexical)
			env.Bind("_dynamic", dynamic)
			eb := NewEnvBuilder(Car(cdr).([]interface{}))
			env = eb.BindAll(List2Arr(args), env)
			env.Bind("exprs", args)

			transformed := EvalExprs(env, List2Arr(Body(cdr)))
			fmt.Println(cdr)
			fmt.Println("--(macro)-->")
			fmt.Println(transformed)
			//Our macro is macor object. it won't expand until see it.
			return Eval(dynamic, transformed)
		}
		return m
	})

	// to be implemented
	env.Bind("macroexpand", func(env Environment, cdr Pair) interface{} {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
	})

	//
	env.Bind("and", and)
	env.Bind("or", or)
	return env
}

func and(env Environment, cdr Pair) interface{} {
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

func or(env Environment, cdr Pair) interface{} {
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
