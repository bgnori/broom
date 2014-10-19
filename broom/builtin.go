package broom

import (
	"fmt"
	"reflect"
	"time"
)

type Injector func(target Sequence) Sequence

func qq(env Environment, x interface{}) interface{} {
	if p, ok := x.(Sequence); ok {
		if s, ok := p.First().(Symbol); ok {
			if s.GetValue() == "uq" {
				return uq(env, Second(p))
			}
		}
		var xs, ys Sequence
		xs = nil
		for ys = p.(Sequence); ys != nil && !ys.IsEmpty(); ys = ys.Rest() {
			v := ys.First()
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
	return func(target Sequence) Sequence {
		v := Eval(env, x)
		if i, ok := v.(Injector); ok {
			return i(target)
		} else {
			return SeqAppend(target, Cons(v, nil))
		}
	}
}

func setupBuiltins(env Environment) Environment {
	env.Bind(sym("true"), true)
	env.Bind(sym("false"), false)
	env.Bind(sym("not"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First()).(bool)
		return !x
	})
	env.Bind(sym("eval"), func(env Environment, cdr Sequence) interface{} {
		given_env := cdr.First().(Environment)
		v := Second(cdr).(interface{})
		return Eval(given_env, v)
	})
	env.Bind(sym("qq"), func(env Environment, cdr Sequence) interface{} {
		// (qq x)
		return qq(env, cdr.First())
	})
	env.Bind(sym("uq"), func(env Environment, cdr Sequence) interface{} {
		// (uq x)
		return uq(env, cdr.First().(Sequence))
	})

	env.Bind(sym("splicing"), func(env Environment, cdr Sequence) interface{} {
		// (splicing xs)
		return Injector(func(target Sequence) Sequence {
			v := Eval(env, cdr.First())
			switch xs := v.(type) {
			case Sequence:
				return SeqAppend(target, xs)
			case []interface{}:
				return SeqAppend(target, Slice2List(xs...))
			default:
				panic(fmt.Sprintf("expected sequence, but got %v", v))
			}
		})
	})

	env.Bind(sym("cons"), func(env Environment, body Sequence) interface{} {
		car := Eval(env, body.First())
		cdr, ok := Eval(env, Second(body)).(Sequence)
		if !ok {
			cdr = nil
		}
		return Cons(car, cdr)
	})
	env.Bind(sym("gensym"), func(env Environment, cdr Sequence) interface{} {
		return GenSym()
	})
	env.Bind(sym("Slice2List"), func(env Environment, args Sequence) interface{} {
		xs := Eval(env, args.First()).([]interface{})
		return Slice2List(xs...)
	})
	env.Bind(sym("Seq2Slice"), func(env Environment, args Sequence) interface{} {
		x := Eval(env, args.First())
		if xs, ok := x.(Sequence); ok {
			return Seq2Slice(xs)
		}
		panic("Non List Value")
	})
	env.Bind(sym("list"), func(env Environment, args Sequence) interface{} {
		var xs Sequence
		ys := make([]interface{}, 0)
		for xs = args.(Sequence); xs != nil && !xs.IsEmpty(); xs = xs.Rest() {
			v := xs.First()
			x := Eval(env, v)
			ys = append(ys, x)
		}
		return Slice2List(ys...)
	})
	env.Bind(sym("abs"), func(env Environment, cdr Sequence) interface{} {
		v := reflect.ValueOf(Eval(env, cdr.First()))
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

	env.Bind(sym("reflect"), MakeReflectPackage())
	env.Bind(sym("os"), MakeOSPackage())
	env.Bind(sym("runtime"), MakeRuntimePackage())
	env.Bind(sym("time"), MakeTimePackage())
	env.Bind(sym("broom"), MakeBroomPackage())
	env.Bind(sym("."), GolangInterop())
	env.Bind(sym("="), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return x == y
	})
	env.Bind(sym("eq?"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return x == y
	})
	env.Bind(sym("equal?"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return reflect.DeepEqual(x, y)
	})
	env.Bind(sym("mod"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First()).(int)
		y := Eval(env, Second(cdr)).(int)
		return x % y
	})
	env.Bind(sym("add"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return BinaryAdd(x, y)
	})
	env.Bind(sym("sub"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return BinarySub(x, y)
	})
	env.Bind(sym("mul"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return BinaryMul(x, y)
	})
	env.Bind(sym("div"), func(env Environment, cdr Sequence) interface{} {
		x := Eval(env, cdr.First())
		y := Eval(env, Second(cdr))
		return BinaryDiv(x, y)
	})
	env.Bind(sym("+"), func(env Environment, cdr Sequence) interface{} {
		return SeqReduce(cdr.First(), func(x, y interface{}) interface{} {
			return BinaryAdd(Eval(env, x), Eval(env, y))
		}, cdr.Rest())
	})
	env.Bind(sym("*"), func(env Environment, cdr Sequence) interface{} {
		return SeqReduce(cdr.First(), func(x, y interface{}) interface{} {
			return BinaryMul(Eval(env, x), Eval(env, y))
		}, cdr.Rest())
	})
	env.Bind(sym("panic"), func(env Environment, cdr Sequence) interface{} {
		panic(Eval(env, cdr.First()))
		return nil
	})
	env.Bind(sym("-"), func(env Environment, cdr Sequence) interface{} {
		return SeqReduce(cdr.First(), func(x, y interface{}) interface{} {
			return BinarySub(Eval(env, x), Eval(env, y))
		}, cdr.Rest())
	})
	env.Bind(sym("/"), func(env Environment, cdr Sequence) interface{} {
		return SeqReduce(cdr.First(), func(x, y interface{}) interface{} {
			return BinaryDiv(Eval(env, x), Eval(env, y))
		}, cdr.Rest())
	})
	env.Bind(sym("sprintf"), func(env Environment, cdr Sequence) interface{} {
		format := Eval(env, cdr.First()).(string)
		xs := Seq2Slice(cdr.Rest())
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		return fmt.Sprintf(format, ys...)
	})
	env.Bind(sym("print"), func(env Environment, cdr Sequence) interface{} {
		xs := Seq2Slice(cdr)
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		fmt.Print(ys...)
		return nil
	})
	env.Bind(sym("println"), func(env Environment, cdr Sequence) interface{} {
		xs := Seq2Slice(cdr)
		ys := make([]interface{}, 0)
		for _, x := range xs {
			ys = append(ys, Eval(env, x))
		}
		fmt.Println(ys...)
		return nil
	})
	env.Bind(sym("<"), func(env Environment, cdr Sequence) interface{} {
		first := Eval(env, cdr.First())
		second := Eval(env, (Second(cdr)))
		return BinaryLessThan(first, second)
	})
	env.Bind(sym(">"), func(env Environment, cdr Sequence) interface{} {
		first := Eval(env, cdr.First())
		second := Eval(env, (Second(cdr)))
		return BinaryGreaterThan(first, second)
	})
	env.Bind(sym("null?"), func(env Environment, cdr Sequence) interface{} {
		return nil == Eval(env, cdr.First())
	})
	env.Bind(sym("boolean?"), func(env Environment, cdr Sequence) interface{} {
		_, ok := Eval(env, cdr.First()).(bool)
		return ok
	})
	env.Bind(sym("rune?"), func(env Environment, cdr Sequence) interface{} {
		v := Eval(env, cdr.First())
		_, ok := v.(rune)
		return ok
	})
	env.Bind(sym("symbol?"), func(env Environment, cdr Sequence) interface{} {
		v := Eval(env, cdr.First())
		_, ok := v.(Symbol)
		return ok
	})
	env.Bind(sym("number?"), func(env Environment, cdr Sequence) interface{} {
		return isNumber(Eval(env, cdr.First()))
	})
	env.Bind(sym("pair?"), func(env Environment, cdr Sequence) interface{} {
		_, ok := Eval(env, cdr.First()).(Sequence)
		return ok
	})
	env.Bind(sym("procedure?"), func(env Environment, cdr Sequence) interface{} {
		v := Eval(env, cdr.First())
		_, ok := v.(func(Environment, Sequence) interface{})
		return ok
	})

	env.Bind(sym("string?"), func(env Environment, cdr Sequence) interface{} {
		_, ok := Eval(env, cdr.First()).(string)
		return ok
	})

	env.Bind(sym("array?"), func(env Environment, cdr Sequence) interface{} {
		return isArray(Eval(env, cdr.First()))
	})

	env.Bind(sym("map?"), func(env Environment, cdr Sequence) interface{} {
		return isMap(Eval(env, cdr.First()))
	})

	env.Bind(sym("go"), func(env Environment, cdr Sequence) interface{} {
		proc := Eval(env, cdr.First()).(func(Environment, Sequence) interface{})
		go proc(env, cdr.Rest())
		return nil
	})

	env.Bind(sym("defer"), func(env Environment, cdr Sequence) interface{} {
		handler := Eval(env, cdr.First()).(func(Environment, Sequence) interface{})
		target := Eval(env, Second(cdr)).(func(Environment, Sequence) interface{})
		return func(dynamic Environment, arg Sequence) interface{} {
			defer func() {
				//fmt.Println("evoking defered", handler)
				handler(dynamic, Cons(1, nil))
				//Eval(dynamic, Cons(handler, nil))
			}()
			return Eval(dynamic, Cons(target, arg))
		}
	})
	env.Bind(sym("bound?"), func(env Environment, cdr Sequence) interface{} {
		if s, ok := cdr.First().(Symbol); ok {
			_, err := env.Resolve(s)
			return err == nil
		}
		return false
	})
	env.Bind(sym("select"), func(env Environment, cdr Sequence) interface{} {
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
					if v, err := env.Resolve(s); err == nil {
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
			benv.Bind(h[0].(Symbol), recv)
		}
		return Eval(benv, b)
	})
	env.Bind(sym("make-chan-bool"), func(env Environment, cdr Sequence) interface{} {
		return make(chan bool)
	})
	env.Bind(sym("make-chan-string"), func(env Environment, cdr Sequence) interface{} {
		return make(chan string)
	})
	env.Bind(sym("recv"), func(env Environment, cdr Sequence) interface{} {
		ch := reflect.ValueOf(Eval(env, cdr.First()))
		if v, ok := ch.Recv(); ok {
			return v
		}
		return nil
	})
	env.Bind(sym("send"), func(env Environment, cdr Sequence) interface{} {
		ch := reflect.ValueOf(Eval(env, cdr.First()))
		to_send := reflect.ValueOf(Eval(env, Second(cdr)))
		ch.Send(to_send)
		return nil
	})
	env.Bind(sym("time/Sleep"), func(env Environment, cdr Sequence) interface{} {
		//d := Eval(env, cdr.First()).(time.Duration)
		time.Sleep(1 * time.Second)
		return nil
	})
	env.Bind(sym("time/Second"), time.Second)

	return env
}

func GolangInterop() func(Environment, Sequence) interface{} {
	return func(env Environment, cdr Sequence) interface{} {
		//see  http://stackoverflow.com/questions/14116840/dynamically-call-method-on-interface-regardless-of-receiver-type

		obj := Eval(env, cdr.First())

		var name string
		var f reflect.Value
		var xs []reflect.Value
		name = Second(cdr).(Symbol).GetValue()

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
					if cdr.Rest().Rest() != nil {
						panic(fmt.Sprintf("args are supplied to member, %s is not method.", name))
					}
					return field
				}
			}
			f = vogus.MethodByName(name)
			//fmt.Println(name, vogus.Interface(), f.Kind())
		}
		xs = helper(env, cdr.Rest().Rest(), nil)

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
			s := fmt.Sprintf("object %v (%v) does not have such field or method: %v", obj, cdr.First(), name)
			panic(s)
		}
	}
}

func helper(env Environment, args Sequence, result []reflect.Value) []reflect.Value {
	if len(result) == 0 {
		result = make([]reflect.Value, 0)
	}
	if args == nil {
		return result
	}
	car := args.First()
	cdr := args.Rest()

	v := reflect.ValueOf(Eval(env, car))
	result = append(result, v)

	return helper(env, cdr, result)
}
