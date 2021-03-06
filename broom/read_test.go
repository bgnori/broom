package broom

import (
	"fmt"
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

func TestReaderQuotedSymbol(t *testing.T) {
	buf := NewBuffered(strings.NewReader("'a"))
	reader := NewReader(buf)
	if tkn := reader.Read(); tkn.id != TOKEN_QUOTE {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderEmptySlice2List(t *testing.T) {
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

func TestReaderSemicolonCRLF(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;\r\n)"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SEMICOLON {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFLINE {
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

func TestReaderSemicolonCR(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;\r)"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SEMICOLON {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFLINE {
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

func TestReaderSemicolonLF(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;\n)"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SEMICOLON {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFLINE {
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

func TestReaderSemicolonEOS(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;"))
	reader := NewReader(buf)

	if tkn := reader.Read(); tkn.id != TOKEN_LEFT_PAREN {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SYMBOL {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_SEMICOLON {
		t.Error("bad token id.")
		println(tkn.id)
	}
	if tkn := reader.Read(); tkn.id != TOKEN_ENDOFINPUT {
		t.Error("bad token id.")
		println(tkn.id)
	}
}

func TestReaderSomeSlice2List(t *testing.T) {
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
	expr := BuildSExpr(buf)
	if expr != 42 {
		t.Error("42 is expected")
	}
}

func TestMakeFloat(t *testing.T) {
	buf := NewBuffered(strings.NewReader("2.7182818284"))
	expr := BuildSExpr(buf)
	if expr != 2.7182818284 {
		t.Error("2.7182818284 is expected")
	}
}

func TestMakeSymbol(t *testing.T) {
	buf := NewBuffered(strings.NewReader("a"))
	expr := BuildSExpr(buf)
	if sym("a") != expr {
		t.Error("a is expected")
	}
}

func TestMakeQuotedSymbol(t *testing.T) {
	buf := NewBuffered(strings.NewReader("'a"))
	expr := BuildSExpr(buf)
	if !Eq(Slice2List(sym("quote"), sym("a")), expr) {
		t.Error("'a is expected")
	}
}

func TestMakeEmptySlice2List(t *testing.T) {
	buf := NewBuffered(strings.NewReader("()"))
	expr := BuildSExpr(buf)

	if expr != nil {
		t.Errorf("nil is expected, but got %v", expr)
	}
}

func TestMakeQuotedEmptySlice2List(t *testing.T) {
	buf := NewBuffered(strings.NewReader("'()"))
	expr := BuildSExpr(buf)

	if !Eq(Slice2List(sym("quote"), Slice2List()), expr) {
		t.Errorf("'() is expected %v", expr)
	}
}

func TestMakeBadQuoteSlice2List(t *testing.T) {
	defer func() {
		recover()
	}()
	buf := NewBuffered(strings.NewReader("(')"))
	BuildSExpr(buf)
	t.Error("Should not be reached")
}

func TestMakeEmptyArray(t *testing.T) {
	buf := NewBuffered(strings.NewReader("[]"))
	expr := BuildSExpr(buf)

	if !Eq([]interface{}{}, expr) {
		t.Error("[] is expected")
	}
}

func TestMakeQuotedEmptyArray(t *testing.T) {
	buf := NewBuffered(strings.NewReader("'[]"))
	expr := BuildSExpr(buf)

	if !Eq(Slice2List(sym("quote"), []interface{}{}), expr) {
		t.Error("[] is expected")
	}
}

func TestMakeSomeSlice2List(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(1 2 3)"))
	expr := BuildSExpr(buf)

	if !Eq(Slice2List(1, 2, 3), expr) {
		t.Error("(1 2 3) is expected")
	}
}

func TestMakeNestedSlice2List(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a b (c d) e)"))
	expr := BuildSExpr(buf)

	if !Eq(Slice2List(sym("a"), sym("b"), Slice2List(sym("c"), sym("d")), sym("e")), expr) {
		t.Error("(a b (c d) e) is expected")
	}
}

func TestMakeSomeArray(t *testing.T) {
	buf := NewBuffered(strings.NewReader("[1 2 3]"))
	expr := BuildSExpr(buf)

	if !Eq([]interface{}{1, 2, 3}, expr) {
		t.Error("[1 2 3] is expected")
		fmt.Println(expr)
	}
}

func TestMakeSomeMap(t *testing.T) {
	buf := NewBuffered(strings.NewReader("{1 \"one\" \"two\" \"二\" 3 \"III\"}"))
	expr := BuildSExpr(buf)

	if !Eq(map[interface{}]interface{}{1: "one", "two": "二", 3: "III"}, expr) {
		t.Error("{1:\"one\" \"two\":\"二\" 3:\"III\"} is expected")
		DumpMap(expr)
	}
}

func TestMakeSemicolonCRLF(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;\r\n)"))
	expr := BuildSExpr(buf)
	if !Eq(Slice2List(sym("a")), expr) {
		t.Error("(a) is expected")
	}
}

func TestMakeSemicolonCRLFwithSomething(t *testing.T) {
	buf := NewBuffered(strings.NewReader("(a ;\r\n) (b c)"))
	reader := NewReader(buf)
	builder := NewSExprBuilder()
	seq, err := builder.Run(reader)
	if err != nil {
		t.Error("got error")
	}
	if !Eq(Slice2List(sym("a")), seq.items[0]) {
		t.Error("(a) is expected")
	}
	if !Eq(Slice2List(sym("b"), sym("c")), seq.items[1]) {
		t.Error("(b c) is expected")
	}
}
