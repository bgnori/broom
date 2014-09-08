package broom

import (
	//"fmt"
	"strings"
	"testing"
)

func TestNullImput(t *testing.T) {
	buf := NewBuffered(strings.NewReader(""))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderInt(t *testing.T) {
	buf := NewBuffered(strings.NewReader("42"))
	reader := NewReader(buf)
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderSymbol(t *testing.T) {
	buf := NewBuffered(strings.NewReader("a"))
	reader := NewReader(buf)
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderEmptyList(t *testing.T) {
	buf := NewBuffered(strings.NewReader("()"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderEmptyArr(t *testing.T) {
	buf := NewBuffered(strings.NewReader("[]"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_BRACKET {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_BRACKET {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderEmptyMap(t *testing.T) {
	buf := NewBuffered(strings.NewReader("{}"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_BRACE {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_BRACE {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderSomeList(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a b (c d) e)"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderSomeMap(t *testing.T) {
	buf := NewBuffered(strings.NewReader("{1 \"one\" \"two\" \"二\" 3 \"III\"}"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_BRACE {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_STRING {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_STRING {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_STRING {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_STRING {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_BRACE {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderString(t *testing.T) {
	buf := NewBuffered(strings.NewReader("\"abc\""))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_STRING || tkn.v != "abc" {
		t.Error("bad token id.")
		println(tkn.id)
		println(tkn.v)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderStringWithEscape(t *testing.T) {
	buf := NewBuffered(strings.NewReader("\"a\\\"bc\""))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_STRING || tkn.v != "a\"bc" {
		t.Error("bad token id.")
		println(tkn.id)
		println(tkn.v)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderStringWithEscape2(t *testing.T) {
	buf := NewBuffered(strings.NewReader("\"a\\\\bc\""))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_STRING || tkn.v != "a\\bc" {
		t.Error("bad token id.")
		println(tkn.id)
		println(tkn.v)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderSharp(t *testing.T) {
	buf := NewBuffered(strings.NewReader("#"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_SHARP {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderArr(t *testing.T) {
	buf := NewBuffered(strings.NewReader("[1 2 3]"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_BRACKET {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_INT {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_RIGHT_BRACKET {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestMakeInt(t *testing.T) {
	buf := NewBuffered(strings.NewReader("42"))
	expr := BuildSExpr(buf)[0]
	if expr != 42 {
		t.Error("42 is expected")
	}
}

func TestMakeSymbol(t *testing.T) {
	buf := NewBuffered(strings.NewReader("a"))
	expr := BuildSExpr(buf)[0]
	if !sym("a").Eq(expr) {
		t.Error("'a is expected")
	}
}

func TestMakeEmptyList(t *testing.T) {
	buf := NewBuffered(strings.NewReader("()"))
	expr := BuildSExpr(buf)[0]

	if expr != nil {
		t.Error("nil is expected")
	}
}

func TestMakeSomeList(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(1 2 3)"))
	expr := BuildSExpr(buf)[0]

	if !Eq(List(1, 2, 3), expr) {
		t.Error("(1 2 3) is expected")
	}
}

func TestMakeNestedList(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a b (c d) e)"))
	expr := BuildSExpr(buf)[0]

	if !Eq(List(sym("a"), sym("b"), List(sym("c"), sym("d")), sym("e")), expr) {
		t.Error("(a b (c d) e) is expected")
	}
}

func TestMakeSomeArray(t *testing.T) {
	buf := NewBuffered(strings.NewReader("[1 2 3]"))
	expr := BuildSExpr(buf)[0]

	if Eq([]int{1, 2, 3}, expr) {
		t.Error("[1 2 3] is expected")
	}
}

func TestMakeSomeMap(t *testing.T) {
	buf := NewBuffered(strings.NewReader("{1 \"one\" \"two\" \"二\" 3 \"III\"}"))
	expr := BuildSExpr(buf)[0]

	if !Eq(map[Value]Value{1: "one", "two": "二", 3: "III"}, expr) {
		t.Error("{1:\"one\" \"two\":\"二\" 3:\"III\"} is expected")
		DumpMap(expr)
	}
}
