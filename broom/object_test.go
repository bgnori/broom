package broom

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Eq(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func Test_Symbol(t *testing.T) {
	var v interface{}
	v = sym("a")
	if _, ok := v.(Symbol) ; !ok {
		t.Error("(symbol? 'a) must be true.")
		fmt.Println(v)
		_, ok := v.(symbolImpl)
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

