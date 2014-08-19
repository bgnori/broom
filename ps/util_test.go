package ps

import (
    "testing"
)

func TestRecEqBothPair(t *testing.T) {
    xs := Cons(Cons(Int(0), nil), Cons(Int(1), Cons(Int(2), nil))) // ((0) 1 2)
    ys := Cons(Cons(Int(0), nil), Cons(Int(1), Cons(Int(2), nil))) // ((0) 1 2)
    zs := Cons(Cons(Int(0), Int(3)), Cons(Int(1), Cons(Int(2), nil))) // ((0 . 3) 1 2)

    if !RecEq(xs, ys) {
        t.Error("Should be true for RecEq(xs, ys)")
    }
    if RecEq(xs, zs) {
        t.Error("Should be false for RecEq(xs, zs)")
    }
}
