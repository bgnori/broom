package catfish

type Value interface{} // Anything.

type Undef interface{} // T.B.D.

type Symbol interface {
	//T.B.D.
	GetValue() string
	Eq(other Value) bool
}

type Pair interface {
	Car() Value
	Cdr() Value
	SetCar(v Value) Undef
	SetCdr(v Value) Undef
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
	case float32:
	case float64:
	case complex64:
	case complex128:
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

type SExprOperator interface {
	Apply(env Enviroment, cdr Value) Value
	LexEnv() Enviroment //Leixical Enviroment
}

type Procedure interface {
	SExprOperator
}

func isProcedure(v Value) bool {
	//procedure?
	_, ok := v.(Procedure)
	return ok
}

type Syntax interface {
	SExprOperator
}

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
// bytevector?
// define-record-type
