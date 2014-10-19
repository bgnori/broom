package broom

import (
	"fmt"
	"strings"
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

type FromSlice struct {
	wrapped []interface{}
}

func (fs *FromSlice) IsEmpty() bool {
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
	return Cons(item, fs)
}

func MakeFromSlice(xs ...interface{}) Sequence {
	return &FromSlice{wrapped: xs}
}

type FromChan struct { /* Kind a lazy, might block */
	wrapped  chan interface{}
	realized Sequence /* cannnot be lazy */
}

func (fc *FromChan) IsEmpty() bool {
	return false
}

func (fc *FromChan) realize() {
	v, more := <-fc.wrapped
	if more {
		fc.realized = Cons(v, MakeFromChan(fc.wrapped))
	} else {
		fc.realized = Cons(v, nil)
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
	return Cons(item, fc)
}

func MakeFromChan(ch chan interface{}) Sequence {
	return &FromChan{wrapped: ch}
}

func Length(s Sequence) int {
	if s == nil || s.IsEmpty() {
		return 0
	}
	return Length(s.Rest()) + 1
}

func SeqString(xs Sequence) string {
	ss := make([]string, 0)
	for ; xs != nil && !xs.IsEmpty(); xs = xs.Rest() {
		ss = append(ss, fmt.Sprintf("%v", xs.First()))
	}
	return "(" + strings.Join(ss, " ") + ")"
}

func SeqTake(n int, s Sequence) Sequence {
	if n == 0 || s == nil || s.IsEmpty() {
		return nil
	}
	v := SeqTake(n-1, s.Rest())
	if v == nil || v.IsEmpty() {
		return Cons(s.First(), nil)
	}
	return v.Cons(s.First())
}

func SeqDrop(n int, s Sequence) Sequence {
	i := 0
	for i < n && s != nil && !s.IsEmpty() {
		s = s.Rest()
		i += 1
	}
	return s
}

func Seq2Slice(s Sequence) []interface{} {
	xs := make([]interface{}, 0)
	for ; s != nil && !s.IsEmpty(); s = s.Rest() {
		xs = append(xs, s.First())
	}
	return xs
}

func SeqAppend(xs, ys Sequence) Sequence {
	if xs == nil || xs.IsEmpty() {
		return ys
	} else {
		zs := SeqAppend(xs.Rest(), ys)
		return zs.Cons(xs.First())
	}
}

type SeqByAppend struct {
	// doubtful.
	first  Sequence
	second Sequence
}

func (seq *SeqByAppend) First() interface{} {
	if seq.first != nil && !seq.first.IsEmpty() {
		return seq.first.First()
	}
	if seq.second != nil && !seq.second.IsEmpty() {
		return seq.second.First()
	}
	panic("tried to First() on Empty Sequence")
}

func (seq *SeqByAppend) Rest() Sequence {
	if seq.first != nil && !seq.first.IsEmpty() {
		return MakeSeqByAppend(seq.first.Rest(), seq.second)
	}
	return seq.second
}

func (seq *SeqByAppend) Cons(item interface{}) Sequence {
	return Cons(item, seq)
}

func (seq *SeqByAppend) IsEmpty() bool {
	if seq.first != nil && !seq.first.IsEmpty() {
		return false
	}
	if seq.second != nil && !seq.second.IsEmpty() {
		return false
	}
	return true
}

func MakeSeqByAppend(xs, ys Sequence) Sequence {
	if xs != nil && !xs.IsEmpty() {
		return &SeqByAppend{first: xs, second: ys}
	}
	return ys
}

func SeqFilter(pred func(interface{}) bool, seq Sequence) Sequence {
	if seq == nil || seq.IsEmpty() {
		return nil
	}
	result := SeqFilter(pred, seq.Rest())
	v := seq.First()
	if pred(v) {
		if result == nil || result.IsEmpty() {
			return Cons(v, nil)
		}
		return result.Cons(v)
	}
	return result
}

func SeqEvens(seq Sequence) Sequence {
	n := 0
	return SeqFilter(func(x interface{}) (v bool) {
		v = (n%2 == 1)
		n += 1
		return
	}, seq)
}

func SeqOdds(seq Sequence) Sequence {
	n := 0
	return SeqFilter(func(x interface{}) (v bool) {
		v = (n%2 == 0)
		n += 1
		return
	}, seq)
}

func SeqZip2(xs, ys Sequence) Sequence {
	tmp := make([]interface{}, 0) //FIXME
	for xs != nil && !xs.IsEmpty() && ys != nil && !ys.IsEmpty() {
		v := make([]interface{}, 2)
		v[0] = xs.First()
		v[1] = ys.First()
		tmp = append(tmp, v)
		xs = xs.Rest()
		ys = ys.Rest()
	}
	return MakeFromSlice(tmp...)
}

func SeqEq(xs, ys Sequence) bool {
	//Both Empty
	if (xs == nil || xs.IsEmpty()) && (ys == nil || ys.IsEmpty()) {
		return true
	}
	// Both has something
	if xs != nil && !xs.IsEmpty() && ys != nil && !ys.IsEmpty() {
		if xs.First() == ys.First() {
			return SeqEq(xs.Rest(), ys.Rest())
		}
	}
	fmt.Println(xs, ys)
	return false
}

type Range struct {
	start int
	stop  int
	step  int
}

func (r *Range) First() interface{} {
	return r.start
}

func (r *Range) Rest() Sequence {
	return &Range{start: r.start + r.step, stop: r.stop, step: r.step}
}

func (r *Range) Cons(item interface{}) Sequence {
	return Cons(item, r)
}

func (r *Range) IsEmpty() bool {
	return r.start >= r.stop
}

func SeqRange(start, stop, step int) Sequence {
	return &Range{start: start, stop: stop, step: step}
}

func SeqReduce(init interface{}, f func(interface{}, interface{}) interface{}, xs Sequence) interface{} {
	for ; xs != nil && !xs.IsEmpty(); xs = xs.Rest() {
		init = f(init, xs.First())
	}
	return init
}
