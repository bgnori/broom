package broom

import (
	"fmt"
	"sync"
)

type Environment interface {
	Bind(name string, v interface{})
	Resolve(name string) (interface{}, error)
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
	case bool:
		return expr
	case int:
		return expr
	case int8:
		return expr
	case int16:
		return expr
	case int32:
		return expr
	case int64:
		return expr
	case uint:
		return expr
	case uint8:
		return expr
	case uint16:
		return expr
	case uint32:
		return expr
	case uint64:
		return expr
	case float32:
		return expr
	case float64:
		return expr
	case complex64:
		return expr
	case complex128:
		return expr
	case string:
		return expr
	case func(Environment, List) interface{}:
		return expr
	case *Recur:
		return expr
	case Symbol: // variables?
		if v, err := env.Resolve(x.GetValue()); err != nil {
			panic(err)
		} else {
			return v
		}
	case List:
		switch v := Eval(env, Car(x)).(type) {
		case *Recur:
			v.Update(Seq2Slice(Cdr(x)), env)
			return v
		case ([]interface{}):
			idx := Car(Cdr(x)).(int)
			return v[idx]
		case func(Environment, List) interface{}:
			return v(env, Cdr(x))
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
	variables map[string]interface{}
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
	e.variables = make(map[string]interface{})
	e.SetOuter(outer)
	e.Bind("_env", e)
	e.Bind("_outer", outer)
	return e
}

func NewGlobalRootFrame() *enviroment {
	e := NewEnvFrame(nil)
	setupSpecialForms(e)
	setupBuiltins(e)
	e.Bind("eval", func(env Environment, cdr List) interface{} {
		given := Eval(env, Car(cdr)).(Environment)
		given.Dump()
		return Eval(given, Car(Cdr(cdr)))
	})
	return e
}

func (env *enviroment) Bind(name string, v interface{}) {
	env.Lock()
	defer env.Unlock()
	env.variables[name] = v
}

func (env *enviroment) Resolve(name string) (interface{}, error) {
	env.RLock()
	defer env.RUnlock()
	if v, ok := env.variables[name]; ok {
		return v, nil
	}
	outer := env.Outer()
	if outer != nil {
		return outer.Resolve(name)
	}
	return nil, EvalError(fmt.Sprintf("Unbound variable %s", name))
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
	fmt.Println("NewEnvBuilder", seq)
	eb := &EnvBuilder{}

	eb.params = make([]Symbol, 0)
	for ; seq != nil && !seq.IsEmpty(); seq = seq.Rest() {
		x := seq.First()
		fmt.Println(x)
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
	fmt.Println("EvalAndBindAll", as)
	if eb.variadic {
		last := eb.Len() - 1
		for i, s := range eb.params[:last] {
			v := Eval(to_eval, as[i])
			to_bind.Bind(s.GetValue(), v)
		}
		vs := make([]interface{}, 0)
		for _, v := range as[last:] {
			vs = append(vs, Eval(to_eval, v))
		}
		to_bind.Bind(eb.params[last].GetValue(), vs)
	} else {
		for i, s := range eb.params {
			v := Eval(to_eval, as[i])
			to_bind.Bind(s.GetValue(), v)
		}
	}
	return to_bind
}

func (eb *EnvBuilder) BindAll(as []interface{}, env Environment) Environment {
	if eb.variadic {
		last := eb.Len() - 1
		for i, s := range eb.params[:last] {
			env.Bind(s.GetValue(), as[i])
		}
		env.Bind(eb.params[last].GetValue(), as[last:])
	} else {
		if len(as) < eb.Len() {
			panic("not enough argument!")
		}
		for i, s := range eb.params {
			env.Bind(s.GetValue(), as[i])
		}
	}
	return env
}

func Body(p List) List {
	return Cdr(p).(List)
}
