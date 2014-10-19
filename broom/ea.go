package broom

import (
	"fmt"
	"sync"
)

type Environment interface {
	Bind(symbol Symbol, v interface{})
	Resolve(symbol Symbol) (interface{}, error)
	SetOuter(outer Environment)
	Outer() Environment
	Dump()
}

type EvalError string

func (e EvalError) Error() string {
	return string(e)
}

func Eval(env Environment, expr interface{}) interface{} {
	switch x := expr.(type) {
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		string:
		return expr
	case func(Environment, Sequence) interface{}:
		return expr
	case *Recur:
		return expr
	case Symbol: // variables?
		if v, err := env.Resolve(x); err != nil {
			panic(err)
		} else {
			return v
		}
	case Sequence:
		switch v := Eval(env, x.First()).(type) {
		case *Recur:
			v.Update(Seq2Slice(x.Rest()), env)
			return v
		case ([]interface{}):
			idx := Second(x).(int)
			return v[idx]
		case func(Environment, Sequence) interface{}:
			return v(env, x.Rest())
		default:
			panic("application error, expected SExprOperator, but got " + fmt.Sprintf("%v", v))
		}
	}
	return nil
}

func EvalExprs(env Environment, seq Sequence) interface{} {
	var x interface{}
	for ; seq != nil && !seq.IsEmpty(); seq = seq.Rest() {
		x = Eval(env, seq.First())
	}
	return x
}

type enviroment struct {
	rwm       sync.RWMutex
	variables map[Symbol]interface{}
	outer     Environment
}

func (env *enviroment) RLock() {
	env.rwm.RLock()
}

func (env *enviroment) RUnlock() {
	env.rwm.RUnlock()
}

func (env *enviroment) Lock() {
	env.rwm.Lock()
}

func (env *enviroment) Unlock() {
	env.rwm.Unlock()
}

func NewEnvFrame(outer Environment) *enviroment {
	e := new(enviroment)
	e.variables = make(map[Symbol]interface{})
	e.SetOuter(outer)
	e.Bind(sym("_env"), e)
	e.Bind(sym("_outer"), outer)
	return e
}

func NewGlobalRootFrame() *enviroment {
	e := NewEnvFrame(nil)
	setupSpecialForms(e)
	setupBuiltins(e)
	e.Bind(sym("eval"), func(env Environment, cdr Sequence) interface{} {
		given := Eval(env, cdr.First()).(Environment)
		given.Dump()
		return Eval(given, Second(cdr))
	})
	return e
}

func (env *enviroment) Bind(symbol Symbol, v interface{}) {
	env.Lock()
	defer env.Unlock()
	env.variables[symbol] = v
}

func (env *enviroment) Resolve(symbol Symbol) (interface{}, error) {
	env.RLock()
	defer env.RUnlock()
	if v, ok := env.variables[symbol]; ok {
		return v, nil
	}
	outer := env.Outer()
	if outer != nil {
		return outer.Resolve(symbol)
	}
	return nil, EvalError(fmt.Sprintf("Unbound variable %s", symbol))
}

func (env *enviroment) SetOuter(outer Environment) {
	env.Lock()
	defer env.Unlock()
	env.outer = outer
}
func (env *enviroment) Outer() Environment {
	env.RLock()
	defer env.RUnlock()
	return env.outer
}

func (env *enviroment) Dump() {
	env.RLock()
	defer env.RUnlock()
	fmt.Println("=====")
	for key, value := range env.variables {
		fmt.Println(key, ":", value)
	}
	if env.outer != nil {
		env.outer.Dump()
	}
}

type EnvBuilder struct {
	params   []Symbol
	variadic bool
}

func NewEnvBuilder(seq Sequence) *EnvBuilder {
	eb := &EnvBuilder{}

	eb.params = make([]Symbol, 0)
	for ; seq != nil && !seq.IsEmpty(); seq = seq.Rest() {
		x := seq.First()
		if s, ok := x.(Symbol); !ok {
			panic("Parameter must be symbol")
		} else {
			switch s.GetValue() {
			case "&":
				eb.variadic = true
			default:
				eb.params = append(eb.params, s)
			}
		}
	}
	return eb
}

func (eb *EnvBuilder) Len() int {
	return len(eb.params)
}

func (eb *EnvBuilder) EvalAndBindAll(as []interface{}, to_bind, to_eval Environment) Environment {
	if eb.variadic {
		last := eb.Len() - 1
		for i, s := range eb.params[:last] {
			v := Eval(to_eval, as[i])
			to_bind.Bind(s, v)
		}
		vs := make([]interface{}, 0)
		for _, v := range as[last:] {
			vs = append(vs, Eval(to_eval, v))
		}
		to_bind.Bind(eb.params[last], vs)
	} else {
		for i, s := range eb.params {
			v := Eval(to_eval, as[i])
			to_bind.Bind(s, v)
		}
	}
	return to_bind
}

func (eb *EnvBuilder) BindAll(as []interface{}, env Environment) Environment {
	if eb.variadic {
		last := eb.Len() - 1
		for i, s := range eb.params[:last] {
			env.Bind(s, as[i])
		}
		env.Bind(eb.params[last], as[last:])
	} else {
		if len(as) < eb.Len() {
			panic("not enough argument!")
		}
		for i, s := range eb.params {
			env.Bind(s, as[i])
		}
	}
	return env
}

func Body(p Sequence) Sequence {
	return p.Rest().(Sequence)
}
