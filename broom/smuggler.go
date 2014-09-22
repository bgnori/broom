package broom

import (
	"reflect"
)

type Package struct {
	objects map[string]reflect.Value
	values map[reflect.Kind]reflect.Value
}

func (p *Package)register(name string, f interface{}) {
	v := reflect.ValueOf(f)
	p.objects[name] = v
}

func (p *Package)setupValues(){
	p.values = make(map[reflect.Kind]reflect.Value, reflect.UnsafePointer + 1)
	p.values[reflect.Bool] = reflect.Zero(reflect.TypeOf(true))
	p.values[reflect.Int] = reflect.Zero(reflect.TypeOf(0))
	p.values[reflect.Int8] = reflect.Zero(reflect.TypeOf(int8(0)))
	p.values[reflect.Int16] = reflect.Zero(reflect.TypeOf(int16(0)))
	p.values[reflect.Int32] = reflect.Zero(reflect.TypeOf(int32(0)))
	p.values[reflect.Int64] = reflect.Zero(reflect.TypeOf(int64(0)))
	p.values[reflect.Uint] = reflect.Zero(reflect.TypeOf(uint(0)))
	p.values[reflect.Uint8] = reflect.Zero(reflect.TypeOf(uint8(0)))
	p.values[reflect.Uint16] = reflect.Zero(reflect.TypeOf(uint16(0)))
	p.values[reflect.Uint32] = reflect.Zero(reflect.TypeOf(uint32(0)))
	p.values[reflect.Uint64] = reflect.Zero(reflect.TypeOf(uint64(0)))
	p.values[reflect.Uintptr] = reflect.Zero(reflect.TypeOf(uintptr(0)))
	p.values[reflect.Float32] = reflect.Zero(reflect.TypeOf(float32(0)))
	p.values[reflect.Float64] = reflect.Zero(reflect.TypeOf(float64(0)))
	p.values[reflect.Complex64] = reflect.Zero(reflect.TypeOf(complex64(0)))
	p.values[reflect.Complex128] = reflect.Zero(reflect.TypeOf(complex128(0)))
}

func MakeReflectPackage() *Package {
	p := new(Package)
	p.setupValues()
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

	p.register("reflect.KindFromString", func(name string) reflect.Kind{
		switch name {
			case "Bool": return reflect.Bool
			case "Int": return reflect.Int
			case "Int8": return reflect.Int8
			case "Int16": return reflect.Int16
			case "Int32": return reflect.Int32
			case "Int64": return reflect.Int64
			case "Uint": return reflect.Uint
			case "Uint8": return reflect.Uint8
			case "Uint16": return reflect.Uint16
			case "Uint32": return reflect.Uint32
			case "Uint64": return reflect.Uint64
			case "Uintptr": return reflect.Uintptr
			case "Float32": return reflect.Float32
			case "Float64": return reflect.Float64
			case "Complex64": return reflect.Complex64
			case "Complex128": return reflect.Complex128
			case "Array": return reflect.Array
			case "Chan": return reflect.Chan
			case "Func": return reflect.Func
			case "Interface": return reflect.Interface
			case "Map": return reflect.Map
			case "Ptr": return reflect.Ptr
			case "Slice": return reflect.Slice
			case "String": return reflect.String
			case "Struct": return reflect.Struct
			case "UnsafePointer": return reflect.UnsafePointer
		default:
			panic("bad name: "+name)
		}
	})
	p.register("reflect.ZeroValueOf", func(k reflect.Kind) reflect.Value {
		return p.values[k]
	})
	return p
}

func (p *Package) Query (name string) reflect.Value {
	return p.objects[name]
}
