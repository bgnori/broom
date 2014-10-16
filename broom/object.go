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
}

type List interface {
	Sequence
	Car() interface{}
	Cdr() List
	SetCar(v interface{}) Undef
	SetCdr(p List) Undef
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

//port?

type Recur struct {
	env Environment
	eb  *EnvBuilder
}

func NewRecur(outer Environment, xs []interface{}) *Recur {
	r := new(Recur)
	seq := MakeFromSlice(xs...)
	r.eb = NewEnvBuilder(SeqEvens(seq))
	e := NewEnvFrame(outer)
	r.env = r.eb.EvalAndBindAll(Seq2Slice(SeqOdds(seq)), e, e)
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
