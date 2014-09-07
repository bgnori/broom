package broom

import (
    "fmt"
)

type Value interface{} // Anything.

type Undef interface{} // T.B.D.

type Symbol interface {
	//T.B.D.
	GetValue() string
	Eq(other Value) bool
}

type Pair interface {
	Car() Value
	Cdr() Pair
	SetCar(v Value) Undef
	SetCdr(p Pair) Undef
}

func isNull(v Value) bool {
	//null?
	return v == nil
}

func isBoolean(v Value) bool {
	//boolean?
	_, ok := v.(bool)
	return ok
}

func isChar(v Value) bool {
	//char?
	_, ok := v.(rune)
	return ok
}

func isSymbol(v Value) bool {
	//symbol?
	_, ok := v.(Symbol) //FIXME
	return ok
}

//eof-object?

func isNumber(v Value) bool {
	//number?
	//see golang builtin
	switch v.(type) {
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case float32:
	case float64:
	case complex64:
	case complex128:
	default:
		return false
	}
	return true
}

func isPair(v Value) bool {
	//pair?
	_, ok := v.(Pair)
	return ok
}

//port?

type Closure func(env Enviroment, cdr Pair) Value

func isProcedure(v Value) bool {
	//procedure?
	_, ok := v.(Closure)
	return ok
}

type Syntax Closure

func isSyntax(v Value) bool {
	//syntax?
	_, ok := v.(Syntax)
	return ok
}

func isString(v Value) bool {
	//string?
	_, ok := v.(string)
	return ok
}

// vector?
func isArray(v Value) bool {
    _, ok := v.([]Value)
    return ok
}


// bytevector?
// define-record-type

func isMap(v Value) bool {
    _, ok := v.(map[Value]Value)
    return ok
}

func DumpMap(x Value) {
    mx, _ := x.(map[Value]Value)
    fmt.Println("Dumping", mx)
    for k, vx := range mx {
        fmt.Println(k, vx)
    }
}

func EqMap(x, y Value) bool {
    mx, _ := x.(map[Value]Value)
    my, _ := y.(map[Value]Value)
    for k, vx := range mx {
        vy, in := my[k]
        if in && vx == vy {
            continue
        } else {
            return false
        }
    }
    for k, vy := range my {
        vx, in := mx[k]
        if in && vx == vy {
            continue
        } else {
            return false
        }
    }
    return true
}

func Eq(x, y Value) bool {
    switch {
    case isMap(x) && isMap(y):
        return EqMap(x,y)
    case isSymbol(x) && isSymbol(y):
        sx, _ := x.(Symbol)
        sy, _ := y.(Symbol)
        return sx.Eq(sy)
    case isPair(x) && isPair(y):
        return Eq(Car(x), Car(y)) && Eq(Cdr(x), Cdr(y))
    default:
        return x==y
    }
    return false
}

