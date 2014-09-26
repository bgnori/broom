package broom

import (
	"fmt"
)

func setupSpecialForms(env Environment) Environment {
	//case sym("quote").Eq(car): //quoted?
	env.Bind("quote", Closure(func(env Environment, cdr Pair) interface{} {
		if cdr == nil {
			println("got nil")
			return nil
		}
		return Car(cdr)
	}))

	//case sym("set!").Eq(car): //assignment?
	env.Bind("set!", Closure(func(env Environment, cdr Pair) interface{} {
		panic("not implemented: set!")
		return nil
	}))

	//case sym("def").Eq(car): //definition?
	env.Bind("def", Closure(func(env Environment, cdr Pair) interface{} {
		s, _ := Car(cdr).(Symbol)
		v := Car(Cdr(cdr))
		u := Eval(env, v)
		env.Bind(s.GetValue(), u)
		return s
	}))

	//case sym("if").Eq(car): //if?
	env.Bind("if", Closure(func(env Environment, cdr Pair) interface{} {
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

	//idea from http://clojuredocs.org/clojure_core/clojure.core/fn
	env.Bind("fn", Closure(func(lexical Environment, cdr Pair) interface{} {
		r := Closure(func(dynamic Environment, args Pair) interface{} {
			e := NewFrameForApply(lexical, dynamic, args, Args(cdr))
			var x interface{}
			for _, b := range List2Arr(Body(cdr)) {
				x = Eval(e, b)
			}
			return x
		})
		return r
	}))

	env.Bind("let", Closure(func(env Environment, cdr Pair) interface{} {
		//(let [x 1] ,body)
		e := NewEnvFrame(env)
		xs := Car(cdr).([]interface{})

		for i:= 0 ; i < len(xs) ; i+=2 {
			name := xs[i].(Symbol)
			value := xs[i+1]
			e.Bind(name.GetValue(), Eval(e, value))
		}
		var x interface{}
		for _, b := range List2Arr(Body(cdr)) {
			x = Eval(e, b)
		}
		return x
	}))

	// (loop [i 1] (recur (+ i 1))) ... never stops
	env.Bind("loop", Closure(func(dynamic Environment, cdr Pair) interface{} {
		r := NewRecur(dynamic, Car(cdr))
		var x interface{}
		for {
			for _, b := range List2Arr(Body(cdr)) {
				x = Eval(r.Env(), b)
			}
			_, recuring := x.(*Recur)
			if !recuring {
				return x
			}
		}
		panic("never reach")
	}))

	//case sym("begin").Eq(car): //begin?
	env.Bind("begin", Closure(func(env Environment, cdr Pair) interface{} {
		var x interface{}
		e := NewEnvFrame(env)
		for _, b := range List2Arr(cdr) {
			x = Eval(e, b)
		}
		return x
	}))

	//case sym("cond").Eq(car): //cond?
	env.Bind("cond", Closure(func(env Environment, cdr Pair) interface{} {
		test := Car(cdr)
		onTrue := Car(Cdr(cdr))
		rest := Cdr(Cdr(cdr))
		if v, ok := test.(Symbol); ok && v.GetValue() == "else" {
			return Eval(env, onTrue)
		}
		if Eval(env, test) == true {
			return Eval(env, onTrue)
		}
		if rest != nil {
			return Eval(env, Cons(sym("cond"), rest))
			// ugh! rewind as for loop!
			// * by using Chop2
			// * by macro
		}
		return nil //undef
	}))

	//macro
	env.Bind("macro", Closure(func(lexical Environment, cdr Pair) interface{} {
		m := Closure(func(dynamic Environment, args Pair) interface{} {
			env = NewEnvFrame(lexical)
			env.Bind("_dynamic", dynamic)
			as := List2Arr(args)
			for i, v := range Args(cdr) {
				s := v.(Symbol)
				env.Bind(s.GetValue(), as[i])
			}
			env.Bind("exprs", args)

			var transformed interface{}
			for _, b := range List2Arr(Body(cdr)) {
				transformed = Eval(env, b)
			}
			fmt.Println(cdr)
			fmt.Println("--(macro)-->")
			fmt.Println(transformed)
			//Our macro is macor object. it won't expand until see it.
			return Eval(dynamic, transformed)
		})
		return m
	}))

	// to be implemented
	env.Bind("macroexpand", Closure(func(env Environment, cdr Pair) interface{} {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
	}))

	//
	env.Bind("and", Closure(and))
	env.Bind("or", Closure(or))
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
