package broom

import (
	"reflect"
)

type PackageProxy struct {
	name  string
	funcs map[string]reflect.Value
}

func (pp *PackageProxy) Name() string {
	return pp.name
}

func (pp *PackageProxy) Find(name string) (reflect.Value, bool) {
	val, ok := pp.funcs[name]
	return val, ok
}

func (pp *PackageProxy) register(name string, f interface{}) {
	v := reflect.ValueOf(f)
	pp.funcs[name] = v
}

func NewPackageProxy(name string) *PackageProxy {
	p := &PackageProxy{name: name}
	p.funcs = make(map[string]reflect.Value)
	return p
}

func MakeBroomPackage() *PackageProxy {
	p := NewPackageProxy("broom")
	p.register("Load", Load)
	p.register("Repl", Repl)

	//From sequence.go
	p.register("MakeFromSlice", MakeFromSlice)
	p.register("MakeFromChan", MakeFromChan)
	p.register("SeqString", SeqString)
	p.register("SeqTake", SeqTake)
	p.register("SeqDrop", SeqDrop)
	p.register("Seq2Slice", Seq2Slice)
	p.register("SeqAppend", SeqAppend)
	p.register("MakeSeqByAppend", MakeSeqByAppend)
	p.register("SeqFilter", SeqFilter)
	p.register("SeqEvens", SeqEvens)
	p.register("SeqOdds", SeqOdds)
	p.register("SeqZip2", SeqZip2)
	p.register("SeqEq", SeqEq)
	p.register("SeqRange", SeqRange)
	p.register("SeqReduce", SeqReduce)



	return p
}
