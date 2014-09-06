package broom

import (
    "strings"
    "testing"
)

func TestNullImput(t *testing.T) {
    buf := NewBuffered(strings.NewReader(""))
    reader := NewReader(buf)

    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderChunk(t *testing.T) {
    buf := NewBuffered(strings.NewReader("42"))
    reader := NewReader(buf)
    if tkn := reader.Read(); tkn.id != TOKEN_CHUNK {
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderEmptyList(t *testing.T) {
    buf := NewBuffered(strings.NewReader("()"))
    reader := NewReader(buf)

    if tkn := reader.Read() ; tkn.id != TOKEN_LEFTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_RIGHTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderSomeList(t *testing.T) {
    buf := NewBuffered(strings.NewReader("(a b (c d) e)"))
    reader := NewReader(buf)

    if tkn := reader.Read() ; tkn.id != TOKEN_LEFTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_CHUNK{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_CHUNK{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_LEFTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_CHUNK{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_CHUNK{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_RIGHTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_CHUNK{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read() ; tkn.id != TOKEN_RIGHTPAREN{
        t.Error("bad token id.")
        println(tkn.id)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderString(t *testing.T) {
    buf := NewBuffered(strings.NewReader("\"abc\""))
    reader := NewReader(buf)

    if tkn := reader.Read() ; tkn.id != TOKEN_STRING || tkn.v !="abc" {
        t.Error("bad token id.")
        println(tkn.id)
        println(tkn.v)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderStringWithEscape(t *testing.T) {
    buf := NewBuffered(strings.NewReader("\"a\\\"bc\""))
    reader := NewReader(buf)

    if tkn := reader.Read() ; tkn.id != TOKEN_STRING || tkn.v !="a\"bc" {
        t.Error("bad token id.")
        println(tkn.id)
        println(tkn.v)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}

func TestReaderStringWithEscape2(t *testing.T) {
    buf := NewBuffered(strings.NewReader("\"a\\\\bc\""))
    reader := NewReader(buf)

    if tkn := reader.Read() ; tkn.id != TOKEN_STRING || tkn.v !="a\\bc" {
        t.Error("bad token id.")
        println(tkn.id)
        println(tkn.v)
    }
    if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT{
        t.Error("bad token id.")
        println(tkn.id)
    }
}
