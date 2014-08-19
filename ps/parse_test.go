package ps

import (
    "strings"
    "testing"
)

func TestParse0Byte(t *testing.T) {
    xs := Parse(strings.NewReader(""))
    if len(xs) != 0 {
        t.Error("empty result is expected")
    }
}

func TestParseInt(t *testing.T) {
    xs := Parse(strings.NewReader("123"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(Int)
    if !ok {
        t.Error("type Int is expected")
    }
    if u != Int(123) {
        t.Error("Wrong Value.")
    }
}

func TestParseString(t *testing.T) {
    xs := Parse(strings.NewReader("あいう"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(String)
    if !ok {
        t.Error("type String is expected")
    }
    if u != String("あいう") {
        t.Error("Wrong Value.")
    }
}

func TestParseListNil(t *testing.T) {
    println("TestParsePairListNil")
    xs := Parse(strings.NewReader("(())"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    if !RecEq(u, Cons(nil, nil)) {
        t.Error("Wrong Value.")
    }
}

func TestParsePairNilNil(t *testing.T) {
    println("TestParsePairNilNil")
    xs := Parse(strings.NewReader("(() . ())"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    if !RecEq(u, Cons(nil, nil)) {
        t.Error("Wrong Value.")
    }
}


func TestParseList(t *testing.T) {
    xs := Parse(strings.NewReader("(1 2 3)"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    if !RecEq(u, Cons(Int(1), Cons(Int(2), Cons(Int(3), nil)))) {
        t.Error("Wrong Value.")
    }
}




