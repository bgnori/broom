package broom

import (
	"fmt"
	"reflect"
//	"fmt"
)

/*
See:
	http://clojure.org/sequences
	for motivation

*/
type Sequence interface {
	First() interface{}
	Rest() Sequence
	Cons(item interface{}) Sequence
	IsEmpty() bool
}

type Base struct {
	first interface{}
	rest  Sequence
}

func (bs *Base) IsEmpty() bool {
	fmt.Println("*Base.IsEmpty", bs, reflect.TypeOf(bs))
	return bs == nil
}

func (bs *Base) First() interface{} {
	return bs.first
}

func (bs *Base) Rest() Sequence {
	return bs.rest
}

func (bs *Base) Cons(item interface{}) Sequence {
	return &Base{first: item, rest: bs}
}

type FromSlice struct {
	wrapped []interface{}
}

func (fs *FromSlice) IsEmpty() bool {
	fmt.Println("*FromSlice.IsEmpty", fs, reflect.TypeOf(fs))
	return fs == nil || len(fs.wrapped) == 0
}

func (fs *FromSlice) First() interface{} {
	return fs.wrapped[0]
}

func (fs *FromSlice) Rest() Sequence {
	if len(fs.wrapped) <= 1 {
		return nil
	}
	return MakeFromSlice(fs.wrapped[1:]...)
}

func (fs *FromSlice) Cons(item interface{}) Sequence {
	return &Base{first: item, rest: fs}
}

func MakeFromSlice(xs... interface{}) Sequence {
	return &FromSlice{wrapped: xs}
}

type FromChan struct { /* Kind a lazy, might block */
	wrapped chan interface{}
	realized Sequence /* cannnot be lazy */
}

func (fc *FromChan) IsEmpty() bool {
	return false
}

func (fc *FromChan) realize() {
	v, more := <-fc.wrapped
	if more {
		fc.realized = &Base{first:v, rest: MakeFromChan(fc.wrapped)}
	} else {
		fc.realized = &Base{first:v, rest: nil}
	}
}

func (fc *FromChan) First() interface{} {
	if fc.realized != nil {
		return fc.realized.First()
	}
	fc.realize()
	return fc.realized.First()
}

func (fc *FromChan) Rest() Sequence {
	if fc.realized != nil {
		return fc.realized.Rest()
	}
	fc.realize()
	return fc.realized.Rest()
}

func (fc *FromChan) Cons(item interface{}) Sequence {
	return &Base{first: item, rest: fc}
}

func MakeFromChan(ch chan interface{}) Sequence {
	return &FromChan{wrapped: ch}
}


func Kons(item interface{}, s Sequence) Sequence {
	if s == nil {
		var b *Base
		return b.Cons(item)
	}
	return s.Cons(item)
}

func Length(s Sequence) int {
	fmt.Println("Length", s, reflect.TypeOf(s))
	if s == nil || s.IsEmpty() {
		return 0
	}
	return Length(s.Rest()) + 1
}

func Take(n int, s Sequence) Sequence {
	if n == 0 || s == nil || s.IsEmpty() {
		return nil
	}
	v := Take(n-1, s.Rest())
	if v == nil {
		return Kons(s.First(), nil)
	}
	return v.Cons(s.First())
}

func Seq2Slice(s Sequence)[]interface{} {
	xs := make([]interface{}, 0)
	for ; s != nil && !s.IsEmpty() ; s = s.Rest() {
		xs = append(xs, s.First())
	}
	return xs
}

func SeqAppend(xs, ys Sequence) Sequence {
	if xs == nil {
		return ys
	} else {
		return SeqAppend(xs.Rest(), ys).Cons(xs.First())
	}
}
