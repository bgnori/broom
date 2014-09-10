package broom

import (
	"fmt"
	"testing"
)

func TestEvalNumber(t *testing.T) {
	e := NewGlobalRootFrame()
	if Eval(e, 1) != 1 {
		t.Error("expected 1")
	}
}

func TestEvalString(t *testing.T) {
	e := NewGlobalRootFrame()
	if Eval(e, "あいう") != "あいう" {
		t.Error("expected あいう")
	}
}

func TestEvalVariable(t *testing.T) {
	e := NewGlobalRootFrame()
	e.Bind("a", 42)
	v := Eval(e, sym("a"))
	if v != 42 {
		t.Error("expected 42")
		fmt.Println(v)
	}
}
