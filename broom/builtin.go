package broom

import (
	"fmt"
	"reflect"
)

func setupBuiltins(env Environment) Environment {
	env.Bind("true", true)
	env.Bind("false", false)
	env.Bind("not", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr)).(bool)
		return !x
	}))
	env.Bind("cons", Closure(func(env Environment, body Pair) interface{} {
		car := Eval(env, Car(body))
		cdr, ok := Eval(env, Car(Cdr(body))).(Pair)
		if !ok {
			cdr = nil
		}
		return Cons(car, cdr)
	}))

	env.Bind(".", MakeMethodInvoker())
	env.Bind("=", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return Eq(x, y)
	}))
	env.Bind("mod", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr)).(int)
		y := Eval(env, Car(Cdr(cdr))).(int)
		return x % y
	}))
	env.Bind("+", Closure(func(env Environment, cdr Pair) interface{} {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc += Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("*", Closure(func(env Environment, cdr Pair) interface{} {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc *= Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("-", Closure(func(env Environment, cdr Pair) interface{} {
		xs := List2Arr(Cdr(cdr))
		acc, ok := Eval(env, Car(cdr)).(int)
		if !ok {
			panic("1st arg is not int")
		}
		for _, x := range xs {
			acc -= Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("/", Closure(func(env Environment, cdr Pair) interface{} {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc /= Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("sprintf", Closure(func(env Environment, cdr Pair) interface{} {
		format := Eval(env, Car(cdr)).(string)
		xs := List2Arr(Cdr(cdr))
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		return fmt.Sprintf(format, ys...)
	}))
	env.Bind("println", Closure(func(env Environment, cdr Pair) interface{} {
		fmt.Println(Eval(env, Car(cdr)))
		return nil
	}))
	env.Bind("<", Closure(func(env Environment, cdr Pair) interface{} {
		first := Eval(env, Car(cdr)).(int)
		second := Eval(env, (Car(Cdr(cdr)))).(int)
		return first < second
	}))
	env.Bind(">", Closure(func(env Environment, cdr Pair) interface{} {
		first := Eval(env, Car(cdr)).(int)
		second := Eval(env, (Car(Cdr(cdr)))).(int)
		return first > second
	}))
	env.Bind("null?", Closure(func(env Environment, cdr Pair) interface{} {
		return isNull(Eval(env, Car(cdr)))
	}))
	env.Bind("boolean?", Closure(func(env Environment, cdr Pair) interface{} {
		return isBoolean(Eval(env, Car(cdr)))
	}))
	env.Bind("char?", Closure(func(env Environment, cdr Pair) interface{} {
		return isChar(Eval(env, Car(cdr)))
	}))
	env.Bind("symbol?", Closure(func(env Environment, cdr Pair) interface{} {
		return isSymbol(Eval(env, Car(cdr)))
	}))
	env.Bind("number?", Closure(func(env Environment, cdr Pair) interface{} {
		return isNumber(Eval(env, Car(cdr)))
	}))
	env.Bind("pair?", Closure(func(env Environment, cdr Pair) interface{} {
		return isPair(Eval(env, Car(cdr)))
	}))
	env.Bind("procedure?", Closure(func(env Environment, cdr Pair) interface{} {
		return isProcedure(Eval(env, Car(cdr)))
	}))
	env.Bind("syntax?", Closure(func(env Environment, cdr Pair) interface{} {
		return isSyntax(Eval(env, Car(cdr)))
	}))

	env.Bind("string?", Closure(func(env Environment, cdr Pair) interface{} {
		return isString(Eval(env, Car(cdr)))
	}))

	env.Bind("array?", Closure(func(env Environment, cdr Pair) interface{} {
		return isArray(Eval(env, Car(cdr)))
	}))

	env.Bind("mpa?", Closure(func(env Environment, cdr Pair) interface{} {
		return isMap(Eval(env, Car(cdr)))
	}))

	env.Bind("go", Closure(func(env Environment, cdr Pair) interface{} {
		proc := Eval(env, Car(cdr)).(Closure)
		go proc(env, Cdr(cdr))
		return nil
	}))

	env.Bind("defer", Closure(func(env Environment, cdr Pair) interface{} {
		handler := Eval(env, Car(cdr)).(Closure)
		target := Eval(env, Car(Cdr(cdr))).(Closure)
		return Closure(func(dynamic Environment, arg Pair) interface{} {
			defer func() {
				//fmt.Println("evoking defered", handler)
				handler(dynamic, Cons(1, nil))
				//Eval(dynamic, Cons(handler, nil))
			}()
			return Eval(dynamic, Cons(target, arg))
		})
	}))

	return env
}

func MakeMethodInvoker() Closure {
	return func(env Environment, cdr Pair) interface{} {
		//see  http://stackoverflow.com/questions/14116840/dynamically-call-method-on-interface-regardless-of-receiver-type
		obj := Eval(env, cdr.Car())
		name := cdr.Cdr().Car().(Symbol).GetValue()
		value := reflect.ValueOf(obj)
		method := value.MethodByName(name)

		xs := helper(env, cdr.Cdr().Cdr(), nil)
		if method.IsValid() {
			vs := method.Call(xs)
			i := len(vs)
			if i == 1 {
				return vs[0].Interface()
			} else {
				ys := make([]interface{}, 0, i)
				for _, v := range vs {
					ys = append(ys, v.Interface())
				}
				return List(ys...)
			}
		} else {
			panic("no such method:" + name)
		}
	}
}

func helper(env Environment, args Pair, result []reflect.Value) []reflect.Value {
	if len(result) == 0 {
		result = make([]reflect.Value, 0)
	}
	if args == nil {
		return result
	}
	car := Car(args)
	cdr := Cdr(args)

	v := reflect.ValueOf(Eval(env, car))
	result = append(result, v)

	return helper(env, cdr, result)
}
