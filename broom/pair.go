package broom

import (
	"fmt"
	"strings"
)

type Pair struct {
	car interface{}
	cdr List
}

func Cons(car interface{}, cdr List) List {
	return &Pair{car: car, cdr: cdr}
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
	return p == nil
}

/* As List */

func (p *Pair) Car() interface{} {
	return p.car
}

func (p *Pair) Cdr() List {
	return p.cdr
}

func (p *Pair) String() string {
	//assume that proper list
	var xs Sequence
	ss := make([]string, 0)
	for xs = p; xs != nil && !xs.IsEmpty(); xs = xs.Rest() {
		ss = append(ss, fmt.Sprint(xs.First()))
	}
	return "(" + strings.Join(ss, " ") + ")"
}

func Slice2List(xs ...interface{}) List {
	//(list obj... )
	// this function supports . cdr, for none proper list
	if len(xs) == 0 {
		return nil
	}
	return Cons(xs[0], Slice2List(xs[1:]...))
}

func isList(v interface{}) bool {
	if nil == v {
		return true
	}
	if xs, ok := v.(List); ok {
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
