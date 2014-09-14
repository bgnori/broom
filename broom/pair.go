package broom

import (
	"fmt"
	"strings"
)

type pairImpl struct {
	car interface{}
	cdr Pair
}

func Cons(car interface{}, cdr Pair) Pair {
	return &pairImpl{car: car, cdr: cdr}
}

func Car(v interface{}) interface{} {
	u, ok := v.(Pair)
	if !ok {
		panic("non pair object for Car()")
	}
	return u.Car()
}

func Cdr(v interface{}) Pair {
	u, ok := v.(Pair)
	if !ok {
		panic("non pair object for Cdr()")
	}
	return u.Cdr()
}

func (p *pairImpl) Car() interface{} {
	return p.car
}

func (p *pairImpl) Cdr() Pair {
	return p.cdr
}

func (p *pairImpl) SetCar(v interface{}) Undef {
	p.car = v
	return nil
}

func (p *pairImpl) SetCdr(cdr Pair) Undef {
	p.cdr = cdr
	return nil
}

func (p *pairImpl) String() string {
	//assume that proper list
	xs := List2Arr(p)
	ss := make([]string, 0)
	for _, x := range xs {
		ss = append(ss, fmt.Sprint(x))
	}
	return "(" + strings.Join(ss, " ") + ")"
}

func sub(v interface{}, xs []interface{}) [](interface{}) {
	if v == nil {
		return xs
	} else {
		xs = append(xs, Car(v))
		return sub(Cdr(v), xs)
	}
}

func List2Arr(v interface{}) []interface{} {
	return sub(v, make([]interface{}, 0))
}

func List(xs ...interface{}) Pair {
	//(list obj... )
	// this function supports . cdr, for none proper list
	if len(xs) == 0 {
		return nil
	}
	return Cons(xs[0], List(xs[1:]...))
}

func Append(xs Pair, cdr Pair) Pair {
	if Cdr(xs) == nil {
		return Cons(Car(xs), cdr)
	} else {
		return Cons(Car(xs), Append(Cdr(xs), cdr))
	}
}

func isList(xs interface{}) bool {
	if isNull(xs) {
		return true
	}
	if isPair(xs) {
		return isList(Cdr(xs))
	}
	return false
}

func Length(xs interface{}) int {
	if isNull(xs) {
		return 0
	}
	if isPair(xs) {
		return Length(Cdr(xs)) + 1
	}
	panic("proper list required")
}
