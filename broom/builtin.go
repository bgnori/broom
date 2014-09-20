package broom

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
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
	env.Bind("List2Arr", Closure(func(env Environment, args Pair) interface{} {
		list := Eval(env, Car(args))
		return List2Arr(list)
	}))
	env.Bind("list", Closure(func(env Environment, args Pair) interface{} {
		var head, tail Pair

		for _, v := range List2Arr(args) {
			x := Eval(env, v)
			if head == nil && tail == nil {
				head = Cons(x, nil)
				tail = head
			} else {
				tail.SetCdr(Cons(x, nil))
				tail = tail.Cdr()
			}
		}
		return head
	}))
	env.Bind("abs", Closure(func(env Environment, cdr Pair) interface{} {
		v := reflect.ValueOf(Eval(env, Car(cdr)))
		t := v.Type()
		switch t.Kind() {
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Int:
			if v.Int() < 0 {
				return -v.Int()
			}
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			if v.Float() < 0 {
				return -v.Float()
			}
		}
		return v.Interface()
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
	env.Bind("add", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryAdd(x, y)
	}))
	env.Bind("sub", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinarySub(x, y)
	}))
	env.Bind("mul", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryMul(x, y)
	}))
	env.Bind("div", Closure(func(env Environment, cdr Pair) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryDiv(x, y)
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
		first := Eval(env, Car(cdr))
		second := Eval(env, (Car(Cdr(cdr))))
		return BinaryLessThan(first, second)
	}))
	env.Bind(">", Closure(func(env Environment, cdr Pair) interface{} {
		first := Eval(env, Car(cdr))
		second := Eval(env, (Car(Cdr(cdr))))
		return BinaryGreaterThan(first, second)
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
	env.Bind("GOMAXPROCS", Closure(func(env Environment, cdr Pair) interface{} {
		n := Eval(env, Car(cdr)).(int)
		return runtime.GOMAXPROCS(n)
	}))
	env.Bind("select", Closure(func(env Environment, cdr Pair) interface{} {
		cases := make([]reflect.SelectCase, 0)
		headers := make([][]interface{}, 0)
		bodies := make([]interface{}, 0)
		for _, chunk := range Chop2(cdr) {
			select_case := chunk.header
			body := chunk.body
			if v, ok := select_case.(Symbol); ok && v.GetValue() == "default" {
				case_obj := reflect.SelectCase{Dir: reflect.SelectDefault} //, Chan: nil, Send: nil}
				cases = append(cases, case_obj)
				headers = append(headers, nil)
				bodies = append(bodies, body)
			} else {
				dst := (select_case.([]interface{}))[0]
				src := (select_case.([]interface{}))[1]
				if s, ok := dst.(Symbol); ok {
					if v, err := env.Resolve(s.GetValue()); err == nil {
						// must be chan. Sending
						case_obj := reflect.SelectCase{
							Dir:  reflect.SelectSend,
							Chan: v.(reflect.Value),
							Send: Eval(env, src).(reflect.Value)}
						cases = append(cases, case_obj)
						headers = append(headers, select_case.([]interface{}))
						bodies = append(bodies, body)
					} else {
						// s is target symbol. read and assign
						ch := reflect.ValueOf(Eval(env, src))
						case_obj := reflect.SelectCase{
							Dir:  reflect.SelectRecv,
							Chan: ch}
						cases = append(cases, case_obj)
						headers = append(headers, select_case.([]interface{}))
						bodies = append(bodies, body)
					}
				} else {
					panic("Non-symbol in destination of select case.")
				}
			}
		}

		chosen, recv, recvOK := reflect.Select(cases)
		h := headers[chosen]
		b := bodies[chosen]
		benv := NewEnvFrame(env)
		if recvOK {
			benv.Bind(h[0].(Symbol).GetValue(), recv)
		}
		return Eval(benv, b)
	}))
	env.Bind("make-chan-bool", Closure(func(env Environment, cdr Pair) interface{} {
		return make(chan bool)
	}))
	env.Bind("make-chan-string", Closure(func(env Environment, cdr Pair) interface{} {
		return make(chan string)
	}))
	env.Bind("recv", Closure(func(env Environment, cdr Pair) interface{} {
		ch := reflect.ValueOf(Eval(env, Car(cdr)))
		if v, ok := ch.Recv(); ok {
			return v
		}
		return nil
	}))
	env.Bind("send", Closure(func(env Environment, cdr Pair) interface{} {
		ch := reflect.ValueOf(Eval(env, Car(cdr)))
		to_send := reflect.ValueOf(Eval(env, Car(Cdr(cdr))))
		ch.Send(to_send)
		return nil
	}))
	env.Bind("time/Sleep", Closure(func(env Environment, cdr Pair) interface{} {
		//d := Eval(env, Car(cdr)).(time.Duration)
		time.Sleep(1 * time.Second)
		return nil
	}))
	env.Bind("time/Second", time.Second)

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
			s := fmt.Sprintf("object %v (%v) does not have such method: %v", obj, cdr.Car(), name)
			panic(s)
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
