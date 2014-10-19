package broom

import (
	"testing"
)

/* fail. won't work
func Test_IsEmpty(t *testing.T) {
	var x Sequence
	x = nil
	if !x.IsEmpty() {
		t.Error("nil must be Empty ")
	}
	var y List
	y = nil
	if !y.IsEmpty() {
		t.Error("nil must be Empty ")
	}
}
*/

func Test_Cons(t *testing.T) {
	s := Cons(1, nil)
	if Length(s) != 1 {
		t.Errorf("expected 1 but got %d", Length(s))
	}
}

func Test_MakeFromSlice(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		xs[i] = 'a' + i
	}
	ys := MakeFromSlice(xs...)
	zs := make([]int, 0)
	for ; ys != nil; ys = ys.Rest() {
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", len(zs))
	}
}

func Test_SeqTake5(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		xs[i] = 'a' + i
	}
	zs := make([]int, 0)
	var ys Sequence
	for ys = SeqTake(5, MakeFromSlice(xs...)); ys!=nil && !ys.IsEmpty(); ys = ys.Rest() {
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 5 {
		t.Errorf("expected 5 items but there is %d", len(zs))
	}
}

func Test_SeqTake20(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		xs[i] = 'a' + i
	}
	zs := make([]int, 0)
	var ys Sequence
	for ys = SeqTake(20, MakeFromSlice(xs...)); ys!=nil && !ys.IsEmpty(); ys = ys.Rest() {
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", len(zs))
	}
}

func Test_Seq2Slice(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		xs[i] = 'a' + i
	}
	ys := MakeFromSlice(xs...)
	zs := Seq2Slice(ys)
	if len(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", len(zs))
	}
}

func Test_SeqAppend_0(t *testing.T) {
	xs := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		xs[i] = 'a' + i
	}
	ys := make([]interface{}, 5)
	for i := 0; i < 5; i++ {
		ys[i] = 'A' + i
	}
	zs := SeqAppend(MakeFromSlice(xs...), MakeFromSlice(ys...))
	if Length(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", Length(zs))
	}
}

func Test_SeqAppend_1(t *testing.T) {
	xs := MakeFromSlice(sym("a"), sym("b"), sym("c"))
	ys := MakeFromSlice(sym("d"), sym("e"), sym("f"))
	got := SeqAppend(xs, ys)
	if !SeqEq(got, MakeFromSlice(sym("a"), sym("b"), sym("c"), sym("d"), sym("e"), sym("f"))) {
		t.Errorf("Append not working right, got %v", got)
	}
}

func Test_SetRange(t *testing.T) {
	sum := 0
	for r := SeqRange(0, 101, 1); r != nil && !r.IsEmpty(); r = r.Rest() {
		v := r.First().(int)
		sum += v
	}
	if sum != 5050 {
		t.Errorf("Expected 5050, but got %v", sum)
	}
}

func Test_SeqByAppend(t *testing.T) {
	xs := MakeFromSlice(sym("a"), sym("b"), sym("c"))
	ys := MakeFromSlice(sym("d"), sym("e"), sym("f"))
	got := MakeSeqByAppend(xs, ys)
	if !SeqEq(got, MakeFromSlice(sym("a"), sym("b"), sym("c"), sym("d"), sym("e"), sym("f"))) {
		t.Errorf("Append not working right, got %v", SeqString(got))
	}
}
