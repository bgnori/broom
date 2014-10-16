package broom

import (
	"fmt"
	"testing"
)

func TestEvalQuotedSymbol(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("quote"), sym("A")))
	if sym("A") != v {
		t.Error("expected sym A")
		fmt.Println(v)
	}
}

func TestEvalDefine(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("def"), sym("A"), 42))
	if found, err := e.Resolve(sym("A")); err != nil || found != 42 {
		t.Error("expected 42")
	}
	if sym("A") != v {
		t.Error("expected sym A")
	}
}

func TestEvalIF(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("if"), true, 42, 123))
	if v != 42 {
		t.Error("expected 42")
	}
	v = Eval(e, Slice2List(sym("if"), false, 42, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvalLambda(t *testing.T) {
	e := NewGlobalRootFrame()
	f := Eval(e, Slice2List(sym("fn"), []interface{}{sym("x")}, sym("x")))
	if _, ok := f.(func(Environment, List) interface{}); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, Slice2List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func TestEvalfn(t *testing.T) {
	e := NewGlobalRootFrame()
	f := Eval(e, Slice2List(sym("fn"), []interface{}{sym("x")}, sym("x")))
	if _, ok := f.(func(Environment, List) interface{}); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, Slice2List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func xTestEvaldefn(t *testing.T) {
	//defn is macro now
	e := NewGlobalRootFrame()
	Eval(e, Slice2List(sym("defn"), sym("foo"), []interface{}{sym("x")}, sym("x")))
	f := Eval(e, sym("foo"))
	if _, ok := f.(func(Environment, List) interface{}); !ok {
		t.Error("expected Procedure")
	}
	v := Eval(e, Slice2List(f, 123))
	if v != 123 {
		t.Error("expected 123")
	}
}

func xTestEvalWhen(t *testing.T) {
	// when is macro
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("when"), true, 1, 2, 3))
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
	v := Eval(e, Slice2List(sym("and"), true, true, false))
	if v != false {
		t.Error("expected false")
		fmt.Println(v)
	}
}

func TestEvalAnd2(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("and"), true, true, true))
	if v != true {
		t.Error("expected true")
		fmt.Println(v)
	}
}

func TestEvalOr1(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("or"), false, false, true))
	if v != true {
		t.Error("expected true")
		fmt.Println(v)
	}
}

func TestEvalOr2(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("or"), false, false, false))
	if v != false {
		t.Error("expected false")
		fmt.Println(v)
	}
}

// cond is defined by macro.
func xTestEvalCond(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("cond"),
		Slice2List(sym(">"), 0, 1), 1,
		Slice2List(sym(">"), 3, 2), 2,
		sym("else"), 3))
	if v != 2 {
		t.Error("expected 2")
		fmt.Println(v)
	}
}

// cond is defined by macro.
func xTestEvalCondElse(t *testing.T) {
	e := NewGlobalRootFrame()
	v := Eval(e, Slice2List(sym("cond"),
		Slice2List(sym(">"), 0, 1), 1,
		Slice2List(sym(">"), 0, 2), 2,
		sym("else"), 3))
	if v != 3 {
		t.Error("expected 2")
		fmt.Println(v)
	}
}
