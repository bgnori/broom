package broom

import (
	"fmt"
	"reflect"
	"strings"
)

type Pair struct {
	car interface{}
	cdr List
}

func Cons(car interface{}, cdr List) List {
	return &Pair{car:car, cdr:cdr}
}

func Car(v List) interface{} {
	return v.Car()
}

func Cdr(v List) List {
	return v.Cdr()
}

/* As Sequence */
func (p *Pair) First() interface{} {
	return p.Car()
}

func (p *Pair) Rest() Sequence {
	v := p.Cdr()
	if v != nil {
		return v.(Sequence)
	}
	return nil
}

func (p *Pair) Cons(item interface{}) Sequence {
	return &Pair{car: item, cdr: p}
}

func (p *Pair) IsEmpty() bool {
	fmt.Println("Pair.IsEmpty", p, reflect.TypeOf(p))
	return p == nil
}


/* As List */

func (p *Pair) Car () interface{} {
	return p.car
}

func (p *Pair) Cdr() List {
	return p.cdr
}

func (p *Pair) SetCar(v interface{}) Undef {
	p.car = v
	return nil
}

func (p *Pair) SetCdr(cdr List) Undef {
	p.cdr = cdr
	return nil
}

func (p *Pair) String() string {
	//assume that proper list
	xs := List2Slice(p)
	ss := make([]string, 0)
	for _, x := range xs {
		ss = append(ss, fmt.Sprint(x))
	}
	return "(" + strings.Join(ss, " ") + ")"
}

func sub(v List, xs []interface{}) [](interface{}) {
	if v == nil {
		return xs
	} else {
		xs = append(xs, Car(v))
		return sub(Cdr(v), xs)
	}
}

func List2Slice(v List) []interface{} {
	return sub(v, make([]interface{}, 0))
}

func Slice2List(xs ...interface{}) List {
	//(list obj... )
	// this function supports . cdr, for none proper list
	if len(xs) == 0 {
		return nil
	}
	return Cons(xs[0], Slice2List(xs[1:]...))
}

func Append(xs List, cdr List) List {
	if xs == nil {
		return cdr
	} else {
		return Cons(Car(xs), Append(Cdr(xs), cdr))
	}
}

func isList(v interface{}) bool {
	if nil == v {
		return true
	}
	if xs, ok := v.(List) ; ok {
		return isList(Cdr(xs))
	}
	return false
}

func Chop2(xs List) []struct{ header, body interface{} } {

	ys := make([]struct{ header, body interface{} }, 0)
	for xs != nil {
		header := Car(xs)
		body := Car(Cdr(xs))
		xs = Cdr(Cdr(xs))
		ys = append(ys, struct{ header, body interface{} }{header: header, body: body})
	}
	return ys
}

func Odds(xs []interface{}) []interface{} {
	ys := make([]interface{}, 0)
	for i, v := range xs {
		if i%2 == 1 {
			ys = append(ys, v)
		}
	}
	return ys
}

func Evens(xs []interface{}) []interface{} {
	ys := make([]interface{}, 0)
	for i, v := range xs {
		if i%2 == 0 {
			ys = append(ys, v)
		}
	}
	return ys
}
