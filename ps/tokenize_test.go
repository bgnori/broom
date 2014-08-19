package ps

import (
    "strings"
    "testing"
)


func TestTokenizeOneTokenEmptyString(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader(""))
    tkn := tknz.OneToken()
    if tkn.tkid != TOKEN_ENDOFINPUT {
        t.Error("bad status.")
    }
}

func TestTokenizeOneTokenInt(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("42"))
    tkn := tknz.OneToken()
    if tkn == nil {
        t.Error("nil token")
    }
    if tkn.start != 0 {
        t.Error("bad start index.")
    }
    if tkn.end != 2 {
        t.Error("bad end index.")
    }
    if tkn.tkid != TOKEN_CHUNK{
        t.Error("bad token id .")
    }
}

func TestTokenizeOneTokenLeftParen(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("("))
    tkn := tknz.OneToken()
    if tkn == nil {
        t.Error("nil token")
    }
    if tkn.start != 0 {
        t.Error("bad start index.")
    }
    if tkn.end != 1 {
        t.Error("bad end index.")
    }
    if tkn.tkid != TOKEN_LEFTPAREN{
        t.Error("bad token id .")
    }
}

func TestTokenizeOneTokenMulti(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("(123)"))
    tkn := tknz.OneToken()
    if tkn == nil {
        t.Error("nil token")
    }
    if tkn.start != 0 {
        t.Error("bad start index.")
    }
    if tkn.end != 1 {
        t.Error("bad end index.")
    }
    if tkn.tkid != TOKEN_LEFTPAREN{
        t.Error("bad token id .")
    }
    tkn = tknz.OneToken()
    if tkn == nil {
        t.Error("nil token")
    }
    if tkn.start != 1 {
        t.Error("bad start index.")
    }
    if tkn.end != 4 {
        t.Error("bad end index.")
    }
    if tkn.tkid != TOKEN_CHUNK{
        t.Error("bad token id .")
    }
    tkn = tknz.OneToken()
    if tkn == nil {
        t.Error("nil token")
    }
    if tkn.start != 4 {
        t.Error("bad start index.")
    }
    if tkn.end != 5 {
        t.Error("bad end index.")
    }
    if tkn.tkid != TOKEN_RIGHTPAREN {
        t.Error("bad token id .")
    }
}




func TestTokenizeEmptyString(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader(""))
    xs := ch2xs(tknz.Tokenize())
    if len(xs) != 1 {
        t.Error("just TOKEN_ENDOFINPUT should be there")
    }
}

func TestTokenizeNumber(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("123"))
    xs := ch2xs(tknz.Tokenize())
    if len(xs) != 2 {
        t.Error("an item and TOKEN_ENDOFINPUT are expected")
    }
    if xs[0].value != Int(123) {
        t.Error("unexpected content")
    }
}

func TestTokenizeUnicode(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("あいう"))
    xs := ch2xs(tknz.Tokenize())
    if len(xs) != 2 {
        t.Error("an item and TOKEN_ENDOFINPUT are expected")
    }
    if xs[0].value != String("あいう") {
        t.Error("unexpected content")
    }
}

func TestTokenizeParens(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("(() . ())"))
    xs := ch2xs(tknz.Tokenize())
    if len(xs) != 8 {
        t.Error("7 items and TOKEN_ENDOFINPUT are expected")
        println(xs)
    }
    if xs[3].tkid != TOKEN_DOT {
        t.Error("unexpected token type for dot")
    }
    if xs[4].tkid == TOKEN_DOT {
        t.Error("unexpected token type for dot")
    }
}

func TestTokenizeCons(t *testing.T) {
    tknz := NewTokenizer(strings.NewReader("(あいう . 123)"))
    xs := ch2xs(tknz.Tokenize())
    if len(xs) != 6 {
        t.Error("5 items and TOKEN_ENDOFINPUT are expected")
    }
    if xs[1].value != String("あいう") {
        t.Error("unexpected content")
    }
    if xs[3].value != Int(123) {
        t.Error("unexpected content")
    }
    if xs[2].tkid != TOKEN_DOT {
        t.Error("unexpected token type for dot")
    }
}
