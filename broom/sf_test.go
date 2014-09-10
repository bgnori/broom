package broom

import (
	"fmt"
	"testing"
)

func TestEvalQuotedSymbol(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("quote"), sym("A")))
	if !sym("A").Eq(v) {
		t.Error("expected sym A")
		fmt.Println(v)
	}
}

func TestEvalDefine(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("define"), sym("A"), 42))
	if e.Resolve("A") != 42 {
		t.Error("expected 42")
	}
	if !sym("A").Eq(v) {
		t.Error("expected sym A")
	}
}

func TestEvalIF(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("if"), true, 42, 123))
	if v != 42 {
		t.Error("expected 42")
	}
	v = Eval(e, List(sym("if"), false, 42, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvalLambda(t *testing.T) {
	e := NewGlobalRootFrame()
	f := Eval(e, List(sym("lambda"), List(sym("x")), sym("x")))
	if _, ok := f.(Closure); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvalfn(t *testing.T) {
	e := NewGlobalRootFrame()
	f := Eval(e, List(sym("fn"), []interface{}{sym("x")}, sym("x")))
	if _, ok := f.(Closure); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvaldefn(t *testing.T) {
	e := NewGlobalRootFrame()
	Eval(e, List(sym("defn"), sym("foo"), []interface{}{sym("x")}, sym("x")))
	f := Eval(e, sym("foo"))
	if _, ok := f.(Closure); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvalWhen(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("when"), true, 1, 2, 3))
	if v != 3 {
		t.Error("expected 3")
		fmt.Println(v)
	}
}

func TestEvalEnvSymbol(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, sym("_env"))
	if v != e {
		t.Error("expected enviroment object")
	}
}

func TestEvalAnd1(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("and"), true, true, false))
	if v != false {
		t.Error("expected false")
		fmt.Println(v)
	}
}

func TestEvalAnd2(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("and"), true, true, true))
	if v != true {
		t.Error("expected true")
		fmt.Println(v)
	}
}

func TestEvalOr1(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("or"), false, false, true))
	if v != true {
		t.Error("expected true")
		fmt.Println(v)
	}
}

func TestEvalOr2(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, List(sym("or"), false, false, false))
	if v != false {
		t.Error("expected false")
		fmt.Println(v)
	}
}
