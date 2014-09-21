package broom

import (
	"reflect"
)

type Package struct {
	objects map[string]reflect.Value
}

func (p *Package)register(name string, f interface{}) {
	v := reflect.ValueOf(f)
	p.objects[name] = v
}

func MakeReflectPackage() *Package {
	p := new(Package)
	p.objects = make(map[string]reflect.Value)

	p.register("reflect.Copy", reflect.Copy)
	p.register("reflect.TypeOf", reflect.TypeOf)
	p.register("reflect.DeepEqual", reflect.DeepEqual)

	p.register("reflect.ChanOf", reflect.ChanOf)
	p.register("reflect.MakeChan", reflect.MakeChan)

	p.register("reflect.MapOf", reflect.MapOf)
	p.register("reflect.MakeMap", reflect.MakeMap)

	p.register("reflect.SliceOf", reflect.SliceOf)
	p.register("reflect.MakeSlice", reflect.MakeSlice)

	p.register("reflect.PtrTo", reflect.PtrTo)
	p.register("reflect.Indirect", reflect.Indirect)
	p.register("reflect.New", reflect.New)
	//NewAt.
	p.register("reflect.Zero", reflect.Zero)

	p.register("reflect.Append", reflect.Append)
	p.register("reflect.AppendSlice", reflect.AppendSlice)

	p.register("reflect.Select", reflect.Select)

	p.register("reflect.KindFromString", func(name string) reflect.Value {
		switch name {
			case "Bool": return reflect.ValueOf(reflect.Bool)
			case "Int": return reflect.ValueOf(reflect.Int)
			case "Int8": return reflect.ValueOf(reflect.Int8)
			case "Int16": return reflect.ValueOf(reflect.Int16)
			case "Int32": return reflect.ValueOf(reflect.Int32)
			case "Int64": return reflect.ValueOf(reflect.Int64)
			case "Uint": return reflect.ValueOf(reflect.Uint)
			case "Uint8": return reflect.ValueOf(reflect.Uint8)
			case "Uint16": return reflect.ValueOf(reflect.Uint16)
			case "Uint32": return reflect.ValueOf(reflect.Uint32)
			case "Uint64": return reflect.ValueOf(reflect.Uint64)
			case "Uintptr": return reflect.ValueOf(reflect.Uintptr)
			case "Float32": return reflect.ValueOf(reflect.Float32)
			case "Float64": return reflect.ValueOf(reflect.Float64)
			case "Complex64": return reflect.ValueOf(reflect.Complex64)
			case "Complex128": return reflect.ValueOf(reflect.Complex128)
			case "Array": return reflect.ValueOf(reflect.Array)
			case "Chan": return reflect.ValueOf(reflect.Chan)
			case "Func": return reflect.ValueOf(reflect.Func)
			case "Interface": return reflect.ValueOf(reflect.Interface)
			case "Map": return reflect.ValueOf(reflect.Map)
			case "Ptr": return reflect.ValueOf(reflect.Ptr)
			case "Slice": return reflect.ValueOf(reflect.Slice)
			case "String": return reflect.ValueOf(reflect.String)
			case "Struct": return reflect.ValueOf(reflect.Struct)
			case "UnsafePointer": return reflect.ValueOf(reflect.UnsafePointer)
		default:
			panic("bad name: "+name)
		}
	})

	return p
}

func (p *Package) Query (name string) reflect.Value {
	return p.objects[name]
}
