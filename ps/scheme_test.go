package ps

import (
	"fmt"
	"strings"
	"testing"
)

func TestBool(t *testing.T) {
	if Bool(true) != Bool(true) {
		t.Error("Bad result ")
	}
	if Bool(false) == Bool(true) {
		t.Error("Bad result ")
	}
	if Bool(false) != Bool(false) {
		t.Error("Bad result ")
	}
}

func TestInt(t *testing.T) {
	if Int(1) != Int(1) {
		t.Error("Bad result ")
	}
	if Int(1) == Int(2) {
		t.Error("Bad result ")
	}
	if Int(1) > Int(2) {
		t.Error("Bad result ")
	}
	if Int(2) < Int(1) {
		t.Error("Bad result ")
	}
}

func TestTypeSanity(t *testing.T) {
	var a, b Value
	a = Bool(false)
	b = Int(0)
	if a == b {
		t.Error("Bad result ")
	}
	a = Bool(true)
	b = Int(1)
	if a == b {
		t.Error("Bad result ")
	}
}

func TestEvalInt(t *testing.T) {
	env := MakeEnv()
	expr := Int(3)
	got := Eval(expr, env)
	v, ok := got.(Int)
	if !ok {
		t.Error("Bad result type")
	}
	if v != 3 {
		t.Error("Bad result value")
	}

}

func TestEvalBool(t *testing.T) {
	env := MakeEnv()
	expr := Bool(false)
	got := Eval(expr, env)
	v, ok := got.(Bool)
	if !ok {
		t.Error("Bad result type")
	}
	if v != false {
		t.Error("Bad result value")
	}
}

func TestEvalPair(t *testing.T) {
	env := MakeEnv()
	expr := Cons(nil, nil)
	got := Eval(expr, env)
	v, ok := got.(*Pair)
	if !ok {
		t.Error("Bad result type")
	}
	if v.Car() != nil {
		t.Error("Bad result value, Car()")
	}
	if v.Cdr() != nil {
		t.Error("Bad result value, Cdr()")
	}

}

func TestEvalName(t *testing.T) {
	env := MakeEnv()
	env.Bind("A", Int(42))
	expr := Name("A")
	got := Eval(expr, env)
	v, ok := got.(Int)
	if !ok {
		t.Error("Bad result type")
	}
	if v != 42 {
		t.Error("Bad result value")
	}
}

func TestEvalQuote(t *testing.T) {
	env := MakeEnv()
	expr := Cons(SFQuote, Cons(Cons(Int(1), Int(2)), nil))
	got := Eval(expr, env)
	v, ok := got.(*Pair)
	if !ok {
		t.Error("Bad result type")
	}
	if v.Car() != Int(1) {
		t.Error("Bad result value")
	}
	if v.Cdr() != Int(2) {
		t.Error("Bad result value")
	}
}

func TestBuiltin(t *testing.T) {
	env := MakeEnv()
	env.Bind("+", BuiltinPlus)
	expr := Cons(Name("+"), Cons(Int(1), Cons(Int(2), nil)))
	got := Eval(expr, env)
	v, ok := got.(Int)
	if !ok {
		t.Error("Bad result type")
	}
	if v != Int(3) {
		t.Error("Bad result value")
	}
}

func TestLambda(t *testing.T) {
	println("TestLambda")
	env := MakeEnv()
	env.Bind("+", BuiltinPlus)
	expr := Cons(
		MakeLambda(
			Cons(Name("x"), nil), //Param List
			Cons( //body
				Name("+"),
				Cons(Name("x"),
					Cons(Int(1), nil))),
			env),
		Cons(Int(42), nil))
	got := Eval(expr, env)
	v, ok := got.(Int)
	if !ok {
		t.Error("Bad result type")
	}
	if v != Int(43) {
		t.Error("Bad result value")
	}
}

func TestWithParse(t *testing.T) {
	println("TestWithParse")
	env := MakeEnv()
	env.Bind("+", BuiltinPlus)
    expr := Parse(strings.NewReader("((lambda (x) (+ x 1)) 42)"))[0]

    expected := Cons(
        Cons(
            SFLambda,
            Cons(
                Cons(Name("x"),nil),
                Cons(
                    Cons(Name("+"), Cons(Name("x"), Cons(Int(1), nil))),
                    nil))),
        Cons(Int(42), nil))

    if !RecEq(expected, expr) {
        fmt.Println("expeted: ", expected)
        fmt.Println("got: ", expr)
		t.Error("Bad expr")
    }

	got := Eval(expr, env)
	v, ok := got.(Int)
	if !ok {
		t.Error("Bad result type")
	}
	if v != Int(43) {
		t.Error("Bad result value")
	}
    println("end of TestWithParse")
}

