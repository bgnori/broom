package ps

import (
    "fmt"
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
    u, ok := xs[0].(Name)
    if !ok {
        t.Error("type Name is expected")
    }
    if u != Name("あいう") {
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

func TestParseQuoteAtom(t *testing.T) {
    xs := Parse(strings.NewReader("(quote a)"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    if !RecEq(u, Cons(SFQuote, Cons(Name("a"), nil))) {
        t.Error("Wrong Value.")
    }
}


func TestParseQuoteList(t *testing.T) {
    xs := Parse(strings.NewReader("(quote (a b))"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    if !RecEq(u, MakeList(nil, SFQuote, MakeList(nil, Name("a"), Name("b")))){
        t.Error("Wrong Value.")
    }
}

func TestParseIF(t *testing.T) {
    println("TestParseIF")
    xs := Parse(strings.NewReader("(if #f 1 2)"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    expected := MakeList(nil, SFIf, Bool(false), Int(1), Int(2))
    // (if #f a b)
    if !RecEq(u, expected) {
        fmt.Println(expected)
        fmt.Println(u)
        t.Error("does not match.")
    }
}

func TestParseIFName(t *testing.T) {
    println("TestParseIF")
    xs := Parse(strings.NewReader("(if #f a b)"))
    if len(xs) != 1 {
        t.Error("a item is expected")
    }
    u, ok := xs[0].(*Pair)
    if !ok {
        t.Error("type *Pair is expected")
    }
    expected := MakeList(nil, SFIf, Bool(false), Name("a"), Name("b"))
    // (if #f a b)
    if !RecEq(u, expected) {
        fmt.Println(expected)
        fmt.Println(u)
        t.Error("does not match.")
    }
}

