package broom

import (
    "strings"
    "testing"
)

func TestNullImput(t *testing.T) {

}

func TestReader(t *testing.T) {
    buf := NewBuffered(strings.NewReader("42"))
    reader := NewReader(buf)
    tkn := reader.Read()

    if tkn.id != TOKEN_CHUNK {
        t.Error("bad token id.")
        println(tkn.id)
    }
}


