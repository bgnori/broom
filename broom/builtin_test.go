package broom

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMethodInvocationA(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(nil), List(sym("quote"), List(1, 2)), sym("Car"))
	r := Eval(e, expr)
	if v, ok := r.(int); !ok || v != 1 {
		t.Error("expected 1")
		fmt.Println(r)
		fmt.Println(ok)
		fmt.Println(v)
	}
}

func TestMethodInvocationB(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(nil), List(sym("quote"), List(1, 2)), sym("String"))
	r := Eval(e, expr)
	if r != "(1 2)" {
		t.Error("expected (1 2)")
		fmt.Println(r)
	}
}

func xTestMethodInvocationC(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(nil), string("abcdef"), sym("At"), 3)
	r := Eval(e, expr)
	if r != 'd' {
		t.Error("expected 'd'")
		fmt.Println(r)
	}
}

func TestPackageFuncInvocation(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(MakeMethodInvoker(MakeReflectPackage()), List(sym("quote"), sym("reflect.TypeOf")), 1)
	r := Eval(e, expr)
	if r != reflect.TypeOf(1) {
		t.Error("expected reflect.Typo(1)")
	}
}

func TestNumPlus(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("+"), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	r := Eval(e, expr)
	if r != 55 {
		t.Error("expected 55")
		fmt.Println(r)
	}
}

func TestNumMul(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("*"), 5, 2)
	r := Eval(e, expr)
	if r != 10 {
		t.Error("expected 10")
		fmt.Println(r)
	}
}

func TestNumMinus(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("-"), 10, 2)
	r := Eval(e, expr)
	if r != 8 {
		t.Error("expected 8")
		fmt.Println(r)
	}
}

func TestNumMinus2(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("-"), 10, 2, 3)
	r := Eval(e, expr)
	if r != 5 {
		t.Error("expected 5")
		fmt.Println(r)
	}
}

func TestNumDiv(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("/"), 10, 2)
	r := Eval(e, expr)
	if r != 5 {
		t.Error("expected 5")
		fmt.Println(r)
	}
}

func TestSprintf(t *testing.T) {
	e := NewGlobalRootFrame()
	expr := List(sym("sprintf"), "Answer is %d", 42)
	r := Eval(e, expr)
	if r != "Answer is 42" {
		t.Error("expected \"Answer is 42\"")
		fmt.Println(r)
	}
}

func TestisNull1(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("null?"), nil)
	r := Eval(e, expr)
	if r != true {
		t.Error("expected true")
		fmt.Println(r)
	}
}

func TestisNull2(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("null?"), 1)
	r := Eval(e, expr)
	if r != false {
		t.Error("expected true")
		fmt.Println(r)
	}
}

func TestisBoolean1(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("boolean?"), false)
	r := Eval(e, expr)
	if r != true {
		t.Error("expected true")
		fmt.Println(r)
	}
}

func TestisBoolean2(t *testing.T) {
	e := NewGlobalRootFrame()

	expr := List(sym("boolean?"), 1)
	r := Eval(e, expr)
	if r != false {
		t.Error("expected true")
		fmt.Println(r)
	}
}
