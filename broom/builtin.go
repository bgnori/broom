package broom

import (
	"fmt"
	"reflect"
	"time"
)

type Injector func(target Sequence) Sequence

func qq(env Environment, x interface{}) interface{} {
	if p, ok := x.(List); ok {
		if s, ok := Car(p).(Symbol); ok {
			if s.GetValue() == "uq" {
				return uq(env, Car(Cdr(p)))
			}
		}
		var xs Sequence
		xs = nil
		for _, v := range List2Slice(p) {
			q := qq(env, v)
			if i, ok := q.(Injector); ok {
				xs = i(xs)
			} else {
				xs = SeqAppend(xs, Cons(q, nil))
			}
		}
		return xs
	} else {
		return x
	}
}

func uq(env Environment, x interface{}) Injector {
	return func(target Sequence) Sequence{
		fmt.Println("uq:", x)
		v := Eval(env, x)
		if i, ok := v.(Injector); ok {
			return i(target)
		} else {
			return SeqAppend(target, Cons(v, nil))
		}
	}
}

func setupBuiltins(env Environment) Environment {
	env.Bind("true", true)
	env.Bind("false", false)
	env.Bind("not", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr)).(bool)
		return !x
	})
	env.Bind("eval", func(env Environment, cdr List) interface{} {
		given_env := Car(cdr).(Environment)
		v := Car(Cdr(cdr)).(interface{})
		fmt.Println(given_env, v)
		return Eval(given_env, v)
	})
	env.Bind("qq", func(env Environment, cdr List) interface{} {
		// (qq x)
		return qq(env, Car(cdr))
	})
	env.Bind("uq", func(env Environment, cdr List) interface{} {
		// (uq x)
		return uq(env, Car(cdr).(List))
	})

	env.Bind("splicing", func(env Environment, cdr List) interface{} {
		// (splicing xs)
		return Injector(func(target Sequence) Sequence{
			v := Eval(env, Car(cdr))
			switch xs := v.(type) {
			case List:
				return SeqAppend(target, xs)
			case []interface{}:
				return SeqAppend(target, Slice2List(xs...))
			default:
				panic(fmt.Sprintf("expected sequence, but got %v", v))
			}
		})
	})

	env.Bind("cons", func(env Environment, body List) interface{} {
		car := Eval(env, Car(body))
		cdr, ok := Eval(env, Car(Cdr(body))).(List)
		if !ok {
			cdr = nil
		}
		return Cons(car, cdr)
	})
	env.Bind("gensym", func(env Environment, cdr List) interface{} {
		return GenSym()
	})
	env.Bind("Slice2List", func(env Environment, args List) interface{} {
		xs := Eval(env, Car(args)).([]interface{})
		return Slice2List(xs...)
	})
	env.Bind("List2Slice", func(env Environment, args List) interface{} {
		x := Eval(env, Car(args))
		if xs, ok := x.(List) ; ok {
			return List2Slice(xs)
		}
		panic("Non List Value")
	})
	env.Bind("list", func(env Environment, args List) interface{} {
		var head, tail List

		for _, v := range List2Slice(args) {
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
	})
	env.Bind("abs", func(env Environment, cdr List) interface{} {
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
	})

	env.Bind("reflect", MakeReflectPackage())
	env.Bind("os", MakeOSPackage())
	env.Bind("runtime", MakeRuntimePackage())
	env.Bind("broom", MakeBroomPackage())
	env.Bind(".", GolangInterop())
	env.Bind("=", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return x == y
	})
	env.Bind("eq?", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return x == y
	})
	env.Bind("equal?", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return reflect.DeepEqual(x, y)
	})
	env.Bind("mod", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr)).(int)
		y := Eval(env, Car(Cdr(cdr))).(int)
		return x % y
	})
	env.Bind("add", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryAdd(x, y)
	})
	env.Bind("sub", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinarySub(x, y)
	})
	env.Bind("mul", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryMul(x, y)
	})
	env.Bind("div", func(env Environment, cdr List) interface{} {
		x := Eval(env, Car(cdr))
		y := Eval(env, Car(Cdr(cdr)))
		return BinaryDiv(x, y)
	})
	env.Bind("+", func(env Environment, cdr List) interface{} {
		xs := List2Slice(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc += Eval(env, x).(int)
		}
		return acc
	})
	env.Bind("*", func(env Environment, cdr List) interface{} {
		xs := List2Slice(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc *= Eval(env, x).(int)
		}
		return acc
	})
	env.Bind("panic", func(env Environment, cdr List) interface{} {
		panic(Eval(env, Car(cdr)))
		return nil
	})
	env.Bind("-", func(env Environment, cdr List) interface{} {
		xs := List2Slice(Cdr(cdr))
		acc, ok := Eval(env, Car(cdr)).(int)
		if !ok {
			panic("1st arg is not int")
		}
		for _, x := range xs {
			acc -= Eval(env, x).(int)
		}
		return acc
	})
	env.Bind("/", func(env Environment, cdr List) interface{} {
		xs := List2Slice(Cdr(cdr))
		acc := Eval(env, Car(cdr)).(int)
		for _, x := range xs {
			acc /= Eval(env, x).(int)
		}
		return acc
	})
	env.Bind("sprintf", func(env Environment, cdr List) interface{} {
		format := Eval(env, Car(cdr)).(string)
		xs := List2Slice(Cdr(cdr))
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		return fmt.Sprintf(format, ys...)
	})
	env.Bind("println", func(env Environment, cdr List) interface{} {
		xs := List2Slice(cdr)
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		fmt.Println(ys...)
		return nil
	})
	env.Bind("<", func(env Environment, cdr List) interface{} {
		first := Eval(env, Car(cdr))
		second := Eval(env, (Car(Cdr(cdr))))
		return BinaryLessThan(first, second)
	})
	env.Bind(">", func(env Environment, cdr List) interface{} {
		first := Eval(env, Car(cdr))
		second := Eval(env, (Car(Cdr(cdr))))
		return BinaryGreaterThan(first, second)
	})
	env.Bind("null?", func(env Environment, cdr List) interface{} {
		return nil == Eval(env, Car(cdr))
	})
	env.Bind("boolean?", func(env Environment, cdr List) interface{} {
		_, ok := Eval(env, Car(cdr)).(bool)
		return ok
	})
	env.Bind("rune?", func(env Environment, cdr List) interface{} {
		v := Eval(env, Car(cdr))
		_, ok := v.(rune)
		return ok
	})
	env.Bind("symbol?", func(env Environment, cdr List) interface{} {
		v := Eval(env, Car(cdr))
		_, ok := v.(Symbol)
		return ok
	})
	env.Bind("number?", func(env Environment, cdr List) interface{} {
		return isNumber(Eval(env, Car(cdr)))
	})
	env.Bind("pair?", func(env Environment, cdr List) interface{} {
		_, ok := Eval(env, Car(cdr)).(List)
		return ok
	})
	env.Bind("procedure?", func(env Environment, cdr List) interface{} {
		v := Eval(env, Car(cdr))
		_, ok := v.(func(Environment, List) interface{})
		return ok
	})

	env.Bind("string?", func(env Environment, cdr List) interface{} {
		_, ok := Eval(env, Car(cdr)).(string)
		return ok
	})

	env.Bind("array?", func(env Environment, cdr List) interface{} {
		return isArray(Eval(env, Car(cdr)))
	})

	env.Bind("map?", func(env Environment, cdr List) interface{} {
		return isMap(Eval(env, Car(cdr)))
	})

	env.Bind("go", func(env Environment, cdr List) interface{} {
		proc := Eval(env, Car(cdr)).(func(Environment, List) interface{})
		go proc(env, Cdr(cdr))
		return nil
	})

	env.Bind("defer", func(env Environment, cdr List) interface{} {
		handler := Eval(env, Car(cdr)).(func(Environment, List) interface{})
		target := Eval(env, Car(Cdr(cdr))).(func(Environment, List) interface{})
		return func(dynamic Environment, arg List) interface{} {
			defer func() {
				//fmt.Println("evoking defered", handler)
				handler(dynamic, Cons(1, nil))
				//Eval(dynamic, Cons(handler, nil))
			}()
			return Eval(dynamic, Cons(target, arg))
		}
	})
	env.Bind("bound?", func(env Environment, cdr List) interface{} {
		if s, ok := Car(cdr).(Symbol); ok {
			_, err := env.Resolve(s.GetValue())
			return err == nil
		}
		return false
	})
	env.Bind("select", func(env Environment, cdr List) interface{} {
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
	})
	env.Bind("make-chan-bool", func(env Environment, cdr List) interface{} {
		return make(chan bool)
	})
	env.Bind("make-chan-string", func(env Environment, cdr List) interface{} {
		return make(chan string)
	})
	env.Bind("recv", func(env Environment, cdr List) interface{} {
		ch := reflect.ValueOf(Eval(env, Car(cdr)))
		if v, ok := ch.Recv(); ok {
			return v
		}
		return nil
	})
	env.Bind("send", func(env Environment, cdr List) interface{} {
		ch := reflect.ValueOf(Eval(env, Car(cdr)))
		to_send := reflect.ValueOf(Eval(env, Car(Cdr(cdr))))
		ch.Send(to_send)
		return nil
	})
	env.Bind("time/Sleep", func(env Environment, cdr List) interface{} {
		//d := Eval(env, Car(cdr)).(time.Duration)
		time.Sleep(1 * time.Second)
		return nil
	})
	env.Bind("time/Second", time.Second)

	return env
}

func GolangInterop() func(Environment, List) interface{} {
	return func(env Environment, cdr List) interface{} {
		//see  http://stackoverflow.com/questions/14116840/dynamically-call-method-on-interface-regardless-of-receiver-type

		obj := Eval(env, cdr.Car())

		var name string
		var f reflect.Value
		var xs []reflect.Value
		name = cdr.Cdr().Car().(Symbol).GetValue()

		if proxy, ok := obj.(*PackageProxy); ok {
			//fmt.Println("getting from proxy, ", name)
			f, ok = proxy.Find(name)
			if !ok {
				panic(fmt.Sprintf("package %s do not have %s", proxy.Name(), name))
			}
		} else {
			// anything else is instance.
			vogus := reflect.ValueOf(obj)

			//method or value?
			if vogus.Kind() == reflect.Struct {
				if field := vogus.FieldByName(name); field.IsValid() {
					if cdr.Cdr().Cdr() != nil {
						panic(fmt.Sprintf("args are supplied to member, %s is not method.", name))
					}
					return field
				}
			}
			f = vogus.MethodByName(name)
			//fmt.Println(name, vogus.Interface(), f.Kind())
		}
		xs = helper(env, cdr.Cdr().Cdr(), nil)

		if f.IsValid() {
			vs := f.Call(xs)
			i := len(vs)
			if i == 1 {
				return vs[0].Interface()
			} else {
				ys := make([]interface{}, 0, i)
				for _, v := range vs {
					ys = append(ys, v.Interface())
				}
				return Slice2List(ys...)
			}
		} else {
			s := fmt.Sprintf("object %v (%v) does not have such field or method: %v", obj, cdr.Car(), name)
			panic(s)
		}
	}
}

func helper(env Environment, args List, result []reflect.Value) []reflect.Value {
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
