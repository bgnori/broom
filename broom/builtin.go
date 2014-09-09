package broom

import (
	"fmt"
	"reflect"
)

func setupBuiltins(env Environment) Environment {
	env.Bind("true", true)
	env.Bind("false", false)
	env.Bind(".", MakeMethodInvoker())
	env.Bind("+", Closure(func(env Environment, cdr Pair) Value {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc += Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("*", Closure(func(env Environment, cdr Pair) Value {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc *= Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("-", Closure(func(env Environment, cdr Pair) Value {
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
	env.Bind("/", Closure(func(env Environment, cdr Pair) Value {
		xs := List2Arr(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc /= Eval(env, x).(int)
		}
		return acc
	}))
	env.Bind("sprintf", Closure(func(env Environment, cdr Pair) Value {
		format := Car(cdr).(string)
		xs := List2Arr(Cdr(cdr))
		return fmt.Sprintf(format, xs...)
	}))
	env.Bind("println", Closure(func(env Environment, cdr Pair) Value {
		fmt.Println(Car(cdr))
		return nil
	}))
	env.Bind("<", Closure(func(env Environment, cdr Pair) Value {
		first := Eval(env, Car(cdr)).(int)
		second := Eval(env, (Car(Cdr(cdr)))).(int)
		return first < second
	}))
	env.Bind(">", Closure(func(env Environment, cdr Pair) Value {
		first := Eval(env, Car(cdr)).(int)
		second := Eval(env, (Car(Cdr(cdr)))).(int)
		return first > second
	}))
	return env
}

func MakeMethodInvoker() Closure {
	return func(env Environment, cdr Pair) Value {
		//see  http://stackoverflow.com/questions/14116840/dynamically-call-method-on-interface-regardless-of-receiver-type
		obj := Eval(env, cdr.Car())
		fmt.Println("obj: ", obj)
		name := cdr.Cdr().Car().(Symbol).GetValue()
		fmt.Println("to invoke:", name)
		xs := helper(cdr.Cdr().Cdr(), nil)

		value := reflect.ValueOf(obj)
		method := value.MethodByName(name)
		if method.IsValid() {
			vs := method.Call(xs)
			i := len(vs)
			if i == 1 {
				return vs[0].Interface()
			} else {
				ys := make([]Value, 0, i)
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

func helper(args Pair, result []reflect.Value) []reflect.Value {
	if len(result) == 0 {
		result = make([]reflect.Value, 0)
	}
	if args == nil {
		return result
	}
	car := Car(args)
	cdr := Cdr(args)

	v := reflect.ValueOf(car)
	result = append(result, v)

	return helper(cdr, result)
}
