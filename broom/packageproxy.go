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
	return p
}
