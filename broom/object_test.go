package broom

import (
	"fmt"
	"testing"
)

func Test_Symbol(t *testing.T) {
	var v Value
	v = sym("a")
	if !isSymbol(v) {
		t.Error("(Symbol? 'a) must be true.")
		fmt.Println(v)
		_, ok := v.(symbolImpl)
		fmt.Println(ok)
		_, ok = v.(Symbol)
		fmt.Println(ok)
	}
	if !sym("a").Eq(v) {
		t.Error("'a and 'a must match.")
	}
	if sym("b").Eq(v) {
		t.Error("'a and 'b must not match.")
	}
	if !sym("quote").Eq(sym("quote")) {
		t.Error(" must match.")
	}
}

func Test_isNull(t *testing.T) {
	var v Value
	v = nil
	if !isNull(v) {
		t.Error("nil must be null? true.")
	}
}

func Test_isBoolean(t *testing.T) {
	var v Value
	v = true
	if !isBoolean(v) {
		t.Error("v must be boolean? true.")
	}
	if v != true {
		t.Error("v must be true.")
	}
	v = false
	if !isBoolean(v) {
		t.Error("v must be boolean? true.")
	}
	if v != false {
		t.Error("v must be false.")
	}
}

func Test_isChar(t *testing.T) {
	var v Value
	v = 'あ'
	if !isChar(v) {
		t.Error("v must be char? true.")
	}
}

func xTest_isSymbol(t *testing.T) {
	var v Value
	v = 'あ'
	if !isChar(v) {
		t.Error("v must be char? true.")
	}
}

func Test_isNumber(t *testing.T) {
	var v Value
	v = 123
	if !isNumber(v) {
		t.Error("v must be number? true.")
	}
	v = 123.4
	if !isNumber(v) {
		t.Error("v must be number? true.")
	}
	v = false
	if isNumber(v) {
		t.Error("v must be number? false.")
	}
}

func Test_isPair(t *testing.T) {
	var v Value
	v = nil
	if isPair(v) {
		t.Error("nil must not be pair")
	}
}

func Test_isString(t *testing.T) {
	var v Value
	v = "abcd"
	if !isString(v) {
		t.Error("string must be string? true")
	}
	v = 1
	if isString(v) {
		t.Error("int must be string? false")
	}
}
