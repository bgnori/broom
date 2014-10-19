package broom

import (
	"fmt"
	"strings"
)

type Pair struct {
	car interface{}
	cdr Sequence
}

func Cons(car interface{}, cdr Sequence) Sequence {
	return &Pair{car: car, cdr: cdr}
}

func Second(v Sequence) interface{} {
	return v.Rest().First()
}

/* As Sequence */
func (p *Pair) First() interface{} {
	return p.car
}

func (p *Pair) Rest() Sequence {
	v := p.cdr
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

func (p *Pair) String() string {
	//assume that proper list
	var xs Sequence
	ss := make([]string, 0)
	for xs = p; xs != nil && !xs.IsEmpty(); xs = xs.Rest() {
		ss = append(ss, fmt.Sprint(xs.First()))
	}
	return "(" + strings.Join(ss, " ") + ")"
}

func Slice2List(xs ...interface{}) Sequence {
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
	if xs, ok := v.(Sequence); ok {
		return isList(xs.Rest())
	}
	return false
}

func Chop2(xs Sequence) []struct{ header, body interface{} } {

	ys := make([]struct{ header, body interface{} }, 0)
	for xs != nil {
		header := xs.First()
		body := Second(xs)
		xs = xs.Rest().Rest()
		ys = append(ys, struct{ header, body interface{} }{header: header, body: body})
	}
	return ys
}
