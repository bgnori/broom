package broom

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

	//case sym("define").Eq(car): //definition?
	env.Bind("define", Closure(func(env Environment, cdr Pair) interface{} {
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

	env.Bind("lambda", Closure(func(env Environment, cdr Pair) interface{} {
		//sort of macro
		xs := List2Arr(Car(cdr).(Pair))
		body := Cdr(cdr).(Pair)
		transformed := Cons(sym("fn"), Cons(xs, body))
		return Eval(env, transformed) //fn
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

	// (loop [i 1] (recur (+ i 1))) ... never stops
	env.Bind("loop", Closure(func(dynamic Environment, cdr Pair) interface{} {
		r := NewRecur(dynamic, Car(cdr))
		var x interface{}
		for recuring := true; recuring; {
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

	//idea from http://clojuredocs.org/clojure_core/clojure.core/defn
	// (defn name [params*] body)
	env.Bind("defn", Closure(func(env Environment, cdr Pair) interface{} {
		//sort of macro
		name := Car(cdr)
		xs := Car(Cdr(cdr))
		body := Cdr(Cdr(cdr))
		transformed := List(sym("define"), name, Cons(sym("fn"), Cons(xs, body)))
		return Eval(env, transformed)
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

	// when macro
	// http://www.shido.info/lisp/scheme_syntax.html
	env.Bind("when", Closure(func(env Environment, cdr Pair) interface{} {
		conv := List(sym("if"), Car(cdr),
			Cons(sym("begin"), Cdr(cdr)))
		return Eval(env, conv)
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
