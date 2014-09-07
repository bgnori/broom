package broom

import (
	"fmt"
	"testing"
)

func TestMethodInvocationA(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(), Cons(1, nil), sym("Car"))
	r := Eval(expr, e)
	if v, ok := r.(int); !ok || v != 1 {
		t.Error("expected 1")
		fmt.Println(r)
		fmt.Println(ok)
		fmt.Println(v)
	}
}

func TestMethodInvocationB(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(), Cons(1, nil), sym("String"))
	r := Eval(expr, e)
	if r != "(1)" {
		t.Error("expected (1)")
		fmt.Println(r)
	}
}

func xTestMethodInvocationC(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(), string("abcdef"), sym("At"), 3)
	r := Eval(expr, e)
	if r != 'd' {
		t.Error("expected 'd'")
		fmt.Println(r)
	}
}

func TestNumPlus(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("+"), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	r := Eval(expr, e)
    if r != 55 {
		t.Error("expected 55")
		fmt.Println(r)
	}
}

func TestNumMul(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("*"), 5, 2)
	r := Eval(expr, e)
    if r != 10 {
		t.Error("expected 10")
		fmt.Println(r)
	}
}

func TestNumMinus(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("-"), 10, 2)
	r := Eval(expr, e)
    if r != 8 {
		t.Error("expected 8")
		fmt.Println(r)
	}
}

func TestNumDiv(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("/"), 10, 2)
	r := Eval(expr, e)
    if r != 5 {
		t.Error("expected 5")
		fmt.Println(r)
	}
}

func TestSprintf(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(sym("sprintf"), "Answer is %d", 42)
	r := Eval(expr, e)
    if r != "Answer is 42" {
		t.Error("expected \"Answer is 42\"")
		fmt.Println(r)
	}
}
