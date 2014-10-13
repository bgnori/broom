package broom

import (
	"fmt"
	"testing"
)

func Test_ConsCarCdr(t *testing.T) {
	var v interface{}
	v = Cons(1, Cons(2, nil))
	if v == nil {
		t.Error("(null? v) must be null? false.")
	}
	p, ok := v.(List)
	if!ok  {
		t.Error("(pair? v) must be true.")
	}
	if !isNumber(Car(p)) {
		t.Error("(car v) must be number.")
	}
	if u, ok := Car(p).(int); ok && u == 1 {
	} else {
		t.Error("(car v) must be 1.")
	}
	if !isNumber(Car(Cdr(p))) {
		t.Error("(cdr v) must be number.")
	}
	if u, ok := Car(Cdr(p)).(int); ok && u == 2 {
	} else {
		t.Error("(cdr v) must be 2.")
	}
}

func Test_ConsNilNil(t *testing.T) {
	p := Cons(nil, nil)
	if _, ok := p.(List) ; !ok  {
		t.Error("(pair? xs) must be true.")
	}
}

func Test_ListNil(t *testing.T) {
	xs := Slice2List()
	if nil != xs {
		t.Error("xs must be null, i.e. '()")
	}
	if _, ok := xs.(List) ; ok  {
		t.Error("(pair? '()) must be false")
	}
}

func Test_isList_01(t *testing.T) {
	//(list? '(a b c)) =) #t
	xs := Cons(sym("A"), Cons(sym("B"), Cons(sym("C"), nil)))
	if !isList(xs) {
		t.Error("expect that (list? '(a b c)) =) #t")
	}
}

func Test_isList_02(t *testing.T) {
	//(list? '()) =) #t
	if !isList(nil) {
		t.Error("expect that (list? '()) =) #t")
	}
}

func xTest_isList_03(t *testing.T) {
	//(list? '(a . b)) =) #f
	// no support for improper list.
}

func Test_isList_04(t *testing.T) {
	//(let ((x (list 'a)))
	//  (set-cdr! x x)
	//  (list? x))
}

func Test_Length_01(t *testing.T) {
	//(length '(a b c)) =) 3
	xs := Cons(sym("a"), Cons(sym("b"), Cons(sym("c"), nil)))
	if Length(xs) != 3 {
		t.Error("expected 3 for (length '(a b c))")
	}
}

func Test_Length_02(t *testing.T) {
	//(length '(a (b) (c d e))) =) 3
	xs := Slice2List(sym("a"), Slice2List(sym("b")), Slice2List(sym("c"), sym("d"), sym("e")))
	if Length(xs) != 3 {
		t.Error("expected 3 for (length '(a (b) (c d e))) ")
	}
}

func Test_Length_03(t *testing.T) {
	//(length '()) =) 0
	if Length(nil) != 0 {
		t.Error("expected 0 for (length '()) ")
	}

}

func Test_String(t *testing.T) {
	xs := Slice2List(sym("a"), Slice2List(sym("b")), Slice2List(sym("c"), sym("d"), sym("e")))
	if fmt.Sprint(xs) != "(a (b) (c d e))" {
		t.Error("expected: (a (b) (c d e))")
		t.Error("actually:", fmt.Sprint(xs))
	}
}

func Test_Append(t *testing.T) {
	xs := Slice2List(sym("a"), sym("b"), sym("c"))
	ys := Slice2List(sym("d"), sym("e"), sym("f"))
	got := Append(xs, ys)
	if !Eq(got, Slice2List(sym("a"), sym("b"), sym("c"), sym("d"), sym("e"), sym("f"))) {
		t.Error("Append not working right")
		println(got)
	}
}
