package broom

import (
	"fmt"
	"reflect"
	"testing"
)

/* fail. won't work
func Test_IsEmpty(t *testing.T) {
	var x Sequence
	x = nil
	if !x.IsEmpty() {
		t.Error("nil must be Empty ")
	}
	var y List
	y = nil
	if !y.IsEmpty() {
		t.Error("nil must be Empty ")
	}
}
*/

func Test_Kons(t *testing.T) {
	s := Kons(1, nil)
	if Length(s) != 1 {
		t.Errorf("expected 1 but got %d", Length(s))
	}
}


func Test_MakeFromSlice(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0 ; i< 10 ; i++ {
		xs[i] = 'a'+i
	}
	ys :=  MakeFromSlice(xs...)
	zs := make([]int, 0)
	for ; ys != nil; ys = ys.Rest() {
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", len(zs))
	}
}

func Test_Take5(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0 ; i< 10 ; i++ {
		xs[i] = 'a'+i
	}
	zs := make([]int, 0)
	for ys := Take(5, MakeFromSlice(xs...)) ; !ys.IsEmpty() ; ys = ys.Rest() {
		fmt.Println(ys, reflect.TypeOf(ys))
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 5 {
		t.Errorf("expected 5 items but there is %d", len(zs))
	}
}

func Test_Take20(t *testing.T) {
	xs := make([]interface{}, 10)
	for i := 0 ; i< 10 ; i++ {
		xs[i] = 'a'+i
	}
	zs := make([]int, 0)
	for ys := Take(20, MakeFromSlice(xs...)) ; !ys.IsEmpty() ; ys = ys.Rest() {
		fmt.Println(ys, reflect.TypeOf(ys))
		zs = append(zs, ys.First().(int))
	}
	if len(zs) != 10 {
		t.Errorf("expected 10 items but there is %d", len(zs))
	}
}




