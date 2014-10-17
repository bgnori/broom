package broom

import (
//	"fmt"
)

func setupSpecialForms(env Environment) Environment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind(sym("quote"), func(env Environment, cdr List) interface{} {
		if cdr == nil {
			println("got nil")
			return nil
		}
		return Car(cdr)
	})

	//case sym("set!").Eq(car): //assignment?
	env.Bind(sym("set!"), func(env Environment, cdr List) interface{} {
		panic("not implemented: set!")
		return nil
	})

	//case sym("def").Eq(car): //definition?
	env.Bind(sym("def"), func(env Environment, cdr List) interface{} {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(env, v)
		env.Bind(s, u)
		return s
	})

	//case sym("if").Eq(car): //if?
	env.Bind(sym("if"), func(env Environment, cdr List) interface{} {
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
	env.Bind(sym("fn"), func(lexical Environment, cdr List) interface{} {
		r := func(dynamic Environment, args List) interface{} {
			car := Car(cdr).([]interface{})
			eb := NewEnvBuilder(MakeFromSlice(car...))
			e := eb.EvalAndBindAll(Seq2Slice(args), NewEnvFrame(lexical), dynamic)
			return EvalExprs(e, Body(cdr))
		}
		return r
	})

	env.Bind(sym("let"), func(env Environment, cdr List) interface{} {
		//(let [x 1] ,body)
		xs := Car(cdr).([]interface{})
		// fmt.Println("let", xs)
		seq := MakeFromSlice(xs...)
		//fmt.Println("made seq from xs", SeqString(seq))

		e := NewEnvFrame(env)
		eb := NewEnvBuilder(SeqEvens(seq))
		x := eb.EvalAndBindAll(Seq2Slice(SeqOdds(seq)), e, e)
		return EvalExprs(x, Body(cdr))
	})

	// (loop [i 1] (recur (+ i 1))) ... never stops
	env.Bind(sym("loop"), func(dynamic Environment, cdr List) interface{} {
		r := NewRecur(dynamic, Car(cdr).([]interface{}))
		for {
			x := EvalExprs(r.Env(), Body(cdr))
			_, recuring := x.(*Recur)
			if !recuring {
				return x
			}
		}
		panic("never reach")
	})

	//case sym("begin").Eq(car): //begin?
	env.Bind(sym("begin"), func(env Environment, cdr List) interface{} {
		e := NewEnvFrame(env)
		return EvalExprs(e, cdr)
	})

	//macro
	env.Bind(sym("macro"), func(lexical Environment, cdr List) interface{} {
		m := func(dynamic Environment, args List) interface{} {
			env = NewEnvFrame(lexical)
			env.Bind(sym("_dynamic"), dynamic)
			eb := NewEnvBuilder(MakeFromSlice(Car(cdr).([]interface{})...))
			env = eb.BindAll(Seq2Slice(args), env)
			env.Bind(sym("exprs"), args)

			transformed := EvalExprs(env, Body(cdr))
			//fmt.Println(cdr)
			//fmt.Println("--(macro)-->")
			//fmt.Println(transformed)
			//Our macro is macor object. it won't expand until see it.
			return Eval(dynamic, transformed)
		}
		return m
	})

	// to be implemented
	env.Bind(sym("macroexpand"), func(env Environment, cdr List) interface{} {
		conv := Slice2List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
	})

	//
	env.Bind(sym("and"), and)
	env.Bind(sym("or"), or)
	return env
}

func and(env Environment, cdr List) interface{} {
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

func or(env Environment, cdr List) interface{} {
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
