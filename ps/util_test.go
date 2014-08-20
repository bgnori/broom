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

func TestMakeListNested(t *testing.T) {
    u := MakeList(nil, MakeList(nil, Int(1)))
    if !RecEq(Cons(Cons(Int(1), nil), nil), u) {
        t.Error("Wrong Value.")
    }
}

func TestMakeList(t *testing.T) {
    u := MakeList(nil, Int(1), Int(2), Int(3))
    if !RecEq(u, Cons(Int(1), Cons(Int(2), Cons(Int(3), nil)))) {
        t.Error("Wrong Value.")
    }
}

func TestMakeListEmpty(t *testing.T) {
    u := MakeList(nil, nil)
    if !RecEq(Cons(nil, nil), u) {
        t.Error("Wrong Value.")
    }
}

func TestReqEqBoolList(t *testing.T) {
    u := MakeList(nil, Bool(false))
    v := MakeList(nil, Bool(false))
    if !RecEq(u, v) {
        t.Error("RecEq has problem over Bool comp.")
    }
}

func TestReqEqListName(t *testing.T) {
    u := MakeList(nil, Name("a"))
    v := MakeList(nil, Name("a"))
    if !RecEq(u, v) {
        t.Error("RecEq has problem over Bool comp.")
    }
}

func TestReqEqNestedListName(t *testing.T) {
    u := MakeList(nil, MakeList(nil, Name("a")), MakeList(nil, Name("b")))
    v := MakeList(nil, MakeList(nil, Name("a")), MakeList(nil, Name("b")))
    if !RecEq(u, v) {
        t.Error("RecEq has problem over Bool comp.")
    }
}


func TestReqEqNestedListName2(t *testing.T) {
    u := MakeList(nil, MakeList(nil, Name("a")), MakeList(nil, Name("b")))
    v := Cons(Cons(Name("a"), nil), Cons(Cons(Name("b"), nil), nil))
    if !RecEq(u, v) {
        t.Error("RecEq has problem over Bool comp.")
    }
}

func TestReqEqNestedList2(t *testing.T) {
    u := MakeList(nil, SFQuote, MakeList(nil, Name("a")), MakeList(nil, Name("b")))
    v := Cons(SFQuote,
              Cons(
                   Cons(Name("a"), nil),
                   Cons(Cons(Name("b"), nil), nil)))
    if !RecEq(u, v) {
        t.Error("RecEq / MakeList has problem over Quote?.")
    }
}

