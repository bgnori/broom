package broom

import (
	"fmt"
	"testing"
	"time"
)

func Test_Symbol(t *testing.T) {
	var v interface{}
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
	var v interface{}
	v = nil
	if !isNull(v) {
		t.Error("nil must be null? true.")
	}
}

func Test_isBoolean(t *testing.T) {
	var v interface{}
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
	var v interface{}
	v = 'あ'
	if !isChar(v) {
		t.Error("v must be char? true.")
	}
}

func xTest_isSymbol(t *testing.T) {
	var v interface{}
	v = 'あ'
	if !isChar(v) {
		t.Error("v must be char? true.")
	}
}

func Test_isNumber(t *testing.T) {
	var v interface{}
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
	var v interface{}
	v = nil
	if isPair(v) {
		t.Error("nil must not be pair")
	}
}

func Test_isString(t *testing.T) {
	var v interface{}
	v = "abcd"
	if !isString(v) {
		t.Error("string must be string? true")
	}
	v = 1
	if isString(v) {
		t.Error("int must be string? false")
	}
}

func TestBinaryAddIntXInt(t *testing.T) {
	var v_int8, v_int16, v_int32, v_int64 interface{}
	v_int8 = int8(1)
	v_int16 = int16(1)
	v_int32 = int32(1)
	v_int64 = int64(1)
	nano := time.Nanosecond

	if v, ok := BinaryAdd(v_int8, v_int8).(int8); !ok || v != 2 {
		t.Errorf("Expected int8 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int16, v_int16).(int16); !ok || v != 2 {
		t.Errorf("Expected int16 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int32, v_int32).(int32); !ok || v != 2 {
		t.Errorf("Expected int32 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int64, v_int64).(int64); !ok || v != 2 {
		t.Errorf("Expected int64 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(nano, nano).(time.Duration); !ok || v != 2 {
		t.Errorf("Expected time.Duration value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int8, v_int16).(int16); !ok || v != 2 {
		t.Errorf("Expected int16 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int16, v_int8).(int16); !ok || v != 2 {
		t.Errorf("Expected int16 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int8, v_int64).(int64); !ok || v != 2 {
		t.Errorf("Expected int64 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int64, v_int8).(int64); !ok || v != 2 {
		t.Errorf("Expected int64 value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(nano, v_int8).(time.Duration); !ok || v != 2 {
		t.Errorf("Expected time.Duration value 2 but got %v %v", v, ok)
	}
	if v, ok := BinaryAdd(v_int8, nano).(time.Duration); !ok || v != 2 {
		t.Errorf("Expected time.Duration value 2 but got %v %v", v, ok)
	}
}

func TestBinaryAddFloat(t *testing.T) {
	var x, y interface{}
	x = 3.0
	y = 3.1

	if v, ok := BinaryAdd(x, 1).(float64); !ok || v != 3.0+1 {
		t.Errorf("Expected float64 value 4.0 but got %v", v)
	}
	if v, ok := BinaryAdd(1, x).(float64); !ok || v != 1+3.0 {
		t.Errorf("Expected float64 value 4.0 but got %v", v)
	}
	if v, ok := BinaryAdd(1, y).(float64); !ok || v != 1+3.1 {
		t.Errorf("Expected float64 value 4.1 but got %v", v)
	}
	if v, ok := BinaryAdd(y, 1).(float64); !ok || v != 3.1+1 {
		t.Errorf("Expected float64 value 4.1 but got %v", v)
	}
	if v, ok := BinaryAdd(x, y).(float64); !ok || v != 3.0+3.1 {
		t.Errorf("Expected float64 value 6.1 but got %v", v)
	}
}

func xTestBinaryAddComplexXFloat(t *testing.T) {
	var x, y interface{}
	x = 3.0
	y = 4i
	if v, ok := BinaryAdd(x, y).(complex128); !ok || v != 3.0+4i {
		t.Errorf("Expected complex128 value 7.0i but got %v", v)
	}
}

func TestBinaryAddComplex(t *testing.T) {
	var x, y interface{}
	x = 3.0i
	y = 4i
	if v, ok := BinaryAdd(x, y).(complex128); !ok || v != 3.0i+4i {
		t.Errorf("Expected complex128 value 7.0i but got %v", v)
	}
}

func NoneOf(env Environment, x interface{}) bool {
	return false
}

func AnyOf(env Environment, x interface{}) bool {
	return true
}

func IfInt(env Environment, x interface{}) bool {
	_, ok := x.(int)
	return ok
}

func Be42(env Environment, x interface{}) interface{} {
	return 42
}

func TestVisit(t *testing.T) {
	env := NewGlobalRootFrame()
	x := Visit(env, NoneOf, nil, 1)
	if x.(int) != 1 {
		t.Errorf("Expected 1 but got %v", x)
	}
	v := Cons(1, nil)
	w := Visit(env, NoneOf, nil, v)
	if !Eq(w, v) {
		t.Errorf("Expected (1) but got %v", w)
	}
	v = Cons(Cons(1, nil), nil)
	w = Visit(env, NoneOf, nil, v)
	if !Eq(w, v) {
		t.Errorf("Expected ((1)) but got %v", w)
	}
	v = Cons(1, nil)
	w = Visit(env, IfInt, Be42, v)
	if !Eq(w, Cons(42, nil)) {
		t.Errorf("Expected (42) but got %v", w)
	}

}
