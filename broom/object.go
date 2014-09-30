package broom

import (
	"fmt"
	"reflect"
)

//type interface{} interface{} // Anything.

type Undef interface{} // T.B.D.

type Symbol interface {
	//T.B.D.
	GetValue() string
	Eq(other interface{}) bool
}

type Pair interface {
	Car() interface{}
	Cdr() Pair
	SetCar(v interface{}) Undef
	SetCdr(p Pair) Undef
}

func isNull(v interface{}) bool {
	//null?
	return v == nil
}

func isBoolean(v interface{}) bool {
	//boolean?
	_, ok := v.(bool)
	return ok
}

func isChar(v interface{}) bool {
	//char?
	_, ok := v.(rune)
	return ok
}

func isSymbol(v interface{}) bool {
	//symbol?
	_, ok := v.(Symbol) //FIXME
	return ok
}

//eof-object?

func isNumber(v interface{}) bool {
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

func isPair(v interface{}) bool {
	//pair?
	_, ok := v.(Pair)
	return ok
}

//port?

type Recur struct {
	env Environment
	eb  *EnvBuilder
}

func NewRecur(outer Environment, xs []interface{}) *Recur {
	r := new(Recur)
	r.eb = NewEnvBuilder(Evens(xs))
	e := NewEnvFrame(outer)
	r.env = r.eb.EvalAndBindAll(Odds(xs), e, e)
	r.Env().Bind("recur", r)
	return r
}

func (r *Recur) Update(xs []interface{}, env Environment) {
	r.env = r.eb.EvalAndBindAll(xs, NewEnvFrame(r.Env().Outer()), env)
	r.Env().Bind("recur", r)
}

func (r *Recur) Env() Environment {
	return r.env
}

func isRecur(v interface{}) bool {
	_, ok := v.(*Recur)
	return ok
}

type Closure func(env Environment, cdr Pair) interface{}

func isProcedure(v interface{}) bool {
	//procedure?
	_, ok := v.(Closure)
	return ok
}

type Syntax Closure

func isSyntax(v interface{}) bool {
	//syntax?
	_, ok := v.(Syntax)
	return ok
}

func isString(v interface{}) bool {
	//string?
	_, ok := v.(string)
	return ok
}

// vector?
func isArray(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}

// bytevector?
// define-record-type

func isMap(v interface{}) bool {
	_, ok := v.(map[interface{}]interface{})
	return ok
}

func DumpMap(x interface{}) {
	mx, _ := x.(map[interface{}]interface{})
	fmt.Println("Dumping", mx)
	for k, vx := range mx {
		fmt.Println(k, vx)
	}
}

func EqMap(x, y interface{}) bool {
	mx, _ := x.(map[interface{}]interface{})
	my, _ := y.(map[interface{}]interface{})
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

func EqArray(x, y interface{}) bool {
	println("EqArray")
	xs, _ := x.([]interface{})
	ys, _ := y.([]interface{})
	if len(xs) != len(ys) {
		return false
	}
	for i, v := range xs {
		if !Eq(ys[i], v) {
			return false
		}
	}
	return true
}

func Eq(x, y interface{}) bool {
	switch {
	case isMap(x) && isMap(y):
		return EqMap(x, y)
	case isSymbol(x) && isSymbol(y):
		sx, _ := x.(Symbol)
		sy, _ := y.(Symbol)
		return sx.Eq(sy)
	case isPair(x) && isPair(y):
		return Eq(Car(x), Car(y)) && Eq(Cdr(x), Cdr(y))
	case isArray(x) && isArray(y):
		return EqArray(x, y)
	default:
		return x == y
	}
	return false
}

func Wider(t1, t2 reflect.Type) int {
	k1 := t1.Kind()
	k2 := t2.Kind()
	if k1 > k2 {
		return 1
	}
	if k1 < k2 {
		return -1
	}
	return 0
}

func CoerceType(xv, yv reflect.Value) (av, bv reflect.Value, t reflect.Type, ok bool) {
	xt := xv.Type()
	yt := yv.Type()

	if xt == yt {
		return xv, yv, xt, true
	}

	switch Wider(xt, yt) {
	case -1:
		if xt.ConvertibleTo(yt) {
			return xv.Convert(yt), yv, yt, true
		}

	case 1:
		if yt.ConvertibleTo(xt) {
			return xv, yv.Convert(xt), xt, true
		}
	}
	return xv, yv, nil, false
}

func BinaryAdd(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return reflect.ValueOf(xv.Int() + yv.Int()).Convert(t).Interface()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return reflect.ValueOf(xv.Float() + yv.Float()).Convert(t).Interface()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return reflect.ValueOf(xv.Complex() + yv.Complex()).Convert(t).Interface()
	default:
		panic("Add is not supported by this type.")
	}
}

func BinarySub(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return reflect.ValueOf(xv.Int() - yv.Int()).Convert(t).Interface()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return reflect.ValueOf(xv.Float() - yv.Float()).Convert(t).Interface()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return reflect.ValueOf(xv.Complex() - yv.Complex()).Convert(t).Interface()
	default:
		panic("Sub is not supported by this type.")
	}
}

func BinaryMul(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return reflect.ValueOf(xv.Int() * yv.Int()).Convert(t).Interface()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return reflect.ValueOf(xv.Float() * yv.Float()).Convert(t).Interface()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return reflect.ValueOf(xv.Complex() * yv.Complex()).Convert(t).Interface()
	default:
		panic("Mul is not supported by this type.")
	}
}

func BinaryDiv(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return reflect.ValueOf(xv.Int() / yv.Int()).Convert(t).Interface()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return reflect.ValueOf(xv.Float() / yv.Float()).Convert(t).Interface()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return reflect.ValueOf(xv.Complex() / yv.Complex()).Convert(t).Interface()
	default:
		panic("Div is not supported by this type.")
	}
}

func BinaryLessThan(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		fmt.Println(xv, xv.Type(), yv, yv.Type())
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return xv.Int() < yv.Int()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return xv.Float() < yv.Float()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	default:
		panic("Less than is not supported by this type.")
	}
}

func BinaryGreaterThan(x, y interface{}) interface{} {
	xv, yv, t, ok := CoerceType(reflect.ValueOf(x), reflect.ValueOf(y))
	if !ok || t == nil {
		panic("Failed to coerce.")
	}

	switch t.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		return xv.Int() > yv.Int()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return xv.Float() > yv.Float()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	default:
		panic("Greater than is not supported by this type.")
	}
}
