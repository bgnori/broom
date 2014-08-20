package ps

import (
	"fmt"
)

type Value interface {
	Pair() bool
	String() string
}

type Bool bool

func (x Bool) Pair() bool {
	return false
}

func (x Bool) String() string {
	if bool(x) {
		return "#t"
	} else {
		return "#f"
	}
}

type Int int64

func (x Int) Pair() bool {
	return false
}

func (x Int) String() string {
	return fmt.Sprint(int64(x))
}

type Rune rune

func (x Rune) Pair() bool {
	return false
}

type String string

func (x String) Pair() bool {
	return false
}

func (x String) String() string {
	return fmt.Sprint(string(x))
}

type Lambda struct {
	arg  Value
	body Value
	env  *EnvFrame
}

func MakeLambda(arg Value, body Value, env *EnvFrame) *Lambda {
    fmt.Println("made lambda, ", arg)
	return &Lambda{arg: arg, body: body, env: env}
}

func (x *Lambda) Pair() bool {
	return false
}

func (x *Lambda) Apply(cdr Value, e *EnvFrame) Value {
    fmt.Println("evaled lambda", cdr)
	argenv := MakeEnv()
	for z := range Zip(x.arg, cdr) {
		name := z.x.(Name)
		v := Eval(z.y, e)
		argenv.Bind(string(name), v)
	}
	argenv.SetOuter(x.env) //Lexical Scope
    fmt.Println(x.body)
	return Eval(x.body, argenv)
}

func (x *Lambda) String() string {
	return fmt.Sprintf("%x", x)
}

const (
	sfid_if     = iota
	sfid_quote  = iota
	sfid_lambda = iota
)

type SpecialForm struct {
	sfid int
}

func (sf *SpecialForm) String() string {
	switch sf.sfid {
	case sfid_quote:
		return "quote"
	case sfid_if:
		return "if"
	case sfid_lambda:
		return "lambda"
	}
	return "undefined SpecialForm id"
}

func (sf *SpecialForm) Pair() bool {
	return false
}

func (sf *SpecialForm) Eval(cdr Value, e *EnvFrame) Value {
	switch sf.sfid {
	case sfid_quote:
        if Cdr(cdr) != nil {
            panic("got non nil cdr for quote")
        }
		return Car(cdr)
	case sfid_if:
		return nil
	case sfid_lambda:
		return MakeLambda(Car(cdr), Car(Cdr(cdr)), e)
    }
	return nil
}

var SFQuote = &SpecialForm{sfid: sfid_quote}
var SFIf = &SpecialForm{sfid: sfid_if}
var SFLambda = &SpecialForm{sfid: sfid_lambda}

type Pair struct {
	car Value
	cdr Value
}

func Cons(car, cdr Value) *Pair {
	return &Pair{car: car, cdr: cdr}
}

func (p *Pair) Pair() bool {
	return true
}

func (p *Pair) Car() Value {
	return p.car
}

func Car(v Value) Value {
    if p, ok := v.(*Pair) ; ok {
        return p.Car()
    } else {
        panic("not pair")
    }
}

func (p *Pair) Cdr() Value {
	return p.cdr
}

func Cdr(v Value) Value {
    if p, ok := v.(*Pair) ; ok {
        return p.Cdr()
    } else {
        panic("not pair")
    }
}

func MakeList(xs []Value, cdr Value) Value {
    if xs == nil || len(xs) == 0 {
        return cdr
    }
    return Cons(xs[0], MakeList(xs[1:], cdr))
}


func (p *Pair) String() string {
	return fmt.Sprintf("(%v . %v)", p.Car(), p.Cdr())
}

func (pair *Pair) Eval(e *EnvFrame) Value {
	car := pair.Car()
	cdr := pair.Cdr()
	if car == nil && cdr == nil {
		return pair
	}
	if car == nil && cdr != nil {
		panic("car is nil while evaluating pair")
	}
	x := Eval(car, e)
	switch v := x.(type) {
	case *SpecialForm:
		return v.Eval(cdr, e)
	case *Lambda:
		return v.Apply(cdr, e)
	case *Builtin:
		return v.Do(cdr, e)
	}
    fmt.Println("got ", x)
	panic("Bad object in car")
}

type EnvFrame struct {
	bindings map[string]Value
	outer    *EnvFrame
}

func MakeEnv() *EnvFrame {
	return &EnvFrame{bindings: make(map[string]Value), outer: nil}
}

func (e *EnvFrame) Resolve(name string) Value {
	v, ok := e.bindings[name]
	if ok {
		return v
	}
	if e.outer == nil {
		panic(fmt.Sprintf("no such name:%s", name))
	}
	return e.outer.Resolve(name)
}

func (e *EnvFrame) Bind(name string, v Value) {
	e.bindings[name] = v
}

func (e *EnvFrame) SetOuter(outer *EnvFrame) {
	e.outer = outer
}

type Name string

func (n Name) Pair() bool {
	return false
}

func (n Name) String() string {
	return string(n)
}

func Eval(x Value, e *EnvFrame) Value {
	switch v := x.(type) {
	case Name:
		return e.Resolve(string(v))
	case Int:
		return v
	case Bool:
		return v
	case String:
		return v
	case *Pair:
		return v.Eval(e)
    case *SpecialForm:
        return v // expect: "Syntax Special Form"
    case *Lambda:
        return v // expect: "Clousure"
	}
    fmt.Println(x)
    panic("bad object")
	return nil
}
