package broom

import (
//	"fmt"
)

func setupSpecialForms(env Environment) Environment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind(sym("quote"), func(env Environment, cdr Sequence) interface{} {
		if cdr == nil {
			println("got nil")
			return nil
		}
		return cdr.First()
	})

	//case sym("set!").Eq(car): //assignment?
	env.Bind(sym("set!"), func(env Environment, cdr Sequence) interface{} {
		panic("not implemented: set!")
		return nil
	})

	//case sym("def").Eq(car): //definition?
	env.Bind(sym("def"), func(env Environment, cdr Sequence) interface{} {
		s, _ := cdr.First().(Symbol)
		v := Second(cdr)
		u := Eval(env, v)
		env.Bind(s, u)
		return s
	})

	//case sym("if").Eq(car): //if?
	env.Bind(sym("if"), func(env Environment, cdr Sequence) interface{} {
		cond := cdr.First()
		if Eval(env, cond) == true {
			clauseThen := Second(cdr)
			return Eval(env, clauseThen)
		} else {
			clauseElse := Second(cdr.Rest())
			return Eval(env, clauseElse)
		}
		return nil
	})

	//idea from http://clojuredocs.org/clojure_core/clojure.core/fn
	env.Bind(sym("fn"), func(lexical Environment, cdr Sequence) interface{} {
		r := func(dynamic Environment, args Sequence) interface{} {
			car := cdr.First().([]interface{})
			eb := NewEnvBuilder(MakeFromSlice(car...))
			e := eb.EvalAndBindAll(Seq2Slice(args), NewEnvFrame(lexical), dynamic)
			return EvalExprs(e, Body(cdr))
		}
		return r
	})

	env.Bind(sym("let"), func(env Environment, cdr Sequence) interface{} {
		//(let [x 1] ,body)
		xs := cdr.First().([]interface{})
		// fmt.Println("let", xs)
		seq := MakeFromSlice(xs...)
		//fmt.Println("made seq from xs", SeqString(seq))

		e := NewEnvFrame(env)
		eb := NewEnvBuilder(SeqEvens(seq))
		x := eb.EvalAndBindAll(Seq2Slice(SeqOdds(seq)), e, e)
		return EvalExprs(x, Body(cdr))
	})

	// (loop [i 1] (recur (+ i 1))) ... never stops
	env.Bind(sym("loop"), func(dynamic Environment, cdr Sequence) interface{} {
		r := NewRecur(dynamic, cdr.First().([]interface{}))
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
	env.Bind(sym("begin"), func(env Environment, cdr Sequence) interface{} {
		e := NewEnvFrame(env)
		return EvalExprs(e, cdr)
	})

	//macro
	env.Bind(sym("macro"), func(lexical Environment, cdr Sequence) interface{} {
		m := func(dynamic Environment, args Sequence) interface{} {
			env = NewEnvFrame(lexical)
			env.Bind(sym("_dynamic"), dynamic)
			eb := NewEnvBuilder(MakeFromSlice(cdr.First().([]interface{})...))
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
	env.Bind(sym("macroexpand"), func(env Environment, cdr Sequence) interface{} {
		conv := Slice2List(sym("if"), cdr.First(),
			Cons(sym("begin"), cdr.Rest()))
		return Eval(env, conv)
	})

	//
	env.Bind(sym("and"), and)
	env.Bind(sym("or"), or)
	return env
}

func and(env Environment, cdr Sequence) interface{} {
	v := Eval(env, cdr.First()).(bool)
	if v {
		next := cdr.Rest()
		if next == nil {
			return true
		} else {
			return and(env, next)
		}
	}
	return false
}

func or(env Environment, cdr Sequence) interface{} {
	v := Eval(env, cdr.First()).(bool)
	if !v {
		next := cdr.Rest()
		if next == nil {
			return false
		} else {
			return or(env, next)
		}
	}
	return true
}
