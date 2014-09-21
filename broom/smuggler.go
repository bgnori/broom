package broom

import (
	"reflect"
)

type Package struct {
	Functions map[string]reflect.Value
}

func (p *Package)register(name string, f interface{}) {
	v := reflect.ValueOf(f)
	p.Functions[name] = v
}

func MakeReflectPackage() *Package {
	p := new(Package)
	p.Functions = make(map[string]reflect.Value)

	p.register("reflect.Copy", reflect.Copy)
	p.register("reflect.TypeOf", reflect.TypeOf)

	return p
}

func (p *Package) Query (name string) reflect.Value {
	return p.Functions[name]
}
