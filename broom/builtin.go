package broom

import (
	"fmt"
	"reflect"
)


func setupBuiltins(env Enviroment) Enviroment {
	env.Bind(".", MakeMethodInvoker)
	env.Bind("+", Closure(func (env Enviroment, cdr Pair) Value {
        xs := List2Arr(Cdr(cdr))
        acc := Car(cdr).(int)
        for _, x := range xs {
            acc += Eval(x, env).(int)
        }
        return acc
    }))
	env.Bind("*", Closure(func (env Enviroment, cdr Pair) Value {
        xs := List2Arr(Cdr(cdr))
        acc := Car(cdr).(int)
        for _, x := range xs {
            acc *= Eval(x, env).(int)
        }
        return acc
    }))
	env.Bind("-", Closure(func (env Enviroment, cdr Pair) Value {
        xs := List2Arr(Cdr(cdr))
        acc := Car(cdr).(int)
        for _, x := range xs {
            acc -= Eval(x, env).(int)
        }
        return acc
    }))
	env.Bind("/", Closure(func (env Enviroment, cdr Pair) Value {
        xs := List2Arr(Cdr(cdr))
        acc := Car(cdr).(int)
        for _, x := range xs {
            acc /= Eval(x, env).(int)
        }
        return acc
    }))
    return env
}


func MakeMethodInvoker() Closure {
	return func(env Enviroment, cdr Pair) Value {
		//see  http://stackoverflow.com/questions/14116840/dynamically-call-method-on-interface-regardless-of-receiver-type
		obj := cdr.Car()
		fmt.Println(obj)
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


func GetFunc() Closure {
    return func(env Enviroment, cdr Pair) Value {
        name := Car(cdr)
        fmt.Println(name)

        return nil
    }
}
